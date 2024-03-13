package smr

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/nekomeowww/insights-bot/internal/configs"
	"github.com/nekomeowww/insights-bot/internal/datastore"
	"github.com/nekomeowww/insights-bot/internal/models/smr"
	"github.com/nekomeowww/insights-bot/internal/services/smr/smrqueue"
	"github.com/nekomeowww/insights-bot/pkg/bots/discordbot"
	"github.com/nekomeowww/insights-bot/pkg/bots/slackbot"
	"github.com/nekomeowww/insights-bot/pkg/bots/tgbot"
	"github.com/nekomeowww/insights-bot/pkg/i18n"
	"github.com/nekomeowww/insights-bot/pkg/logger"
	"github.com/redis/rueidis"
	"github.com/samber/lo"
	"go.uber.org/fx"
	"go.uber.org/zap"
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
	I18n   *i18n.I18n

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
	i18n   *i18n.I18n

	ent *datastore.Ent

	model *smr.Model

	tgBot      *tgbot.BotService
	slackBot   *slackbot.BotService
	discordBot *discordbot.BotService

	started       bool
	closeFunc     context.CancelFunc
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
			i18n:       param.I18n,
			config:     param.Config,
			ent:        param.Ent,
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
	var ctx context.Context
	ctx, s.closeFunc = context.WithCancel(context.Background())
	s.started = true

	s.logger.Info("smr service started")

	needToClose := false

	for {
		select {
		case <-ctx.Done():
			s.logger.Info("smr service: received stop signal, waiting for all tasks done", zap.Int("last_tasks_count", s.queue.Count()))

			needToClose = true
		default:
		}

		info, err := s.queue.GetTask()
		if err != nil {
			if errors.Is(err, rueidis.Nil) || errors.Is(err, smrqueue.ErrQueueFull) {
				continue
			}

			s.logger.Warn("smr service: failed to get task", zap.Error(err))

			continue
		}

		go func() {
			defer func() {
				err2 := recover()
				if err2 != nil {
					s.logger.Error("smr service: task failed with panic",
						zap.Any("panic", err2),
						zap.String("task", string(lo.Must(json.Marshal(info)))),
					)
				}
			}()

			s.processor(info)
			s.queue.FinishTask()
		}()

		if needToClose && s.queue.Count() == 0 {
			break
		}

		time.Sleep(time.Second * 2)
	}

	s.alreadyClosed = true
}

func (s *Service) stop() {
	if s.alreadyClosed {
		return
	}

	s.closeFunc()
}

func Run() func(s *Service) {
	return func(s *Service) {
		go s.run()
	}
}
