package telegram

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/fx"

	"github.com/nekomeowww/insights-bot/internal/bots/telegram/handlers"
	"github.com/nekomeowww/insights-bot/internal/bots/telegram/middlewares"
	"github.com/nekomeowww/insights-bot/internal/configs"
	"github.com/nekomeowww/insights-bot/internal/models/chat_histories"
	"github.com/nekomeowww/insights-bot/internal/models/tgchats"
	"github.com/nekomeowww/insights-bot/pkg/bots/tgbot"
	"github.com/nekomeowww/insights-bot/pkg/logger"
	"github.com/nekomeowww/insights-bot/pkg/utils"
)

func NewModules() fx.Option {
	return fx.Options(
		fx.Provide(NewBot()),
		fx.Provide(tgbot.NewDispatcher()),
		fx.Options(handlers.NewModules()),
	)
}

type NewBotParam struct {
	fx.In

	Lifecycle fx.Lifecycle

	Config     *configs.Config
	Logger     *logger.Logger
	Dispatcher *tgbot.Dispatcher
	Handlers   *handlers.Handlers

	ChatHistories *chat_histories.Model
	TgChats       *tgchats.Model
}

type Bot struct {
	*tgbotapi.BotAPI

	Config     *configs.Config
	Logger     *logger.Logger
	Dispatcher *tgbot.Dispatcher

	ChatHistories *chat_histories.Model

	alreadyClose bool
	closeChan    chan struct{}
}

func NewBot() func(param NewBotParam) (*Bot, error) {
	return func(param NewBotParam) (*Bot, error) {
		if param.Config.TelegramBotToken == "" {
			param.Logger.Fatal("must supply a valid telegram bot token in configs or environment variable")
		}

		b, err := tgbotapi.NewBotAPI(param.Config.TelegramBotToken)
		if err != nil {
			return nil, err
		}

		bot := &Bot{
			BotAPI:        b,
			Logger:        param.Logger,
			Dispatcher:    param.Dispatcher,
			ChatHistories: param.ChatHistories,
			closeChan:     make(chan struct{}, 1),
		}

		param.Lifecycle.Append(fx.Hook{
			OnStop: func(ctx context.Context) error {
				bot.stopPull(ctx)
				return nil
			},
		})

		param.Dispatcher.Use(middlewares.RecordMessage(param.ChatHistories, param.TgChats))
		param.Handlers.InstallAll()
		param.Logger.Infof("Authorized as bot @%s", bot.Self.UserName)
		return bot, nil
	}
}

func (b *Bot) stopPull(ctx context.Context) {
	if b.alreadyClose {
		return
	}

	_ = utils.Invoke0(func() error {
		b.alreadyClose = true
		b.StopReceivingUpdates()
		b.closeChan <- struct{}{}
		close(b.closeChan)

		return nil
	}, utils.WithContext(ctx))
}

func (b *Bot) pullUpdates() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.GetUpdatesChan(u)
	for {
		if b.alreadyClose {
			b.Logger.Info("stopped to receiving updates")
			return
		}

		select {
		case update := <-updates:
			b.Dispatcher.Dispatch(b.BotAPI, update)
		case <-b.closeChan:
			b.Logger.Info("stopped to receiving updates")
			return
		}
	}
}

func Run() func(bot *Bot) {
	return func(bot *Bot) {
		go bot.pullUpdates()
	}
}
