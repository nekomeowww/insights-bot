package openai

import (
	"context"
	"fmt"
	"math"

	"github.com/pandodao/tokenizer-go"
	"github.com/sashabaranov/go-openai"
)

type Client struct {
	OpenAIClient *openai.Client
}

func NewClient(apiSecret string) *Client {
	return &Client{
		OpenAIClient: openai.NewClient(apiSecret),
	}
}

// truncateContentBasedOnTokens 基于 token 计算的方式截断文本
func (c *Client) TruncateContentBasedOnTokens(textContent string) (string, error) {
	tokens, err := tokenizer.CalToken(textContent)
	if err != nil {
		return "", err
	}
	if tokens > 3900 {
		return string([]rune(textContent)[:int(math.Min(3900, float64(len([]rune(textContent)))))]), nil
	}

	return textContent, nil
}

// SplitContentBasedByTokenLimitations 基于 token 计算的方式分割文本
func (c *Client) SplitContentBasedByTokenLimitations(textContent string) ([]string, error) {
	slices := make([]string, 0)
	slices, err := appendSplitTextByTokenLimitations(slices, textContent)
	if err != nil {
		return make([]string, 0), err
	}

	return slices, nil
}

func appendSplitTextByTokenLimitations(slices []string, textContent string) ([]string, error) {
	tokens, err := tokenizer.CalToken(textContent)
	if err != nil {
		return make([]string, 0), err
	}
	if tokens > 3900 {
		sliceFrom := math.Min(3900, float64(len([]rune(textContent))))
		slices = append(slices, string([]rune(textContent)[:int(sliceFrom)]))
		return appendSplitTextByTokenLimitations(slices, string([]rune(textContent)[int(sliceFrom):]))
	}

	slices = append(slices, textContent)
	return slices, nil
}

// SummarizeWithQuestionsAsSimplifiedChinese 通过 OpenAI 的 Chat API 来为文章生成摘要和联想问题
func (c *Client) SummarizeWithQuestionsAsSimplifiedChinese(ctx context.Context, title, by, content string) (*openai.ChatCompletionResponse, error) {
	resp, err := c.OpenAIClient.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
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

	return &resp, nil
}

func (c *Client) SummarizeWithOneChatHistory(ctx context.Context, llmFriendlyChatHistory string) (*openai.ChatCompletionResponse, error) {
	resp, err := c.OpenAIClient.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
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

	return &resp, nil
}

func (c *Client) SummarizeWithChatHistories(ctx context.Context, llmFriendlyChatHistories string) (*openai.ChatCompletionResponse, error) {
	resp, err := c.OpenAIClient.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role: openai.ChatMessageRoleSystem,
					Content: "" +
						"你是我的聊天记录总结和回顾主力。我将为你提供一份不完整的、在过去一个小时中" +
						"的、包含了人物名称、人物用户名、消息发送时间、消息内容等信息的聊天记录，这" +
						"些聊天记录条目每条一行，我需要你通过这些聊天记录总结并以 Markdown 的语法" +
						"输出一个列表，这个列表中包含了你发现的聊天主题，参与人和内容。不需要输出总" +
						"结的大标题。" +
						"",
				},
				{
					Role: openai.ChatMessageRoleUser,
					Content: "" +
						"聊天记录：\n" +
						llmFriendlyChatHistories + "\n" +
						"请你帮我总结一下。",
				},
			},
		},
	)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}
