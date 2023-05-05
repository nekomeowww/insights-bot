package slack

import (
	"context"
	"errors"
	"net"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nekomeowww/insights-bot/internal/configs"
	"github.com/nekomeowww/insights-bot/internal/datastore"
	"github.com/nekomeowww/insights-bot/internal/models/smr"
	"github.com/nekomeowww/insights-bot/pkg/logger"
	"github.com/samber/lo"
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

	Ent *datastore.Ent
}

type SlackBot struct {
	config *configs.Config
	logger *logger.Logger

	smrModel *smr.Model

	server *http.Server

	ent *datastore.Ent

	alreadyClosed bool
	closeChan     chan struct{}

	processChan chan smrRequestInfo
}

func NewSlackBot() func(param NewSlackBotParam) *SlackBot {
	return func(param NewSlackBotParam) *SlackBot {
		slackConfig := param.Config.Slack

		if slackConfig.ClientID == "" || slackConfig.ClientSecret == "" {
			param.Logger.Warn("slack client id or secret not provided, will not create bot instance")
			return nil
		}

		slackBot := &SlackBot{
			config:      param.Config,
			logger:      param.Logger,
			closeChan:   make(chan struct{}, 1),
			processChan: make(chan smrRequestInfo, 10),
			smrModel:    param.SMR,
			ent:         param.Ent,
		}

		engine := gin.Default()
		engine.POST("/slack/command/smr", slackBot.postCommandInfo)
		engine.GET("/slack/install/auth", slackBot.getInstallAuth)
		slackBot.server = &http.Server{
			Addr:              lo.Ternary(slackConfig.Port == "", ":7070", slackConfig.Port),
			Handler:           engine,
			ReadHeaderTimeout: time.Second * 10,
		}

		param.Lifecycle.Append(fx.Hook{
			OnStop: func(ctx context.Context) error {
				err := slackBot.server.Shutdown(ctx)

				if !errors.Is(err, context.Canceled) {
					param.Logger.WithField("error", err.Error()).Error("slack bot server shutdown failed")
					return err
				}
				slackBot.logger.Info("stopped to receiving new requests")

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
				bot.logger.WithField("error", err.Error()).Fatal("slack bot server error")
			}
		}()

		go bot.runSmr()

		return nil
	}
}
