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
	Point     string `json:"point"`
	MessageID int64  `json:"messageId"`
}

type ChatHistorySummarizationOutputs struct {
	TopicName                        string                                      `json:"topicName"`
	ParticipantsNamesWithoutUsername []string                                    `json:"participantsNamesWithoutUsername"`
	Discussion                       []ChatHistorySummarizationOutputsDiscussion `json:"discussion"`
	Conclusion                       string                                      `json:"conclusion"`
}

var ChatHistorySummarizationPrompt = lo.Must(template.New(uuid.New().String()).Parse("" +
	`你是我的聊天记录总结和回顾助理。我将为你提供一份不完整的、在过去一个小时中的、包含了人物名称、人物用户名、消息发送时间、消息内容等信息的聊天记录，这些聊天记录条目每条一行，我需要你总结这些聊天记录，并在有结论的时候提供结论总结。
并请你使用下面的 JSON 格式进行输出，并且不需要提供额外的解释和说明。

输出格式："""
[
  {
    "topicName": "..",
    "participantsNamesWithoutUsername": [ "..", ".." ],
    "discussion": [ { "point": "..", "messageId": 123456789 }, { "point": "..", "messageId": 123456789 } ],
    "conclusion": ".."
  },
  {
    "topicName": "..",
    "participantsNamesWithoutUsername": [ "..", ".." ],
    "discussion": [ { "point": "..", "messageId": 123456789 }, { "point": "..", "messageId": 123456789 } ],
    "conclusion": ".."
  }
]
"""

聊天记录："""
{{ .ChatHistory }}
"""`))
