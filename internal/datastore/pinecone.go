package datastore

import (
	pinecone "github.com/nekomeowww/go-pinecone"
	"go.uber.org/fx"

	"github.com/nekomeowww/insights-bot/internal/configs"
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
		client, err := pinecone.New(
			pinecone.WithAPIKey(param.Config.Pinecone.APIKey),
			pinecone.WithEnvironment(param.Config.Pinecone.Environment),
			pinecone.WithProjectName(param.Config.Pinecone.ProjectName),
		)
		if err != nil {
			return nil, err
		}

		return &Pinecone{
			Client: client,
		}, nil
	}
}
