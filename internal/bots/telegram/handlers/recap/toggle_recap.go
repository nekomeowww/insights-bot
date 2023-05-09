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

func checkTogglingRecapPermission(chatID, userID int64, update tgbotapi.Update, bot *tgbot.Bot) error {
	member, err := bot.GetChatMember(tgbotapi.GetChatMemberConfig{
		ChatConfigWithUser: tgbotapi.ChatConfigWithUser{
			ChatID: chatID,
			UserID: userID,
		},
	})
	if err != nil {
		return tgbot.NewExceptionError(err).WithMessage("聊天记录回顾功能开启失败，请稍后再试！").WithReply(update.Message)
	}
	if !lo.Contains([]telegram.MemberStatus{
		telegram.MemberStatusCreator,
		telegram.MemberStatusAdministrator,
	}, telegram.MemberStatus(member.Status)) {
		return tgbot.NewMessageError("你没有权限开启聊天记录回顾功能哦！").WithReply(update.Message)
	}

	return nil
}

func checkBotMember(chatID int64, update tgbotapi.Update, bot *tgbot.Bot) error {
	botMember, err := bot.GetChatMember(tgbotapi.GetChatMemberConfig{
		ChatConfigWithUser: tgbotapi.ChatConfigWithUser{
			ChatID: chatID,
			UserID: bot.Self.ID,
		},
	})
	if err != nil {
		return tgbot.NewExceptionError(err).WithMessage("聊天记录回顾功能开启失败，请稍后再试！").WithReply(update.Message)
	}
	if !lo.Contains([]telegram.MemberStatus{
		telegram.MemberStatusAdministrator,
	}, telegram.MemberStatus(botMember.Status)) {
		return tgbot.NewMessageError("现在机器人不是群组管理员，已经不会记录任何聊天记录了。如果需要打开聊天记录回顾功能，请先将机器人设为群组管理员。").WithReply(update.Message)
	}

	return nil
}

func checkToggle(update tgbotapi.Update, bot *tgbot.Bot) error {
	chatType := telegram.ChatType(update.Message.Chat.Type)
	if !lo.Contains([]telegram.ChatType{telegram.ChatTypeGroup, telegram.ChatTypeSuperGroup}, chatType) {
		return tgbot.NewMessageError("聊天记录回顾功能只有群组和超级群组可以配置开启哦！").WithReply(update.Message)
	}

	err := checkTogglingRecapPermission(update.Message.Chat.ID, update.Message.From.ID, update, bot)
	if err != nil {
		return err
	}

	err = checkBotMember(update.Message.Chat.ID, update, bot)
	if err != nil {
		return err
	}

	return nil
}

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

func (h *EnableRecapCommandHandler) Handle(c *tgbot.Context) (tgbot.Response, error) {
	chatType := telegram.ChatType(c.Update.Message.Chat.Type)

	err := checkToggle(c.Update, c.Bot)
	if err != nil {
		return nil, err
	}

	err = h.tgchats.EnableChatHistoriesRecap(c.Update.Message.Chat.ID, chatType)
	if err != nil {
		return nil, tgbot.NewExceptionError(err).WithMessage("聊天记录回顾功能开启失败，请稍后再试！").WithReply(c.Update.Message)
	}

	return c.NewMessageReplyTo("聊天记录回顾功能已开启，开启后将会自动收集群组中的聊天记录并定时发送聊天回顾快报！", c.Update.Message.MessageID), nil
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

func (h *DisableRecapCommandHandler) Handle(c *tgbot.Context) (tgbot.Response, error) {
	chatType := telegram.ChatType(c.Update.Message.Chat.Type)

	err := checkToggle(c.Update, c.Bot)
	if err != nil {
		return nil, err
	}

	err = h.tgchats.DisableChatHistoriesRecap(c.Update.Message.Chat.ID, chatType)
	if err != nil {
		return nil, tgbot.NewExceptionError(err).WithMessage("聊天记录回顾功能关闭失败，请稍后再试！").WithReply(c.Update.Message)
	}

	return c.NewMessageReplyTo("聊天记录回顾功能已关闭，关闭后将不会自动收集群组中的聊天记录并定时发送聊天回顾快报了。", c.Update.Message.MessageID), nil
}
