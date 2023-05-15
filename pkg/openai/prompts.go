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
你是我的总结助手。我将为你提供一段话，我需要你在不丢失原文主旨和情感的情况下帮我用不超过100字总结一下这段话说了什么。`))

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
	TopicName                        string                                       `json:"topicName"`
	SinceID                          int64                                        `json:"sinceId"`
	ParticipantsNamesWithoutUsername []string                                     `json:"participantsNamesWithoutUsername"`
	Discussion                       []*ChatHistorySummarizationOutputsDiscussion `json:"discussion"`
	Conclusion                       string                                       `json:"conclusion"`
}

var ChatHistorySummarizationPrompt = lo.Must(template.New("chat histories summarization prompt").Parse("" +
	`Chat histories:"""
{{ .ChatHistory }}
"""

You are my chat histories summary and review assistant. Above are chat histories, each message starts with msgId, please summarize these chats as 1 to 5 topics, each topic should contain the following fields: sinceId (msgId at the beginning of the topic)、keyIds (key msgId in the discussion, max 5 msgIds) and conclusion (ignore this field if no clear conclusion). Please output as the following JSON format without additional explanation using {{ .Language }}:"""
[{"topicName":"..","sinceId":123456789,"participantsNamesWithoutUsername":[".."],"discussion":[{"point":"..","keyIds":[123456789]}],"conclusion":".."}]"""`))
