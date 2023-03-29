package handlers

import (
	"go.uber.org/fx"

	"github.com/nekomeowww/insights-bot/internal/bots/telegram/dispatcher"
	"github.com/nekomeowww/insights-bot/internal/bots/telegram/handlers/chat_with_chat_history"
	"github.com/nekomeowww/insights-bot/internal/bots/telegram/handlers/summarize"
	"github.com/nekomeowww/insights-bot/pkg/handler"
)

func NewModules() fx.Option {
	return fx.Options(
		fx.Provide(NewHandlers()),
		fx.Provide(summarize.NewHandler()),
		fx.Provide(chat_with_chat_history.NewHandler()),
	)
}

type NewHandlersParam struct {
	fx.In

	Dispatcher                 *dispatcher.Dispatcher
	SummarizeHandler           *summarize.Handler
	ChatWithChatHistoryHandler *chat_with_chat_history.Handler
}

type Handlers struct {
	Dispatcher *dispatcher.Dispatcher

	CommandHandlers     map[string]handler.HandleFunc
	MessageHandlers     []handler.HandleFunc
	ChannelPostHandlers []handler.HandleFunc
}

func NewHandlers() func(param NewHandlersParam) *Handlers {
	return func(param NewHandlersParam) *Handlers {
		return &Handlers{
			Dispatcher: param.Dispatcher,
			CommandHandlers: map[string]handler.HandleFunc{
				"smr":           param.SummarizeHandler.HandleSMRCommand,
				"recap":         param.ChatWithChatHistoryHandler.HandleRecapCommand,
				"enable_recap":  param.ChatWithChatHistoryHandler.HandleEnableRecapCommand,
				"disable_recap": param.ChatWithChatHistoryHandler.HandleDisableRecapCommand,
			},
			MessageHandlers: []handler.HandleFunc{
				param.ChatWithChatHistoryHandler.HandleRecordMessage,
			},
			ChannelPostHandlers: []handler.HandleFunc{
				param.SummarizeHandler.HandleChannelPost,
			},
		}
	}
}

func (h *Handlers) RegisterHandlers() {
	for cmd, ch := range h.CommandHandlers {
		h.Dispatcher.RegisterOneCommandHandler(cmd, ch)
	}
	for _, mh := range h.MessageHandlers {
		h.Dispatcher.RegisterOneMessageHandler(mh)
	}
	for _, cph := range h.ChannelPostHandlers {
		h.Dispatcher.RegisterOneChannelPostHandler(cph)
	}
}
