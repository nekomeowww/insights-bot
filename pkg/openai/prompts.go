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
	`你是我的聊天记录总结和回顾助理。我将为你提供一份不完整的、在过去一个小时中的、包含了人物名称、人物用户名、消息发送时间、消息内容等信息的聊天记录，这些聊天记录条目每条一行，我需要你总结这些聊天记录，并在有结论的时候提供结论总结。
请你使用下面的 JSON 格式进行输出，不需要提供额外的解释和说明。为了方便理解，下面的 JSON 中 sinceMsgId 代表了话题开始的消息 ID，而 criticalMsgIds 代表了讨论过程中出现的「关键消息」，因此你不需要罗列出所有的相关 criticalMsgIds。

输出时所使用的 JSON 格式："""
[
  { "topicName": "..", "sinceMsgId": 123456789, "participantsNamesWithoutUsername": [ "..", ".." ], "discussion": [ { "point": "..", "criticalMsgIds": [ 123456789, 123456789 ] }, { "point": "..", "criticalMsgIds": [ 123456789, 123456789 ] } ], "conclusion": ".." },
  { "topicName": "..", "sinceMsgId": 123456789, "participantsNamesWithoutUsername": [ "..", ".." ], "discussion": [ { "point": "..", "criticalMsgIds": [ 123456789, 123456789 ] }, { "point": "..", "criticalMsgIds": [ 123456789, 123456789 ] } ], "conclusion": ".." }
]
"""

聊天记录："""
{{ .ChatHistory }}
"""`))
