package discordbot

import (
	"context"
	"crypto/ed25519"
	"errors"
	"net"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/httpserver"
	"github.com/nekomeowww/insights-bot/internal/configs"
	"github.com/nekomeowww/insights-bot/pkg/logger"
	"github.com/samber/lo"
)

type BotService struct {
	bot.Client

	logger *logger.Logger

	webhookStarted bool
}

func NewBotService[E bot.Event](
	f func(e E),
	cfg configs.SectionDiscord,
	logger *logger.Logger,
) *BotService {
	discordBot := &BotService{
		logger: logger,
	}

	port := lo.Ternary(cfg.Port == "", "7072", cfg.Port)

	client, err := disgo.New(
		cfg.Token,
		bot.WithHTTPServerConfigOpts(
			cfg.PublicKey,
			httpserver.WithAddress(net.JoinHostPort("", port)),
			httpserver.WithURL("/discord/command/smr"),
		),
		bot.WithEventListenerFunc(f),
	)
	if err != nil {
		logger.WithField("error", err.Error()).Fatal("discord: failed to create bot instance")
	}

	discordBot.Client = client

	return discordBot
}

func (b *BotService) SetLogger(logger *logger.Logger) {
	b.logger = logger
}

func (b *BotService) Check(ctx context.Context) error {
	return lo.Ternary(b.webhookStarted, nil, errors.New("discord bot service is not started yet"))
}

func (b *BotService) Run() error {
	// use custom ed25519 verify implementation.
	// this code is from examples of disgoorg/disgo.
	httpserver.Verify = func(publicKey httpserver.PublicKey, message, sig []byte) bool {
		return ed25519.Verify(publicKey, message, sig)
	}

	b.logger.Info("discord: registering commands...")

	_, err := b.Rest().SetGlobalCommands(b.ApplicationID(), commands)
	if err != nil {
		return err
	}

	b.logger.Info("discord: starting webhook server...")

	err = b.OpenHTTPServer()
	if err != nil {
		return err
	}

	b.webhookStarted = true

	return nil
}

func (b *BotService) Stop(ctx context.Context) {
	b.logger.Info("discord: shutting down...")
	b.Close(ctx)
	b.webhookStarted = false
}
