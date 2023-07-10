package handlers

import (
	"go.uber.org/fx"

	"github.com/nekomeowww/insights-bot/internal/bots/telegram/handlers/recap"
	"github.com/nekomeowww/insights-bot/internal/bots/telegram/handlers/summarize"
	"github.com/nekomeowww/insights-bot/internal/bots/telegram/handlers/welcome"
	"github.com/nekomeowww/insights-bot/pkg/bots/tgbot"
)

func NewModules() fx.Option {
	return fx.Options(
		fx.Provide(NewHandlers()),
		fx.Options(recap.NewModules()),
		fx.Options(summarize.NewModules()),
		fx.Options(welcome.NewModules()),
	)
}

type NewHandlersParam struct {
	fx.In

	Dispatcher *tgbot.Dispatcher

	RecapHandlers     *recap.Handlers
	SummarizeHandlers *summarize.Handlers
	WelcomeHandlers   *welcome.Handlers
}

type Handlers struct {
	Dispatcher    *tgbot.Dispatcher
	HandlerGroups []tgbot.HandlerGroup
}

func NewHandlers() func(param NewHandlersParam) *Handlers {
	return func(param NewHandlersParam) *Handlers {
		return &Handlers{
			Dispatcher: param.Dispatcher,
			HandlerGroups: []tgbot.HandlerGroup{
				param.SummarizeHandlers,
				param.RecapHandlers,
				param.WelcomeHandlers,
			},
		}
	}
}

func (h *Handlers) InstallAll() {
	for _, g := range h.HandlerGroups {
		g.Install(h.Dispatcher)
	}
}
