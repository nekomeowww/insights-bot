package summarize

import (
	"github.com/nekomeowww/insights-bot/internal/models/smr"
	"github.com/nekomeowww/insights-bot/internal/services/smr/smrqueue"
	"github.com/nekomeowww/insights-bot/pkg/bots/tgbot"
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
	SMR      *smr.Model
	SmrQueue *smrqueue.Queue
}

var _ tgbot.HandlerGroup = (*Handlers)(nil)

type Handlers struct {
	logger   *logger.Logger
	smr      *smr.Model
	smrQueue *smrqueue.Queue
}

func NewHandlers() func(NewHandlersParams) *Handlers {
	return func(param NewHandlersParams) *Handlers {
		handler := &Handlers{
			logger:   param.Logger,
			smrQueue: param.SmrQueue,
			smr:      param.SMR,
		}

		return handler
	}
}

func (h *Handlers) Install(dispatcher *tgbot.Dispatcher) {
	dispatcher.OnCommandGroup("量子速读", []tgbot.Command{
		{Command: "smr", HelpMessage: "量子速读网页文章（也支持在频道中使用） 用法：/smr <code>&lt;链接&gt;</code>", Handler: tgbot.NewHandler(h.Handle)},
	})

	dispatcher.OnChannelPost(tgbot.NewHandler(h.HandleChannelPost))
	dispatcher.OnCallbackQuery("smr/summarization/feedback/react", tgbot.NewHandler(h.handleCallbackQueryReact))
}
