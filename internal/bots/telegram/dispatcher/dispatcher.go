package dispatcher

import (
	"runtime/debug"

	"github.com/nekomeowww/insights-bot/pkg/handler"
	"github.com/nekomeowww/insights-bot/pkg/logger"
	"go.uber.org/fx"
)

func NewModules() fx.Option {
	return fx.Options(
		fx.Provide(NewDispatcher()),
	)
}

type NewDispatcherParam struct {
	fx.In

	Logger *logger.Logger
}

type Dispatcher struct {
	Logger              *logger.Logger
	MessageHandlers     []handler.HandleFunc
	ChannelPostHandlers []handler.HandleFunc
}

func NewDispatcher() func(param NewDispatcherParam) *Dispatcher {
	return func(param NewDispatcherParam) *Dispatcher {
		return &Dispatcher{
			Logger:              param.Logger,
			MessageHandlers:     make([]handler.HandleFunc, 0),
			ChannelPostHandlers: make([]handler.HandleFunc, 0),
		}
	}
}

func (d *Dispatcher) RegisterOneMessageHandler(handler handler.HandleFunc) {
	d.MessageHandlers = append(d.MessageHandlers, handler)
}

func (d *Dispatcher) DispatchMessage(c *handler.Context) {
	for _, h := range d.MessageHandlers {
		go func(handlerFunc handler.HandleFunc) {
			defer func() {
				if err := recover(); err != nil {
					d.Logger.Errorf("Panic recovered from message handler, %v\n%s", err, debug.Stack())
					return
				}
			}()

			handlerFunc(c)
		}(h)
	}
}

func (d *Dispatcher) RegisterOneChannelPostHandler(handler handler.HandleFunc) {
	d.ChannelPostHandlers = append(d.ChannelPostHandlers, handler)
}

func (d *Dispatcher) DispatchChannelPost(c *handler.Context) {
	for _, h := range d.ChannelPostHandlers {
		go func(handlerFunc handler.HandleFunc) {
			defer func() {
				if err := recover(); err != nil {
					d.Logger.Errorf("Panic recovered from channel post handler, %v\n%s", err, debug.Stack())
					return
				}
			}()

			handlerFunc(c)
		}(h)
	}
}
