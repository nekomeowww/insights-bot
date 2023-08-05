package thirdparty

import (
	"go.uber.org/fx"

	"github.com/nekomeowww/insights-bot/internal/thirdparty/openai"
)

func NewModules() fx.Option {
	return fx.Options(
		fx.Provide(openai.NewClient(true)),
	)
}
