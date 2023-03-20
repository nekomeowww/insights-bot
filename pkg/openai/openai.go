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
						"我想让你担任我的网页文章阅读助理。我将为你提供文章的标题、作" +
						"者、所抓取的网页中的正文等信息，你将从我提供的文章信息中总结" +
						"出一个简体中文版本的文章的摘要，并且，你将利用你已有的知识和" +
						"经验，对我提供的文章信息提出至少 3 个具有创造性和发散思维的" +
						"问题，在所提出的 3 个问题的前面使用「关联提问：」作为" +
						"开头，且在问题与问题都使用换行进行分割。",
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
