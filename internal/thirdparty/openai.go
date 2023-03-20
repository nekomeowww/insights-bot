package thirdparty

import (
	"go.uber.org/fx"

	"github.com/nekomeowww/insights-bot/internal/configs"
	"github.com/nekomeowww/insights-bot/pkg/openai"
)

type NewOpenAIClientParam struct {
	fx.In

	Config *configs.Config
}

func NewOpenAIClient() func(NewOpenAIClientParam) *openai.Client {
	return func(param NewOpenAIClientParam) *openai.Client {
		return openai.NewClient(param.Config.OpenAIAPISecret)
	}
}
