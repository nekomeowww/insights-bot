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

Read through the provided chat history and identify all distinct discussion topics that took place. Summarize each topic by extracting the most relevant points and key message IDs.

Output topics in the following JSON format in language {{ .Language }}:"""
[{"topicName":"Topic 1 Name","sinceId":123456789,"participantsNamesWithoutUsername":["John","Mary"],"discussion":[{"point":"Key point 1","keyIds":[123456789,987654321]},{"point":"Key point 2","keyIds":[456789123]}],"conclusion":"Optional conclusion"},{"topicName":"Topic 2 Name","sinceId":987654321,"participantsNamesWithoutUsername":["Bob","Alice"],"discussion":[{"point":"Key point 1","keyIds":[987654321]}],"conclusion":"Optional conclusion"}]
"""

Extract all distinct topics discussed in the chat history, including the topic name, starting message ID, participating members, key discussion points, and optional conclusion. Only include the most relevant points and message IDs. Be concise yet comprehensive.`))
