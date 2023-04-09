package chat_with_chat_history

import (
	"github.com/nekomeowww/insights-bot/internal/models/chat_histories"
	"github.com/nekomeowww/insights-bot/internal/models/telegram_chat_feature_flags"
	"github.com/nekomeowww/insights-bot/pkg/logger"
	"go.uber.org/fx"
)

type NewHandlerParam struct {
	fx.In

	Logger *logger.Logger

	ChatHistories            *chat_histories.ChatHistoriesModel
	TelegramChatFeatureFlags *telegram_chat_feature_flags.TelegramChatFeatureFlagsModel
}

type Handler struct {
	Logger                   *logger.Logger
	ChatHistories            *chat_histories.ChatHistoriesModel
	TelegramChatFeatureFlags *telegram_chat_feature_flags.TelegramChatFeatureFlagsModel
}

func NewHandler() func(NewHandlerParam) *Handler {
	return func(param NewHandlerParam) *Handler {
		return &Handler{
			Logger:                   param.Logger,
			ChatHistories:            param.ChatHistories,
			TelegramChatFeatureFlags: param.TelegramChatFeatureFlags,
		}
	}
}
