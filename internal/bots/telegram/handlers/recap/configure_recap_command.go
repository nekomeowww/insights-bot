package recap

import (
	"errors"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/samber/lo"

	"github.com/nekomeowww/insights-bot/ent"
	"github.com/nekomeowww/insights-bot/pkg/bots/tgbot"
	"github.com/nekomeowww/insights-bot/pkg/types/bot/handlers/recap"
	"github.com/nekomeowww/insights-bot/pkg/types/telegram"
	"github.com/nekomeowww/insights-bot/pkg/types/tgchat"
)

var (
	errOperationCanNotBeDone                                   = errors.New("抱歉，此操作无法进行")
	errAdministratorPermissionRequired                         = errors.New("<b>管理员</b>")
	errCreatorPermissionRequired                               = errors.New("<b>群组创建者</b>")
	errToggleRecapPermissionDeniedDueToAdministratorIsRequired = fmt.Errorf("%w，只有%w或%w角色可以开启/关闭聊天记录回顾功能。", errOperationCanNotBeDone, errAdministratorPermissionRequired, errCreatorPermissionRequired)
	errAssignModePermissionDeniedDueToAdministratorIsRequired  = fmt.Errorf("%w，只有%w角色可以配置聊天记录回顾的模式。", errOperationCanNotBeDone, errCreatorPermissionRequired)
)

func checkBotIsAdmin(ctx *tgbot.Context) error {
	is, err := ctx.IsBotAdministrator()
	if err != nil {
		return err
	}
	if !is {
		return fmt.Errorf("%w，%s", errOperationCanNotBeDone, "现在机器人不是<b>群组管理员</b>，已经不会记录任何聊天记录了。如果需要配置聊天记录回顾功能，<b>请先将机器人设为群组管理员</b>，然后再次执行命令后再试")
	}

	return nil
}

func checkToggle(ctx *tgbot.Context, _ int64, user *tgbotapi.User) error {
	err := checkBotIsAdmin(ctx)
	if err != nil {
		return err
	}

	if !lo.Contains([]telegram.ChatType{telegram.ChatTypeGroup, telegram.ChatTypeSuperGroup}, telegram.ChatType(ctx.Update.FromChat().Type)) {
		return fmt.Errorf("%w，%s", errOperationCanNotBeDone, "聊天记录回顾功能只有<b>群组</b>和<b>超级群组</b>的管理员可以配置哦！\n请将 Bot 添加到群组中，并配置 Bot 为管理员后使用管理员权限的用户账户为 Bot 进行配置吧。")
	}
	if user == nil {
		return fmt.Errorf("%s，只有%w角色可以进行此操作", errOperationCanNotBeDone, errAdministratorPermissionRequired)
	}

	is, err := ctx.IsUserMemberStatus(user.ID, []telegram.MemberStatus{
		telegram.MemberStatusCreator,
		telegram.MemberStatusAdministrator,
	})
	if err != nil {
		return err
	}
	if !is && !ctx.Bot.IsGroupAnonymousBot(user) {
		return fmt.Errorf("%s，%w", errOperationCanNotBeDone, errToggleRecapPermissionDeniedDueToAdministratorIsRequired)
	}

	return nil
}

func checkAssignMode(ctx *tgbot.Context, _ int64, user *tgbotapi.User) error {
	err := checkBotIsAdmin(ctx)
	if err != nil {
		return err
	}

	if !lo.Contains([]telegram.ChatType{telegram.ChatTypeGroup, telegram.ChatTypeSuperGroup}, telegram.ChatType(ctx.Update.FromChat().Type)) {
		return fmt.Errorf("%w，%s", errOperationCanNotBeDone, "聊天记录回顾功能只有<b>群组</b>和<b>超级群组</b>的管理员可以配置哦！\n请将 Bot 添加到群组中，并配置 Bot 为管理员后使用管理员权限的用户账户为 Bot 进行配置吧。")
	}
	if user == nil {
		return fmt.Errorf("%s，只有%w角色可以进行此操作", errOperationCanNotBeDone, errAdministratorPermissionRequired)
	}

	is, err := ctx.IsUserMemberStatus(user.ID, []telegram.MemberStatus{telegram.MemberStatusCreator})
	if err != nil {
		return err
	}
	if !is {
		isAdmin, err := ctx.IsUserMemberStatus(user.ID, []telegram.MemberStatus{telegram.MemberStatusAdministrator})
		if err != nil {
			return err
		}
		if !isAdmin && !ctx.Bot.IsGroupAnonymousBot(user) {
			return fmt.Errorf("%s，只有%w角色可以进行此操作", errOperationCanNotBeDone, errAdministratorPermissionRequired)
		}

		return fmt.Errorf("%w，%w", errOperationCanNotBeDone, errAssignModePermissionDeniedDueToAdministratorIsRequired)
	}

	return nil
}

func newRecapInlineKeyboardMarkup(
	c *tgbot.Context,
	chatID int64,
	fromID int64,
	currentRecapStatusOn bool,
	currentRecapMode tgchat.AutoRecapSendMode,
	currentAutoRecapRatesPerDay int,
	currentPinStatusOn bool,
) (tgbotapi.InlineKeyboardMarkup, error) {
	nopData, err := c.Bot.AssignOneNopCallbackQueryData()
	if err != nil {
		return tgbotapi.InlineKeyboardMarkup{}, err
	}

	toggleOnData, err := c.Bot.AssignOneCallbackQueryData("recap/configure/toggle", recap.ConfigureRecapToggleActionData{Status: true, ChatID: chatID, FromID: fromID})
	if err != nil {
		return tgbotapi.InlineKeyboardMarkup{}, err
	}

	toggleOffData, err := c.Bot.AssignOneCallbackQueryData("recap/configure/toggle", recap.ConfigureRecapToggleActionData{Status: false, ChatID: chatID, FromID: fromID})
	if err != nil {
		return tgbotapi.InlineKeyboardMarkup{}, err
	}

	publicData, err := c.Bot.AssignOneCallbackQueryData("recap/configure/assign_mode", recap.ConfigureRecapAssignModeActionData{Mode: tgchat.AutoRecapSendModePublicly, ChatID: chatID, FromID: fromID})
	if err != nil {
		return tgbotapi.InlineKeyboardMarkup{}, err
	}

	privateData, err := c.Bot.AssignOneCallbackQueryData("recap/configure/assign_mode", recap.ConfigureRecapAssignModeActionData{Mode: tgchat.AutoRecapSendModeOnlyPrivateSubscriptions, ChatID: chatID, FromID: fromID})
	if err != nil {
		return tgbotapi.InlineKeyboardMarkup{}, err
	}

	completeData, err := c.Bot.AssignOneCallbackQueryData("recap/configure/complete", recap.ConfigureRecapCompleteActionData{ChatID: chatID, FromID: fromID})
	if err != nil {
		return tgbotapi.InlineKeyboardMarkup{}, err
	}

	togglePinData, err := c.Bot.AssignOneCallbackQueryData("recap/configure/pin", recap.ConfigureRecapPinMessageData{Status: true, ChatID: chatID})
	if err != nil {
		return tgbotapi.InlineKeyboardMarkup{}, err
	}

	toggleUnpinData, err := c.Bot.AssignOneCallbackQueryData("recap/configure/pin", recap.ConfigureRecapPinMessageData{Status: false, ChatID: chatID})
	if err != nil {
		return tgbotapi.InlineKeyboardMarkup{}, err
	}

	if !currentRecapStatusOn {
		return tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("🔈 聊天记录回顾", nopData),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(lo.Ternary(currentRecapStatusOn, "🔘 开启", "开启"), toggleOnData),
				tgbotapi.NewInlineKeyboardButtonData(lo.Ternary(!currentRecapStatusOn, "🔘 关闭", "关闭"), toggleOffData),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("📩 聊天记录回顾投递方式", nopData),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(lo.Ternary(currentRecapMode == tgchat.AutoRecapSendModePublicly, "🔘 "+tgchat.AutoRecapSendModePublicly.String(), tgchat.AutoRecapSendModePublicly.String()), publicData),
				tgbotapi.NewInlineKeyboardButtonData(lo.Ternary(currentRecapMode == tgchat.AutoRecapSendModeOnlyPrivateSubscriptions, "🔘 "+tgchat.AutoRecapSendModeOnlyPrivateSubscriptions.String(), tgchat.AutoRecapSendModeOnlyPrivateSubscriptions.String()), privateData),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("✅ 完成", completeData),
			),
		), nil
	}

	twoTimePerDayData, err := c.Bot.AssignOneCallbackQueryData("recap/configure/auto_recap_rates_per_day", recap.ConfigureAutoRecapRatesPerDayActionData{Rates: 2, ChatID: chatID, FromID: fromID})
	if err != nil {
		return tgbotapi.InlineKeyboardMarkup{}, err
	}

	threeTimePerDayData, err := c.Bot.AssignOneCallbackQueryData("recap/configure/auto_recap_rates_per_day", recap.ConfigureAutoRecapRatesPerDayActionData{Rates: 3, ChatID: chatID, FromID: fromID})
	if err != nil {
		return tgbotapi.InlineKeyboardMarkup{}, err
	}

	fourTimePerDayData, err := c.Bot.AssignOneCallbackQueryData("recap/configure/auto_recap_rates_per_day", recap.ConfigureAutoRecapRatesPerDayActionData{Rates: 4, ChatID: chatID, FromID: fromID})
	if err != nil {
		return tgbotapi.InlineKeyboardMarkup{}, err
	}

	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔈 聊天记录回顾", nopData),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(lo.Ternary(currentRecapStatusOn, "🔘 开启", "开启"), toggleOnData),
			tgbotapi.NewInlineKeyboardButtonData(lo.Ternary(!currentRecapStatusOn, "🔘 关闭", "关闭"), toggleOffData),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📩 聊天记录回顾投递方式", nopData),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(lo.Ternary(currentRecapMode == tgchat.AutoRecapSendModePublicly, "🔘 "+tgchat.AutoRecapSendModePublicly.String(), tgchat.AutoRecapSendModePublicly.String()), publicData),
			tgbotapi.NewInlineKeyboardButtonData(lo.Ternary(currentRecapMode == tgchat.AutoRecapSendModeOnlyPrivateSubscriptions, "🔘 "+tgchat.AutoRecapSendModeOnlyPrivateSubscriptions.String(), tgchat.AutoRecapSendModeOnlyPrivateSubscriptions.String()), privateData),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🛎️ 每天自动创建回顾次数", nopData),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(lo.Ternary(currentAutoRecapRatesPerDay == 2, "🔘 2 次", "2 次"), twoTimePerDayData),
			tgbotapi.NewInlineKeyboardButtonData(lo.Ternary(currentAutoRecapRatesPerDay == 3, "🔘 3 次", "3 次"), threeTimePerDayData),
			tgbotapi.NewInlineKeyboardButtonData(lo.Ternary(currentAutoRecapRatesPerDay == 4, "🔘 4 次", "4 次"), fourTimePerDayData),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🪧 置顶聊天记录回顾", nopData),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(lo.Ternary(currentPinStatusOn, "🔘 开启", "开启"), togglePinData),
			tgbotapi.NewInlineKeyboardButtonData(lo.Ternary(!currentPinStatusOn, "🔘 关闭", "关闭"), toggleUnpinData),
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
		return nil, tgbot.NewMessageError("只有在群组和超级群组内才可以配置聊天记录回顾功能哦！").WithReply(c.Update.Message)
	}

	err := checkBotIsAdmin(c)
	if err != nil {
		if errors.Is(err, errOperationCanNotBeDone) || errors.Is(err, errCreatorPermissionRequired) {
			return nil, tgbot.
				NewMessageError(err.Error()).
				WithReply(c.Update.Message).
				WithParseModeHTML()
		}

		return nil, tgbot.
			NewExceptionError(err).
			WithMessage("暂时无法配置聊天记录回顾功能，请稍后再试！").
			WithReply(c.Update.Message)
	}

	is, err := c.IsUserMemberStatus(c.Update.Message.From.ID, []telegram.MemberStatus{
		telegram.MemberStatusCreator,
		telegram.MemberStatusAdministrator,
	})
	if err != nil {
		return nil, tgbot.
			NewExceptionError(err).
			WithMessage("暂时无法配置聊天记录回顾功能，请稍后再试！").
			WithReply(c.Update.Message)
	}
	if !is && !c.Bot.IsGroupAnonymousBot(c.Update.Message.From) {
		return nil, tgbot.
			NewMessageError(fmt.Errorf("%w，%s", errOperationCanNotBeDone, "需要<b>管理员</b>权限才能配置聊天记录回顾功能。").Error()).
			WithReply(c.Update.Message).
			WithParseModeHTML()
	}

	chatID := c.Update.Message.Chat.ID

	has, err := h.tgchats.HasChatHistoriesRecapEnabledForGroups(chatID, c.Update.Message.Chat.Title)
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

	markup, err := newRecapInlineKeyboardMarkup(
		c,
		chatID,
		c.Update.Message.From.ID,
		has,
		tgchat.AutoRecapSendMode(options.AutoRecapSendMode),
		lo.Ternary(options.AutoRecapRatesPerDay == 0, 4, options.AutoRecapRatesPerDay),
		options.PinAutoRecapMessage,
	)
	if err != nil {
		return nil, tgbot.NewExceptionError(err).WithMessage("暂时无法配置聊天记录回顾功能，请稍后再试！").WithReply(c.Update.Message)
	}

	return c.
		NewMessageReplyTo(configureRecapGeneralInstructionMessage, c.Update.Message.MessageID).
		WithReplyMarkup(markup), nil
}
