package tgchats

import (
	"context"
	"time"

	"github.com/samber/lo"
	"go.uber.org/zap"

	"github.com/nekomeowww/insights-bot/ent"
	"github.com/nekomeowww/insights-bot/ent/telegramchatfeatureflags"
	"github.com/nekomeowww/insights-bot/pkg/types/telegram"
	"github.com/nekomeowww/insights-bot/pkg/types/timecapsules"
)

func (m *Model) findOneFeatureFlag(chatID int64, chatTitle string) (*ent.TelegramChatFeatureFlags, error) {
	featureFlags, err := m.ent.TelegramChatFeatureFlags.
		Query().
		Where(
			telegramchatfeatureflags.ChatID(chatID),
			telegramchatfeatureflags.ChatTypeIn(string(telegram.ChatTypeGroup), string(telegram.ChatTypeSuperGroup)),
		).
		First(context.Background())
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		}

		return nil, err
	}

	if featureFlags.ChatTitle == "" && chatTitle != "" {
		_, err = m.ent.TelegramChatFeatureFlags.
			UpdateOne(featureFlags).
			SetChatTitle(chatTitle).
			Save(context.Background())
		if err != nil {
			return nil, err
		}
	}

	return featureFlags, nil
}

func (m *Model) EnableChatHistoriesRecap(chatID int64, chatType telegram.ChatType, chatTitle string) error {
	if !lo.Contains([]telegram.ChatType{telegram.ChatTypeGroup, telegram.ChatTypeSuperGroup}, chatType) {
		return nil
	}

	featureFlags, err := m.findOneFeatureFlag(chatID, chatTitle)
	if err != nil {
		return err
	}
	if featureFlags == nil {
		_, err = m.ent.TelegramChatFeatureFlags.
			Create().
			SetChatID(chatID).
			SetChatTitle(chatTitle).
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

	m.logger.Info("enabled chat histories recap",
		zap.Int64("chat_id", chatID),
		zap.String("chat_title", chatTitle),
		zap.String("chat_type", string(chatType)),
	)

	return nil
}

func (m *Model) DisableChatHistoriesRecap(chatID int64, chatType telegram.ChatType, chatTitle string) error {
	if !lo.Contains([]telegram.ChatType{telegram.ChatTypeGroup, telegram.ChatTypeSuperGroup}, chatType) {
		return nil
	}

	featureFlags, err := m.findOneFeatureFlag(chatID, chatTitle)
	if err != nil {
		return err
	}
	if featureFlags == nil {
		_, err = m.ent.TelegramChatFeatureFlags.
			Create().
			SetChatID(chatID).
			SetChatTitle(chatTitle).
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

	m.logger.Info("disabled chat histories recap",
		zap.Int64("chat_id", chatID),
		zap.String("chat_title", chatTitle),
		zap.String("chat_type", string(chatType)),
	)

	return nil
}

func (m *Model) HasChatHistoriesRecapEnabled(chatID int64, chatTitle string) (bool, error) {
	featureFlags, err := m.findOneFeatureFlag(chatID, chatTitle)
	if err != nil {
		return false, err
	}
	if featureFlags == nil {
		return false, nil
	}

	return featureFlags.FeatureChatHistoriesRecap, nil
}

func (m *Model) ListChatHistoriesRecapEnabledChats() ([]*ent.TelegramChatFeatureFlags, error) {
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

	return featureFlagsChats, nil
}

func (m *Model) QueueSendChatHistoriesRecapTask() {
	chats, err := m.ListChatHistoriesRecapEnabledChats()
	if err != nil {
		m.logger.Error("failed to list chat histories recap enabled chats", zap.Error(err))
		return
	}

	for _, chat := range chats {
		err = m.QueueOneSendChatHistoriesRecapTaskForChatID(chat.ChatID)
		if err != nil {
			m.logger.Error("failed to queue send chat histories recap task", zap.Error(err))
			continue
		}
	}
}

func (m *Model) QueueOneSendChatHistoriesRecapTaskForChatID(chatID int64) error {
	location := time.UTC
	if m.config.TimezoneShiftSeconds != 0 {
		location = time.FixedZone("Local", int(m.config.TimezoneShiftSeconds))
	}

	now := time.
		Now().       // Current time.
		UTC().       // Resets to UTC.
		In(location) // Align current timezone with the configured offset (if any) for later calculation.

	scheduleTargets := []int64{2, 8, 14, 20} // queue for 02:00, 08:00, 14:00, 20:00
	scheduleSets := make([]time.Time, 0, len(scheduleTargets))

	for _, target := range scheduleTargets {
		if now.Hour() < int(target) {
			scheduleSets = append(scheduleSets, time.Date(now.Year(), now.Month(), now.Day(), int(target), 0, 0, 0, location))
			break
		}
	}
	if len(scheduleSets) == 0 {
		scheduleSets = append(scheduleSets, time.Date(now.Year(), now.Month(), now.Day()+1, int(scheduleTargets[0]), 0, 0, 0, location))
	}

	for _, schedule := range scheduleSets {
		m.logger.Info("scheduled one send chat histories recap task for chat", zap.Int64("chat_id", chatID), zap.Time("schedule", schedule))

		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)

		err := m.digger.BuryUtil(ctx, timecapsules.AutoRecapCapsule{
			ChatID: chatID,
		}, schedule.UnixMilli())
		if err != nil {
			m.logger.Error("failed to bury one send chat histories recap task for chat", zap.Int64("chat_id", chatID), zap.Error(err))
		}

		cancel()
	}

	return nil
}
