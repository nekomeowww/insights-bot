package tgchat

import "github.com/nekomeowww/insights-bot/pkg/types/telegram"

type FeatureFlag struct {
	ID                        string            `clover:"_id"`
	ChatID                    int64             `clover:"chat_id"`
	ChatType                  telegram.ChatType `clover:"chat_type"`
	FeatureChatHistoriesRecap bool              `clover:"feature_chat_histories_recap"`
}

func (FeatureFlag) CollectionName() string {
	return "telegram_chat_feature_flags"
}
