package smr

import (
	"context"
	"errors"
	"github.com/nekomeowww/insights-bot/internal/configs"
	"github.com/nekomeowww/insights-bot/internal/datastore"
	"github.com/nekomeowww/insights-bot/internal/models/smr"
	"github.com/nekomeowww/insights-bot/internal/services/smr/smrqueue"
	"github.com/nekomeowww/insights-bot/pkg/bots/discordbot"
	"github.com/nekomeowww/insights-bot/pkg/bots/slackbot"
	"github.com/nekomeowww/insights-bot/pkg/bots/tgbot"
	"github.com/nekomeowww/insights-bot/pkg/logger"
	"github.com/redis/rueidis"
	"github.com/samber/lo"
	"github.com/sourcegraph/conc/pool"
	"go.uber.org/fx"
)

func NewModules() fx.Option {
	return fx.Options(
		fx.Options(smrqueue.NewModules()),
		fx.Provide(NewService()),
	)
}

type NewServiceParam struct {
	fx.In

	LifeCycle fx.Lifecycle

	Config *configs.Config
	Logger *logger.Logger

	RedisClient *datastore.Redis
	Ent         *datastore.Ent
	Queue       *smrqueue.Queue

	Model *smr.Model

	TgBot      *tgbot.BotService
	SlackBot   *slackbot.BotService
	DiscordBot *discordbot.BotService
}

type Service struct {
	logger *logger.Logger
	config *configs.Config

	ent *datastore.Ent

	model *smr.Model

	tgBot      *tgbot.BotService
	slackBot   *slackbot.BotService
	discordBot *discordbot.BotService

	started       bool
	closeChan     chan struct{}
	alreadyClosed bool

	queue *smrqueue.Queue
}

func (s *Service) Check(ctx context.Context) error {
	return lo.Ternary(s.started, nil, errors.New("smr service not started yet"))
}

func NewService() func(param NewServiceParam) (*Service, error) {
	return func(param NewServiceParam) (*Service, error) {
		s := &Service{
			logger:     param.Logger,
			ent:        param.Ent,
			config:     param.Config,
			model:      param.Model,
			queue:      param.Queue,
			tgBot:      param.TgBot,
			slackBot:   param.SlackBot,
			discordBot: param.DiscordBot,
		}

		param.LifeCycle.Append(fx.Hook{OnStop: func(ctx context.Context) error {
			s.stop()

			return nil
		}})

		return s, nil
	}
}

func (s *Service) run() {
	s.closeChan = make(chan struct{})
	s.started = true

	s.logger.Info("smr service started")

	needToClose := false

	taskRunner := pool.New().WithMaxGoroutines(10)

	for {
		select {
		case <-s.closeChan:
			s.logger.WithField("last tasks count", s.queue.Count()).Info("smr service: received stop signal, waiting for all tasks done")

			needToClose = true
		default:
		}

		info, err := s.queue.GetTask()
		if err != nil {
			if errors.Is(err, rueidis.Nil) || errors.Is(err, smrqueue.ErrQueueFull) {
				continue
			}

			s.logger.WithError(err).Warn("smr service: failed to get task")

			continue
		}

		taskRunner.Go(func() {
			defer func() {
				err2 := recover()
				if err2 != nil {
					s.logger.
						WithField("err", err2).
						WithField("task", info).
						Error("smr service: task failed with panic")
				}
			}()

			s.processor(info)
			s.queue.FinishTask()
		})

		if needToClose {
			break
		}
	}

	s.alreadyClosed = true

	taskRunner.Wait()
}

func (s *Service) stop() {
	if s.alreadyClosed {
		return
	}
	s.closeChan <- struct{}{}
	close(s.closeChan)
}

func Run() func(s *Service) error {
	return func(s *Service) error {
		s.run()
		return nil
	}
}
