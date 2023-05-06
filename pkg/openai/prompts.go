package openai

import (
	"text/template"

	"github.com/google/uuid"
	"github.com/samber/lo"
)

type ChatHistorySummarizationPromptInputs struct {
	ChatHistory string
}

type ChatHistorySummarizationOutputsDiscussion struct {
	Point       string  `json:"point"`
	CriticalIDs []int64 `json:"criticalIds"`
}

type ChatHistorySummarizationOutputs struct {
	TopicName                        string                                       `json:"topicName"`
	SinceID                          int64                                        `json:"sinceId"`
	ParticipantsNamesWithoutUsername []string                                     `json:"participantsNamesWithoutUsername"`
	Discussion                       []*ChatHistorySummarizationOutputsDiscussion `json:"discussion"`
	Conclusion                       string                                       `json:"conclusion"`
}

var ChatHistorySummarizationPrompt = lo.Must(template.New(uuid.New().String()).Parse("" +
	`聊天记录："""
{{ .ChatHistory }}
"""

你是我的聊天记录总结和回顾助理。以上是一份聊天记录，每条消息以 msgId 开头，请总结这些聊天记录为1~5个话题，每个话题需包含以下字段 sinceId(话题开始的 msgId)、criticalIds(讨论过程中的关键 msgId，最多5条)和 conclusion(结论，若无明确结论则该字段为空)。请使用以下 JSON 格式输出，无需额外解释说明："""
[{"topicName":"..","sinceId":123456789,"participantsNamesWithoutUsername":[".."],"discussion":[{"point":"..","criticalIds":[123456789]}],"conclusion":".."}]"""`))
