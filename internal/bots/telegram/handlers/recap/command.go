package recap

import (
	"github.com/nekomeowww/insights-bot/internal/configs"
	"github.com/nekomeowww/insights-bot/internal/datastore"
	"github.com/nekomeowww/insights-bot/internal/models/chathistories"
	"github.com/nekomeowww/insights-bot/internal/models/tgchats"
	"github.com/nekomeowww/insights-bot/pkg/logger"
	"go.uber.org/fx"
)

type NewCommandHandlerParams struct {
	fx.In

	Config        *configs.Config
	Logger        *logger.Logger
	TgChats       *tgchats.Model
	ChatHistories *chathistories.Model
	Redis         *datastore.Redis
}

type CommandHandler struct {
	config        *configs.Config
	logger        *logger.Logger
	tgchats       *tgchats.Model
	chathistories *chathistories.Model
	redis         *datastore.Redis
}

func NewRecapCommandHandler() func(NewCommandHandlerParams) *CommandHandler {
	return func(param NewCommandHandlerParams) *CommandHandler {
		return &CommandHandler{
			config:        param.Config,
			logger:        param.Logger,
			tgchats:       param.TgChats,
			chathistories: param.ChatHistories,
			redis:         param.Redis,
		}
	}
}
