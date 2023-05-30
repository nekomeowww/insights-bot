package tgbot

import (
	"crypto/sha256"
	"fmt"
	"runtime/debug"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/gookit/color"
	"github.com/redis/rueidis"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"

	"github.com/nekomeowww/insights-bot/pkg/logger"
	"github.com/nekomeowww/insights-bot/pkg/types/telegram"
)

type Dispatcher struct {
	Logger *logger.Logger

	helpCommand                *helpCommandHandler
	cancelCommand              *cancelCommandHandler
	startCommandHandler        *startCommandHandler
	middlewares                []MiddlewareFunc
	commandHandlers            map[string]HandleFunc
	channelPostHandlers        []Handler
	callbackQueryHandlers      map[string]HandleFunc
	callbackQueryHandlersRoute map[string]string
}

func NewDispatcher() func(logger *logger.Logger) *Dispatcher {
	return func(logger *logger.Logger) *Dispatcher {
		d := &Dispatcher{
			Logger:                     logger,
			helpCommand:                newHelpCommandHandler(),
			cancelCommand:              newCancelCommandHandler(),
			startCommandHandler:        newStartCommandHandler(),
			middlewares:                make([]MiddlewareFunc, 0),
			commandHandlers:            make(map[string]HandleFunc),
			channelPostHandlers:        make([]Handler, 0),
			callbackQueryHandlers:      make(map[string]HandleFunc),
			callbackQueryHandlersRoute: make(map[string]string),
		}

		d.startCommandHandler.helpCommandHandler = d.helpCommand

		d.OnCommandGroup("基础命令", []Command{
			{Command: d.helpCommand.Command(), HelpMessage: d.helpCommand.CommandHelp(), Handler: NewHandler(d.helpCommand.handle)},
			{Command: d.cancelCommand.Command(), HelpMessage: d.cancelCommand.CommandHelp(), Handler: NewHandler(d.cancelCommand.handle)},
			{Command: d.startCommandHandler.Command(), HelpMessage: d.startCommandHandler.CommandHelp(), Handler: NewHandler(d.startCommandHandler.handle)},
		})

		return d
	}
}

func (d *Dispatcher) Use(middleware MiddlewareFunc) {
	d.middlewares = append(d.middlewares, middleware)
}

func (d *Dispatcher) OnCommand(cmd, commandHelp string, h Handler) {
	d.helpCommand.defaultGroup.commands = append(d.helpCommand.defaultGroup.commands, Command{
		Command:     cmd,
		HelpMessage: commandHelp,
	})

	d.commandHandlers[cmd] = h.Handle
}

func (d *Dispatcher) OnCommandGroup(groupName string, group []Command) {
	d.helpCommand.commandGroups = append(d.helpCommand.commandGroups, commandGroup{name: groupName, commands: group})

	for _, c := range group {
		d.commandHandlers[c.Command] = c.Handler.Handle
	}
}

func (d *Dispatcher) OnCancelCommand(cancelHandler func(c *Context) (bool, error), handler Handler) {
	d.cancelCommand.cancellableCommands = append(d.cancelCommand.cancellableCommands, cancellableCommand{
		shouldCancelFunc: cancelHandler,
		handler:          handler,
	})
}

func (d *Dispatcher) OnStartCommand(h Handler) {
	d.startCommandHandler.startCommandHandlers = append(d.startCommandHandler.startCommandHandlers, h)
}

func (d *Dispatcher) dispatchMessage(c *Context) {
	identityStrings := make([]string, 0)
	identityStrings = append(identityStrings, FullNameFromFirstAndLastName(c.Update.Message.From.FirstName, c.Update.Message.From.LastName))

	if c.Update.Message.From.UserName != "" {
		identityStrings = append(identityStrings, "@"+c.Update.Message.From.UserName)
	}
	if c.Update.Message.Chat.Type == "private" {
		d.Logger.Tracef("[消息｜%s] %s (%s): %s",
			MapChatTypeToChineseText(telegram.ChatType(c.Update.Message.Chat.Type)),
			strings.Join(identityStrings, " "),
			color.FgYellow.Render(c.Update.Message.From.ID),
			lo.Ternary(c.Update.Message.Text == "", "<empty or contains medias>", c.Update.Message.Text),
		)
	} else {
		d.Logger.Tracef("[消息｜%s] [%s (%s)] %s (%s): %s",
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
	}
}

func (d *Dispatcher) OnChannelPost(handler Handler) {
	d.channelPostHandlers = append(d.channelPostHandlers, handler)
}

func (d *Dispatcher) dispatchChannelPost(c *Context) {
	d.Logger.Tracef("[频道消息｜%s] [%s (%s)]: %s",
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

func (d *Dispatcher) OnCallbackQuery(route string, h Handler) {
	routeHash := fmt.Sprintf("%x", sha256.Sum256([]byte(route)))[0:16]
	d.callbackQueryHandlersRoute[routeHash] = route
	d.callbackQueryHandlers[routeHash] = h.Handle
}

func (d *Dispatcher) dispatchCallbackQuery(c *Context) {
	defer func() {
		identityStrings := make([]string, 0)
		identityStrings = append(identityStrings, FullNameFromFirstAndLastName(c.Update.CallbackQuery.From.FirstName, c.Update.CallbackQuery.From.LastName))

		if c.Update.CallbackQuery.From.UserName != "" {
			identityStrings = append(identityStrings, "@"+c.Update.CallbackQuery.From.UserName)
		}

		if c.callbackQueryHandlerRoute == "" || c.callbackQueryHandlerActionData == "" {
			d.Logger.WithFields(logrus.Fields{
				"route":            c.callbackQueryHandlerRoute,
				"route_hash":       c.callbackQueryHandlerRouteHash,
				"action_data_hash": c.callbackQueryHandlerActionDataHash,
			}).Tracef("[回调查询｜%s] [%s (%s)] %s (%s) : %s (Raw Data, missing route or action data)",
				MapChatTypeToChineseText(telegram.ChatType(c.Update.CallbackQuery.Message.Chat.Type)),
				color.FgGreen.Render(c.Update.CallbackQuery.Message.Chat.Title),
				color.FgYellow.Render(c.Update.CallbackQuery.Message.Chat.ID),
				strings.Join(identityStrings, " "),
				color.FgYellow.Render(c.Update.CallbackQuery.From.ID),
				c.Update.CallbackData(),
			)
		} else {
			d.Logger.WithFields(logrus.Fields{
				"route":            c.callbackQueryHandlerRoute,
				"route_hash":       c.callbackQueryHandlerRouteHash,
				"action_data_hash": c.callbackQueryHandlerActionDataHash,
			}).Tracef("[回调查询｜%s] [%s (%s)] %s (%s): %s: %s",
				MapChatTypeToChineseText(telegram.ChatType(c.Update.CallbackQuery.Message.Chat.Type)),
				color.FgGreen.Render(c.Update.CallbackQuery.Message.Chat.Title),
				color.FgYellow.Render(c.Update.CallbackQuery.Message.Chat.ID),
				strings.Join(identityStrings, " "),
				color.FgYellow.Render(c.Update.CallbackQuery.From.ID),
				c.callbackQueryHandlerRoute, c.callbackQueryHandlerActionData,
			)
		}
	}()

	callbackQueryActionInvalidErrMessage := tgbotapi.NewEditMessageText(c.Update.CallbackQuery.Message.Chat.ID, c.Update.CallbackQuery.Message.MessageID, "抱歉，因为操作无效，此操作无法进行，请重新发起操作后再试。")

	routeHash, actionDataHash := c.Bot.routeHashAndActionHashFromData(c.Update.CallbackQuery.Data)
	if routeHash == "" || actionDataHash == "" {
		c.Bot.MayRequest(callbackQueryActionInvalidErrMessage)
		return
	}

	route, ok := d.callbackQueryHandlersRoute[routeHash]
	if !ok || route == "" {
		return
	}

	handler, ok := d.callbackQueryHandlers[routeHash]
	if !ok || handler == nil {
		c.Bot.MayRequest(callbackQueryActionInvalidErrMessage)
		return
	}

	c.initForCallbackQuery(route, routeHash, actionDataHash)

	err := c.fetchActionDataForCallbackQueryHandler()
	if err != nil {
		d.Logger.Errorf("failed to fetch the callback query action data for handler %s: %v", c.callbackQueryHandlerRoute, err)
		return
	}
	if c.callbackQueryHandlerActionDataIsEmpty {
		c.Bot.MayRequest(callbackQueryActionInvalidErrMessage)
		return
	}

	d.dispatchInGoroutine(func() {
		_, _ = handler(c)
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

	d.Logger.Tracef("[我的成员信息更新｜%s] [%s (%s)] %s (%s): 成员状态自 %s 变更为 %s",
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
			d.Logger.Tracef("已退出频道 %s (%d)", c.Update.MyChatMember.Chat.Title, c.Update.MyChatMember.Chat.ID)
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

		d.Logger.Tracef("已加入频道 %s (%d)", c.Update.MyChatMember.Chat.Title, c.Update.MyChatMember.Chat.ID)
	}
}

func (d *Dispatcher) Dispatch(bot *tgbotapi.BotAPI, update tgbotapi.Update, rueidisClient rueidis.Client) {
	for _, m := range d.middlewares {
		m(NewContext(bot, update, d.Logger, rueidisClient), func() {})
	}

	ctx := NewContext(bot, update, d.Logger, rueidisClient)
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
