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

func (m *Model) findOneFeatureFlagForGroups(chatID int64, chatTitle string) (*ent.TelegramChatFeatureFlags, error) {
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

func (m *Model) findOrCreateFeatureFlagForGroups(chatID int64, chatType telegram.ChatType, chatTitle string) (*ent.TelegramChatFeatureFlags, error) {
	featureFlags, err := m.findOneFeatureFlagForGroups(chatID, chatTitle)
	if err != nil {
		return nil, err
	}
	if featureFlags == nil {
		createdFeatureFlags, err := m.ent.TelegramChatFeatureFlags.
			Create().
			SetChatID(chatID).
			SetChatTitle(chatTitle).
			SetChatType(string(chatType)).
			SetFeatureChatHistoriesRecap(false).
			SetFeatureLanguage("en").
			Save(context.Background())
		if err != nil {
			return nil, err
		}

		return createdFeatureFlags, nil
	}

	return featureFlags, nil
}

func (m *Model) FindLanguageForGroups(chatID int64, chatTitle string) (string, error) {
	featureFlags, err := m.findOneFeatureFlagForGroups(chatID, chatTitle)
	if err != nil {
		return "en", err
	}
	if featureFlags == nil {
		return "en", nil
	}

	return featureFlags.FeatureLanguage, nil
}

func (m *Model) SetLanguageForGroups(chatID int64, chatType telegram.ChatType, chatTitle string, language string) error {
	if !lo.Contains([]telegram.ChatType{telegram.ChatTypeGroup, telegram.ChatTypeSuperGroup}, chatType) {
		return nil
	}

	featureFlags, err := m.findOrCreateFeatureFlagForGroups(chatID, chatType, chatTitle)
	if err != nil {
		return err
	}
	if featureFlags.FeatureLanguage == language {
		return nil
	}

	_, err = m.ent.TelegramChatFeatureFlags.
		UpdateOne(featureFlags).
		SetFeatureLanguage(language).
		Save(context.Background())
	if err != nil {
		return err
	}

	m.logger.Info("set language for chat",
		zap.Int64("chat_id", chatID),
		zap.String("chat_title", chatTitle),
		zap.String("chat_type", string(chatType)),
		zap.String("language", language),
	)

	return nil
}

func (m *Model) EnableChatHistoriesRecapForGroups(chatID int64, chatType telegram.ChatType, chatTitle string) error {
	if !lo.Contains([]telegram.ChatType{telegram.ChatTypeGroup, telegram.ChatTypeSuperGroup}, chatType) {
		return nil
	}

	featureFlags, err := m.findOrCreateFeatureFlagForGroups(chatID, chatType, chatTitle)
	if err != nil {
		return err
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

func (m *Model) DisableChatHistoriesRecapForGroups(chatID int64, chatType telegram.ChatType, chatTitle string) error {
	if !lo.Contains([]telegram.ChatType{telegram.ChatTypeGroup, telegram.ChatTypeSuperGroup}, chatType) {
		return nil
	}

	featureFlags, err := m.findOneFeatureFlagForGroups(chatID, chatTitle)
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

func (m *Model) HasChatHistoriesRecapEnabledForGroups(chatID int64, chatTitle string) (bool, error) {
	featureFlags, err := m.findOneFeatureFlagForGroups(chatID, chatTitle)
	if err != nil {
		return false, err
	}
	if featureFlags == nil {
		return false, nil
	}

	return featureFlags.FeatureChatHistoriesRecap, nil
}

func (m *Model) HasJoinedGroupsBefore(chatID int64, chatTitle string) (bool, error) {
	featureFlags, err := m.findOneFeatureFlagForGroups(chatID, chatTitle)
	if err != nil {
		return false, err
	}

	return featureFlags != nil, nil
}

func (m *Model) ListChatHistoriesRecapEnabledChatsForGroups() ([]*ent.TelegramChatFeatureFlags, error) {
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
	chats, err := m.ListChatHistoriesRecapEnabledChatsForGroups()
	if err != nil {
		m.logger.Error("failed to list chat histories recap enabled chats", zap.Error(err))
		return
	}

	mChatOptions := make(map[int64]*ent.TelegramChatRecapsOptions)

	for _, chat := range chats {
		options, err := m.FindOneOrCreateRecapsOption(chat.ChatID)
		if err != nil {
			m.logger.Error("failed to find one or create recaps option",
				zap.Int64("chat_id", chat.ChatID),
				zap.Error(err),
			)

			continue
		}

		mChatOptions[chat.ChatID] = options
	}

	for _, chat := range chats {
		options, ok := mChatOptions[chat.ChatID]
		if !ok {
			continue
		}

		err := m.QueueOneSendChatHistoriesRecapTaskForChatID(chat.ChatID, options)
		if err != nil {
			continue
		}
	}
}

var MapScheduleHours = map[int][]int64{
	2: {8, 20},        // queue for 08:00, 20:00
	3: {0, 8, 16},     // queue for 00:00, 08:00, 16:00
	4: {2, 8, 14, 20}, // queue for 02:00, 08:00, 14:00, 20:00
}

func (m *Model) newNextScheduleTimeForChatHistoriesRecapTasksForChatID(_ int64, rate int) time.Time {
	location := time.UTC
	if m.config.TimezoneShiftSeconds != 0 {
		location = time.FixedZone("Local", int(m.config.TimezoneShiftSeconds))
	}

	now := time.
		Now().       // Current time.
		UTC().       // Resets to UTC.
		In(location) // Align current timezone with the configured offset (if any) for later calculation.

	scheduleTargets, ok := MapScheduleHours[rate]
	if !ok {
		scheduleTargets = MapScheduleHours[4]
	}

	var nextScheduleTimeSet bool
	var nextScheduleTime time.Time

	for _, target := range scheduleTargets {
		if now.Hour() < int(target) {
			nextScheduleTime = time.Date(now.Year(), now.Month(), now.Day(), int(target), 0, 0, 0, location)
			nextScheduleTimeSet = true

			break
		}
	}

	if !nextScheduleTimeSet {
		nextScheduleTime = time.Date(now.Year(), now.Month(), now.Day()+1, int(scheduleTargets[0]), 0, 0, 0, location)
	}

	return nextScheduleTime
}

func (m *Model) queueOneSendChatHistoriesRecapTaskForChatIDBasedOnScheduleSets(chatID int64, nextScheduleTime time.Time) error {
	m.logger.Info("scheduled one send chat histories recap task for chat",
		zap.Int64("chat_id", chatID),
		zap.Time("schedule", nextScheduleTime),
	)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	err := m.digger.BuryUtil(ctx, timecapsules.AutoRecapCapsule{
		ChatID: chatID,
	}, nextScheduleTime.UnixMilli())
	if err != nil {
		m.logger.Error("failed to bury one send chat histories recap task for chat",
			zap.Int64("chat_id", chatID),
			zap.Time("schedule", nextScheduleTime),
			zap.Error(err),
		)
	}

	return nil
}

func (m *Model) QueueOneSendChatHistoriesRecapTaskForChatID(chatID int64, options *ent.TelegramChatRecapsOptions) error {
	if !lo.Contains([]int{2, 3, 4}, options.AutoRecapRatesPerDay) {
		m.logger.Error("invalid auto recap rates per day, fallbacks, to 4 times a day",
			zap.Int64("chat_id", chatID),
			zap.Int("auto_recap_rates", options.AutoRecapRatesPerDay),
		)

		options.AutoRecapRatesPerDay = 4
	}

	nextScheduleTime := m.newNextScheduleTimeForChatHistoriesRecapTasksForChatID(chatID, options.AutoRecapRatesPerDay)

	err := m.queueOneSendChatHistoriesRecapTaskForChatIDBasedOnScheduleSets(chatID, nextScheduleTime)
	if err != nil {
		m.logger.Error("failed to queue send chat histories recap task",
			zap.Int64("chat_id", chatID),
			zap.Int("auto_recap_rates", options.AutoRecapRatesPerDay),
			zap.Error(err),
		)

		return err
	}

	return nil
}

func (m *Model) DeleteOneFeatureFlagByChatID(chatID int64) error {
	_, err := m.ent.TelegramChatFeatureFlags.
		Delete().
		Where(telegramchatfeatureflags.ChatID(chatID)).
		Exec(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func (m *Model) MigrateFeatureFlagsOfChatFromChatIDToChatID(fromChatID int64, toChatID int64) error {
	affectedRows, err := m.ent.TelegramChatFeatureFlags.
		Update().
		Where(
			telegramchatfeatureflags.ChatIDEQ(fromChatID),
		).
		SetChatID(toChatID).
		SetChatType(string(telegram.ChatTypeSuperGroup)).
		Save(context.Background())
	if err != nil {
		return err
	}

	m.logger.Info("successfully migrated feature flags of chat",
		zap.Int64("from_chat_id", fromChatID),
		zap.Int64("to_chat_id", toChatID),
		zap.Int("affected_rows", affectedRows),
	)

	return nil
}
