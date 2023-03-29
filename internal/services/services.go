package services

import (
	"github.com/nekomeowww/insights-bot/internal/services/chat_history_recap"
	"go.uber.org/fx"
)

func NewModules() fx.Option {
	return fx.Options(
		fx.Provide(chat_history_recap.NewChatHistoryRecapService()),
	)
}
