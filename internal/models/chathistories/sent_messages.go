package chathistories

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/nekomeowww/insights-bot/ent"
	"github.com/nekomeowww/insights-bot/ent/sentmessages"
	"go.uber.org/zap"
	"strings"
)

type messageType int

const (
	autoRecapMessage messageType = iota
)

func (m *Model) SaveOneTelegramSentMessage(message *tgbotapi.Message, isPinned bool) error {
	telegramSentMessageCreate := m.ent.SentMessages.
		Create().
		SetMessageID(message.MessageID).
		SetChatID(message.Chat.ID).
		SetText(message.Text).
		SetFromPlatform(int(FromPlatformTelegram)).
		SetIsPinned(isPinned).
		SetMessageType(int(autoRecapMessage))

	telegramSentMessage, err := telegramSentMessageCreate.Save(context.TODO())
	if err != nil {
		return err
	}

	m.logger.Debug("saved one telegram sent message",
		zap.String("id", telegramSentMessage.ID.String()),
		zap.Int64("chat_id", telegramSentMessage.ChatID),
		zap.Int("message_id", telegramSentMessage.MessageID),
		zap.String("text", strings.ReplaceAll(telegramSentMessage.Text, "\n", " ")),
	)

	return nil
}

func (m *Model) FindLastTelegramPinnedMessage(chatID int64) (*ent.SentMessages, error) {
	telegramSentMessage, err := m.ent.SentMessages.
		Query().
		Where(
			sentmessages.ChatID(chatID),
			sentmessages.IsPinned(true),
		).
		Order(ent.Desc(sentmessages.FieldCreatedAt)).
		First(context.Background())
	if err != nil {
		return nil, err
	}

	return telegramSentMessage, nil
}

func (m *Model) UpdatePinnedMessage(chatID int64, messageID int, isPinned bool) error {
	_, err := m.ent.SentMessages.
		Update().
		SetIsPinned(isPinned).
		Where(
			sentmessages.ChatID(chatID),
			sentmessages.MessageID(messageID),
		).
		Save(context.Background())
	if err != nil {
		return err
	}

	m.logger.Debug("updated one pinned message",
		zap.Int64("chat_id", chatID),
		zap.Int("message_id", messageID),
		zap.Bool("is_pinned", isPinned),
	)

	return nil
}
