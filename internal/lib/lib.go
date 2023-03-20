package lib

import "go.uber.org/fx"

func NewModules() fx.Option {
	return fx.Options(
		fx.Provide(NewLogger()),
	)
}
