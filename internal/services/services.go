package services

import (
	"github.com/nekomeowww/insights-bot/internal/services/autorecap"
	"github.com/nekomeowww/insights-bot/internal/services/health"
	"github.com/nekomeowww/insights-bot/internal/services/pprof"
	"github.com/nekomeowww/insights-bot/internal/services/smr"
	"go.uber.org/fx"
)

func NewModules() fx.Option {
	return fx.Options(
		fx.Provide(health.NewHealth()),
		fx.Provide(pprof.NewPprof()),
		fx.Provide(autorecap.NewAutoRecapService()),
		fx.Provide(smr.NewService()),
	)
}
