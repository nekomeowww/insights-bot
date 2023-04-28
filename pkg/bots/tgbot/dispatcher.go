package tgbot

import (
	"net/url"
	"runtime/debug"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/gookit/color"
	"github.com/samber/lo"
	"go.uber.org/fx"

	"github.com/nekomeowww/insights-bot/pkg/logger"
	"github.com/nekomeowww/insights-bot/pkg/types/telegram"
	"github.com/nekomeowww/insights-bot/pkg/utils"
)

type NewDispatcherParam struct {
	fx.In

	Logger *logger.Logger
}

type Dispatcher struct {
	Logger *logger.Logger

	helpCommand           *helpCommandHandler
	middlewares           []MiddlewareFunc
	commandHandlers       map[string]HandleFunc
	messageHandlers       map[string]HandleFunc
	channelPostHandlers   []Handler
	callbackQueryHandlers map[string]HandleFunc
}

func NewDispatcher() func(param NewDispatcherParam) *Dispatcher {
	return func(param NewDispatcherParam) *Dispatcher {
		d := &Dispatcher{
			Logger:                param.Logger,
			helpCommand:           newHelpCommandHandler(),
			middlewares:           make([]MiddlewareFunc, 0),
			commandHandlers:       make(map[string]HandleFunc),
			messageHandlers:       make(map[string]HandleFunc),
			channelPostHandlers:   make([]Handler, 0),
			callbackQueryHandlers: make(map[string]HandleFunc),
		}

		d.OnCommand(d.helpCommand)
		return d
	}
}

func (d *Dispatcher) Use(middleware MiddlewareFunc) {
	d.middlewares = append(d.middlewares, middleware)
}

func (d *Dispatcher) OnCommand(h CommandHandler) {
	d.helpCommand.commands = append(d.helpCommand.commands, h)
	d.commandHandlers[h.Command()] = NewHandler(h.Handle).Handle
}

func (d *Dispatcher) OnMessage(h MessageHandler) {
	d.messageHandlers[h.Message()] = NewHandler(h.Handle).Handle
}

func (d *Dispatcher) dispatchMessage(c *Context) {
	identityStrings := make([]string, 0)
	identityStrings = append(identityStrings, FullNameFromFirstAndLastName(c.Update.Message.From.FirstName, c.Update.Message.From.LastName))
	if c.Update.Message.From.UserName != "" {
		identityStrings = append(identityStrings, "@"+c.Update.Message.From.UserName)
	}
	if c.Update.Message.Chat.Type == "private" {
		d.Logger.Infof("[消息｜%s] %s (%s): %s",
			MapChatTypeToChineseText(telegram.ChatType(c.Update.Message.Chat.Type)),
			strings.Join(identityStrings, " "),
			color.FgYellow.Render(c.Update.Message.From.ID),
			lo.Ternary(c.Update.Message.Text == "", "<empty or contains medias>", c.Update.Message.Text),
		)
	} else {
		d.Logger.Infof("[消息｜%s] [%s (%s)] %s (%s): %s",
			MapChatTypeToChineseText(telegram.ChatType(c.Update.Message.Chat.Type)),
			color.FgGreen.Render(c.Update.Message.Chat.Title),
			color.FgYellow.Render(c.Update.Message.Chat.ID),
			strings.Join(identityStrings, " "),
			color.FgYellow.Render(c.Update.Message.From.ID),
			lo.Ternary(c.Update.Message.Text == "", "<empty or contains medias>", c.Update.Message.Text),
		)
	}
	if c.Update.Message.Command() != "" {
		d.dispatchInGoroutine(func() {
			for cmd, f := range d.commandHandlers {
				if c.Update.Message.Command() == cmd {
					_, _ = f(c)
				}
			}
		})

	} else {
		d.dispatchInGoroutine(func() {
			for msg, f := range d.messageHandlers {
				if c.Update.Message.Text == msg {
					_, _ = f(c)
				}
				if c.Update.Message.Caption == msg {
					_, _ = f(c)
				}
			}
		})
	}
}

func (d *Dispatcher) OnChannelPost(handler Handler) {
	d.channelPostHandlers = append(d.channelPostHandlers, handler)
}

func (d *Dispatcher) dispatchChannelPost(c *Context) {
	d.Logger.Infof("[频道消息｜%s] [%s (%s)]: %s",
		MapChatTypeToChineseText(telegram.ChatType(c.Update.ChannelPost.Chat.Type)),
		color.FgGreen.Render(c.Update.ChannelPost.Chat.Title),
		color.FgYellow.Render(c.Update.ChannelPost.Chat.ID),
		lo.Ternary(c.Update.ChannelPost.Text == "", "<empty or contains medias>", c.Update.ChannelPost.Text),
	)

	d.dispatchInGoroutine(func() {
		for _, h := range d.channelPostHandlers {
			_, _ = h.Handle(c)
		}
	})
}

func (d *Dispatcher) OnCallbackQuery(h CallbackQueryHandler) {
	d.callbackQueryHandlers[h.CallbackQueryRoute()] = NewHandler(h.Handle).Handle
}

func (d *Dispatcher) dispatchCallbackQuery(c *Context) {
	d.Logger.Infof("[回调查询｜%s] [%s (%s)]: %s",
		MapChatTypeToChineseText(telegram.ChatType(c.Update.CallbackQuery.Message.Chat.Type)),
		color.FgGreen.Render(c.Update.CallbackQuery.Message.Chat.Title),
		color.FgYellow.Render(c.Update.CallbackQuery.Message.Chat.ID),
		c.Update.CallbackQuery.Data,
	)

	d.dispatchInGoroutine(func() {
		for route, h := range d.callbackQueryHandlers {
			parsedRoute, err := url.Parse(c.Update.CallbackQuery.Data)
			if err != nil {
				d.Logger.Errorf("failed to parse callback query data, err: %v, data: %v", err, utils.SprintJSON(parsedRoute))
				continue
			}
			if parsedRoute.Host+parsedRoute.Path == route {
				_, _ = h(c)
			}
		}
	})
}

func (d *Dispatcher) dispatchMyChatMember(c *Context) {
	identityStrings := make([]string, 0)
	identityStrings = append(identityStrings, FullNameFromFirstAndLastName(c.Update.MyChatMember.From.FirstName, c.Update.MyChatMember.From.LastName))
	if c.Update.MyChatMember.From.UserName != "" {
		identityStrings = append(identityStrings, "@"+c.Update.MyChatMember.From.UserName)
	}

	oldMemberStatus := telegram.MemberStatus(c.Update.MyChatMember.OldChatMember.Status)
	newMemberStatus := telegram.MemberStatus(c.Update.MyChatMember.NewChatMember.Status)

	d.Logger.Infof("[我的成员信息更新｜%s] [%s (%s)] %s (%s): 成员状态自 %s 变更为 %s",
		MapChatTypeToChineseText(telegram.ChatType(c.Update.MyChatMember.Chat.Type)),
		color.FgGreen.Render(c.Update.MyChatMember.Chat.Title),
		color.FgYellow.Render(c.Update.MyChatMember.Chat.ID),
		strings.Join(identityStrings, " "),
		color.FgYellow.Render(c.Update.MyChatMember.From.ID),
		MapMemberStatusToChineseText(oldMemberStatus),
		MapMemberStatusToChineseText(newMemberStatus),
	)
	switch c.Update.MyChatMember.Chat.Type {
	case "channel":
		if newMemberStatus != "administrator" {
			d.Logger.Infof("已退出频道 %s (%d)", c.Update.MyChatMember.Chat.Title, c.Update.MyChatMember.Chat.ID)
			return
		}

		_, err := c.Bot.GetChat(tgbotapi.ChatInfoConfig{
			ChatConfig: tgbotapi.ChatConfig{
				ChatID: c.Update.MyChatMember.Chat.ID,
			},
		})
		if err != nil {
			d.Logger.Error(err)
			return
		}

		d.Logger.Infof("已加入频道 %s (%d)", c.Update.MyChatMember.Chat.Title, c.Update.MyChatMember.Chat.ID)
	}
}

func (d *Dispatcher) Dispatch(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	for _, m := range d.middlewares {
		m(NewContext(bot, update, d.Logger), func() {})
	}

	ctx := NewContext(bot, update, d.Logger)
	switch ctx.UpdateType() {
	case UpdateTypeMessage:
		d.dispatchMessage(ctx)
	case UpdateTypeEditedMessage:
		d.Logger.Warn("edited message is not supported yet")
	case UpdateTypeChannelPost:
		d.dispatchChannelPost(ctx)
	case UpdateTypeEditedChannelPost:
		d.Logger.Warn("edited channel post is not supported yet")
	case UpdateTypeInlineQuery:
		d.Logger.Warn("inline query is not supported yet")
	case UpdateTypeChosenInlineResult:
		d.Logger.Warn("chosen inline result is not supported yet")
	case UpdateTypeCallbackQuery:
		d.dispatchCallbackQuery(ctx)
	case UpdateTypeShippingQuery:
		d.Logger.Warn("shipping query is not supported yet")
	case UpdateTypePreCheckoutQuery:
		d.Logger.Warn("pre checkout query is not supported yet")
	case UpdateTypePoll:
		d.Logger.Warn("poll is not supported yet")
	case UpdateTypePollAnswer:
		d.Logger.Warn("poll answer is not supported yet")
	case UpdateTypeMyChatMember:
		d.dispatchMyChatMember(ctx)
	case UpdateTypeChatMember:
		d.Logger.Warn("chat member is not supported yet")
	case UpdateTypeChatJoinRequest:
		d.Logger.Warn("chat join request is not supported yet")
	case UpdateTypeUnknown:
		d.Logger.Warn("unable to dispatch update due to unknown update type")
	}
}

func (d *Dispatcher) dispatchInGoroutine(f func()) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				d.Logger.Errorf("Panic recovered from command dispatcher, %v\n%s", err, debug.Stack())
				return
			}
		}()

		f()
	}()
}
