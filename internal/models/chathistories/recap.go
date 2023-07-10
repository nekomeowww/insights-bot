package chathistories

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"strconv"
	"strings"
	"time"

	"github.com/nekomeowww/fo"
	"github.com/samber/lo"
	goopenai "github.com/sashabaranov/go-openai"
	"go.uber.org/zap"

	"github.com/nekomeowww/insights-bot/internal/thirdparty/openai"
	"github.com/nekomeowww/insights-bot/pkg/bots/tgbot"
)

type RecapOutputTemplateInputs struct {
	ChatID string
	Recap  *openai.ChatHistorySummarizationOutputs
}

func formatChatID(chatID int64) string {
	chatIDStr := strconv.FormatInt(chatID, 10)
	if strings.HasPrefix(chatIDStr, "-100") {
		return strings.TrimPrefix(chatIDStr, "-100")
	}

	return chatIDStr
}

var RecapOutputTemplate = lo.Must(template.
	New("recap output markdown template").
	Funcs(template.FuncMap{
		"join":   strings.Join,
		"sub":    func(a, b int) int { return a - b },
		"add":    func(a, b int) int { return a + b },
		"escape": tgbot.EscapeHTMLSymbols,
	}).
	Parse(`{{ $chatID := .ChatID }}{{ if .Recap.SinceID }}## <a href="https://t.me/c/{{ $chatID }}/{{ .Recap.SinceID }}">{{ escape .Recap.TopicName }}</a>{{ else }}## {{ escape .Recap.TopicName }}{{ end }}
参与人：{{ join .Recap.ParticipantsNamesWithoutUsername "，" }}
讨论：{{ range $di, $d := .Recap.Discussion }}
 - {{ escape $d.Point }}{{ if len $d.KeyIDs }} {{ range $cIndex, $c := $d.KeyIDs }}<a href="https://t.me/c/{{ $chatID }}/{{ $c }}">[{{ add $cIndex 1 }}]</a>{{ if not (eq $cIndex (sub (len $d.KeyIDs) 1)) }} {{ end }}{{ end }}{{ end }}{{ end }}{{ if .Recap.Conclusion }}
结论：{{ escape .Recap.Conclusion }}{{ end }}`))

func (m *Model) summarizeChatHistoriesSlice(s string) ([]*openai.ChatHistorySummarizationOutputs, goopenai.Usage, error) {
	if s == "" {
		return make([]*openai.ChatHistorySummarizationOutputs, 0), goopenai.Usage{}, nil
	}

	m.logger.Info(fmt.Sprintf("✍️ summarizing chat histories:\n%s", s))

	resp, err := m.openAI.SummarizeChatHistories(context.Background(), s)
	if err != nil {
		return nil, goopenai.Usage{}, err
	}
	if len(resp.Choices) == 0 {
		return nil, goopenai.Usage{}, nil
	}

	m.logger.Info("✅ summarized chat histories")
	if resp.Choices[0].Message.Content == "" {
		return nil, goopenai.Usage{}, nil
	}

	var outputs []*openai.ChatHistorySummarizationOutputs

	err = json.Unmarshal([]byte(resp.Choices[0].Message.Content), &outputs)
	if err != nil {
		m.logger.Error("failed to unmarshal chat history summarization output", zap.String("content", resp.Choices[0].Message.Content))
		return nil, resp.Usage, err
	}

	m.logger.Info(fmt.Sprintf("✅ unmarshaled chat history summarization output: %s", fo.May(json.Marshal(outputs))))

	return outputs, resp.Usage, nil
}

func (m *Model) summarizeChatHistories(messageIDs []int64, llmFriendlyChatHistories string) ([]*openai.ChatHistorySummarizationOutputs, goopenai.Usage, error) {
	chatHistoriesSlices := m.openAI.SplitContentBasedByTokenLimitations(llmFriendlyChatHistories, 15000)
	chatHistoriesSummarizations := make([]*openai.ChatHistorySummarizationOutputs, 0, len(chatHistoriesSlices))

	var statusUsage goopenai.Usage

	for _, s := range chatHistoriesSlices {
		var outputs []*openai.ChatHistorySummarizationOutputs

		_, _, err := lo.AttemptWithDelay(3, time.Second, func(tried int, delay time.Duration) error {
			o, usage, err := m.summarizeChatHistoriesSlice(s)
			statusUsage.CompletionTokens += usage.CompletionTokens
			statusUsage.PromptTokens += usage.PromptTokens
			statusUsage.TotalTokens += usage.TotalTokens

			if err != nil {
				m.logger.Error(fmt.Sprintf("failed to summarize chat histories slice: %s, tried %d...", s, tried))
				return err
			}

			outputs = o
			return nil
		})
		if err != nil {
			return make([]*openai.ChatHistorySummarizationOutputs, 0), goopenai.Usage{}, err
		}
		if outputs == nil {
			continue
		}

		// filter out empty outputs
		outputs = lo.Filter(outputs, func(item *openai.ChatHistorySummarizationOutputs, _ int) bool {
			return item != nil &&
				item.TopicName != "" && // filter out empty topic name
				item.SinceID != 0 && // filter out empty since id
				len(item.ParticipantsNamesWithoutUsername) > 0 // filter out empty participants
		})

		// filter out non-exist message ids
		for i := range outputs {
			for j := range outputs[i].Discussion {
				outputs[i].Discussion[j].KeyIDs = lo.Filter(outputs[i].Discussion[j].KeyIDs, func(item int64, _ int) bool {
					return lo.Contains(messageIDs, item)
				})
			}

			outputs[i].Discussion = lo.Filter(outputs[i].Discussion, func(item *openai.ChatHistorySummarizationOutputsDiscussion, _ int) bool {
				return len(item.KeyIDs) > 0 && item.Point != ""
			})
		}

		// limit key ids to 5
		for _, o := range outputs {
			for _, d := range o.Discussion {
				d.KeyIDs = lo.UniqBy(d.KeyIDs, func(item int64) int64 {
					return item
				})
				d.KeyIDs = lo.Filter(d.KeyIDs, func(item int64, _ int) bool {
					return item != 0
				})

				if len(d.KeyIDs) > 5 {
					d.KeyIDs = d.KeyIDs[:5]
				}
			}
		}

		chatHistoriesSummarizations = append(chatHistoriesSummarizations, outputs...)
	}

	return chatHistoriesSummarizations, statusUsage, nil
}

func (m *Model) renderRecapTemplates(chatID int64, summarizations []*openai.ChatHistorySummarizationOutputs) ([]string, error) {
	ss := make([]string, 0)

	for _, r := range summarizations {
		sb := new(strings.Builder)

		err := RecapOutputTemplate.Execute(sb, RecapOutputTemplateInputs{
			ChatID: formatChatID(chatID),
			Recap:  r,
		})
		if err != nil {
			return make([]string, 0), err
		}

		ss = append(ss, sb.String())
	}

	return ss, nil
}
