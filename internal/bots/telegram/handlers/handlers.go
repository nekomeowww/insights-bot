package handlers

import (
	"go.uber.org/fx"

	"github.com/nekomeowww/insights-bot/internal/bots/telegram/dispatcher"
	"github.com/nekomeowww/insights-bot/internal/bots/telegram/handlers/summarize"
	"github.com/nekomeowww/insights-bot/pkg/handler"
)

func NewModules() fx.Option {
	return fx.Options(
		fx.Provide(NewHandlers()),
		fx.Provide(summarize.NewHandler()),
	)
}

type NewHandlersParam struct {
	fx.In

	Dispatcher       *dispatcher.Dispatcher
	SummarizeHandler *summarize.Handler
}

type Handlers struct {
	Dispatcher *dispatcher.Dispatcher

	MessageHandlers     []handler.HandleFunc
	ChannelPostHandlers []handler.HandleFunc
}

func NewHandlers() func(param NewHandlersParam) *Handlers {
	return func(param NewHandlersParam) *Handlers {
		return &Handlers{
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
		h.Dispatcher.RegisterOneMessageHandler(mh)
	}
	for _, cph := range h.ChannelPostHandlers {
		h.Dispatcher.RegisterOneChannelPostHandler(cph)
	}
}
