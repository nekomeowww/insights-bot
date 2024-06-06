package chathistories

import (
	"context"
	"encoding/json"
	"errors"
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
	"github.com/nekomeowww/insights-bot/pkg/types/telegram"
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
参与人：{{ join .Recap.Participants "，" }}
讨论：{{ range $di, $d := .Recap.Discussion }}
 - {{ escape $d.Point }}{{ if len $d.KeyIDs }} {{ range $cIndex, $c := $d.KeyIDs }}<a href="https://t.me/c/{{ $chatID }}/{{ $c }}">[{{ add $cIndex 1 }}]</a>{{ if not (eq $cIndex (sub (len $d.KeyIDs) 1)) }} {{ end }}{{ end }}{{ end }}{{ end }}{{ if .Recap.Conclusion }}
结论：{{ escape .Recap.Conclusion }}{{ end }}`))

var RecapWithoutLinksOutputTemplate = lo.Must(template.
	New("recap output markdown template").
	Funcs(template.FuncMap{
		"join":   strings.Join,
		"sub":    func(a, b int) int { return a - b },
		"add":    func(a, b int) int { return a + b },
		"escape": tgbot.EscapeHTMLSymbols,
	}).
	Parse(`{{ $chatID := .ChatID }}{{ if .Recap.SinceID }}## {{ escape .Recap.TopicName }}{{ else }}## {{ escape .Recap.TopicName }}{{ end }}
参与人：{{ join .Recap.Participants "，" }}
讨论：{{ range $di, $d := .Recap.Discussion }}
 - {{ escape $d.Point }}{{ end }}{{ if .Recap.Conclusion }}
结论：{{ escape .Recap.Conclusion }}{{ end }}`))

func (m *Model) summarizeChatHistoriesSlice(chatID int64, s string) ([]*openai.ChatHistorySummarizationOutputs, goopenai.Usage, error) {
	if s == "" {
		return make([]*openai.ChatHistorySummarizationOutputs, 0), goopenai.Usage{}, nil
	}

	m.logger.Info(fmt.Sprintf("✍️ summarizing chat histories:\n%s", s),
		zap.Int64("chat_id", chatID),
		zap.String("model_name", m.openAI.GetModelName()),
	)

	resp, err := m.openAI.SummarizeChatHistories(context.Background(), s)
	if err != nil {
		return nil, goopenai.Usage{}, err
	}
	if len(resp.Choices) == 0 {
		return nil, goopenai.Usage{}, nil
	}

	m.logger.Info("✅ summarized chat histories",
		zap.Int64("chat_id", chatID),
		zap.String("model_name", m.openAI.GetModelName()),
	)
	if resp.Choices[0].Message.Content == "" {
		return nil, goopenai.Usage{}, nil
	}

	var outputs []*openai.ChatHistorySummarizationOutputs
	resp.Choices[0].Message.Content = strings.TrimPrefix(resp.Choices[0].Message.Content, "```json")
	resp.Choices[0].Message.Content = strings.TrimPrefix(resp.Choices[0].Message.Content, "```")
	resp.Choices[0].Message.Content = strings.TrimSuffix(resp.Choices[0].Message.Content, "```")

	err = json.Unmarshal([]byte(resp.Choices[0].Message.Content), &outputs)
	if err != nil {
		m.logger.Error("failed to unmarshal chat history summarization output",
			zap.String("content", resp.Choices[0].Message.Content),
			zap.Int64("chat_id", chatID),
			zap.String("model_name", m.openAI.GetModelName()),
		)

		return nil, resp.Usage, err
	}

	m.logger.Info(fmt.Sprintf("✅ unmarshaled chat history summarization output: %s", fo.May(json.Marshal(outputs))),
		zap.Int64("chat_id", chatID),
		zap.String("model_name", m.openAI.GetModelName()),
	)

	return outputs, resp.Usage, nil
}

func filterOutInvalidFields(messageIDs []int64) func(output *openai.ChatHistorySummarizationOutputs, _ int) *openai.ChatHistorySummarizationOutputs {
	return func(output *openai.ChatHistorySummarizationOutputs, _ int) *openai.ChatHistorySummarizationOutputs {
		// limit key ids to 5
		output.Participants = lo.Uniq(output.Participants)

		// filter out non-exist message ids
		for _, d := range output.Discussion {
			d.KeyIDs = lo.Filter(d.KeyIDs, func(item int64, _ int) bool {
				return lo.Contains(messageIDs, item) && item != 0
			})
			d.KeyIDs = lo.UniqBy(d.KeyIDs, func(item int64) int64 {
				return item
			})

			if len(d.KeyIDs) > 5 {
				d.KeyIDs = d.KeyIDs[:5]
			}
		}

		output.Discussion = lo.Filter(output.Discussion, func(item *openai.ChatHistorySummarizationOutputsDiscussion, _ int) bool {
			return len(item.KeyIDs) > 0 && item.Point != ""
		})

		return output
	}
}

func filterOutInvalidOutputFilterFunc(output *openai.ChatHistorySummarizationOutputs, _ int) bool {
	return output != nil &&
		output.TopicName != "" && // filter out empty topic name
		output.SinceID != 0 && // filter out empty since id
		len(output.Participants) > 0 && // filter out empty participants
		len(output.Discussion) > 0 // filter out empty discussion
}

func filterOutMention(output *openai.ChatHistorySummarizationOutputs, _ int) *openai.ChatHistorySummarizationOutputs {
	output.Participants = lo.Map(output.Participants, func(item string, _ int) string {
		return strings.TrimPrefix(item, "@")
	})

	return output
}

func (m *Model) summarizeChatHistories(chatID int64, messageIDs []int64, llmFriendlyChatHistories string) ([]*openai.ChatHistorySummarizationOutputs, goopenai.Usage, error) {
	tokenLimit := m.config.OpenAI.TokenLimit - m.config.OpenAI.ChatHistoriesRecapTokenLimit
	chatHistoriesSlices := m.openAI.SplitContentBasedByTokenLimitations(llmFriendlyChatHistories, int(tokenLimit))
	chatHistoriesSummarizations := make([]*openai.ChatHistorySummarizationOutputs, 0, len(chatHistoriesSlices))

	var statusUsage goopenai.Usage

	for _, s := range chatHistoriesSlices {
		var outputs []*openai.ChatHistorySummarizationOutputs

		_, _, err := lo.AttemptWithDelay(5, time.Second, func(tried int, delay time.Duration) error {
			o, usage, err := m.summarizeChatHistoriesSlice(chatID, s)
			statusUsage.CompletionTokens += usage.CompletionTokens
			statusUsage.PromptTokens += usage.PromptTokens
			statusUsage.TotalTokens += usage.TotalTokens

			if err != nil {
				m.logger.Error(fmt.Sprintf("failed to summarize chat histories slice: %s, tried %d...", s, tried),
					zap.Int64("chat_id", chatID),
					zap.String("model_name", m.openAI.GetModelName()),
					zap.Int64("configured_chat_histories_recap_token_limit", m.config.OpenAI.ChatHistoriesRecapTokenLimit),
					zap.Int64("configured_token_limit", m.config.OpenAI.TokenLimit),
					zap.Int64("calculated_token_limit", tokenLimit),
				)

				return err
			}

			// filter out invalid fields
			o = lo.Map(o, filterOutInvalidFields(messageIDs))
			// filter out empty outputs
			o = lo.Filter(o, filterOutInvalidOutputFilterFunc)
			// filter out mentions
			o = lo.Map(o, filterOutMention)

			if len(o) == 0 {
				m.logger.Error(fmt.Sprintf("no valid outputs from chat histories slice: %s, tried %d...", s, tried),
					zap.Int64("chat_id", chatID),
					zap.String("model_name", m.openAI.GetModelName()),
					zap.Int64("configured_chat_histories_recap_token_limit", m.config.OpenAI.ChatHistoriesRecapTokenLimit),
					zap.Int64("configured_token_limit", m.config.OpenAI.TokenLimit),
					zap.Int64("calculated_token_limit", tokenLimit),
				)

				return errors.New("no valid outputs")
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

		chatHistoriesSummarizations = append(chatHistoriesSummarizations, outputs...)
	}

	return chatHistoriesSummarizations, statusUsage, nil
}

func (m *Model) renderRecapTemplates(chatID int64, chatType telegram.ChatType, summarizations []*openai.ChatHistorySummarizationOutputs) ([]string, error) {
	ss := make([]string, 0)

	for _, r := range summarizations {
		sb := new(strings.Builder)

		switch chatType {
		case telegram.ChatTypeSuperGroup:
			err := RecapOutputTemplate.Execute(sb, RecapOutputTemplateInputs{
				ChatID: formatChatID(chatID),
				Recap:  r,
			})
			if err != nil {
				return make([]string, 0), err
			}

			ss = append(ss, sb.String())
		case telegram.ChatTypePrivate, telegram.ChatTypeGroup, telegram.ChatTypeChannel:
			err := RecapWithoutLinksOutputTemplate.Execute(sb, RecapOutputTemplateInputs{
				ChatID: formatChatID(chatID),
				Recap:  r,
			})
			if err != nil {
				return make([]string, 0), err
			}

			ss = append(ss, sb.String())
		}
	}

	return ss, nil
}
