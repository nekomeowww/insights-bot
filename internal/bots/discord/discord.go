package discord

import (
	"context"
	"github.com/nekomeowww/insights-bot/internal/bots/discord/listeners"
	"github.com/nekomeowww/insights-bot/internal/configs"
	"github.com/nekomeowww/insights-bot/internal/services/smr"
	"github.com/nekomeowww/insights-bot/pkg/bots/discordbot"
	"github.com/nekomeowww/insights-bot/pkg/logger"
	"go.uber.org/fx"
)

func NewModules() fx.Option {
	return fx.Options(
		fx.Provide(NewDiscordBot()),
		fx.Options(listeners.NewModules()),
	)
}

type NewDiscordBotParam struct {
	fx.In

	Lifecycle fx.Lifecycle

	Logger *logger.Logger
	Config *configs.Config

	Smr *smr.Service

	Listeners *listeners.Listeners
}

func NewDiscordBot() func(p NewDiscordBotParam) *discordbot.BotService {
	return func(p NewDiscordBotParam) *discordbot.BotService {
		cfg := p.Config.Discord

		if cfg.PublicKey == "" || cfg.Token == "" {
			p.Logger.Warn("discordbot: public key or bot token not provided, will not create bot instance")
			return nil
		}

		s := discordbot.NewBotService(p.Listeners.CommandListener, cfg, p.Logger)

		p.Lifecycle.Append(fx.Hook{
			OnStop: func(ctx context.Context) error {
				s.Stop(ctx)
				return nil
			},
		})

		p.Smr.SetDiscordBot(s)

		return s
	}
}

func Run() func(b *discordbot.BotService) error {
	return func(b *discordbot.BotService) error {
		if b == nil {
			return nil
		}

		return b.Run()
	}
}
