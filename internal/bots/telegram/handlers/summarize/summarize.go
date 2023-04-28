package summarize

import (
	"github.com/nekomeowww/insights-bot/internal/models/smr"
	"github.com/nekomeowww/insights-bot/pkg/bots/tgbot"
	"go.uber.org/fx"
)

func NewModules() fx.Option {
	return fx.Options(
		fx.Provide(NewHandlers()),
	)
}

type NewHandlersParams struct {
	fx.In

	SMR *smr.Model
}

var _ tgbot.HandlerGroup = (*Handlers)(nil)

type Handlers struct {
	smr *smr.Model
}

func NewHandlers() func(NewHandlersParams) *Handlers {
	return func(param NewHandlersParams) *Handlers {
		handler := &Handlers{
			smr: param.SMR,
		}

		return handler
	}
}

func (h *Handlers) Install(dispatcher *tgbot.Dispatcher) {
	dispatcher.OnCommand(h)
	dispatcher.OnChannelPost(tgbot.NewHandler(h.HandleChannelPost))
}
