package tgchats

import (
	"context"

	"github.com/nekomeowww/insights-bot/ent"
	"github.com/nekomeowww/insights-bot/ent/telegramchatautorecapssubscribers"
	"go.uber.org/zap"
)

func (m *Model) FindOneAutoRecapsSubscriber(chatID int64, userID int64) (*ent.TelegramChatAutoRecapsSubscribers, error) {
	subscriber, err := m.ent.TelegramChatAutoRecapsSubscribers.
		Query().
		Where(
			telegramchatautorecapssubscribers.ChatID(chatID),
			telegramchatautorecapssubscribers.UserID(userID),
		).
		First(context.Background())
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		}

		return nil, err
	}

	return subscriber, nil
}

func (m *Model) FindAutoRecapsSubscribers(chatID int64) ([]*ent.TelegramChatAutoRecapsSubscribers, error) {
	subscribers, err := m.ent.TelegramChatAutoRecapsSubscribers.
		Query().
		Where(telegramchatautorecapssubscribers.ChatID(chatID)).
		All(context.Background())
	if err != nil {
		if ent.IsNotFound(err) {
			return make([]*ent.TelegramChatAutoRecapsSubscribers, 0), nil
		}

		return make([]*ent.TelegramChatAutoRecapsSubscribers, 0), err
	}

	return subscribers, nil
}

func (m *Model) SubscribeToAutoRecaps(chatID int64, userID int64) error {
	subscriber, err := m.FindOneAutoRecapsSubscriber(chatID, userID)
	if err != nil {
		return err
	}
	if subscriber != nil {
		return nil
	}

	return m.ent.TelegramChatAutoRecapsSubscribers.
		Create().
		SetChatID(chatID).
		SetUserID(userID).
		Exec(context.Background())
}

func (m *Model) UnsubscribeToAutoRecaps(chatID int64, userID int64) error {
	subscriber, err := m.FindOneAutoRecapsSubscriber(chatID, userID)
	if err != nil {
		return err
	}
	if subscriber == nil {
		return nil
	}

	return m.ent.TelegramChatAutoRecapsSubscribers.
		DeleteOne(subscriber).
		Exec(context.Background())
}

func (m *Model) DeleteAllSubscribersByChatID(chatID int64) error {
	_, err := m.ent.TelegramChatAutoRecapsSubscribers.
		Delete().
		Where(telegramchatautorecapssubscribers.ChatID(chatID)).
		Exec(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func (m *Model) MigrateSubscribersOfChatFromChatIDToChatID(fromChatID int64, toChatID int64) error {
	affectedRows, err := m.ent.TelegramChatAutoRecapsSubscribers.
		Update().
		Where(telegramchatautorecapssubscribers.ChatID(fromChatID)).
		SetChatID(toChatID).
		Save(context.Background())
	if err != nil {
		return err
	}

	m.logger.Info("successfully migrated options of chat",
		zap.Int64("from_chat_id", fromChatID),
		zap.Int64("to_chat_id", toChatID),
		zap.Int("affected_rows", affectedRows),
	)

	return nil
}
