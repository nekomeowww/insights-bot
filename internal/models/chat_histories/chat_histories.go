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
	"github.com/nekomeowww/insights-bot/pkg/logger"
	"github.com/nekomeowww/insights-bot/pkg/openai"
	"github.com/nekomeowww/insights-bot/pkg/types/chat_history"
	"github.com/nekomeowww/insights-bot/pkg/utils"
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

func FullNameFromFirstAndLastName(firstName, lastName string) string {
	if lastName == "" {
		return firstName
	}
	if firstName == "" {
		return lastName
	}
	if utils.ContainsCJKChar(firstName) || utils.ContainsCJKChar(lastName) {
		return lastName + firstName
	}

	return firstName + " " + lastName
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
		FullName:  FullNameFromFirstAndLastName(message.From.FirstName, message.From.LastName),
		ChattedAt: time.Unix(int64(message.Date), 0).UnixMilli(),
		CreatedAt: time.Now().UnixMilli(),
		UpdatedAt: time.Now().UnixMilli(),
	}
	if message.ForwardFrom != nil {
		telegramChatHistory.Text = "转发了来自" + FullNameFromFirstAndLastName(message.ForwardFrom.FirstName, message.ForwardFrom.LastName) + "的消息：" + message.Text
	} else if message.ForwardFromChat != nil {
		telegramChatHistory.Text = "转发了来自" + message.ForwardFromChat.Title + "的消息：" + message.Text
	} else if message.Caption != "" {
		telegramChatHistory.Text = message.Caption
	} else {
		telegramChatHistory.Text = message.Text
	}
	if utf8.RuneCountInString(telegramChatHistory.Text) >= 200 {
		resp, err := m.OpenAI.SummarizeWithOneChatHistory(context.Background(), telegramChatHistory.Text)
		if err != nil {
			return err
		}
		if len(resp.Choices) == 0 {
			return nil
		}

		telegramChatHistory.Text = resp.Choices[0].Message.Content
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
		"text":       strings.ReplaceAll(telegramChatHistory.Text, "\n", "\\n"),
	}).Info("saved one telegram chat history")
	return nil
}

func (m *ChatHistoriesModel) FindLastOneHourChatHistories(chatID int64) ([]*chat_history.TelegramChatHistory, error) {
	query := clover.
		NewQuery(chat_history.TelegramChatHistory{}.CollectionName()).
		Where(clover.Field("chat_id").Eq(chatID)).
		Where(clover.Field("chatted_at").Gt(time.Now().Add(-time.Hour).UnixMilli())).
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

func (c *ChatHistoriesModel) SummarizeLastOneHourChatHistories(chatID int64) (string, error) {
	histories, err := c.FindLastOneHourChatHistories(chatID)
	if err != nil {
		return "", err
	}
	if len(histories) <= 5 {
		return "", nil
	}

	historiesLLMFriendly := make([]string, 0, len(histories))
	for _, message := range histories {
		historiesLLMFriendly = append(historiesLLMFriendly, fmt.Sprintf("%s (用户名：%s) 于 %s 发送：%s", message.FullName, message.Username, time.UnixMilli(message.ChattedAt).Format("2006-01-02 15:04:05"), message.Text))
	}

	chatHistories := strings.Join(historiesLLMFriendly, "\n")
	chatHistoriesContent, err := c.OpenAI.TruncateContentBasedOnTokens(chatHistories)
	if err != nil {
		return "", fmt.Errorf("failed to truncate content based on tokens... %w", err)
	}

	c.Logger.Info("✍️ summarizing last one hour chat histories")
	resp, err := c.OpenAI.SummarizeWithChatHistories(context.Background(), chatHistoriesContent)
	if err != nil {
		return "", err
	}
	if len(resp.Choices) == 0 {
		return "", nil
	}

	c.Logger.WithFields(logrus.Fields{
		"chat_id":                chatID,
		"prompt_token_usage":     resp.Usage.PromptTokens,
		"completion_token_usage": resp.Usage.CompletionTokens,
		"total_token_usage":      resp.Usage.TotalTokens,
	}).Info("✅ summarized last one hour chat histories")
	return resp.Choices[0].Message.Content, nil
}
