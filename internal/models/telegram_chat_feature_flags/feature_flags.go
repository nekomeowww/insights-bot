package telegram_chat_feature_flags

import (
	"github.com/ostafen/clover/v2"
	"github.com/samber/lo"
	"go.uber.org/fx"

	"github.com/nekomeowww/insights-bot/internal/datastore"
	"github.com/nekomeowww/insights-bot/pkg/types/telegram"
	"github.com/nekomeowww/insights-bot/pkg/types/telegram_chat_feature_flag"
)

type NewTelegramChatFeatureFlagsModelParam struct {
	fx.In

	Clover *datastore.Clover
}

type TelegramChatFeatureFlagsModel struct {
	Clover *datastore.Clover
}

func NewFeatureFlagsModel() func(NewTelegramChatFeatureFlagsModelParam) (*TelegramChatFeatureFlagsModel, error) {
	return func(param NewTelegramChatFeatureFlagsModelParam) (*TelegramChatFeatureFlagsModel, error) {
		hasCollection, err := param.Clover.HasCollection(telegram_chat_feature_flag.TelegramChatFeatureFlag{}.CollectionName())
		if err != nil {
			return nil, err
		}
		if !hasCollection {
			err = param.Clover.CreateCollection(telegram_chat_feature_flag.TelegramChatFeatureFlag{}.CollectionName())
			if err != nil {
				return nil, err
			}
		}

		return &TelegramChatFeatureFlagsModel{
			Clover: param.Clover,
		}, nil
	}
}

func (m *TelegramChatFeatureFlagsModel) findOneFeatureFlag(chatID int64, chatType telegram.ChatType) (*telegram_chat_feature_flag.TelegramChatFeatureFlag, error) {
	query := clover.
		NewQuery(telegram_chat_feature_flag.TelegramChatFeatureFlag{}.CollectionName()).
		Where(clover.Field("chat_id").Eq(chatID)).
		Where(clover.Field("chat_type").In(string(telegram.ChatTypeGroup), string(telegram.ChatTypeSuperGroup)))

	doc, err := m.Clover.FindFirst(query)
	if err != nil {
		return nil, err
	}
	if doc == nil {
		return nil, nil
	}

	var featureFlags telegram_chat_feature_flag.TelegramChatFeatureFlag
	err = doc.Unmarshal(&featureFlags)
	if err != nil {
		return nil, err
	}

	return &featureFlags, nil
}

func (m *TelegramChatFeatureFlagsModel) EnableChatHistoriesRecap(chatID int64, chatType telegram.ChatType) error {
	if !lo.Contains([]telegram.ChatType{telegram.ChatTypeGroup, telegram.ChatTypeSuperGroup}, chatType) {
		return nil
	}

	featureFlags, err := m.findOneFeatureFlag(chatID, chatType)
	if err != nil {
		return err
	}
	if featureFlags == nil {
		_, err = m.Clover.InsertOne(
			telegram_chat_feature_flag.TelegramChatFeatureFlag{}.CollectionName(),
			clover.NewDocumentOf(telegram_chat_feature_flag.TelegramChatFeatureFlag{
				ID:                        clover.NewObjectId(),
				ChatID:                    chatID,
				ChatType:                  chatType,
				FeatureChatHistoriesRecap: true,
			}),
		)
		if err != nil {
			return err
		}

		return nil
	}
	if featureFlags.FeatureChatHistoriesRecap {
		return nil
	}

	updates := make(map[string]any)
	updates["feature_chat_histories_recap"] = true
	err = m.Clover.UpdateById(
		telegram_chat_feature_flag.TelegramChatFeatureFlag{}.CollectionName(),
		featureFlags.ID,
		updates,
	)
	if err != nil {
		return err
	}

	return nil
}

func (m *TelegramChatFeatureFlagsModel) DisableChatHistoriesRecap(chatID int64, chatType telegram.ChatType) error {
	if !lo.Contains([]telegram.ChatType{telegram.ChatTypeGroup, telegram.ChatTypeSuperGroup}, chatType) {
		return nil
	}

	featureFlags, err := m.findOneFeatureFlag(chatID, chatType)
	if err != nil {
		return err
	}
	if featureFlags == nil {
		_, err = m.Clover.InsertOne(
			telegram_chat_feature_flag.TelegramChatFeatureFlag{}.CollectionName(),
			clover.NewDocumentOf(telegram_chat_feature_flag.TelegramChatFeatureFlag{
				ID:                        clover.NewObjectId(),
				ChatID:                    chatID,
				ChatType:                  chatType,
				FeatureChatHistoriesRecap: false,
			}),
		)
		if err != nil {
			return err
		}

		return nil
	}
	if !featureFlags.FeatureChatHistoriesRecap {
		return nil
	}

	updates := make(map[string]any)
	updates["feature_chat_histories_recap"] = false
	err = m.Clover.UpdateById(
		telegram_chat_feature_flag.TelegramChatFeatureFlag{}.CollectionName(),
		featureFlags.ID,
		updates,
	)
	if err != nil {
		return err
	}

	return nil
}

func (m *TelegramChatFeatureFlagsModel) HasChatHistoriesRecapEnabled(chatID int64, chatType telegram.ChatType) (bool, error) {
	featureFlags, err := m.findOneFeatureFlag(chatID, chatType)
	if err != nil {
		return false, err
	}
	if featureFlags == nil {
		return false, nil
	}

	return featureFlags.FeatureChatHistoriesRecap, nil
}

func (m *TelegramChatFeatureFlagsModel) ListChatHistoriesRecapEnabledChats() ([]int64, error) {
	query := clover.
		NewQuery(telegram_chat_feature_flag.TelegramChatFeatureFlag{}.CollectionName()).
		Where(clover.Field("chat_type").In(string(telegram.ChatTypeGroup), string(telegram.ChatTypeSuperGroup))).
		Where(clover.Field("feature_chat_histories_recap").Eq(true))

	docs, err := m.Clover.FindAll(query)
	if err != nil {
		return make([]int64, 0), err
	}
	if len(docs) == 0 {
		return make([]int64, 0), nil
	}

	featureFlagsChats := make([]*telegram_chat_feature_flag.TelegramChatFeatureFlag, 0, len(docs))
	for _, doc := range docs {
		var featureFlags telegram_chat_feature_flag.TelegramChatFeatureFlag
		err = doc.Unmarshal(&featureFlags)
		if err != nil {
			return make([]int64, 0), err
		}

		featureFlagsChats = append(featureFlagsChats, &featureFlags)
	}

	return lo.Map(featureFlagsChats, func(featureFlags *telegram_chat_feature_flag.TelegramChatFeatureFlag, _ int) int64 {
		return featureFlags.ChatID
	}), nil
}
