package handlers

import (
	"runtime/debug"

	"go.uber.org/fx"

	"github.com/nekomeowww/insights-bot/internal/bots/telegram/dispatcher"
	"github.com/nekomeowww/insights-bot/internal/bots/telegram/handlers/summarize"
	"github.com/nekomeowww/insights-bot/pkg/handler"
	"github.com/nekomeowww/insights-bot/pkg/logger"
)

func NewModules() fx.Option {
	return fx.Options(
		fx.Provide(NewHandlers()),
		fx.Provide(summarize.NewHandler()),
	)
}

type NewHandlersParam struct {
	fx.In

	Logger *logger.Logger

	Dispatcher       *dispatcher.Dispatcher
	SummarizeHandler *summarize.Handler
}

type Handlers struct {
	Logger     *logger.Logger
	Dispatcher *dispatcher.Dispatcher

	MessageHandlers     []handler.HandleFunc
	ChannelPostHandlers []handler.HandleFunc
}

func NewHandlers() func(param NewHandlersParam) *Handlers {
	return func(param NewHandlersParam) *Handlers {
		return &Handlers{
			Logger:     param.Logger,
			Dispatcher: param.Dispatcher,
			MessageHandlers: []handler.HandleFunc{
				param.SummarizeHandler.HandleMessage,
			},
			ChannelPostHandlers: []handler.HandleFunc{
				param.SummarizeHandler.HandleChannelPost,
			},
		}
	}
}

func (h *Handlers) RegisterHandlers() {
	for _, mh := range h.MessageHandlers {
		go func(handlerFunc handler.HandleFunc) {
			defer func() {
				if err := recover(); err != nil {
					h.Logger.Errorf("Panic recovered, %v\n%s", err, debug.Stack())
					return
				}
			}()

			h.Dispatcher.RegisterOneMessageHandler(handlerFunc)
		}(mh)
	}
	for _, cph := range h.ChannelPostHandlers {
		go func(handlerFunc handler.HandleFunc) {
			defer func() {
				if err := recover(); err != nil {
					h.Logger.Errorf("Panic recovered, %v\n%s", err, debug.Stack())
					return
				}
			}()

			h.Dispatcher.RegisterOneChannelPostHandler(handlerFunc)
		}(cph)
	}
}
