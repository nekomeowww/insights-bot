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
	Point              string  `json:"point"`
	CriticalMessageIDs []int64 `json:"criticalMsgIds"`
}

type ChatHistorySummarizationOutputs struct {
	TopicName                        string                                       `json:"topicName"`
	SinceMsgID                       int64                                        `json:"sinceMsgId"`
	ParticipantsNamesWithoutUsername []string                                     `json:"participantsNamesWithoutUsername"`
	Discussion                       []*ChatHistorySummarizationOutputsDiscussion `json:"discussion"`
	Conclusion                       string                                       `json:"conclusion"`
}

var ChatHistorySummarizationPrompt = lo.Must(template.New(uuid.New().String()).Parse("" +
	`聊天记录："""
{{ .ChatHistory }}
"""

你是我的聊天记录总结和回顾助理。上文是我提供的不完整的、在过去一段时间内、包含了人物、消息内容等信息的聊天记录，这些记录条目每条都以 msgId 为开头，你需要总结这些聊天记录，并在有结论的时候提供结论总结。请你使用下面的 JSON 格式进行输出，不需要提供额外的解释和说明，在 JSON 中 sinceMsgId 代表了话题开始的消息 ID，criticalMsgIds 代表了讨论过程中出现的关键消息（列出最多 5 条即可）："""
[{"topicName":"..","sinceMsgId":123456789,"participantsNamesWithoutUsername":[".."],"discussion":[{"point":"..","criticalMsgIds":[123456789]}],"conclusion":".."}]
"""`))
