package datastore

import (
	"context"

	"entgo.io/ent/dialect"
	_ "github.com/lib/pq"
	"github.com/nekomeowww/insights-bot/ent"
	"github.com/nekomeowww/insights-bot/internal/configs"
	"go.uber.org/fx"
)

type NewEntParams struct {
	fx.In

	Lifecycle fx.Lifecycle
	Configs   *configs.Config
}

type Ent struct {
	*ent.Client
}

func NewEnt() func(NewEntParams) (*Ent, error) {
	return func(param NewEntParams) (*Ent, error) {
		opts := make([]ent.Option, 0)
		client, err := ent.Open(dialect.Postgres, param.Configs.DB.ConnectionString, opts...)
		if err != nil {
			return nil, err
		}

		err = client.Schema.Create(context.Background())
		if err != nil {
			return nil, err
		}

		param.Lifecycle.Append(fx.Hook{
			OnStop: func(ctx context.Context) error {
				return client.Close()
			},
		})

		return &Ent{client}, nil
	}
}
