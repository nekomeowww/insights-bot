package recap

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/nekomeowww/insights-bot/internal/models/tgchats"
	"github.com/nekomeowww/insights-bot/pkg/bots/tgbot"
	"github.com/nekomeowww/insights-bot/pkg/logger"
	"github.com/nekomeowww/insights-bot/pkg/types/telegram"
	"github.com/samber/lo"
	"go.uber.org/fx"
)

var (
	_ tgbot.CommandHandler = (*EnableRecapCommandHandler)(nil)
	_ tgbot.CommandHandler = (*DisableRecapCommandHandler)(nil)
)

type NewEnableRecapCommandHandlerParams struct {
	fx.In

	Logger *logger.Logger
	TgChat *tgchats.Model
}

type EnableRecapCommandHandler struct {
	logger  *logger.Logger
	tgchats *tgchats.Model
}

func NewEnableRecapCommandHandler() func(NewEnableRecapCommandHandlerParams) *EnableRecapCommandHandler {
	return func(param NewEnableRecapCommandHandlerParams) *EnableRecapCommandHandler {
		return &EnableRecapCommandHandler{
			logger:  param.Logger,
			tgchats: param.TgChat,
		}
	}
}

func (h EnableRecapCommandHandler) Command() string {
	return "enable_recap"
}

func (h EnableRecapCommandHandler) CommandHelp() string {
	return "开启聊天记录回顾（需要管理权限）"
}

func (h *EnableRecapCommandHandler) Handle(c *tgbot.Context) error {
	chatType := telegram.ChatType(c.Update.Message.Chat.Type)
	if !lo.Contains([]telegram.ChatType{telegram.ChatTypeGroup, telegram.ChatTypeSuperGroup}, chatType) {
		return tgbot.NewMessageError("聊天记录回顾功能只有群组和超级群组可以配置开启哦！").WithReply(c.Update.Message)
	}

	hasTogglingRecapPermission, err := checkTogglingRecapPermission(c.Bot, c.Update.Message.Chat.ID, c.Update.Message.From.ID)
	if err != nil {
		return tgbot.NewExceptionError(err).WithMessage("聊天记录回顾功能开启失败，请稍后再试！").WithReply(c.Update.Message)
	}
	if !hasTogglingRecapPermission {
		return tgbot.NewMessageError("你没有权限开启聊天记录回顾功能哦！").WithReply(c.Update.Message)
	}

	botMember, err := c.Bot.GetChatMember(tgbotapi.GetChatMemberConfig{
		ChatConfigWithUser: tgbotapi.ChatConfigWithUser{
			ChatID: c.Update.Message.Chat.ID,
			UserID: c.Bot.Self.ID,
		},
	})
	if err != nil {
		return tgbot.NewExceptionError(err).WithMessage("聊天记录回顾功能开启失败，请稍后再试！").WithReply(c.Update.Message)
	}
	if !lo.Contains([]telegram.MemberStatus{
		telegram.MemberStatusAdministrator,
	}, telegram.MemberStatus(botMember.Status)) {
		return tgbot.NewMessageError("请先将机器人设为群组管理员，然后再开启聊天记录回顾功能哦！").WithReply(c.Update.Message)
	}

	err = h.tgchats.EnableChatHistoriesRecap(c.Update.Message.Chat.ID, chatType)
	if err != nil {
		return tgbot.NewExceptionError(err).WithMessage("聊天记录回顾功能开启失败，请稍后再试！").WithReply(c.Update.Message)
	}

	message := tgbotapi.NewMessage(c.Update.Message.Chat.ID, "聊天记录回顾功能已开启，开启后将会自动收集群组中的聊天记录并定时发送聊天回顾快报！")
	message.ReplyToMessageID = c.Update.Message.MessageID
	c.Bot.MustSend(message)
	return nil
}

type NewDisableRecapCommandHandlerParams struct {
	fx.In

	Logger *logger.Logger
	TgChat *tgchats.Model
}

type DisableRecapCommandHandler struct {
	logger  *logger.Logger
	tgchats *tgchats.Model
}

func NewDisableRecapCommandHandler() func(NewDisableRecapCommandHandlerParams) *DisableRecapCommandHandler {
	return func(param NewDisableRecapCommandHandlerParams) *DisableRecapCommandHandler {
		return &DisableRecapCommandHandler{
			logger:  param.Logger,
			tgchats: param.TgChat,
		}
	}
}

func (h DisableRecapCommandHandler) Command() string {
	return "disable_recap"
}

func (h DisableRecapCommandHandler) CommandHelp() string {
	return "关闭聊天记录回顾（需要管理权限）"
}

func (h *DisableRecapCommandHandler) Handle(c *tgbot.Context) error {
	chatType := telegram.ChatType(c.Update.Message.Chat.Type)
	if !lo.Contains([]telegram.ChatType{telegram.ChatTypeGroup, telegram.ChatTypeSuperGroup}, chatType) {
		return tgbot.NewMessageError("聊天记录回顾功能只有群组和超级群组可以配置关闭哦！").WithReply(c.Update.Message)
	}

	hasTogglingRecapPermission, err := checkTogglingRecapPermission(c.Bot, c.Update.Message.Chat.ID, c.Update.Message.From.ID)
	if err != nil {
		return tgbot.NewExceptionError(err).WithMessage("聊天记录回顾功能开启失败，请稍后再试！").WithReply(c.Update.Message)
	}
	if !hasTogglingRecapPermission {
		return tgbot.NewMessageError("你没有权限关闭聊天记录回顾功能哦！").WithReply(c.Update.Message)
	}

	botMember, err := c.Bot.GetChatMember(tgbotapi.GetChatMemberConfig{
		ChatConfigWithUser: tgbotapi.ChatConfigWithUser{
			ChatID: c.Update.Message.Chat.ID,
			UserID: c.Bot.Self.ID,
		},
	})
	if err != nil {
		return tgbot.NewExceptionError(err).WithMessage("聊天记录回顾功能开启失败，请稍后再试！").WithReply(c.Update.Message)
	}
	if !lo.Contains([]telegram.MemberStatus{
		telegram.MemberStatusAdministrator,
	}, telegram.MemberStatus(botMember.Status)) {
		return tgbot.NewMessageError("现在机器人不是群组管理员，已经不会记录任何聊天记录了。如果需要打开聊天记录回顾功能，请先将机器人设为群组管理员。").WithReply(c.Update.Message)
	}

	err = h.tgchats.DisableChatHistoriesRecap(c.Update.Message.Chat.ID, chatType)
	if err != nil {
		return tgbot.NewExceptionError(err).WithMessage("聊天记录回顾功能关闭失败，请稍后再试！").WithReply(c.Update.Message)
	}

	message := tgbotapi.NewMessage(c.Update.Message.Chat.ID, "聊天记录回顾功能已关闭，关闭后将不会自动收集群组中的聊天记录并定时发送聊天回顾快报了。")
	message.ReplyToMessageID = c.Update.Message.MessageID
	c.Bot.MustSend(message)
	return nil
}

func checkTogglingRecapPermission(bot *tgbot.Bot, chatID, userID int64) (bool, error) {
	member, err := bot.GetChatMember(tgbotapi.GetChatMemberConfig{
		ChatConfigWithUser: tgbotapi.ChatConfigWithUser{
			ChatID: chatID,
			UserID: userID,
		},
	})
	if err != nil {
		return false, err
	}
	if !lo.Contains([]telegram.MemberStatus{
		telegram.MemberStatusCreator,
		telegram.MemberStatusAdministrator,
	}, telegram.MemberStatus(member.Status)) {
		return false, nil
	}

	return true, nil
}
