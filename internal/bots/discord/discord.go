package discord

import (
	"context"
	"crypto/ed25519"
	"errors"
	"github.com/samber/lo"
	"net"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/httpserver"
	"github.com/nekomeowww/insights-bot/internal/configs"
	"github.com/nekomeowww/insights-bot/internal/models/smr"
	"github.com/nekomeowww/insights-bot/pkg/logger"
	"go.uber.org/fx"
)

type NewDiscordBotParam struct {
	fx.In

	Lifecycle fx.Lifecycle

	Logger *logger.Logger
	Config *configs.Config

	Smr *smr.Model
}

type DiscordBot struct {
	logger *logger.Logger
	config *configs.Config

	smr       *smr.Model
	botClient bot.Client

	webhookStarted bool
}

func (b *DiscordBot) Check(ctx context.Context) error {
	return lo.Ternary(b.webhookStarted, nil, errors.New("discord bot service is not started yet"))
}

func NewModules() fx.Option {
	return fx.Options(
		fx.Provide(NewDiscordBot()),
	)
}

func NewDiscordBot() func(p NewDiscordBotParam) *DiscordBot {
	return func(p NewDiscordBotParam) *DiscordBot {
		cfg := p.Config.Discord

		if cfg.PublicKey == "" || cfg.Token == "" {
			p.Logger.Warn("discord: public key or bot token not provided, will not create bot instance")
			return nil
		}

		discordBot := &DiscordBot{
			logger: p.Logger,
			config: p.Config,
			smr:    p.Smr,
		}

		port := lo.Ternary(cfg.Port == "", "7072", cfg.Port)

		client, err := disgo.New(
			cfg.Token,
			bot.WithHTTPServerConfigOpts(
				cfg.PublicKey,
				httpserver.WithAddress(net.JoinHostPort("", port)),
				httpserver.WithURL("/discord/command/smr"),
			),
			bot.WithEventListenerFunc(discordBot.commandListener),
		)
		if err != nil {
			p.Logger.WithField("error", err.Error()).Fatal("discord: failed to create bot instance")
		}

		p.Lifecycle.Append(fx.Hook{
			OnStop: func(ctx context.Context) error {
				p.Logger.Info("discord: shutting down...")
				client.Close(ctx)
				discordBot.webhookStarted = false
				return nil
			},
		})
		discordBot.botClient = client

		return discordBot
	}
}

func Run() func(b *DiscordBot) error {
	return func(b *DiscordBot) error {
		// use custom ed25519 verify implementation.
		// this code is from examples of disgoorg/disgo.
		httpserver.Verify = func(publicKey httpserver.PublicKey, message, sig []byte) bool {
			return ed25519.Verify(publicKey, message, sig)
		}

		b.logger.Info("discord: registering commands...")

		//_, err := b.botClient.Rest().SetGuildCommands(b.botClient.ApplicationID(), b.config.Discord.GuildID, commands)
		_, err := b.botClient.Rest().SetGlobalCommands(b.botClient.ApplicationID(), commands)
		if err != nil {
			return err
		}

		b.logger.Info("discord: starting webhook server...")

		err = b.botClient.OpenHTTPServer()
		if err != nil {
			return err
		}

		b.webhookStarted = true

		return nil
	}
}
