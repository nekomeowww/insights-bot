package chat_histories

import (
	"context"
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/ostafen/clover/v2"
	"github.com/sirupsen/logrus"
	"go.uber.org/fx"

	"github.com/nekomeowww/insights-bot/internal/datastore"
	telegram_bot "github.com/nekomeowww/insights-bot/pkg/bots/telegram"
	"github.com/nekomeowww/insights-bot/pkg/logger"
	"github.com/nekomeowww/insights-bot/pkg/openai"
	"github.com/nekomeowww/insights-bot/pkg/types/chat_history"
)

type NewChatHistoriesModelParam struct {
	fx.In

	Logger *logger.Logger
	Clover *datastore.Clover
	OpenAI *openai.Client
}

type ChatHistoriesModel struct {
	Logger *logger.Logger
	Clover *datastore.Clover
	OpenAI *openai.Client
}

func NewChatHistoriesModel() func(NewChatHistoriesModelParam) (*ChatHistoriesModel, error) {
	return func(param NewChatHistoriesModelParam) (*ChatHistoriesModel, error) {
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

		return &ChatHistoriesModel{
			Logger: param.Logger,
			Clover: param.Clover,
			OpenAI: param.OpenAI,
		}, nil
	}
}

func (m *ChatHistoriesModel) extractTextWithSummarization(message *tgbotapi.Message) (string, error) {
	text := telegram_bot.ExtractTextFromMessage(message)
	if text == "" {
		return "", nil
	}
	if utf8.RuneCountInString(text) >= 200 {
		resp, err := m.OpenAI.SummarizeWithOneChatHistory(context.Background(), text)
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

func (m *ChatHistoriesModel) SaveOneTelegramChatHistory(message *tgbotapi.Message) error {
	if message.Text == "" && message.Caption == "" {
		m.Logger.Warn("message text is empty")
		return nil
	}

	telegramChatHistory := chat_history.TelegramChatHistory{
		ID:        clover.NewObjectId(),
		ChatID:    message.Chat.ID,
		MessageID: message.MessageID,
		UserID:    message.From.ID,
		Username:  message.From.UserName,
		FullName:  telegram_bot.FullNameFromFirstAndLastName(message.From.FirstName, message.From.LastName),
		ChattedAt: time.Unix(int64(message.Date), 0).UnixMilli(),
		CreatedAt: time.Now().UnixMilli(),
		UpdatedAt: time.Now().UnixMilli(),
	}

	text, err := m.extractTextWithSummarization(message)
	if err != nil {
		return err
	}
	if text == "" {
		m.Logger.Warn("message text is empty")
		return nil
	}
	if message.ForwardFrom != nil {
		telegramChatHistory.Text = "转发了来自" + telegram_bot.FullNameFromFirstAndLastName(message.ForwardFrom.FirstName, message.ForwardFrom.LastName) + "的消息：" + text
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
			telegramChatHistory.RepliedToFullName = telegram_bot.FullNameFromFirstAndLastName(message.ReplyToMessage.From.FirstName, message.ReplyToMessage.From.LastName)
			telegramChatHistory.RepliedToUsername = message.ReplyToMessage.From.UserName
			telegramChatHistory.RepliedToText = repliedToText
		}
	}

	id, err := m.Clover.InsertOne(
		chat_history.TelegramChatHistory{}.CollectionName(),
		clover.NewDocumentOf(telegramChatHistory),
	)
	if err != nil {
		return err
	}

	m.Logger.WithFields(logrus.Fields{
		"id":         id,
		"chat_id":    telegramChatHistory.ChatID,
		"message_id": telegramChatHistory.MessageID,
		"text":       strings.ReplaceAll(telegramChatHistory.Text, "\n", " "),
	}).Debug("saved one telegram chat history")
	return nil
}

func (m *ChatHistoriesModel) FindLastOneHourChatHistories(chatID int64) ([]*chat_history.TelegramChatHistory, error) {
	return m.FindChatHistoriesByTimeBefore(chatID, time.Hour)
}

func (m *ChatHistoriesModel) FindLastSixHourChatHistories(chatID int64) ([]*chat_history.TelegramChatHistory, error) {
	return m.FindChatHistoriesByTimeBefore(chatID, 6*time.Hour)
}

func (m *ChatHistoriesModel) FindChatHistoriesByTimeBefore(chatID int64, before time.Duration) ([]*chat_history.TelegramChatHistory, error) {
	query := clover.
		NewQuery(chat_history.TelegramChatHistory{}.CollectionName()).
		Where(clover.Field("chat_id").Eq(chatID)).
		Where(clover.Field("chatted_at").Gt(time.Now().Add(-before).UnixMilli())).
		Sort(clover.SortOption{
			Field:     "message_id",
			Direction: 1,
		})

	m.Logger.Infof("querying chat histories for %d", chatID)
	docs, err := m.Clover.FindAll(query)
	if err != nil {
		return make([]*chat_history.TelegramChatHistory, 0), err
	}

	m.Logger.Infof("found %d chat histories", len(docs))
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
	if username == "" {
		return fullName
	}
	if utf8.RuneCountInString(fullName) >= 10 {
		return fmt.Sprintf("%s (用户名：%s)", username, username)
	}

	return fmt.Sprintf("%s (用户名：%s)", fullName, username)
}

func (c *ChatHistoriesModel) SummarizeChatHistories(histories []*chat_history.TelegramChatHistory) (string, error) {
	historiesLLMFriendly := make([]string, 0, len(histories))
	for _, message := range histories {
		chattedAt := time.UnixMilli(message.ChattedAt).Format("15:04:05")
		partialContextMessage := fmt.Sprintf("%s 于 %s", formatFullNameAndUsername(message.FullName, message.Username), chattedAt)
		if message.RepliedToMessageID == 0 {
			historiesLLMFriendly = append(historiesLLMFriendly, fmt.Sprintf(
				"%s 发送：%s",
				partialContextMessage,
				message.Text,
			))
		} else {
			historiesLLMFriendly = append(historiesLLMFriendly, fmt.Sprintf(
				"%s 回复 %s ：%s",
				partialContextMessage,
				formatFullNameAndUsername(message.RepliedToFullName, message.RepliedToUsername),
				message.Text,
			))
		}
	}

	chatHistories := strings.Join(historiesLLMFriendly, "\n")
	chatHistoriesSlices := c.OpenAI.SplitContentBasedByTokenLimitations(chatHistories, 3000)
	chatHistoriesSummarizations := make([]string, 0, len(chatHistoriesSlices))
	for _, s := range chatHistoriesSlices {
		c.Logger.Infof("✍️ summarizing last one hour chat histories:\n%s", s)
		resp, err := c.OpenAI.SummarizeWithChatHistories(context.Background(), s)
		if err != nil {
			return "", err
		}
		if len(resp.Choices) == 0 {
			return "", nil
		}

		c.Logger.WithFields(logrus.Fields{
			"prompt_token_usage":     resp.Usage.PromptTokens,
			"completion_token_usage": resp.Usage.CompletionTokens,
			"total_token_usage":      resp.Usage.TotalTokens,
		}).Info("✅ summarized last one hour chat histories")
		if resp.Choices[0].Message.Content == "" {
			continue
		}

		chatHistoriesSummarizations = append(chatHistoriesSummarizations, resp.Choices[0].Message.Content)
	}

	return strings.Join(chatHistoriesSummarizations, "\n\n"), nil
}
