package telegram

import (
	"context"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/gookit/color"
	"github.com/samber/lo"
	"go.uber.org/fx"

	"github.com/nekomeowww/insights-bot/internal/bots/telegram/dispatcher"
	"github.com/nekomeowww/insights-bot/internal/bots/telegram/handlers"
	"github.com/nekomeowww/insights-bot/internal/configs"
	"github.com/nekomeowww/insights-bot/pkg/handler"
	"github.com/nekomeowww/insights-bot/pkg/logger"
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
}

type Bot struct {
	*tgbotapi.BotAPI

	Config     *configs.Config
	Logger     *logger.Logger
	Dispatcher *dispatcher.Dispatcher

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
			BotAPI:     b,
			Logger:     param.Logger,
			Dispatcher: param.Dispatcher,
			closeChan:  make(chan struct{}, 1),
		}

		param.Lifecycle.Append(fx.Hook{
			OnStop: func(ctx context.Context) error {
				bot.StopPull(ctx)
				return nil
			},
		})

		param.Logger.Infof("Authorized as bot @%s", bot.Self.UserName)
		param.Handlers.RegisterHandlers()
		return bot, nil
	}
}

func (b *Bot) StopPull(ctx context.Context) {
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

func (b *Bot) MapChatTypeToChineseText(chatType string) string {
	switch chatType {
	case "private":
		return "私聊"
	case "group":
		return "群组"
	case "supergroup":
		return "超级群组"
	case "channel":
		return "频道"
	default:
		return "未知"
	}
}

func (b *Bot) MapMemberStatusToChineseText(memberStatus string) string {
	switch memberStatus {
	case "creator":
		return "创建者"
	case "administrator":
		return "管理员"
	case "member":
		return "成员"
	case "restricted":
		return "受限成员"
	case "left":
		return "已离开"
	case "kicked":
		return "已被踢出"
	default:
		return "未知"
	}
}

func (b *Bot) PullUpdates() {
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
				identityStrings := make([]string, 0)
				if update.Message.From.FirstName != "" {
					identityStrings = append(identityStrings, update.Message.From.FirstName)
				}
				if update.Message.From.LastName != "" {
					identityStrings = append(identityStrings, update.Message.From.LastName)
				}
				if update.Message.From.UserName != "" {
					identityStrings = append(identityStrings, "@"+update.Message.From.UserName)
				}

				if update.Message.Chat.Type == "private" {
					b.Logger.Infof("[消息｜%s] %s (%s): %s",
						b.MapChatTypeToChineseText(update.Message.Chat.Type),
						strings.Join(identityStrings, " "),
						color.FgYellow.Render(update.Message.From.ID),
						lo.Ternary(update.Message.Text == "", "<empty or contains medias>", update.Message.Text),
					)
				} else {
					b.Logger.Infof("[消息｜%s] [%s (%s)] %s (%s): %s",
						b.MapChatTypeToChineseText(update.Message.Chat.Type),
						color.FgGreen.Render(update.Message.Chat.Title),
						color.FgYellow.Render(update.Message.Chat.ID),
						strings.Join(identityStrings, " "),
						color.FgYellow.Render(update.Message.From.ID),
						lo.Ternary(update.Message.Text == "", "<empty or contains medias>", update.Message.Text),
					)
				}

				b.Dispatcher.DispatchMessage(handler.NewContext(b.BotAPI, update))
			}
			if update.MyChatMember != nil {
				identityStrings := make([]string, 0)
				if update.MyChatMember.From.FirstName != "" {
					identityStrings = append(identityStrings, update.MyChatMember.From.FirstName)
				}
				if update.MyChatMember.From.LastName != "" {
					identityStrings = append(identityStrings, update.MyChatMember.From.LastName)
				}
				if update.MyChatMember.From.UserName != "" {
					identityStrings = append(identityStrings, "@"+update.MyChatMember.From.UserName)
				}

				oldMemberStatus := update.MyChatMember.OldChatMember.Status
				newMemberStatus := update.MyChatMember.NewChatMember.Status

				b.Logger.Infof("[我的成员信息更新｜%s] [%s (%s)] %s (%s): 成员状态自 %s 变更为 %s",
					b.MapChatTypeToChineseText(update.MyChatMember.Chat.Type),
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
						continue
					}

					_, err := b.BotAPI.GetChat(tgbotapi.ChatInfoConfig{
						ChatConfig: tgbotapi.ChatConfig{
							ChatID: update.MyChatMember.Chat.ID,
						},
					})
					if err != nil {
						b.Logger.Error(err)
						continue
					}

					b.Logger.Infof("已加入频道 %s (%d)", update.MyChatMember.Chat.Title, update.MyChatMember.Chat.ID)
				}
			}
			if update.ChannelPost != nil {
				b.Logger.Infof("[频道消息｜%s] [%s (%s)]: %s",
					b.MapChatTypeToChineseText(update.ChannelPost.Chat.Type),
					color.FgGreen.Render(update.ChannelPost.Chat.Title),
					color.FgYellow.Render(update.ChannelPost.Chat.ID),
					lo.Ternary(update.ChannelPost.Text == "", "<empty or contains medias>", update.ChannelPost.Text),
				)
				b.Dispatcher.DispatchChannelPost(handler.NewContext(b.BotAPI, update))
			}
		case <-b.closeChan:
			b.Logger.Info("stopped to receiving updates")
			return
		}
	}
}

func Run() func(bot *Bot) {
	return func(bot *Bot) {
		go bot.PullUpdates()
	}
}
