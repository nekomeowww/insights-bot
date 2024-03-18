package slack

import (
	"context"
	"net/http"

	"github.com/nekomeowww/insights-bot/internal/bots/slack/handlers"
	"github.com/nekomeowww/insights-bot/internal/configs"
	"github.com/nekomeowww/insights-bot/pkg/bots/slackbot"
	"github.com/nekomeowww/insights-bot/pkg/bots/slackbot/services"
	"github.com/nekomeowww/insights-bot/pkg/logger"
	"go.uber.org/fx"
)

func NewModules() fx.Option {
	return fx.Options(
		fx.Provide(services.NewServices()),
		fx.Provide(NewSlackBot()),
		fx.Options(handlers.NewModules()),
	)
}

type NewSlackBotParam struct {
	fx.In

	Lifecycle fx.Lifecycle

	Config *configs.Config
	Logger *logger.Logger

	Handlers *handlers.Handlers
	Services *services.Services
}

func NewSlackBot() func(param NewSlackBotParam) *slackbot.BotService {
	return func(param NewSlackBotParam) *slackbot.BotService {
		slackConfig := param.Config.Slack

		if slackConfig.ClientID == "" || slackConfig.ClientSecret == "" {
			param.Logger.Warn("slack client id or secret not provided, will not create bot instance")
			return nil
		}

		bot := slackbot.NewBotService(slackConfig)
		bot.Handle(http.MethodPost, "/slack/command/smr", param.Handlers.PostCommandInfo)
		bot.Handle(http.MethodGet, "/slack/install/auth", param.Handlers.GetInstallAuth)
		bot.SetService(param.Services)
		bot.SetLogger(param.Logger)

		param.Lifecycle.Append(fx.Hook{
			OnStop: func(ctx context.Context) error {
				return bot.Stop(ctx)
			},
		})

		return bot
	}
}

func Run() func(bot *slackbot.BotService) error {
	return func(bot *slackbot.BotService) error {
		if bot == nil {
			return nil
		}

		return bot.Run()
	}
}
