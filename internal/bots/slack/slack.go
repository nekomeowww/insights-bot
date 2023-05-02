package slack

import (
	"context"
	"errors"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nekomeowww/insights-bot/internal/configs"
	"github.com/nekomeowww/insights-bot/internal/models/smr"
	"github.com/nekomeowww/insights-bot/pkg/logger"
	"github.com/slack-go/slack"
	"go.uber.org/fx"
)

func NewModules() fx.Option {
	return fx.Options(
		fx.Provide(NewSlackBot()),
	)
}

type NewSlackBotParam struct {
	fx.In

	Lifecycle fx.Lifecycle

	Config *configs.Config
	Logger *logger.Logger

	SMR *smr.Model
}

type SlackBot struct {
	Config *configs.Config
	Logger *logger.Logger

	smrModel *smr.Model

	server   *http.Server
	slackCli *slack.Client

	alreadyClosed bool
	closeChan     chan struct{}

	processChan chan recivedCommandInfo
}

func NewSlackBot() func(param NewSlackBotParam) *SlackBot {
	return func(param NewSlackBotParam) *SlackBot {
		if param.Config.SlackBotToken == "" {
			param.Logger.Warn("slack bot token not provided, will not create bot instance")
			return nil
		}

		slackBot := &SlackBot{
			Config:      param.Config,
			Logger:      param.Logger,
			closeChan:   make(chan struct{}, 1),
			processChan: make(chan recivedCommandInfo, 10),
			slackCli:    slack.New(param.Config.SlackBotToken),
			smrModel:    param.SMR,
		}

		_, err := slackBot.slackCli.AuthTest()
		if err != nil {
			param.Logger.WithError(err).Fatalf("slack bot token auth test failed")
			return nil
		}

		engine := gin.Default()
		engine.POST("/slack/command/smr", slackBot.postCommandInfo)
		slackBot.server = &http.Server{Addr: ":7070", Handler: engine}

		param.Lifecycle.Append(fx.Hook{
			OnStop: func(ctx context.Context) error {
				err := slackBot.server.Shutdown(ctx)

				if !errors.Is(err, context.Canceled) {
					param.Logger.WithField("error", err.Error()).Error("slack bot server shutdown failed")
					return err
				}
				slackBot.Logger.Info("stopped to receiving new requests")

				slackBot.alreadyClosed = true
				slackBot.closeChan <- struct{}{}

				return nil
			},
		})

		return slackBot
	}
}

func Run() func(bot *SlackBot) error {
	return func(bot *SlackBot) error {
		if bot == nil {
			return nil
		}

		listener, err := net.Listen("tcp", bot.server.Addr)
		if err != nil {
			return err
		}

		go func() {
			err = bot.server.Serve(listener)
			if err != nil && err != http.ErrServerClosed {
				bot.Logger.WithField("error", err.Error()).Fatal("slack bot server error")
			}
		}()

		go bot.runSmr()
		return nil
	}
}
