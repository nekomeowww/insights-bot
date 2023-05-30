package recap

import (
	"errors"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/samber/lo"

	"github.com/nekomeowww/insights-bot/ent"
	"github.com/nekomeowww/insights-bot/pkg/bots/tgbot"
	"github.com/nekomeowww/insights-bot/pkg/types/bots/handlers/recap"
	"github.com/nekomeowww/insights-bot/pkg/types/telegram"
	"github.com/nekomeowww/insights-bot/pkg/types/tgchat"
)

var (
	errOperationCanNotBeDone = errors.New("抱歉，此操作无法进行")
)

func checkToggle(ctx *tgbot.Context, chatID int64, fromID int64) error {
	if !lo.Contains([]telegram.ChatType{telegram.ChatTypeGroup, telegram.ChatTypeSuperGroup}, telegram.ChatType(ctx.Update.FromChat().Type)) {
		return fmt.Errorf("%w，%s", errOperationCanNotBeDone, "聊天记录回顾功能只有<b>群组</b>和<b>超级群组</b>的管理员可以配置哦！\n请将 Bot 添加到群组中，并配置 Bot 为管理员后使用管理员权限的用户账户为 Bot 进行配置吧。")
	}

	is, err := ctx.IsUserMemberStatus(fromID, []telegram.MemberStatus{
		telegram.MemberStatusCreator,
		telegram.MemberStatusAdministrator,
	})
	if err != nil {
		return err
	}
	if !is {
		return fmt.Errorf("%w，%s", errOperationCanNotBeDone, "开启/关闭聊天记录回顾功能需要<b>管理员</b>权限执行 /configure_recap 命令。")
	}

	is, err = ctx.IsBotAdministrator()
	if err != nil {
		return err
	}
	if !is {
		return fmt.Errorf("%w，%s", errOperationCanNotBeDone, "现在机器人不是<b>群组管理员</b>，已经不会记录任何聊天记录了。如果需要配置聊天记录回顾功能，请先将机器人设为群组管理员。然后再次执行命令后再试")
	}

	return nil
}

func checkAssignMode(ctx *tgbot.Context, chatID int64, fromID int64) error {
	if !lo.Contains([]telegram.ChatType{telegram.ChatTypeGroup, telegram.ChatTypeSuperGroup}, telegram.ChatType(ctx.Update.FromChat().Type)) {
		return fmt.Errorf("%w，%s", errOperationCanNotBeDone, "聊天记录回顾功能只有<b>群组</b>和<b>超级群组</b>的管理员可以配置哦！\n请将 Bot 添加到群组中，并配置 Bot 为管理员后使用管理员权限的用户账户为 Bot 进行配置吧。")
	}

	is, err := ctx.IsUserMemberStatus(fromID, []telegram.MemberStatus{
		telegram.MemberStatusCreator,
	})
	if err != nil {
		return err
	}
	if !is {
		return fmt.Errorf("%w，%s", errOperationCanNotBeDone, "配置聊天记录回顾功能的模式需要<b>群组创建者</b>权限执行 /configure_recap 命令。")
	}

	is, err = ctx.IsBotAdministrator()
	if err != nil {
		return err
	}
	if !is {
		return fmt.Errorf("%w，%s", errOperationCanNotBeDone, "现在机器人不是<b>群组管理员</b>，已经不会记录任何聊天记录了。如果需要配置聊天记录回顾功能，请先将机器人设为群组管理员。然后再次执行命令后再试")
	}

	return nil
}

func newRecapInlineKeyboardMarkup(ctx *tgbot.Context, chatID int64, fromID int64, currentRecapStatus bool, currentRecapMode tgchat.AutoRecapSendMode) (tgbotapi.InlineKeyboardMarkup, error) {
	data, err := ctx.Bot.AssignOneCallbackQueryData("recap/configure/toggle", recap.ConfigureRecapToggleActionData{Status: !currentRecapStatus, ChatID: chatID, FromID: fromID})
	if err != nil {
		return tgbotapi.InlineKeyboardMarkup{}, err
	}

	publicData, err := ctx.Bot.AssignOneCallbackQueryData("recap/configure/assign_mode", recap.ConfigureRecapAssignModeActionData{Mode: tgchat.AutoRecapSendModePublicly, ChatID: chatID, FromID: fromID})
	if err != nil {
		return tgbotapi.InlineKeyboardMarkup{}, err
	}

	privateData, err := ctx.Bot.AssignOneCallbackQueryData("recap/configure/assign_mode", recap.ConfigureRecapAssignModeActionData{Mode: tgchat.AutoRecapSendModeOnlyPrivateSubscriptions, ChatID: chatID, FromID: fromID})
	if err != nil {
		return tgbotapi.InlineKeyboardMarkup{}, err
	}

	completeData, err := ctx.Bot.AssignOneCallbackQueryData("recap/configure/complete", recap.ConfigureRecapCompleteActionData{ChatID: chatID, FromID: fromID})
	if err != nil {
		return tgbotapi.InlineKeyboardMarkup{}, err
	}

	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(lo.Ternary(!currentRecapStatus, "🔈 开启聊天记录回顾", "🔇 关闭聊天记录回顾"), data),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(lo.Ternary(currentRecapMode == tgchat.AutoRecapSendModePublicly, "🔘 公开模式", "公开模式"), publicData),
			tgbotapi.NewInlineKeyboardButtonData(lo.Ternary(currentRecapMode == tgchat.AutoRecapSendModeOnlyPrivateSubscriptions, "🔘 私聊订阅模式", "私聊订阅模式"), privateData),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("✅ 完成", completeData),
		),
	), nil
}

const (
	configureRecapGeneralInstructionMessage = "好的。请在下面点击你想配置的选项进行操作吧。"
)

func (h *CommandHandler) handleConfigureRecapCommand(c *tgbot.Context) (tgbot.Response, error) {
	chatType := telegram.ChatType(c.Update.Message.Chat.Type)
	if !lo.Contains([]telegram.ChatType{telegram.ChatTypeGroup, telegram.ChatTypeSuperGroup}, chatType) {
		return nil, tgbot.NewMessageError("只有在群组和超级群组内猜可以配置聊天记录回顾功能哦！").WithReply(c.Update.Message)
	}

	is, err := c.IsUserMemberStatus(c.Update.Message.From.ID, []telegram.MemberStatus{telegram.MemberStatusCreator, telegram.MemberStatusAdministrator})
	if err != nil {
		return nil, tgbot.
			NewExceptionError(err).
			WithMessage("暂时无法配置聊天记录回顾功能，请稍后再试！").
			WithReply(c.Update.Message)
	}
	if !is {
		return nil, tgbot.
			NewMessageError(fmt.Errorf("%w，%s", errOperationCanNotBeDone, "开启/关闭聊天记录回顾功能需要<b>管理员</b>权限执行 /configure_recap 命令。").Error()).
			WithReply(c.Update.Message).
			WithParseModeHTML()
	}

	chatID := c.Update.Message.Chat.ID

	has, err := h.tgchats.HasChatHistoriesRecapEnabled(chatID, c.Update.Message.Chat.Title)
	if err != nil {
		return nil, tgbot.NewExceptionError(err).WithMessage("暂时无法配置聊天记录回顾功能，请稍后再试！").WithReply(c.Update.Message)
	}

	options, err := h.tgchats.FindOneRecapsOption(chatID)
	if err != nil {
		return nil, tgbot.NewExceptionError(err).WithMessage("暂时无法配置聊天记录回顾功能，请稍后再试！").WithReply(c.Update.Message)
	}
	if options == nil {
		options = &ent.TelegramChatRecapsOptions{AutoRecapSendMode: int(tgchat.AutoRecapSendModePublicly)}
	}

	markup, err := newRecapInlineKeyboardMarkup(c, chatID, c.Update.Message.From.ID, has, tgchat.AutoRecapSendMode(options.AutoRecapSendMode))
	if err != nil {
		return nil, tgbot.NewExceptionError(err).WithMessage("暂时无法配置聊天记录回顾功能，请稍后再试！").WithReply(c.Update.Message)
	}

	return c.
		NewMessageReplyTo(configureRecapGeneralInstructionMessage, c.Update.Message.MessageID).
		WithReplyMarkup(markup), nil
}
