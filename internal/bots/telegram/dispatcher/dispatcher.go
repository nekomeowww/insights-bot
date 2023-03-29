package dispatcher

import (
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
	CommandHandlers     map[string]handler.HandleFunc
	MessageHandlers     []handler.HandleFunc
	ChannelPostHandlers []handler.HandleFunc
}

func NewDispatcher() func(param NewDispatcherParam) *Dispatcher {
	return func(param NewDispatcherParam) *Dispatcher {
		return &Dispatcher{
			Logger:              param.Logger,
			CommandHandlers:     make(map[string]handler.HandleFunc),
			MessageHandlers:     make([]handler.HandleFunc, 0),
			ChannelPostHandlers: make([]handler.HandleFunc, 0),
		}
	}
}

func (d *Dispatcher) RegisterOneCommandHandler(cmd string, handler handler.HandleFunc) {
	d.CommandHandlers[cmd] = handler
}

func (d *Dispatcher) DispatchCommand(c *handler.Context) {
	for cmd, h := range d.CommandHandlers {
		if c.Update.Message.Command() == cmd {
			h(c)
			continue
		}
	}
}

func (d *Dispatcher) RegisterOneMessageHandler(handler handler.HandleFunc) {
	d.MessageHandlers = append(d.MessageHandlers, handler)
}

func (d *Dispatcher) DispatchMessage(c *handler.Context) {
	for _, h := range d.MessageHandlers {
		h(c)
	}
}

func (d *Dispatcher) RegisterOneChannelPostHandler(handler handler.HandleFunc) {
	d.ChannelPostHandlers = append(d.ChannelPostHandlers, handler)
}

func (d *Dispatcher) DispatchChannelPost(c *handler.Context) {
	for _, h := range d.ChannelPostHandlers {
		h(c)
	}
}
