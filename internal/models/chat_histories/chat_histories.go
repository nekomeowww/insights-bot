package chat_histories

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"text/template"
	"time"
	"unicode/utf8"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	"github.com/ostafen/clover/v2"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	"go.uber.org/fx"

	"github.com/nekomeowww/insights-bot/internal/datastore"
	"github.com/nekomeowww/insights-bot/pkg/bots/tgbot"
	"github.com/nekomeowww/insights-bot/pkg/logger"
	"github.com/nekomeowww/insights-bot/pkg/openai"
	"github.com/nekomeowww/insights-bot/pkg/types/chat_history"
)

type NewModelParams struct {
	fx.In

	Logger *logger.Logger
	Clover *datastore.Clover
	OpenAI *openai.Client
}

type Model struct {
	logger *logger.Logger
	clover *datastore.Clover
	openAI *openai.Client
}

func NewModel() func(NewModelParams) (*Model, error) {
	return func(param NewModelParams) (*Model, error) {
		hasCollection, err := param.Clover.HasCollection(chat_history.TelegramChatHistory{}.CollectionName())
		if err != nil {
			return nil, err
		}
		if !hasCollection {
			err = param.Clover.CreateCollection(chat_history.TelegramChatHistory{}.CollectionName())
			if err != nil {
				return nil, err
			}
		}

		return &Model{
			logger: param.Logger,
			clover: param.Clover,
			openAI: param.OpenAI,
		}, nil
	}
}

func (m *Model) extractTextWithSummarization(message *tgbotapi.Message) (string, error) {
	text := tgbot.ExtractTextFromMessage(message)
	if text == "" {
		return "", nil
	}
	if utf8.RuneCountInString(text) >= 200 {
		resp, err := m.openAI.SummarizeWithOneChatHistory(context.Background(), text)
		if err != nil {
			return "", err
		}
		if len(resp.Choices) == 0 {
			return "", nil
		}

		return resp.Choices[0].Message.Content, nil
	}

	return text, nil
}

func (m *Model) SaveOneTelegramChatHistory(message *tgbotapi.Message) error {
	if message.Text == "" && message.Caption == "" {
		m.logger.Warn("message text is empty")
		return nil
	}

	telegramChatHistory := chat_history.TelegramChatHistory{
		ID:        clover.NewObjectId(),
		ChatID:    message.Chat.ID,
		MessageID: message.MessageID,
		UserID:    message.From.ID,
		Username:  message.From.UserName,
		FullName:  tgbot.FullNameFromFirstAndLastName(message.From.FirstName, message.From.LastName),
		ChattedAt: time.Unix(int64(message.Date), 0).UnixMilli(),
		CreatedAt: time.Now().UnixMilli(),
		UpdatedAt: time.Now().UnixMilli(),
	}

	text, err := m.extractTextWithSummarization(message)
	if err != nil {
		return err
	}
	if text == "" {
		m.logger.Warn("message text is empty")
		return nil
	}
	if message.ForwardFrom != nil {
		telegramChatHistory.Text = "转发了来自" + tgbot.FullNameFromFirstAndLastName(message.ForwardFrom.FirstName, message.ForwardFrom.LastName) + "的消息：" + text
	} else if message.ForwardFromChat != nil {
		telegramChatHistory.Text = "转发了来自" + message.ForwardFromChat.Title + "的消息：" + text
	} else {
		telegramChatHistory.Text = text
	}
	if message.ReplyToMessage != nil {
		repliedToText, err := m.extractTextWithSummarization(message.ReplyToMessage)
		if err != nil {
			return err
		}
		if repliedToText != "" {
			telegramChatHistory.RepliedToMessageID = message.ReplyToMessage.MessageID
			telegramChatHistory.RepliedToUserID = message.ReplyToMessage.From.ID
			telegramChatHistory.RepliedToFullName = tgbot.FullNameFromFirstAndLastName(message.ReplyToMessage.From.FirstName, message.ReplyToMessage.From.LastName)
			telegramChatHistory.RepliedToUsername = message.ReplyToMessage.From.UserName
			telegramChatHistory.RepliedToText = repliedToText
		}
	}

	id, err := m.clover.InsertOne(
		chat_history.TelegramChatHistory{}.CollectionName(),
		clover.NewDocumentOf(telegramChatHistory),
	)
	if err != nil {
		return err
	}

	m.logger.WithFields(logrus.Fields{
		"id":         id,
		"chat_id":    telegramChatHistory.ChatID,
		"message_id": telegramChatHistory.MessageID,
		"text":       strings.ReplaceAll(telegramChatHistory.Text, "\n", " "),
	}).Debug("saved one telegram chat history")
	return nil
}

func (m *Model) FindLastOneHourChatHistories(chatID int64) ([]*chat_history.TelegramChatHistory, error) {
	return m.FindChatHistoriesByTimeBefore(chatID, time.Hour)
}

func (m *Model) FindLastSixHourChatHistories(chatID int64) ([]*chat_history.TelegramChatHistory, error) {
	return m.FindChatHistoriesByTimeBefore(chatID, 6*time.Hour)
}

func (m *Model) FindChatHistoriesByTimeBefore(chatID int64, before time.Duration) ([]*chat_history.TelegramChatHistory, error) {
	query := clover.
		NewQuery(chat_history.TelegramChatHistory{}.CollectionName()).
		Where(clover.Field("chat_id").Eq(chatID)).
		Where(clover.Field("chatted_at").Gt(time.Now().Add(-before).UnixMilli())).
		Sort(clover.SortOption{
			Field:     "message_id",
			Direction: 1,
		})

	m.logger.Infof("querying chat histories for %d", chatID)
	docs, err := m.clover.FindAll(query)
	if err != nil {
		return make([]*chat_history.TelegramChatHistory, 0), err
	}

	m.logger.Infof("found %d chat histories", len(docs))
	if len(docs) == 0 {
		return make([]*chat_history.TelegramChatHistory, 0), nil
	}

	chatHistories := make([]*chat_history.TelegramChatHistory, 0, len(docs))
	for _, doc := range docs {
		var chatHistory chat_history.TelegramChatHistory
		err = doc.Unmarshal(&chatHistory)
		if err != nil {
			return make([]*chat_history.TelegramChatHistory, 0), err
		}

		chatHistories = append(chatHistories, &chatHistory)
	}

	return chatHistories, nil
}

func formatFullNameAndUsername(fullName, username string) string {
	if utf8.RuneCountInString(fullName) >= 10 {
		return username
	}

	return strings.ReplaceAll(fullName, "#", "")
}

type RecapOutputTemplateInputs struct {
	ChatID string
	Recaps []*openai.ChatHistorySummarizationOutputs
}

func formatChatID(chatID int64) string {
	chatIDStr := strconv.FormatInt(chatID, 10)
	if strings.HasPrefix(chatIDStr, "-100") {
		return strings.TrimPrefix(chatIDStr, "-100")
	}

	return chatIDStr
}

var RecapOutputTemplate = lo.Must(template.
	New(uuid.New().String()).
	Funcs(template.FuncMap{
		"join":   strings.Join,
		"sub":    func(a, b int) int { return a - b },
		"add":    func(a, b int) int { return a + b },
		"escape": tgbot.EscapeHTMLSymbols,
	}).
	Parse(`{{ $chatID := .ChatID }}{{ $recapLen := len .Recaps }}{{ range $i, $r := .Recaps }}{{ if $r.SinceMsgID }}## <a href="https://t.me/c/{{ $chatID }}/{{ $r.SinceMsgID }}">{{ escape $r.TopicName }}</a>{{ else }}## {{ escape $r.TopicName }}{{ end }}
参与人：{{ join $r.ParticipantsNamesWithoutUsername "，" }}
讨论：{{ range $di, $d := $r.Discussion }}
 - {{ escape $d.Point }}{{ if len $d.CriticalMessageIDs }} {{ range $cIndex, $c := $d.CriticalMessageIDs }}<a href="https://t.me/c/{{ $chatID }}/{{ $c }}">[{{ add $cIndex 1 }}]</a>{{ if not (eq $cIndex (sub (len $d.CriticalMessageIDs) 1)) }} {{ end }}{{ end }}{{ end }}{{ end }}{{ if $r.Conclusion }}
结论：{{ escape $r.Conclusion }}{{ end }}{{ if eq $i (sub $recapLen 1) }}{{ else }}

{{ end }}{{ end }}`))

func (c *Model) summarizeChatHistoriesSlice(s string) ([]*openai.ChatHistorySummarizationOutputs, error) {
	c.logger.Infof("✍️ summarizing last one hour chat histories:\n%s", s)
	resp, err := c.openAI.SummarizeWithChatHistories(context.Background(), s)
	if err != nil {
		return nil, err
	}
	if len(resp.Choices) == 0 {
		return nil, nil
	}

	c.logger.WithFields(logrus.Fields{
		"prompt_token_usage":     resp.Usage.PromptTokens,
		"completion_token_usage": resp.Usage.CompletionTokens,
		"total_token_usage":      resp.Usage.TotalTokens,
	}).Info("✅ summarized last one hour chat histories")
	if resp.Choices[0].Message.Content == "" {
		return nil, nil
	}

	var outputs []*openai.ChatHistorySummarizationOutputs
	err = json.Unmarshal([]byte(resp.Choices[0].Message.Content), &outputs)
	if err != nil {
		c.logger.Errorf("failed to unmarshal chat history summarization output: %s", resp.Choices[0].Message.Content)
		return nil, err
	}

	return outputs, nil
}

func (c *Model) SummarizeChatHistories(chatID int64, histories []*chat_history.TelegramChatHistory) (string, error) {
	historiesLLMFriendly := make([]string, 0, len(histories))
	for _, message := range histories {
		if message.RepliedToMessageID == 0 {
			historiesLLMFriendly = append(historiesLLMFriendly, fmt.Sprintf(
				"msgId:%d: %s 发送：%s",
				message.MessageID,
				formatFullNameAndUsername(message.FullName, message.Username),
				message.Text,
			))
		} else {
			repliedToPartialContextMessage := fmt.Sprintf("%s 发送的 msgId:%d 的消息", formatFullNameAndUsername(message.RepliedToFullName, message.RepliedToUsername), message.RepliedToMessageID)
			historiesLLMFriendly = append(historiesLLMFriendly, fmt.Sprintf(
				"msgId:%d: %s 回复 %s：%s",
				message.MessageID,
				formatFullNameAndUsername(message.FullName, message.Username),
				repliedToPartialContextMessage,
				message.Text,
			))
		}
	}

	chatHistories := strings.Join(historiesLLMFriendly, "\n")
	chatHistoriesSlices := c.openAI.SplitContentBasedByTokenLimitations(chatHistories, 2800)
	chatHistoriesSummarizations := make([]*openai.ChatHistorySummarizationOutputs, 0, len(chatHistoriesSlices))
	for _, s := range chatHistoriesSlices {
		var outputs []*openai.ChatHistorySummarizationOutputs
		_, _, err := lo.AttemptWithDelay(3, time.Second, func(tried int, delay time.Duration) error {
			o, err := c.summarizeChatHistoriesSlice(s)
			if err != nil {
				c.logger.Errorf("failed to summarize chat histories slice: %s, tried %d...", s, tried)
				return err
			}

			outputs = o
			return nil
		})
		if err != nil {
			return "", err
		}
		if outputs == nil {
			continue
		}

		for _, o := range outputs {
			for _, d := range o.Discussion {
				d.CriticalMessageIDs = lo.UniqBy(d.CriticalMessageIDs, func(item int64) int64 {
					return item
				})
				d.CriticalMessageIDs = lo.Filter(d.CriticalMessageIDs, func(item int64, _ int) bool {
					return item != 0
				})
				if len(d.CriticalMessageIDs) > 5 {
					d.CriticalMessageIDs = d.CriticalMessageIDs[:5]
				}
			}
		}

		chatHistoriesSummarizations = append(chatHistoriesSummarizations, outputs...)
	}

	sb := new(strings.Builder)
	err := RecapOutputTemplate.Execute(sb, RecapOutputTemplateInputs{
		ChatID: formatChatID(chatID),
		Recaps: chatHistoriesSummarizations,
	})
	if err != nil {
		return "", err
	}

	return sb.String(), nil
}
