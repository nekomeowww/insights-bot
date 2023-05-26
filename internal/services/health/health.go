package health

import (
	"context"
	"fmt"
	"github.com/nekomeowww/insights-bot/internal/services/smr"
	"github.com/nekomeowww/insights-bot/pkg/bots/discordbot"
	"github.com/nekomeowww/insights-bot/pkg/bots/slackbot"
	"net"
	"net/http"
	"time"

	"github.com/alexliesenfeld/health"
	"go.uber.org/fx"

	"github.com/nekomeowww/insights-bot/internal/services/autorecap"
	"github.com/nekomeowww/insights-bot/internal/services/pprof"
	"github.com/nekomeowww/insights-bot/pkg/bots/tgbot"
	"github.com/nekomeowww/insights-bot/pkg/logger"
)

type NewHealthParams struct {
	fx.In

	Lifecycle fx.Lifecycle

	Logger *logger.Logger

	TelegramBot *tgbot.BotService
	SlackBot    *slackbot.BotService
	DiscordBot  *discordbot.BotService
	AutoRecap   *autorecap.AutoRecapService
	Pprof       *pprof.Pprof
	SmrService  *smr.Service
}

type Health struct {
	server *http.Server
	logger *logger.Logger
}

func NewHealth() func(NewHealthParams) (*Health, error) {
	return func(params NewHealthParams) (*Health, error) {
		opts := make([]health.CheckerOption, 0)
		opts = append(opts,
			health.WithCacheDuration(time.Second),
			health.WithTimeout(time.Second*10),
			health.WithCheck(health.Check{
				Name:  "telegram_bot",
				Check: params.TelegramBot.Check,
			}),
			health.WithCheck(health.Check{
				Name:  "auto_recap",
				Check: params.AutoRecap.Check,
			}),
			health.WithCheck(health.Check{
				Name:  "pprof",
				Check: params.Pprof.Check,
			}),
			health.WithCheck(health.Check{
				Name:  "smr_service",
				Check: params.SmrService.Check,
			}),
		)

		if params.SlackBot != nil {
			opts = append(opts, health.WithCheck(health.Check{
				Name:  "slack_bot",
				Check: params.SlackBot.Check,
			}))
		}

		if params.DiscordBot != nil {
			opts = append(opts, health.WithCheck(health.Check{
				Name:  "discord_bot",
				Check: params.DiscordBot.Check,
			}))
		}

		checker := health.NewChecker(opts...)

		srvMux := http.NewServeMux()
		srvMux.HandleFunc("/health", health.NewHandler(checker))

		srvr := &http.Server{
			Addr:              ":7069",
			Handler:           srvMux,
			ReadHeaderTimeout: time.Second * 15,
		}

		srv := &Health{
			server: srvr,
			logger: params.Logger,
		}

		params.Lifecycle.Append(fx.Hook{
			OnStop: func(ctx context.Context) error {
				closeCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()

				if err := srv.server.Shutdown(closeCtx); err != nil && err != http.ErrServerClosed {
					return err
				}

				return nil
			},
		})

		return srv, nil
	}
}

func Run() func(health *Health) error {
	return func(health *Health) error {
		listener, err := net.Listen("tcp", health.server.Addr)
		if err != nil {
			return fmt.Errorf("failed to listen %s: %v", health.server.Addr, err)
		}

		go func() {
			if err := health.server.Serve(listener); err != nil && err != http.ErrServerClosed {
				health.logger.Fatalf("failed to serve health checker: %v", err)
			}
		}()

		time.Sleep(time.Second)

		return nil
	}
}
