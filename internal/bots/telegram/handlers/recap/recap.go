package recap

import (
	"fmt"

	"github.com/nekomeowww/insights-bot/pkg/bots/tgbot"
	"github.com/samber/lo"
	"go.uber.org/fx"
)

func NewModules() fx.Option {
	return fx.Options(
		fx.Provide(NewHandlers()),
		fx.Provide(NewRecapCommandHandler()),
		fx.Provide(NewCallbackQueryHandler()),
	)
}

type NewHandlersParams struct {
	fx.In

	Command       *CommandHandler
	CallbackQuery *CallbackQueryHandler
}

type Handlers struct {
	command       *CommandHandler
	callbackQuery *CallbackQueryHandler
}

func NewHandlers() func(NewHandlersParams) *Handlers {
	return func(param NewHandlersParams) *Handlers {
		return &Handlers{
			command:       param.Command,
			callbackQuery: param.CallbackQuery,
		}
	}
}

func (h *Handlers) Install(dispatcher *tgbot.Dispatcher) {
	dispatcher.OnCommandGroup(func(c *tgbot.Context) string {
		return "聊天回顾"
	}, []tgbot.Command{
		{
			Command: "recap",
			Handler: tgbot.NewHandler(h.command.handleRecapCommand),
			HelpMessage: func(c *tgbot.Context) string {
				return "总结过去的聊天记录并生成回顾快报"
			},
		},
		{
			Command: "configure_recap",
			Handler: tgbot.NewHandler(h.command.handleConfigureRecapCommand),
			HelpMessage: func(c *tgbot.Context) string {
				return "配置聊天记录回顾（需要管理权限，<b>请在配置的时候尽量避免使用匿名用户身份或者其他群组的身份进行配置，可能会导致权限检查异常而配置失败。</b>）"
			},
		},
		{
			Command: "recap_forwarded_start",
			Handler: tgbot.NewHandler(h.command.handleRecapForwardedStartCommand),
			HelpMessage: func(c *tgbot.Context) string {
				return "使 Bot 接收在私聊中转发给 Bot 的消息，并在发送 /recap_forwarded 后开始总结"
			},
		},
		{
			Command: "recap_forwarded",
			Handler: tgbot.NewHandler(h.command.handleRecapForwardedCommand),
			HelpMessage: func(c *tgbot.Context) string {
				return "使 Bot 停止接收在私聊中转发给 Bot 的消息，对已经转发过的消息进行总结"
			},
		},
		{
			Command: "subscribe_recap",
			Handler: tgbot.NewHandler(h.command.handleSubscribeRecapCommand),
			HelpMessage: func(c *tgbot.Context) string {
				return "订阅当前群组的定时聊天回顾"
			},
		},
		{
			Command: "unsubscribe_recap",
			Handler: tgbot.NewHandler(h.command.handleUnsubscribeRecapCommand),
			HelpMessage: func(c *tgbot.Context) string {
				return "取消订阅当前群组的定时聊天回顾"
			},
		},
	})

	dispatcher.OnCancelCommand(h.command.handleRecapForwardedStartShouleCancel, tgbot.NewHandler(h.command.handleRecapForwardedStartCancelCommand))
	dispatcher.OnStartCommand(tgbot.NewHandler(h.command.handleStartCommandWithPrivateSubscriptionsRecap))
	dispatcher.OnStartCommand(tgbot.NewHandler(h.command.handleStartCommandWithRecapSubscription))

	dispatcher.OnCallbackQuery("recap/recap/select_hours", tgbot.NewHandler(h.callbackQuery.handleCallbackQuerySelectHours))
	dispatcher.OnCallbackQuery("recap/configure/toggle", tgbot.NewHandler(h.callbackQuery.handleCallbackQueryToggle))
	dispatcher.OnCallbackQuery("recap/configure/assign_mode", tgbot.NewHandler(h.callbackQuery.handleCallbackQueryAssignMode))
	dispatcher.OnCallbackQuery("recap/configure/complete", tgbot.NewHandler(h.callbackQuery.handleCallbackQueryComplete))
	dispatcher.OnCallbackQuery("recap/unsubscribe_recap", tgbot.NewHandler(h.callbackQuery.handleCallbackQueryUnsubscribe))
	dispatcher.OnCallbackQuery("recap/recap/feedback/react", tgbot.NewHandler(h.callbackQuery.handleCallbackQueryReact))
	dispatcher.OnCallbackQuery("recap/configure/auto_recap_rates_per_day", tgbot.NewHandler(h.callbackQuery.handleAutoRecapRatesPerDaySelect))
	dispatcher.OnCallbackQuery("recap/configure/pin", tgbot.NewHandler(h.callbackQuery.handleCallbackQueryPin))

	dispatcher.OnLeftChatMember(tgbot.NewHandler(h.command.handleChatMemberLeft))
}

var (
	RecapSelectHourAvailable = []int64{
		1, 2, 4, 6, 12,
	}
	RecapSelectHourAvailableText = lo.SliceToMap(RecapSelectHourAvailable, func(item int64) (int64, string) {
		return item, fmt.Sprintf("%d 小时", item)
	})
)
