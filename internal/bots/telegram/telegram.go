package telegram

import (
	"context"
	"runtime/debug"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/gookit/color"
	"github.com/samber/lo"
	"go.uber.org/fx"

	"github.com/nekomeowww/insights-bot/internal/bots/telegram/dispatcher"
	"github.com/nekomeowww/insights-bot/internal/bots/telegram/handlers"
	"github.com/nekomeowww/insights-bot/internal/configs"
	"github.com/nekomeowww/insights-bot/internal/models/chat_histories"
	telegram_bot "github.com/nekomeowww/insights-bot/pkg/bots/telegram"
	"github.com/nekomeowww/insights-bot/pkg/handler"
	"github.com/nekomeowww/insights-bot/pkg/logger"
	"github.com/nekomeowww/insights-bot/pkg/types/telegram"
	"github.com/nekomeowww/insights-bot/pkg/utils"
)

func NewModules() fx.Option {
	return fx.Options(
		fx.Provide(NewBot()),
		fx.Options(dispatcher.NewModules()),
		fx.Options(handlers.NewModules()),
	)
}

type NewBotParam struct {
	fx.In

	Lifecycle fx.Lifecycle

	Config     *configs.Config
	Logger     *logger.Logger
	Dispatcher *dispatcher.Dispatcher
	Handlers   *handlers.Handlers

	ChatHistories *chat_histories.ChatHistoriesModel
}

type Bot struct {
	*tgbotapi.BotAPI

	Config     *configs.Config
	Logger     *logger.Logger
	Dispatcher *dispatcher.Dispatcher

	ChatHistories *chat_histories.ChatHistoriesModel

	alreadyClose bool
	closeChan    chan struct{}
}

func NewBot() func(param NewBotParam) (*Bot, error) {
	return func(param NewBotParam) (*Bot, error) {
		if param.Config.TelegramBotToken == "" {
			param.Logger.Fatal("must supply a valid telegram bot token in configs or environment variable")
		}

		b, err := tgbotapi.NewBotAPI(param.Config.TelegramBotToken)
		if err != nil {
			return nil, err
		}

		bot := &Bot{
			BotAPI:        b,
			Logger:        param.Logger,
			Dispatcher:    param.Dispatcher,
			ChatHistories: param.ChatHistories,
			closeChan:     make(chan struct{}, 1),
		}

		param.Lifecycle.Append(fx.Hook{
			OnStop: func(ctx context.Context) error {
				bot.stopPull(ctx)
				return nil
			},
		})

		param.Logger.Infof("Authorized as bot @%s", bot.Self.UserName)
		param.Handlers.RegisterHandlers()
		return bot, nil
	}
}

func (b *Bot) stopPull(ctx context.Context) {
	if b.alreadyClose {
		return
	}

	_ = utils.Invoke0(func() error {
		b.alreadyClose = true
		b.StopReceivingUpdates()
		b.closeChan <- struct{}{}
		close(b.closeChan)

		return nil
	}, utils.WithContext(ctx))
}

func (b *Bot) MapChatTypeToChineseText(chatType telegram.ChatType) string {
	switch chatType {
	case telegram.ChatTypePrivate:
		return "私聊"
	case telegram.ChatTypeGroup:
		return "群组"
	case telegram.ChatTypeSuperGroup:
		return "超级群组"
	case telegram.ChatTypeChannel:
		return "频道"
	default:
		return "未知"
	}
}

func (b *Bot) MapMemberStatusToChineseText(memberStatus telegram.MemberStatus) string {
	switch memberStatus {
	case telegram.MemberStatusCreator:
		return "创建者"
	case telegram.MemberStatusAdministrator:
		return "管理员"
	case telegram.MemberStatusMember:
		return "成员"
	case telegram.MemberStatusRestricted:
		return "受限成员"
	case telegram.MemberStatusLeft:
		return "已离开"
	case telegram.MemberStatusKicked:
		return "已被踢出"
	default:
		return "未知"
	}
}

func (b *Bot) handleMessage(update tgbotapi.Update) {
	identityStrings := make([]string, 0)
	identityStrings = append(identityStrings, telegram_bot.FullNameFromFirstAndLastName(update.Message.From.FirstName, update.Message.From.LastName))
	if update.Message.From.UserName != "" {
		identityStrings = append(identityStrings, "@"+update.Message.From.UserName)
	}
	if update.Message.Chat.Type == "private" {
		b.Logger.Infof("[消息｜%s] %s (%s): %s",
			b.MapChatTypeToChineseText(telegram.ChatType(update.Message.Chat.Type)),
			strings.Join(identityStrings, " "),
			color.FgYellow.Render(update.Message.From.ID),
			lo.Ternary(update.Message.Text == "", "<empty or contains medias>", update.Message.Text),
		)
	} else {
		b.Logger.Infof("[消息｜%s] [%s (%s)] %s (%s): %s",
			b.MapChatTypeToChineseText(telegram.ChatType(update.Message.Chat.Type)),
			color.FgGreen.Render(update.Message.Chat.Title),
			color.FgYellow.Render(update.Message.Chat.ID),
			strings.Join(identityStrings, " "),
			color.FgYellow.Render(update.Message.From.ID),
			lo.Ternary(update.Message.Text == "", "<empty or contains medias>", update.Message.Text),
		)
	}
	if update.Message.Command() != "" {
		go func() {
			defer func() {
				if err := recover(); err != nil {
					b.Logger.Errorf("Panic recovered from command dispatcher, %v\n%s", err, debug.Stack())
					return
				}
			}()

			b.Dispatcher.DispatchCommand(handler.NewContext(b.BotAPI, update, b.Logger))
		}()

	} else {
		go func() {
			defer func() {
				if err := recover(); err != nil {
					b.Logger.Errorf("Panic recovered from message dispatcher, %v\n%s", err, debug.Stack())
					return
				}
			}()

			b.Dispatcher.DispatchMessage(handler.NewContext(b.BotAPI, update, b.Logger))
		}()
	}
}

func (b *Bot) handleChatMember(update tgbotapi.Update) {
	identityStrings := make([]string, 0)
	identityStrings = append(identityStrings, telegram_bot.FullNameFromFirstAndLastName(update.MyChatMember.From.FirstName, update.MyChatMember.From.LastName))
	if update.MyChatMember.From.UserName != "" {
		identityStrings = append(identityStrings, "@"+update.MyChatMember.From.UserName)
	}

	oldMemberStatus := telegram.MemberStatus(update.MyChatMember.OldChatMember.Status)
	newMemberStatus := telegram.MemberStatus(update.MyChatMember.NewChatMember.Status)

	b.Logger.Infof("[我的成员信息更新｜%s] [%s (%s)] %s (%s): 成员状态自 %s 变更为 %s",
		b.MapChatTypeToChineseText(telegram.ChatType(update.MyChatMember.Chat.Type)),
		color.FgGreen.Render(update.MyChatMember.Chat.Title),
		color.FgYellow.Render(update.MyChatMember.Chat.ID),
		strings.Join(identityStrings, " "),
		color.FgYellow.Render(update.MyChatMember.From.ID),
		b.MapMemberStatusToChineseText(oldMemberStatus),
		b.MapMemberStatusToChineseText(newMemberStatus),
	)
	switch update.MyChatMember.Chat.Type {
	case "channel":
		if newMemberStatus != "administrator" {
			b.Logger.Infof("已退出频道 %s (%d)", update.MyChatMember.Chat.Title, update.MyChatMember.Chat.ID)
			return
		}

		_, err := b.BotAPI.GetChat(tgbotapi.ChatInfoConfig{
			ChatConfig: tgbotapi.ChatConfig{
				ChatID: update.MyChatMember.Chat.ID,
			},
		})
		if err != nil {
			b.Logger.Error(err)
			return
		}

		b.Logger.Infof("已加入频道 %s (%d)", update.MyChatMember.Chat.Title, update.MyChatMember.Chat.ID)
	}
}

func (b *Bot) handleChannelPost(update tgbotapi.Update) {
	b.Logger.Infof("[频道消息｜%s] [%s (%s)]: %s",
		b.MapChatTypeToChineseText(telegram.ChatType(update.ChannelPost.Chat.Type)),
		color.FgGreen.Render(update.ChannelPost.Chat.Title),
		color.FgYellow.Render(update.ChannelPost.Chat.ID),
		lo.Ternary(update.ChannelPost.Text == "", "<empty or contains medias>", update.ChannelPost.Text),
	)

	go func() {
		defer func() {
			if err := recover(); err != nil {
				b.Logger.Errorf("Panic recovered from channel post dispatcher, %v\n%s", err, debug.Stack())
				return
			}
		}()

		b.Dispatcher.DispatchChannelPost(handler.NewContext(b.BotAPI, update, b.Logger))
	}()
}

func (b *Bot) handleCallbackQuery(update tgbotapi.Update) {
	b.Logger.Infof("[回调查询｜%s] [%s (%s)]: %s",
		b.MapChatTypeToChineseText(telegram.ChatType(update.CallbackQuery.Message.Chat.Type)),
		color.FgGreen.Render(update.CallbackQuery.Message.Chat.Title),
		color.FgYellow.Render(update.CallbackQuery.Message.Chat.ID),
		update.CallbackQuery.Data,
	)

	go func() {
		defer func() {
			if err := recover(); err != nil {
				b.Logger.Errorf("Panic recovered from callback query dispatcher, %v\n%s", err, debug.Stack())
				return
			}
		}()

		b.Dispatcher.DispatchCallbackQuery(handler.NewContext(b.BotAPI, update, b.Logger))
	}()
}

func (b *Bot) pullUpdates() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.GetUpdatesChan(u)
	for {
		if b.alreadyClose {
			b.Logger.Info("stopped to receiving updates")
			return
		}

		select {
		case update := <-updates:
			if update.Message != nil {
				b.handleMessage(update)
				continue
			}
			if update.MyChatMember != nil {
				b.handleChatMember(update)
				continue
			}
			if update.ChannelPost != nil {
				b.handleChannelPost(update)
				continue
			}
			if update.CallbackQuery != nil {
				b.handleCallbackQuery(update)
				continue
			}
		case <-b.closeChan:
			b.Logger.Info("stopped to receiving updates")
			return
		}
	}
}

func Run() func(bot *Bot) {
	return func(bot *Bot) {
		go bot.pullUpdates()
	}
}
