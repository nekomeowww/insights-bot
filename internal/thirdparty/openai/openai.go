package openai

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"unicode/utf8"

	"github.com/pkoukk/tiktoken-go"
	"github.com/sashabaranov/go-openai"
	"github.com/sourcegraph/conc/pool"
	"go.uber.org/fx"
	"go.uber.org/ratelimit"
	"go.uber.org/zap"

	"github.com/nekomeowww/insights-bot/internal/configs"
	"github.com/nekomeowww/insights-bot/internal/datastore"
	"github.com/nekomeowww/insights-bot/pkg/logger"
)

//counterfeiter:generate -o openaimock/mock_client.go --fake-name MockClient . Client
type Client interface {
	SplitContentBasedByTokenLimitations(textContent string, limits int) []string
	SummarizeAny(ctx context.Context, content string) (*openai.ChatCompletionResponse, error)
	SummarizeChatHistories(ctx context.Context, llmFriendlyChatHistories string) (*openai.ChatCompletionResponse, error)
	SummarizeOneChatHistory(ctx context.Context, llmFriendlyChatHistory string) (*openai.ChatCompletionResponse, error)
	SummarizeWithQuestionsAsSimplifiedChinese(ctx context.Context, title string, by string, content string) (*openai.ChatCompletionResponse, error)
	TruncateContentBasedOnTokens(textContent string, limits int) string
}

var _ Client = (*OpenAIClient)(nil)

type OpenAIClient struct {
	tiktokenEncoding *tiktoken.Tiktoken
	client           *openai.Client
	ent              *datastore.Ent
	logger           *logger.Logger
	limiter          ratelimit.Limiter
}

func parseOpenAIAPIHost(apiHost string) (string, error) {
	if !strings.HasPrefix(apiHost, "https://") && !strings.HasPrefix(apiHost, "http://") {
		apiHost = "http://" + apiHost
	}

	parsedURL, err := url.Parse(apiHost)
	if err != nil {
		return "", err
	}

	host := fmt.Sprintf("%s://%s", parsedURL.Scheme, parsedURL.Host)
	if host != "" {
		return host, nil
	}

	return "", fmt.Errorf("invalid API host: %s", apiHost)
}

type NewClientParams struct {
	fx.In

	Config *configs.Config
	Logger *logger.Logger
	Ent    *datastore.Ent
}

func NewClient() func(NewClientParams) (Client, error) {
	return func(params NewClientParams) (Client, error) {
		tokenizer, err := tiktoken.EncodingForModel(openai.GPT3Dot5Turbo)
		if err != nil {
			return nil, err
		}

		apiHost := params.Config.OpenAIAPIHost

		config := openai.DefaultConfig(params.Config.OpenAIAPISecret)
		if apiHost != "" {
			apiHost, err = parseOpenAIAPIHost(apiHost)
			if err != nil {
				return nil, err
			}

			config.BaseURL = fmt.Sprintf("%s/v1", apiHost)
		}

		client := openai.NewClientWithConfig(config)

		limiter := ratelimit.New(1)
		limiter.Take()

		p := pool.New().WithMaxGoroutines(10)
		p.Go(func() {

		})

		return &OpenAIClient{
			client:           client,
			tiktokenEncoding: tokenizer,
			ent:              params.Ent,
			logger:           params.Logger,
			limiter:          ratelimit.New(5),
		}, nil
	}
}

// truncateContentBasedOnTokens 基于 token 计算的方式截断文本。
func (c *OpenAIClient) TruncateContentBasedOnTokens(textContent string, limits int) string {
	tokens := c.tiktokenEncoding.Encode(textContent, nil, nil)
	if len(tokens) <= limits {
		return textContent
	}

	truncated := c.tiktokenEncoding.Decode(tokens[:limits])

	for len(truncated) > 0 {
		// 假设 textContent = "小溪河水清澈见底", Encode 结果为 "[31809,36117,103,31106,111,53610,80866,162,122,230,90070,11795,243]"
		// 当 limits = 4, 那么 tokens[:limits] = "[31809,36117,103,31106]", Decode 结果为 "小溪\xe6\xb2"
		// 这里的 \xe6\xb2 是一个不完整的 UTF-8 编码，无法正确解析为一个完整的字符。下面得代码处理这种情况把它去掉。
		r, size := utf8.DecodeLastRuneInString(truncated)
		if r != utf8.RuneError {
			break
		}
		truncated = truncated[:len(truncated)-size]
	}

	return truncated
}

// SplitContentBasedByTokenLimitations 基于 token 计算的方式分割文本。
func (c *OpenAIClient) SplitContentBasedByTokenLimitations(textContent string, limits int) []string {
	slices := make([]string, 0)

	for {
		s := c.TruncateContentBasedOnTokens(textContent, limits)
		slices = append(slices, s)
		textContent = textContent[len(s):]

		if textContent == "" {
			break
		}
	}

	return slices
}

// SummarizeWithQuestionsAsSimplifiedChinese 通过 OpenAI 的 Chat API 来为文章生成摘要和联想问题。
func (c *OpenAIClient) SummarizeWithQuestionsAsSimplifiedChinese(ctx context.Context, title, by, content string) (*openai.ChatCompletionResponse, error) {
	c.limiter.Take()

	resp, err := c.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo16K,
			Messages: []openai.ChatCompletionMessage{
				{
					Role: openai.ChatMessageRoleSystem,
					Content: "" +
						"你是我的网页文章阅读助理。我将为你提供文章的标题、作" +
						"者、所抓取的网页中的正文等信息，然后你将对文章做出总结。\n请你在总结时满足以下要求：" +
						"1. 首先如果文章的标题不是中文的请依据上下文将标题信达雅的翻译为简体中文并放在第一行" +
						"2. 然后从我提供的文章信息中总结出一个三百字以内的文章的摘要" +
						"3. 最后，你将利用你已有的知识和经验，对我提供的文章信息提出 3 个具有创造性和发散思维的问题" +
						"4. 请用简体中文进行回复" +
						"最终你回复的消息格式应像这个例句一样（例句中的双花括号为需要替换的内容）：\n" +
						"{{简体中文标题，可省略}}\n\n摘要：{{文章的摘要}}\n\n关联提问：\n1. {{关联提问 1}}\n2. {{关联提问 2}}\n2. {{关联提问 3}}",
				},
				{
					Role: openai.ChatMessageRoleUser,
					Content: "" +
						"我的第一个要求相关的信息如下：" +
						fmt.Sprintf("文章标题：%s；", title) +
						fmt.Sprintf("文章作者：%s；", by) +
						fmt.Sprintf("文章正文：%s；", content) +
						"接下来请你完成我所要求的任务。",
				},
			},
		},
	)
	if err != nil {
		return nil, err
	}

	err = c.ent.MetricOpenAIChatCompletionTokenUsage.
		Create().
		SetPromptOperation("Summarize With Questions As Simplified Chinese").
		SetPromptTokenUsage(resp.Usage.PromptTokens).
		SetCompletionTokenUsage(resp.Usage.CompletionTokens).
		SetTotalTokenUsage(resp.Usage.TotalTokens).
		Exec(ctx)
	if err != nil {
		c.logger.Error("failed to create metric openai chat completion token usage", zap.Error(err),
			zap.String("prompt_operation", "Summarize With Questions As Simplified Chinese"),
			zap.Int("prompt_token_usage", resp.Usage.PromptTokens),
			zap.Int("completion_token_usage", resp.Usage.CompletionTokens),
			zap.Int("total_token_usage", resp.Usage.TotalTokens),
		)
	}

	return &resp, nil
}

func (c *OpenAIClient) SummarizeOneChatHistory(ctx context.Context, llmFriendlyChatHistory string) (*openai.ChatCompletionResponse, error) {
	c.limiter.Take()

	resp, err := c.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo16K,
			Messages: []openai.ChatCompletionMessage{
				{
					Role: openai.ChatMessageRoleSystem,
					Content: "" +
						"你是我的聊天消息总结助手。我将为你提供一条包含了人物名称、人物用户名、消息" +
						"发送时间、消息内容等信息的消息，因为这条聊天消息有些过长了，我需要你帮我总" +
						"结一下这条消息说了什么。最好一句话概括，如果这条消息有标题的话你可以直接返" +
						"回标题。" +
						"",
				},
				{
					Role: openai.ChatMessageRoleUser,
					Content: "" +
						"消息：\n" +
						llmFriendlyChatHistory + "\n" +
						"请你帮我总结一下。",
				},
			},
		},
	)
	if err != nil {
		return nil, err
	}

	err = c.ent.MetricOpenAIChatCompletionTokenUsage.
		Create().
		SetPromptOperation("Summarize One Chat History").
		SetPromptTokenUsage(resp.Usage.PromptTokens).
		SetCompletionTokenUsage(resp.Usage.CompletionTokens).
		SetTotalTokenUsage(resp.Usage.TotalTokens).
		Exec(ctx)
	if err != nil {
		c.logger.Error("failed to create metric openai chat completion token usage",
			zap.Error(err),
			zap.String("prompt_operation", "Summarize One Chat History"),
			zap.Int("prompt_token_usage", resp.Usage.PromptTokens),
			zap.Int("completion_token_usage", resp.Usage.CompletionTokens),
			zap.Int("total_token_usage", resp.Usage.TotalTokens),
		)
	}

	return &resp, nil
}

// SummarizeAny 通过 OpenAI 的 Chat API 来为任意内容生成摘要。
func (c *OpenAIClient) SummarizeAny(ctx context.Context, content string) (*openai.ChatCompletionResponse, error) {
	c.limiter.Take()

	sb := new(strings.Builder)

	err := AnySummarizationPrompt.Execute(sb, AnySummarizationInputs{
		Content: content,
	})
	if err != nil {
		return nil, err
	}

	resp, err := c.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo16K,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: sb.String(),
				},
			},
		},
	)
	if err != nil {
		return nil, err
	}

	err = c.ent.MetricOpenAIChatCompletionTokenUsage.
		Create().
		SetPromptOperation("Summarize Any").
		SetPromptTokenUsage(resp.Usage.PromptTokens).
		SetCompletionTokenUsage(resp.Usage.CompletionTokens).
		SetTotalTokenUsage(resp.Usage.TotalTokens).
		Exec(ctx)
	if err != nil {
		c.logger.Error("failed to create metric openai chat completion token usage",
			zap.Error(err),
			zap.String("prompt_operation", "Summarize Any"),
			zap.Int("prompt_token_usage", resp.Usage.PromptTokens),
			zap.Int("completion_token_usage", resp.Usage.CompletionTokens),
			zap.Int("total_token_usage", resp.Usage.TotalTokens),
		)
	}

	return &resp, nil
}

func (c *OpenAIClient) SummarizeChatHistories(ctx context.Context, llmFriendlyChatHistories string) (*openai.ChatCompletionResponse, error) {
	c.limiter.Take()

	sb := new(strings.Builder)

	err := ChatHistorySummarizationPrompt.Execute(
		sb,
		NewChatHistorySummarizationPromptInputs(
			llmFriendlyChatHistories,
			"Simplified Chinese",
		),
	)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo16K,
			Messages: []openai.ChatCompletionMessage{{
				Role:    openai.ChatMessageRoleSystem,
				Content: sb.String(),
			}},
		},
	)
	if err != nil {
		return nil, err
	}

	err = c.ent.MetricOpenAIChatCompletionTokenUsage.
		Create().
		SetPromptOperation("Summarize Chat Histories").
		SetPromptTokenUsage(resp.Usage.PromptTokens).
		SetCompletionTokenUsage(resp.Usage.CompletionTokens).
		SetTotalTokenUsage(resp.Usage.TotalTokens).
		Exec(ctx)
	if err != nil {
		c.logger.Error("failed to create metric openai chat completion token usage",
			zap.Error(err),
			zap.String("prompt_operation", "Summarize Chat Histories"),
			zap.Int("prompt_token_usage", resp.Usage.PromptTokens),
			zap.Int("completion_token_usage", resp.Usage.CompletionTokens),
			zap.Int("total_token_usage", resp.Usage.TotalTokens),
		)
	}

	return &resp, nil
}
