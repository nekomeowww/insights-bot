package smr

import (
	"context"
	"errors"
	"github.com/nekomeowww/insights-bot/internal/configs"
	"github.com/nekomeowww/insights-bot/internal/datastore"
	"github.com/nekomeowww/insights-bot/internal/models/smr"
	"github.com/nekomeowww/insights-bot/internal/services/smr/types"
	"github.com/nekomeowww/insights-bot/pkg/bots/discordbot"
	"github.com/nekomeowww/insights-bot/pkg/bots/slackbot"
	"github.com/nekomeowww/insights-bot/pkg/bots/tgbot"
	"github.com/nekomeowww/insights-bot/pkg/logger"
	"github.com/samber/lo"
	"go.uber.org/fx"
)

type NewServiceParam struct {
	fx.In

	LifeCycle fx.Lifecycle

	Config *configs.Config
	Logger *logger.Logger

	RedisClient *datastore.Redis
	Ent         *datastore.Ent

	Model *smr.Model
}

type Service struct {
	logger *logger.Logger
	config *configs.Config

	redisClient *datastore.Redis
	ent         *datastore.Ent

	model *smr.Model

	tgBot      *tgbot.BotService
	slackBot   *slackbot.BotService
	discordBot *discordbot.BotService

	started       bool
	closeChan     chan struct{}
	alreadyClosed bool

	queue *types.TaskQueue
}

func (s *Service) SetTelegramBot(bot *tgbot.BotService) {
	s.tgBot = bot
}

func (s *Service) SetSlackBot(bot *slackbot.BotService) {
	s.slackBot = bot
}

func (s *Service) SetDiscordBot(bot *discordbot.BotService) {
	s.discordBot = bot
}

func (s *Service) Check(ctx context.Context) error {
	return lo.Ternary(s.started, nil, errors.New("smr service not started yet"))
}

func NewService() func(param NewServiceParam) (*Service, error) {
	return func(param NewServiceParam) (*Service, error) {
		s := &Service{
			logger:      param.Logger,
			redisClient: param.RedisClient,
			ent:         param.Ent,
			config:      param.Config,
			model:       param.Model,
		}

		param.LifeCycle.Append(fx.Hook{OnStop: func(ctx context.Context) error {
			s.stop()

			return nil
		}})

		return s, nil
	}
}

func Run() func(s *Service) error {
	return func(s *Service) error {
		s.run()
		return nil
	}
}
