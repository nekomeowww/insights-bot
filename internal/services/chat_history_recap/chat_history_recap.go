package chat_history_recap

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	cron "github.com/robfig/cron/v3"
	"go.uber.org/fx"

	"github.com/nekomeowww/insights-bot/internal/bots/telegram"
	"github.com/nekomeowww/insights-bot/internal/models/chat_histories"
	"github.com/nekomeowww/insights-bot/internal/models/telegram_chat_feature_flags"
	"github.com/nekomeowww/insights-bot/pkg/logger"
	"github.com/nekomeowww/insights-bot/pkg/openai"
)

type NewChatHistoryRecapServiceParam struct {
	fx.In

	Lifecycle fx.Lifecycle

	Bot                      *telegram.Bot
	Logger                   *logger.Logger
	ChatHistories            *chat_histories.ChatHistoriesModel
	TelegramChatFeatureFlags *telegram_chat_feature_flags.TelegramChatFeatureFlagsModel
	OpenAI                   *openai.Client
}

type ChatHistoryRecapService struct {
	Cron *cron.Cron

	Bot                      *telegram.Bot
	Logger                   *logger.Logger
	ChatHistories            *chat_histories.ChatHistoriesModel
	TelegramChatFeatureFlags *telegram_chat_feature_flags.TelegramChatFeatureFlagsModel
	OpenAI                   *openai.Client
}

func NewChatHistoryRecapService() func(NewChatHistoryRecapServiceParam) *ChatHistoryRecapService {
	return func(param NewChatHistoryRecapServiceParam) *ChatHistoryRecapService {
		service := &ChatHistoryRecapService{
			Cron:                     cron.New(),
			Bot:                      param.Bot,
			Logger:                   param.Logger,
			ChatHistories:            param.ChatHistories,
			TelegramChatFeatureFlags: param.TelegramChatFeatureFlags,
			OpenAI:                   param.OpenAI,
		}

		service.Cron.AddFunc("@hourly", service.SendChatHistoriesRecap)
		param.Lifecycle.Append(fx.Hook{
			OnStop: func(context.Context) error {
				service.Cron.Stop()
				return nil
			},
		})

		return service
	}
}

func Run() func(service *ChatHistoryRecapService) {
	return func(service *ChatHistoryRecapService) {
		service.Cron.Start()
	}
}

func (s *ChatHistoryRecapService) SendChatHistoriesRecap() {
	chatIDs, err := s.TelegramChatFeatureFlags.ListChatHistoriesRecapEnabledChats()
	if err != nil {
		s.Logger.Errorf("failed to list chat histories recap enabled chats: %v", err)
		return
	}

	for _, chatID := range chatIDs {
		s.Logger.Infof("generating chat histories recap for chat %d", chatID)

		summarization, err := s.ChatHistories.SummarizeLastOneHourChatHistories(chatID)
		if err != nil {
			s.Logger.Errorf("failed to summarize last one hour chat histories: %v", err)
			continue
		}
		if summarization == "" {
			s.Logger.Warn("summarization is empty")
			continue
		}

		s.Logger.Info("sending chat histories recap for chat %d", chatID)
		message := tgbotapi.NewMessage(chatID, summarization)
		_, err = s.Bot.Send(message)
		if err != nil {
			s.Logger.Errorf("failed to send chat histories recap: %v", err)
			continue
		}
	}
}
