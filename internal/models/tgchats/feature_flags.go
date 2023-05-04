package tgchats

import (
	"context"

	"github.com/samber/lo"

	"github.com/nekomeowww/insights-bot/ent"
	"github.com/nekomeowww/insights-bot/ent/telegramchatfeatureflags"
	"github.com/nekomeowww/insights-bot/pkg/types/telegram"
)

func (m *Model) findOneFeatureFlag(chatID int64, chatType telegram.ChatType) (*ent.TelegramChatFeatureFlags, error) {
	featureFlags, err := m.ent.TelegramChatFeatureFlags.
		Query().
		Where(
			telegramchatfeatureflags.ChatID(chatID),
			telegramchatfeatureflags.ChatType(string(chatType)),
		).
		First(context.Background())
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		}

		return nil, err
	}

	return featureFlags, nil
}

func (m *Model) EnableChatHistoriesRecap(chatID int64, chatType telegram.ChatType) error {
	if !lo.Contains([]telegram.ChatType{telegram.ChatTypeGroup, telegram.ChatTypeSuperGroup}, chatType) {
		return nil
	}

	featureFlags, err := m.findOneFeatureFlag(chatID, chatType)
	if err != nil {
		return err
	}
	if featureFlags == nil {
		_, err = m.ent.TelegramChatFeatureFlags.
			Create().
			SetChatID(chatID).
			SetChatType(string(chatType)).
			SetFeatureChatHistoriesRecap(true).
			Save(context.Background())
		if err != nil {
			return err
		}

		return nil
	}
	if featureFlags.FeatureChatHistoriesRecap {
		return nil
	}

	_, err = m.ent.TelegramChatFeatureFlags.
		UpdateOne(featureFlags).
		SetFeatureChatHistoriesRecap(true).
		Save(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func (m *Model) DisableChatHistoriesRecap(chatID int64, chatType telegram.ChatType) error {
	if !lo.Contains([]telegram.ChatType{telegram.ChatTypeGroup, telegram.ChatTypeSuperGroup}, chatType) {
		return nil
	}

	featureFlags, err := m.findOneFeatureFlag(chatID, chatType)
	if err != nil {
		return err
	}
	if featureFlags == nil {
		_, err = m.ent.TelegramChatFeatureFlags.
			Create().
			SetChatID(chatID).
			SetChatType(string(chatType)).
			SetFeatureChatHistoriesRecap(false).
			Save(context.Background())
		if err != nil {
			return err
		}

		return nil
	}
	if !featureFlags.FeatureChatHistoriesRecap {
		return nil
	}

	_, err = m.ent.TelegramChatFeatureFlags.
		UpdateOne(featureFlags).
		SetFeatureChatHistoriesRecap(false).
		Save(context.Background())
	if err != nil {
		return err
	}
	if err != nil {
		return err
	}

	return nil
}

func (m *Model) HasChatHistoriesRecapEnabled(chatID int64, chatType telegram.ChatType) (bool, error) {
	featureFlags, err := m.findOneFeatureFlag(chatID, chatType)
	if err != nil {
		return false, err
	}
	if featureFlags == nil {
		return false, nil
	}

	return featureFlags.FeatureChatHistoriesRecap, nil
}

func (m *Model) ListChatHistoriesRecapEnabledChats() ([]int64, error) {
	featureFlagsChats, err := m.ent.TelegramChatFeatureFlags.
		Query().
		Where(
			telegramchatfeatureflags.ChatTypeIn(string(telegram.ChatTypeGroup), string(telegram.ChatTypeSuperGroup)),
			telegramchatfeatureflags.FeatureChatHistoriesRecap(true),
		).
		All(context.Background())
	if err != nil {
		return nil, err
	}

	return lo.Map(featureFlagsChats, func(featureFlags *ent.TelegramChatFeatureFlags, _ int) int64 {
		return featureFlags.ChatID
	}), nil
}
