package summarize

import (
	"github.com/nekomeowww/insights-bot/internal/models/smr"
	"github.com/nekomeowww/insights-bot/internal/services/smr/smrqueue"
	"github.com/nekomeowww/insights-bot/pkg/bots/tgbot"
	"github.com/nekomeowww/insights-bot/pkg/i18n"
	"github.com/nekomeowww/insights-bot/pkg/logger"
	"go.uber.org/fx"
)

func NewModules() fx.Option {
	return fx.Options(
		fx.Provide(NewHandlers()),
	)
}

type NewHandlersParams struct {
	fx.In

	Logger   *logger.Logger
	I18n     *i18n.I18n
	SMR      *smr.Model
	SmrQueue *smrqueue.Queue
}

var _ tgbot.HandlerGroup = (*Handlers)(nil)

type Handlers struct {
	logger   *logger.Logger
	i18n     *i18n.I18n
	smr      *smr.Model
	smrQueue *smrqueue.Queue
}

func NewHandlers() func(NewHandlersParams) *Handlers {
	return func(param NewHandlersParams) *Handlers {
		handler := &Handlers{
			logger:   param.Logger,
			i18n:     param.I18n,
			smrQueue: param.SmrQueue,
			smr:      param.SMR,
		}

		return handler
	}
}

func (h *Handlers) Install(dispatcher *tgbot.Dispatcher) {
	dispatcher.OnCommandGroup(func(c *tgbot.Context) string {
		return c.T("commands.groups.summarization.name")
	}, []tgbot.Command{
		{
			Command: "smr",
			Handler: tgbot.NewHandler(h.Handle),
			HelpMessage: func(c *tgbot.Context) string {
				return c.T("commands.groups.summarization.commands.smr.help")
			},
		},
	})

	dispatcher.OnChannelPost(tgbot.NewHandler(h.HandleChannelPost))
	dispatcher.OnCallbackQuery("smr/summarization/feedback/react", tgbot.NewHandler(h.handleCallbackQueryReact))
	dispatcher.OnCallbackQuery("smr/summarization/retry", tgbot.NewHandler(h.handleCallbackQueryRetry))
}
