package dispatcher

import (
	"net/url"

	"github.com/nekomeowww/insights-bot/pkg/handler"
	"github.com/nekomeowww/insights-bot/pkg/logger"
	"github.com/nekomeowww/insights-bot/pkg/utils"
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
	Logger                *logger.Logger
	CommandHandlers       map[string]handler.HandleFunc
	MessageHandlers       []handler.HandleFunc
	ChannelPostHandlers   []handler.HandleFunc
	CallbackQueryHandlers map[string]handler.HandleFunc
}

func NewDispatcher() func(param NewDispatcherParam) *Dispatcher {
	return func(param NewDispatcherParam) *Dispatcher {
		return &Dispatcher{
			Logger:                param.Logger,
			CommandHandlers:       make(map[string]handler.HandleFunc),
			MessageHandlers:       make([]handler.HandleFunc, 0),
			ChannelPostHandlers:   make([]handler.HandleFunc, 0),
			CallbackQueryHandlers: make(map[string]handler.HandleFunc),
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

func (d *Dispatcher) RegisterOneCallbackQueryHandler(route string, handler handler.HandleFunc) {
	d.CallbackQueryHandlers[route] = handler
}

func (d *Dispatcher) DispatchCallbackQuery(c *handler.Context) {
	for route, h := range d.CallbackQueryHandlers {
		parsedRoute, err := url.Parse(c.Update.CallbackQuery.Data)
		if err != nil {
			d.Logger.Errorf("failed to parse callback query data, err: %v, data: %v", err, utils.SprintJSON(parsedRoute))
			continue
		}
		if parsedRoute.Host+parsedRoute.Path == route {
			h(c)
		}
	}
}
