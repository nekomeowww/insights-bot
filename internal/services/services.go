package services

import (
	"github.com/nekomeowww/insights-bot/internal/services/autorecap"
	"github.com/nekomeowww/insights-bot/internal/services/pprof"
	"go.uber.org/fx"
)

func NewModules() fx.Option {
	return fx.Options(
		fx.Provide(autorecap.NewChatHistoryRecapService()),
		fx.Provide(pprof.NewPprof()),
	)
}
