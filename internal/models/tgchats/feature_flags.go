package tgchats

import (
	"github.com/ostafen/clover/v2"
	"github.com/samber/lo"

	"github.com/nekomeowww/insights-bot/pkg/types/telegram"
	"github.com/nekomeowww/insights-bot/pkg/types/telegram/tgchat"
)

func (m *Model) findOneFeatureFlag(chatID int64, chatType telegram.ChatType) (*tgchat.FeatureFlag, error) {
	query := clover.
		NewQuery(tgchat.FeatureFlag{}.CollectionName()).
		Where(clover.Field("chat_id").Eq(chatID)).
		Where(clover.Field("chat_type").In(string(telegram.ChatTypeGroup), string(telegram.ChatTypeSuperGroup)))

	doc, err := m.Clover.FindFirst(query)
	if err != nil {
		return nil, err
	}
	if doc == nil {
		return nil, nil
	}

	var featureFlags tgchat.FeatureFlag
	err = doc.Unmarshal(&featureFlags)
	if err != nil {
		return nil, err
	}

	return &featureFlags, nil
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
		_, err = m.Clover.InsertOne(
			tgchat.FeatureFlag{}.CollectionName(),
			clover.NewDocumentOf(tgchat.FeatureFlag{
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
		tgchat.FeatureFlag{}.CollectionName(),
		featureFlags.ID,
		updates,
	)
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
		_, err = m.Clover.InsertOne(
			tgchat.FeatureFlag{}.CollectionName(),
			clover.NewDocumentOf(tgchat.FeatureFlag{
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
		tgchat.FeatureFlag{}.CollectionName(),
		featureFlags.ID,
		updates,
	)
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
	query := clover.
		NewQuery(tgchat.FeatureFlag{}.CollectionName()).
		Where(clover.Field("chat_type").In(string(telegram.ChatTypeGroup), string(telegram.ChatTypeSuperGroup))).
		Where(clover.Field("feature_chat_histories_recap").Eq(true))

	docs, err := m.Clover.FindAll(query)
	if err != nil {
		return make([]int64, 0), err
	}
	if len(docs) == 0 {
		return make([]int64, 0), nil
	}

	featureFlagsChats := make([]*tgchat.FeatureFlag, 0, len(docs))
	for _, doc := range docs {
		var featureFlags tgchat.FeatureFlag
		err = doc.Unmarshal(&featureFlags)
		if err != nil {
			return make([]int64, 0), err
		}

		featureFlagsChats = append(featureFlagsChats, &featureFlags)
	}

	return lo.Map(featureFlagsChats, func(featureFlags *tgchat.FeatureFlag, _ int) int64 {
		return featureFlags.ChatID
	}), nil
}
