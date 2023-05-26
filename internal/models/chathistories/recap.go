package chathistories

import (
	"context"
	"encoding/json"
	"html/template"
	"strconv"
	"strings"
	"time"

	"github.com/nekomeowww/fo"
	"github.com/nekomeowww/insights-bot/internal/thirdparty/openai"
	"github.com/nekomeowww/insights-bot/pkg/bots/tgbot"
	"github.com/samber/lo"
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

func (m *Model) summarizeChatHistoriesSlice(s string) ([]*openai.ChatHistorySummarizationOutputs, int, int, int, error) {
	if s == "" {
		return make([]*openai.ChatHistorySummarizationOutputs, 0), 0, 0, 0, nil
	}

	m.logger.Infof("✍️ summarizing chat histories:\n%s", s)

	resp, err := m.openAI.SummarizeChatHistories(context.Background(), s)
	if err != nil {
		return nil, 0, 0, 0, err
	}
	if len(resp.Choices) == 0 {
		return nil, 0, 0, 0, nil
	}

	m.logger.Info("✅ summarized chat histories")
	if resp.Choices[0].Message.Content == "" {
		return nil, 0, 0, 0, nil
	}

	var outputs []*openai.ChatHistorySummarizationOutputs

	err = json.Unmarshal([]byte(resp.Choices[0].Message.Content), &outputs)
	if err != nil {
		m.logger.Errorf("failed to unmarshal chat history summarization output: %s", resp.Choices[0].Message.Content)
		return nil, resp.Usage.CompletionTokens, resp.Usage.PromptTokens, resp.Usage.TotalTokens, err
	}

	m.logger.Infof("✅ unmarshaled chat history summarization output: %s", fo.May(json.Marshal(outputs)))

	return outputs, resp.Usage.CompletionTokens, resp.Usage.PromptTokens, resp.Usage.TotalTokens, nil
}

func (m *Model) summarizeChatHistories(llmFriendlyChatHistories string) ([]*openai.ChatHistorySummarizationOutputs, int, int, int, error) {
	chatHistoriesSlices := m.openAI.SplitContentBasedByTokenLimitations(llmFriendlyChatHistories, 2800)
	chatHistoriesSummarizations := make([]*openai.ChatHistorySummarizationOutputs, 0, len(chatHistoriesSlices))

	statsCompletionTokenUsage := 0
	statsPromptTokenUsage := 0
	statsTotalTokenUsage := 0

	for _, s := range chatHistoriesSlices {
		var outputs []*openai.ChatHistorySummarizationOutputs

		_, _, err := lo.AttemptWithDelay(3, time.Second, func(tried int, delay time.Duration) error {
			o, completionTokenUsage, promptTokenUsage, totalTokenUsage, err := m.summarizeChatHistoriesSlice(s)
			statsCompletionTokenUsage += completionTokenUsage
			statsPromptTokenUsage += promptTokenUsage
			statsTotalTokenUsage += totalTokenUsage

			if err != nil {
				m.logger.Errorf("failed to summarize chat histories slice: %s, tried %d...", s, tried)
				return err
			}

			outputs = o
			return nil
		})
		if err != nil {
			return make([]*openai.ChatHistorySummarizationOutputs, 0), 0, 0, 0, err
		}
		if outputs == nil {
			continue
		}

		// filter out empty outputs
		outputs = lo.Filter(outputs, func(item *openai.ChatHistorySummarizationOutputs, _ int) bool {
			return item != nil &&
				item.TopicName != "" && // filter out empty topic name
				item.SinceID != 0 && // filter out empty since id
				len(item.ParticipantsNamesWithoutUsername) > 0 && // filter out empty participants
				len(item.Discussion) > 0 // filter out empty discussion
		})

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

	return chatHistoriesSummarizations, statsCompletionTokenUsage, statsPromptTokenUsage, statsTotalTokenUsage, nil
}

func (m *Model) fillIntoRecapTemplates(chatID int64, summarizations []*openai.ChatHistorySummarizationOutputs) ([]string, error) {
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
