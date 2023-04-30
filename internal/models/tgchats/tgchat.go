package tgchats

import (
	"github.com/nekomeowww/insights-bot/internal/datastore"
	"go.uber.org/fx"
)

type NewModelParams struct {
	fx.In

	Ent *datastore.Ent
}

type Model struct {
	ent *datastore.Ent
}

func NewModel() func(NewModelParams) (*Model, error) {
	return func(param NewModelParams) (*Model, error) {
		return &Model{
			ent: param.Ent,
		}, nil
	}
}
