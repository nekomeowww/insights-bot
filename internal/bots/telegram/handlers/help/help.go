package help

import (
	"github.com/nekomeowww/insights-bot/pkg/logger"
	"go.uber.org/fx"
)

type NewHandlerParam struct {
	fx.In

	Logger *logger.Logger
}

type Handler struct {
	Logger *logger.Logger
}

func NewHandler() func(NewHandlerParam) *Handler {
	return func(param NewHandlerParam) *Handler {
		return &Handler{
			Logger:                   param.Logger,
		}
	}
}
