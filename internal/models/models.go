package models

import (
	"github.com/nekomeowww/insights-bot/internal/models/chat_histories"
	"github.com/nekomeowww/insights-bot/internal/models/telegram_chat_feature_flags"
	"go.uber.org/fx"
)

func NewModules() fx.Option {
	return fx.Options(
		fx.Provide(chat_histories.NewChatHistoriesModel()),
		fx.Provide(telegram_chat_feature_flags.NewFeatureFlagsModel()),
	)
}
