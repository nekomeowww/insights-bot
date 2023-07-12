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

You are a expert in summarizing the refined outlines from documents and dialogues. Please read through the provided chat history and identify the 1-10 distinct discussion topics that discussed and talked about. Summarize only those key topics by extracting the most relevant points and message IDs.

Output topics in the following JSON format in language {{ .Language }}:"""
[{"topicName":"Most Important Topic 1","sinceId":123456789,"participantsNamesWithoutUsername":["John","Mary"],"discussion":[{"point":"Most relevant key point","keyIds":[123456789,987654321]}],"conclusion":"Optional brief conclusion"},{"topicName":"Most Important Topic 2","sinceId":987654321,"participantsNamesWithoutUsername":["Bob","Alice"],"discussion":[{"point":"Most relevant key point","keyIds":[987654321]}],"conclusion":"Optional brief conclusion"}]
"""

Only summarize the 1-10 distinct topics from the chat history. For each topic, extract just the most relevant point and key message IDs. Be very concise and focused on the key essence of each topic.`))
