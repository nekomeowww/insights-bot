package openai

import (
	"context"
	"fmt"

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

// SummarizeWithQuestionsAsSimplifiedChinese 通过 OpenAI 的 Chat API 来为文章生成摘要和联想问题
func (c *Client) SummarizeWithQuestionsAsSimplifiedChinese(title, by, content string) (*openai.ChatCompletionResponse, error) {
	resp, err := c.OpenAIClient.CreateChatCompletion(
		context.Background(),
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
