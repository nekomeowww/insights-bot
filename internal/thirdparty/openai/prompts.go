package openai

import (
	"text/template"

	"github.com/samber/lo"
)

type AnySummarizationInputs struct {
	Content string
}

var AnySummarizationPrompt = lo.Must(template.New("anything summarization prompt").Parse("" +
	`内容：{{ .Content }}
你是我的总结助手。我将为你提供一段话，我需要你在不丢失原文主旨和情感、不做更多的解释和说明的情况下帮我用不超过100字总结一下这段话说了什么。`))

type ChatHistorySummarizationPromptInputs struct {
	ChatHistory string
	Language    string
}

func NewChatHistorySummarizationPromptInputs(chatHistory string, language string) *ChatHistorySummarizationPromptInputs {
	return &ChatHistorySummarizationPromptInputs{
		ChatHistory: chatHistory,
		Language:    lo.Ternary(language != "", language, "Simplified Chinese"),
	}
}

type ChatHistorySummarizationOutputsDiscussion struct {
	Point  string  `json:"point"`
	KeyIDs []int64 `json:"keyIds"`
}

type ChatHistorySummarizationOutputs struct {
	TopicName    string                                       `json:"topicName"`
	SinceID      int64                                        `json:"sinceId"`
	Participants []string                                     `json:"participants"`
	Discussion   []*ChatHistorySummarizationOutputsDiscussion `json:"discussion"`
	Conclusion   string                                       `json:"conclusion"`
}

var ChatHistorySummarizationPrompt = lo.Must(template.New("chat histories summarization prompt").Parse("" +
	`聊天记录："""
{{ .ChatHistory }}
"""

你是一位擅长总结文档和对话大纲的专家。请仔细阅读提供的聊天记录，识别出其中讨论的 1-10 个不同主题。

请按照以下 JSON Schema 格式输出结果，并使用{{ .Language }}语言："""
{"$schema":"http://json-schema.org/draft-07/schema#","title":"聊天记录总结模式","type":"array","items":{"type":"object","properties":{"topicName":{"type":"string","description":"在聊天记录中讨论的话题的标题或简短标题。"},"sinceId":{"type":"number","description":"该话题最初开始的消息ID。"},"participants":{"type":"array","description":"参与该话题讨论的用户名称列表。","items":{"type":"string"}},"discussion":{"type":"array","description":"该话题讨论期间的要点列表。","items":{"type":"object","properties":{"point":{"type":"string","description":"在该话题中谈论、表达、提到或讨论的关键点。"},"keyIds":{"type":"array","description":"包含该关键点的消息ID列表。","items":{"type":"number"}}},"required":["point","keyIds"]},"minItems": 1,"maxItems": 5},"conclusion":{"type":"string","description":"该话题的结论，可选。"}},"required":["topicName","sinceId","participants","discussion"]}}
"""

例如："""
[{"topicName":"最重要的话题1","sinceId":123456789,"participants":["张三","李四"],"discussion":[{"point":"最相关的关键点","keyIds":[123456789,987654321]}],"conclusion":"可选的简短结论"},{"topicName":"最重要的话题2","sinceId":987654321,"participants":["王五","赵六"],"discussion":[{"point":"最相关的关键点","keyIds":[987654321]}],"conclusion":"可选的简短结论"}]
"""

请注意，话题可能会并行讨论，所以请考虑整个聊天记录中出现的相关关键词。从聊天记录中总结出不同的话题。对于每个话题，提取 1-5 个最相关的要点和关键消息ID。请保持简洁，专注于每个话题的核心要点。`))
