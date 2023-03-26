package datastore

import (
	"github.com/nekomeowww/insights-bot/internal/configs"
	"github.com/nekomeowww/insights-bot/pkg/datastore/pinecone"
	"go.uber.org/fx"
)

type NewPineconeParam struct {
	fx.In

	Config *configs.Config
}

type Pinecone struct {
	*pinecone.Client
}

func NewPinecone() func(NewPineconeParam) (*Pinecone, error) {
	return func(param NewPineconeParam) (*Pinecone, error) {
		client, err := pinecone.NewClient(
			param.Config.Pinecone.IndexName,
			param.Config.Pinecone.ProjectName,
			param.Config.Pinecone.Environment,
			param.Config.Pinecone.APIKey,
		)
		if err != nil {
			return nil, err
		}

		return &Pinecone{
			Client: client,
		}, nil
	}
}
