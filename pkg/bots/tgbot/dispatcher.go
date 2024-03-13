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
	"go.uber.org/zap"

	"github.com/nekomeowww/insights-bot/pkg/i18n"
	"github.com/nekomeowww/insights-bot/pkg/logger"
	"github.com/nekomeowww/insights-bot/pkg/types/telegram"
)

type Dispatcher struct {
	Logger *logger.Logger
	I18n   *i18n.I18n

	helpCommand                *helpCommandHandler
	cancelCommand              *cancelCommandHandler
	startCommandHandler        *startCommandHandler
	middlewares                []MiddlewareFunc
	commandHandlers            map[string]HandleFunc
	channelPostHandlers        []Handler
	callbackQueryHandlers      map[string]HandleFunc
	callbackQueryHandlersRoute map[string]string
	leftChatMemberHandlers     []Handler
	newChatMembersHandlers     []Handler
	myChatMemberHandlers       []Handler
	chatMigrationFromHandlers  []Handler
}

func NewDispatcher() func(logger *logger.Logger, i18n *i18n.I18n) *Dispatcher {
	return func(logger *logger.Logger, i18n *i18n.I18n) *Dispatcher {
		d := &Dispatcher{
			Logger:                     logger,
			I18n:                       i18n,
			helpCommand:                newHelpCommandHandler(),
			cancelCommand:              newCancelCommandHandler(),
			startCommandHandler:        newStartCommandHandler(),
			middlewares:                make([]MiddlewareFunc, 0),
			commandHandlers:            make(map[string]HandleFunc),
			channelPostHandlers:        make([]Handler, 0),
			callbackQueryHandlers:      make(map[string]HandleFunc),
			callbackQueryHandlersRoute: make(map[string]string),
			leftChatMemberHandlers:     make([]Handler, 0),
			newChatMembersHandlers:     make([]Handler, 0),
			myChatMemberHandlers:       make([]Handler, 0),
			chatMigrationFromHandlers:  make([]Handler, 0),
		}

		d.startCommandHandler.helpCommandHandler = d.helpCommand

		d.OnCommandGroup(func(c *Context) string {
			return c.T("system.commands.groups.basic.name")
		}, []Command{
			{Command: d.helpCommand.Command(), HelpMessage: d.helpCommand.CommandHelp, Handler: NewHandler(d.helpCommand.handle)},
			{Command: d.cancelCommand.Command(), HelpMessage: d.cancelCommand.CommandHelp, Handler: NewHandler(d.cancelCommand.handle)},
			{Command: d.startCommandHandler.Command(), HelpMessage: d.startCommandHandler.CommandHelp, Handler: NewHandler(d.startCommandHandler.handle)},
		})
		d.OnCallbackQuery("nop", NewHandler(func(ctx *Context) (Response, error) {
			return nil, nil
		}))

		return d
	}
}

func (d *Dispatcher) Use(middleware MiddlewareFunc) {
	d.middlewares = append(d.middlewares, middleware)
}

func (d *Dispatcher) OnCommand(cmd string, commandHelp func(c *Context) string, h Handler) {
	d.helpCommand.defaultGroup.commands = append(d.helpCommand.defaultGroup.commands, Command{
		Command:     cmd,
		HelpMessage: commandHelp,
	})

	d.commandHandlers[cmd] = h.Handle
}

func (d *Dispatcher) OnCommandGroup(groupName func(*Context) string, group []Command) {
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
		d.Logger.Debug(fmt.Sprintf("[消息｜%s] %s (%s): %s",
			MapChatTypeToChineseText(telegram.ChatType(c.Update.Message.Chat.Type)),
			strings.Join(identityStrings, " "),
			color.FgYellow.Render(c.Update.Message.From.ID),
			lo.Ternary(c.Update.Message.Text == "", "<empty or contains medias>", c.Update.Message.Text)),
		)
	} else {
		d.Logger.Debug(fmt.Sprintf("[消息｜%s] [%s (%s)] %s (%s): %s",
			MapChatTypeToChineseText(telegram.ChatType(c.Update.Message.Chat.Type)),
			color.FgGreen.Render(c.Update.Message.Chat.Title),
			color.FgYellow.Render(c.Update.Message.Chat.ID),
			strings.Join(identityStrings, " "),
			color.FgYellow.Render(c.Update.Message.From.ID),
			lo.Ternary(c.Update.Message.Text == "", "<empty or contains medias>", c.Update.Message.Text)),
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
	d.Logger.Debug(fmt.Sprintf("[频道消息｜%s] [%s (%s)]: %s",
		MapChatTypeToChineseText(telegram.ChatType(c.Update.ChannelPost.Chat.Type)),
		color.FgGreen.Render(c.Update.ChannelPost.Chat.Title),
		color.FgYellow.Render(c.Update.ChannelPost.Chat.ID),
		lo.Ternary(c.Update.ChannelPost.Text == "", "<empty or contains medias>", c.Update.ChannelPost.Text),
	))

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

		if c.callbackQueryHandlerRoute == "" {
			d.Logger.Error(fmt.Sprintf("[回调查询｜%s] [%s (%s)] %s (%s) : %s (Raw Data) \n%s\n%s\n\n%s\n%s\n",
				MapChatTypeToChineseText(telegram.ChatType(c.Update.CallbackQuery.Message.Chat.Type)),
				color.FgGreen.Render(c.Update.CallbackQuery.Message.Chat.Title),
				color.FgYellow.Render(c.Update.CallbackQuery.Message.Chat.ID),
				strings.Join(identityStrings, " "),
				color.FgYellow.Render(c.Update.CallbackQuery.From.ID),
				c.Update.CallbackData(),
				color.FgRed.Render("无法调度 Callback Query，检测到缺少路由。"),
				color.FgRed.Render("Unable to dispatch Callback Query due to missing route DETECTED."),
				color.FgRed.Render("大多数情况下，发生这种情况的原因是相应的处理程序没有通过 OnCallbackQuery(...) 方法正确注册，或者内部派发器未能与之匹配，请检查已注册的处理程序及其路由，然后再试一次。"),
				color.FgRed.Render("For most of the time, this happens when the corresponding handler wasn't registered properly through OnCallbackQuery(...) method or internal dispatcher failed to match it, please check registered handlers and the route of them and then try again."),
			),
				zap.String("route", c.callbackQueryHandlerRoute),
				zap.String("route_hash", c.callbackQueryHandlerRouteHash),
				zap.String("action_data_hash", c.callbackQueryHandlerActionDataHash),
			)
		} else if c.callbackQueryHandlerActionData == "" {
			d.Logger.Error(fmt.Sprintf("[回调查询｜%s] [%s (%s)] %s (%s) : %s (Raw Data) \n%s\n%s\n\n%s\n%s\n",
				MapChatTypeToChineseText(telegram.ChatType(c.Update.CallbackQuery.Message.Chat.Type)),
				color.FgGreen.Render(c.Update.CallbackQuery.Message.Chat.Title),
				color.FgYellow.Render(c.Update.CallbackQuery.Message.Chat.ID),
				strings.Join(identityStrings, " "),
				color.FgYellow.Render(c.Update.CallbackQuery.From.ID),
				c.Update.CallbackData(),
				color.FgRed.Render("无法调度 Callback Query，检测到缺少操作数据。"),
				color.FgRed.Render("Unable to dispatch Callback Query due to missing action data DETECTED."),
				color.FgRed.Render("大多数情况下，当存储在回调查询数据中的操作数据为空、不存在于缓存中或无法从缓存中获取时会出现这种情况，请尝试刷新相应的缓存键并重试。"),
				color.FgRed.Render("For most of the time, this happens when the action data that stored into callback query data is either empty, not exist on cache, or failed to fetch from cache, please try to flush any corresponding cache keys and try again."),
			),
				zap.String("route", c.callbackQueryHandlerRoute),
				zap.String("route_hash", c.callbackQueryHandlerRouteHash),
				zap.String("action_data_hash", c.callbackQueryHandlerActionDataHash),
			)
		} else {
			d.Logger.Debug(fmt.Sprintf("[回调查询｜%s] [%s (%s)] %s (%s): %s: %s",
				MapChatTypeToChineseText(telegram.ChatType(c.Update.CallbackQuery.Message.Chat.Type)),
				color.FgGreen.Render(c.Update.CallbackQuery.Message.Chat.Title),
				color.FgYellow.Render(c.Update.CallbackQuery.Message.Chat.ID),
				strings.Join(identityStrings, " "),
				color.FgYellow.Render(c.Update.CallbackQuery.From.ID),
				c.callbackQueryHandlerRoute, c.callbackQueryHandlerActionData,
			),
				zap.String("route", c.callbackQueryHandlerRoute),
				zap.String("route_hash", c.callbackQueryHandlerRouteHash),
				zap.String("action_data_hash", c.callbackQueryHandlerActionDataHash),
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
		d.Logger.Error("failed to fetch the callback query action data for handler", zap.String("route", route), zap.Error(err))
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

func (d *Dispatcher) OnMyChatMember(handler Handler) {
	d.myChatMemberHandlers = append(d.myChatMemberHandlers, handler)
}

func (d *Dispatcher) dispatchMyChatMember(c *Context) {
	identityStrings := make([]string, 0)
	identityStrings = append(identityStrings, FullNameFromFirstAndLastName(c.Update.MyChatMember.From.FirstName, c.Update.MyChatMember.From.LastName))

	if c.Update.MyChatMember.From.UserName != "" {
		identityStrings = append(identityStrings, "@"+c.Update.MyChatMember.From.UserName)
	}

	oldMemberStatus := telegram.MemberStatus(c.Update.MyChatMember.OldChatMember.Status)
	newMemberStatus := telegram.MemberStatus(c.Update.MyChatMember.NewChatMember.Status)

	d.Logger.Debug(fmt.Sprintf("[我的成员信息更新｜%s] [%s (%s)] %s (%s): 成员状态自 %s 变更为 %s",
		MapChatTypeToChineseText(telegram.ChatType(c.Update.MyChatMember.Chat.Type)),
		color.FgGreen.Render(c.Update.MyChatMember.Chat.Title),
		color.FgYellow.Render(c.Update.MyChatMember.Chat.ID),
		strings.Join(identityStrings, " "),
		color.FgYellow.Render(c.Update.MyChatMember.From.ID),
		MapMemberStatusToChineseText(oldMemberStatus),
		MapMemberStatusToChineseText(newMemberStatus),
	))

	switch c.Update.MyChatMember.Chat.Type {
	case "channel":
		if newMemberStatus != "administrator" {
			d.Logger.Debug(fmt.Sprintf("已退出频道 %s (%d)", c.Update.MyChatMember.Chat.Title, c.Update.MyChatMember.Chat.ID))
			return
		}

		_, err := c.Bot.GetChat(tgbotapi.ChatInfoConfig{
			ChatConfig: tgbotapi.ChatConfig{
				ChatID: c.Update.MyChatMember.Chat.ID,
			},
		})
		if err != nil {
			d.Logger.Error(err.Error())
			return
		}

		d.Logger.Debug(fmt.Sprintf("已加入频道 %s (%d)", c.Update.MyChatMember.Chat.Title, c.Update.MyChatMember.Chat.ID))
	}

	d.dispatchInGoroutine(func() {
		for _, h := range d.myChatMemberHandlers {
			_, _ = h.Handle(c)
		}
	})
}

func (d *Dispatcher) OnLeftChatMember(h Handler) {
	d.leftChatMemberHandlers = append(d.leftChatMemberHandlers, h)
}

func (d *Dispatcher) dispatchLeftChatMember(c *Context) {
	identityStrings := make([]string, 0)
	identityStrings = append(identityStrings, FullNameFromFirstAndLastName(c.Update.Message.LeftChatMember.FirstName, c.Update.Message.LeftChatMember.LastName))

	if c.Update.Message.LeftChatMember.UserName != "" {
		identityStrings = append(identityStrings, "@"+c.Update.Message.LeftChatMember.UserName)
	}

	d.Logger.Debug(fmt.Sprintf("[成员信息更新｜%s] [%s (%s)] %s (%s) 离开了聊天",
		MapChatTypeToChineseText(telegram.ChatType(c.Update.Message.Chat.Type)),
		color.FgGreen.Render(c.Update.Message.Chat.Title),
		color.FgYellow.Render(c.Update.Message.Chat.ID),
		strings.Join(identityStrings, " "),
		color.FgYellow.Render(c.Update.Message.LeftChatMember.ID),
	))

	d.dispatchInGoroutine(func() {
		for _, h := range d.leftChatMemberHandlers {
			_, _ = h.Handle(c)
		}
	})
}

func (d *Dispatcher) OnNewChatMember(h Handler) {
	d.leftChatMemberHandlers = append(d.leftChatMemberHandlers, h)
}

func (d *Dispatcher) dispatchNewChatMember(c *Context) {
	identities := make([]string, len(c.Update.Message.NewChatMembers))

	for _, identity := range c.Update.Message.NewChatMembers {
		identityStrings := make([]string, 0)
		identityStrings = append(identityStrings, FullNameFromFirstAndLastName(identity.FirstName, identity.LastName))

		if identity.UserName != "" {
			identityStrings = append(identityStrings, "@"+identity.UserName)
		}

		identityStrings = append(identityStrings, fmt.Sprintf("(%s)", color.FgYellow.Render(identity.ID)))
		identities = append(identities, strings.Join(identityStrings, " "))
	}

	d.Logger.Debug(fmt.Sprintf("[成员信息更新｜%s] [%s (%s)] %s 加入了聊天",
		MapChatTypeToChineseText(telegram.ChatType(c.Update.Message.Chat.Type)),
		color.FgGreen.Render(c.Update.Message.Chat.Title),
		color.FgYellow.Render(c.Update.Message.Chat.ID),
		strings.Join(identities, ", "),
	))

	d.dispatchInGoroutine(func() {
		for _, h := range d.newChatMembersHandlers {
			_, _ = h.Handle(c)
		}
	})
}

func (d *Dispatcher) OnChatMigrationFrom(h Handler) {
	d.chatMigrationFromHandlers = append(d.chatMigrationFromHandlers, h)
}

func (d *Dispatcher) dispatchChatMigrationFrom(c *Context) {
	d.Logger.Debug(fmt.Sprintf("[群组迁移] 超级群组 [%s (%s)] 已迁移自群组 [%s (%s)]",
		color.FgGreen.Render(c.Update.Message.Chat.Title),
		color.FgYellow.Render(c.Update.Message.Chat.ID),
		color.FgGreen.Render(c.Update.Message.Chat.Title),
		color.FgYellow.Render(c.Update.Message.MigrateFromChatID),
	))

	d.dispatchInGoroutine(func() {
		for _, h := range d.chatMigrationFromHandlers {
			_, _ = h.Handle(c)
		}
	})
}

func (d *Dispatcher) dispatchChatMigrationTo(c *Context) {
	d.Logger.Debug(fmt.Sprintf("[群组迁移] 群组 [%s (%s)] 已迁移至超级群组 [%s (%s)]",
		color.FgGreen.Render(c.Update.Message.Chat.Title),
		color.FgYellow.Render(c.Update.Message.Chat.ID),
		color.FgGreen.Render(c.Update.Message.Chat.Title),
		color.FgYellow.Render(c.Update.Message.MigrateToChatID),
	))
}

func (d *Dispatcher) Dispatch(bot *tgbotapi.BotAPI, update tgbotapi.Update, rueidisClient rueidis.Client) {
	for _, m := range d.middlewares {
		m(NewContext(bot, update, d.Logger, d.I18n, rueidisClient), func() {})
	}

	ctx := NewContext(bot, update, d.Logger, d.I18n, rueidisClient)
	switch ctx.UpdateType() {
	case UpdateTypeMessage:
		d.dispatchMessage(ctx)
	case UpdateTypeEditedMessage:
		d.Logger.Debug("edited message is not supported yet")
	case UpdateTypeChannelPost:
		d.dispatchChannelPost(ctx)
	case UpdateTypeEditedChannelPost:
		d.Logger.Debug("edited channel post is not supported yet")
	case UpdateTypeInlineQuery:
		d.Logger.Debug("inline query is not supported yet")
	case UpdateTypeChosenInlineResult:
		d.Logger.Debug("chosen inline result is not supported yet")
	case UpdateTypeCallbackQuery:
		d.dispatchCallbackQuery(ctx)
	case UpdateTypeShippingQuery:
		d.Logger.Debug("shipping query is not supported yet")
	case UpdateTypePreCheckoutQuery:
		d.Logger.Debug("pre checkout query is not supported yet")
	case UpdateTypePoll:
		d.Logger.Debug("poll is not supported yet")
	case UpdateTypePollAnswer:
		d.Logger.Debug("poll answer is not supported yet")
	case UpdateTypeMyChatMember:
		d.dispatchMyChatMember(ctx)
	case UpdateTypeChatMember:
		d.Logger.Debug("chat member is not supported yet")
	case UpdateTypeLeftChatMember:
		d.dispatchLeftChatMember(ctx)
	case UpdateTypeNewChatMembers:
		d.dispatchNewChatMember(ctx)
	case UpdateTypeChatJoinRequest:
		d.Logger.Debug("chat join request is not supported yet")
	case UpdateTypeChatMigrationFrom:
		d.dispatchChatMigrationFrom(ctx)
	case UpdateTypeChatMigrationTo:
		d.dispatchChatMigrationTo(ctx)
	case UpdateTypeUnknown:
		d.Logger.Debug("unable to dispatch update due to unknown update type")
	default:
		d.Logger.Debug("unable to dispatch update due to unknown update type", zap.String("update_type", string(ctx.UpdateType())))
	}
}

func (d *Dispatcher) dispatchInGoroutine(f func()) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				d.Logger.Error("Panic recovered from command dispatcher",
					zap.Error(fmt.Errorf("panic error: %v", err)),
					zap.Stack("stack"),
				)
				fmt.Println("Panic recovered from command dispatcher: " + string(debug.Stack()))

				return
			}
		}()

		f()
	}()
}
