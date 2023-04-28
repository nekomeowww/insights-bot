package tgchats

import (
	"github.com/nekomeowww/insights-bot/internal/datastore"
	"github.com/nekomeowww/insights-bot/pkg/types/telegram/tgchat"
	"go.uber.org/fx"
)

type NewModelParams struct {
	fx.In

	Clover *datastore.Clover
}

type Model struct {
	Clover *datastore.Clover
}

func NewModel() func(NewModelParams) (*Model, error) {
	return func(param NewModelParams) (*Model, error) {
		hasCollection, err := param.Clover.HasCollection(tgchat.FeatureFlag{}.CollectionName())
		if err != nil {
			return nil, err
		}
		if !hasCollection {
			err = param.Clover.CreateCollection(tgchat.FeatureFlag{}.CollectionName())
			if err != nil {
				return nil, err
			}
		}

		return &Model{
			Clover: param.Clover,
		}, nil
	}
}
