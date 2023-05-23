package tgchats

import (
	"github.com/nekomeowww/insights-bot/internal/configs"
	"github.com/nekomeowww/insights-bot/internal/datastore"
	"github.com/nekomeowww/insights-bot/pkg/logger"
	"go.uber.org/fx"
)

type NewModelParams struct {
	fx.In

	Config *configs.Config
	Ent    *datastore.Ent
	Digger *datastore.AutoRecapTimeCapsuleDigger
	Logger *logger.Logger
}

type Model struct {
	config *configs.Config
	ent    *datastore.Ent
	logger *logger.Logger
	digger *datastore.AutoRecapTimeCapsuleDigger
}

func NewModel() func(NewModelParams) (*Model, error) {
	return func(param NewModelParams) (*Model, error) {
		return &Model{
			config: param.Config,
			ent:    param.Ent,
			logger: param.Logger,
			digger: param.Digger,
		}, nil
	}
}
