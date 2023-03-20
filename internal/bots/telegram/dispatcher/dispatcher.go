package dispatcher

import (
	"github.com/nekomeowww/insights-bot/pkg/handler"
	"go.uber.org/fx"
)

func NewModules() fx.Option {
	return fx.Options(
		fx.Provide(NewDispatcher()),
	)
}

type NewDispatcherParam struct {
	fx.In
}

type Dispatcher struct {
	MessageHandlers     []handler.HandleFunc
	ChannelPostHandlers []handler.HandleFunc
}

func NewDispatcher() func(param NewDispatcherParam) *Dispatcher {
	return func(param NewDispatcherParam) *Dispatcher {
		return &Dispatcher{
			MessageHandlers:     make([]handler.HandleFunc, 0),
			ChannelPostHandlers: make([]handler.HandleFunc, 0),
		}
	}
}

func (d *Dispatcher) RegisterOneMessageHandler(handler handler.HandleFunc) {
	d.MessageHandlers = append(d.MessageHandlers, handler)
}

func (d *Dispatcher) DispatchMessage(c *handler.Context) {
	for _, handler := range d.MessageHandlers {
		go handler(c)
	}
}

func (d *Dispatcher) RegisterOneChannelPostHandler(handler handler.HandleFunc) {
	d.ChannelPostHandlers = append(d.ChannelPostHandlers, handler)
}

func (d *Dispatcher) DispatchChannelPost(c *handler.Context) {
	for _, handler := range d.ChannelPostHandlers {
		go handler(c)
	}
}
