package datastore

import "go.uber.org/fx"

func NewModules() fx.Option {
	return fx.Options(
		fx.Provide(NewPinecone()),
		fx.Provide(NewEnt()),
	)
}
