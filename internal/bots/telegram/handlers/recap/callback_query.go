package recap

import (
	"errors"
	"fmt"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/samber/lo"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/nekomeowww/insights-bot/ent"
	"github.com/nekomeowww/insights-bot/internal/models/chathistories"
	"github.com/nekomeowww/insights-bot/internal/models/tgchats"
	"github.com/nekomeowww/insights-bot/pkg/bots/tgbot"
	"github.com/nekomeowww/insights-bot/pkg/logger"
	"github.com/nekomeowww/insights-bot/pkg/types/bot/handlers/recap"
	"github.com/nekomeowww/insights-bot/pkg/types/telegram"
	"github.com/nekomeowww/insights-bot/pkg/types/tgchat"
)

type NewCallbackQueryHandlerParams struct {
	fx.In

	Logger        *logger.Logger
	ChatHistories *chathistories.Model
	TgChats       *tgchats.Model
}

type CallbackQueryHandler struct {
	logger        *logger.Logger
	chatHistories *chathistories.Model
	tgchats       *tgchats.Model
}

func NewCallbackQueryHandler() func(NewCallbackQueryHandlerParams) *CallbackQueryHandler {
	return func(param NewCallbackQueryHandlerParams) *CallbackQueryHandler {
		return &CallbackQueryHandler{
			logger:        param.Logger,
			chatHistories: param.ChatHistories,
			tgchats:       param.TgChats,
		}
	}
}

func shouldSkipCallbackQueryHandlingByCheckingActionData[
	D recap.ConfigureRecapToggleActionData | recap.ConfigureRecapAssignModeActionData | recap.ConfigureRecapCompleteActionData,
](c *tgbot.Context, actionData D, chatID, fromID int64) bool {
	var actionDataChatID int64
	var actionDataFromID int64

	switch val := any(actionData).(type) {
	case recap.ConfigureRecapToggleActionData:
		actionDataChatID = val.ChatID
		actionDataFromID = val.FromID
	case recap.ConfigureRecapAssignModeActionData:
		actionDataChatID = val.ChatID
		actionDataFromID = val.FromID
	case recap.ConfigureRecapCompleteActionData:
		actionDataChatID = val.ChatID
		actionDataFromID = val.FromID
	}

	// same chat
	if actionDataChatID != chatID {
		c.Logger.Debug("callback query is not from the same chat",
			zap.Int64("chat_id", chatID),
			zap.Int64("action_data_chat_id", actionDataChatID),
		)

		return true
	}
	// same actor or the original command should sent by Group Anonymous Bot
	callbackQueryMessageFromGroupAnonymousBot := c.Update.CallbackQuery.Message.ReplyToMessage != nil && c.Bot.IsGroupAnonymousBot(c.Update.CallbackQuery.Message.ReplyToMessage.From)
	if !(actionDataFromID == fromID || callbackQueryMessageFromGroupAnonymousBot) {
		c.Logger.Debug("action skipped, because callback query is either not from the same actor or the original command should sent by Group Anonymous Bot",
			zap.Int64("from_id", fromID),
			zap.Int64("action_data_from_id", actionDataFromID),
			zap.Bool("has_reply_to_message", c.Update.CallbackQuery.Message.ReplyToMessage != nil),
			zap.Bool("is_group_anonymous_bot", c.Update.CallbackQuery.Message.ReplyToMessage != nil && c.Bot.IsGroupAnonymousBot(c.Update.CallbackQuery.Message.ReplyToMessage.From)),
		)

		return true
	}

	return false
}

func (h *CallbackQueryHandler) handleCallbackQuerySelectHours(c *tgbot.Context) (tgbot.Response, error) {
	messageID := c.Update.CallbackQuery.Message.MessageID

	replyToMessage := c.Update.CallbackQuery.Message.ReplyToMessage

	var data recap.SelectHourCallbackQueryData

	err := c.BindFromCallbackQueryData(&data)
	if err != nil {
		return nil, tgbot.
			NewExceptionError(err).
			WithMessage("聊天记录回顾生成失败，请稍后再试！").
			WithReply(replyToMessage)
	}
	if !lo.Contains(RecapSelectHourAvailable, data.Hour) {
		return nil, tgbot.
			NewExceptionError(fmt.Errorf("invalid hour: %d", data.Hour)).
			WithReply(replyToMessage)
	}

	var inProgressText string

	switch data.RecapMode {
	case tgchat.AutoRecapSendModePublicly:
		inProgressText = fmt.Sprintf("正在为过去 %d 个小时的聊天记录生成回顾，请稍等...", data.Hour)
	case tgchat.AutoRecapSendModeOnlyPrivateSubscriptions:
		inProgressText = fmt.Sprintf("正在为 <b>%s</b> 过去 %d 个小时的聊天记录生成回顾，请稍等...", tgbot.EscapeHTMLSymbols(data.ChatTitle), data.Hour)
	default:
		inProgressText = fmt.Sprintf("正在为过去 %d 个小时的聊天记录生成回顾，请稍等...", data.Hour)
	}

	editConfig := tgbotapi.NewEditMessageTextAndMarkup(
		c.Update.CallbackQuery.Message.Chat.ID,
		messageID,
		inProgressText,
		tgbotapi.NewInlineKeyboardMarkup([]tgbotapi.InlineKeyboardButton{}),
	)

	editConfig.ParseMode = tgbotapi.ModeHTML

	_, err = c.Bot.Request(editConfig)
	if err != nil {
		h.logger.Error("failed to edit message", zap.Error(err))
	}

	histories, err := h.chatHistories.FindChatHistoriesByTimeBefore(data.ChatID, time.Duration(data.Hour)*time.Hour)
	if err != nil {
		return nil, tgbot.
			NewExceptionError(err).WithMessage("聊天记录回顾生成失败，请稍后再试！").
			WithReply(replyToMessage)
	}
	if len(histories) <= 5 {
		var errMessage string

		switch data.RecapMode {
		case tgchat.AutoRecapSendModePublicly:
			errMessage = fmt.Sprintf("最近 %d 小时内暂时没有超过 5 条的聊天记录可以生成聊天回顾哦，要再多聊点之后再试试吗？", data.Hour)
		case tgchat.AutoRecapSendModeOnlyPrivateSubscriptions:
			errMessage = fmt.Sprintf("最近 %d 小时内暂时没有超过 5 条的聊天记录可以生成聊天回顾哦，要再等待群内成员多聊点之后再试试吗？", data.Hour)
		default:
			errMessage = fmt.Sprintf("最近 %d 小时内暂时没有超过 5 条的聊天记录可以生成聊天回顾哦，要再多聊点之后再试试吗？", data.Hour)
		}

		return nil, tgbot.
			NewMessageError(errMessage).
			WithReply(replyToMessage)
	}

	summarizations, err := h.chatHistories.SummarizeChatHistories(data.ChatID, histories)
	if err != nil {
		return nil, tgbot.
			NewExceptionError(err).WithMessage("聊天记录回顾生成失败，请稍后再试！").
			WithReply(replyToMessage)
	}

	summarizations = lo.Filter(summarizations, func(item string, _ int) bool { return item != "" })
	if len(summarizations) == 0 {
		return nil, tgbot.
			NewMessageError("聊天记录回顾生成失败，请稍后再试！").
			WithReply(replyToMessage)
	}

	for i, s := range summarizations {
		summarizations[i] = tgbot.ReplaceMarkdownTitlesToTelegramBoldElement(s)
	}

	summarizationBatches := tgbot.SplitMessagesAgainstLengthLimitIntoMessageGroups(summarizations)
	for i, s := range summarizationBatches {
		var content string
		if len(summarizationBatches) > 1 {
			content = fmt.Sprintf("%s\n\n(%d/%d)\n#recap\n<em>🤖️ Generated by chatGPT</em>", strings.Join(s, "\n\n"), i+1, len(summarizationBatches))
		} else {
			content = fmt.Sprintf("%s\n\n#recap\n<em>🤖️ Generated by chatGPT</em>", strings.Join(s, "\n\n"))
		}

		msg := tgbotapi.NewMessage(c.Update.CallbackQuery.Message.Chat.ID, content)
		msg.ParseMode = tgbotapi.ModeHTML

		if c.Update.CallbackQuery.Message.ReplyToMessage != nil {
			msg.ReplyToMessageID = c.Update.CallbackQuery.Message.ReplyToMessage.MessageID
		}

		h.logger.Info("sending chat histories recap for chat",
			zap.Int64("chat_id", c.Update.CallbackQuery.Message.Chat.ID),
			zap.String("text", msg.Text),
		)

		c.Bot.MaySend(msg)
	}

	return nil, nil
}

func (h *CallbackQueryHandler) handleCallbackQueryToggle(c *tgbot.Context) (tgbot.Response, error) {
	msg := c.Update.CallbackQuery.Message

	generalErrorMessage := configureRecapGeneralInstructionMessage + "\n\n" + "应用聊天记录回顾功能的配置时出现了问题，请稍后再试！"

	fromID := c.Update.CallbackQuery.From.ID
	chatID := msg.Chat.ID
	chatTitle := msg.Chat.Title
	chatType := msg.Chat.Type
	messageID := msg.MessageID

	var actionData recap.ConfigureRecapToggleActionData

	err := c.BindFromCallbackQueryData(&actionData)
	if err != nil {
		return nil, tgbot.
			NewExceptionError(err).
			WithMessage(generalErrorMessage).
			WithEdit(msg).
			WithReplyMarkup(tgbotapi.NewInlineKeyboardMarkup(msg.ReplyMarkup.InlineKeyboard...))
	}

	shouldSkip := shouldSkipCallbackQueryHandlingByCheckingActionData(c, actionData, chatID, fromID)
	if shouldSkip {
		return nil, nil
	}

	// check whether the actor is admin or creator, and whether the bot is admin
	err = checkToggle(c, chatID, c.Update.CallbackQuery.From)
	if err != nil {
		if errors.Is(err, errAdministratorPermissionRequired) {
			h.logger.Debug("action, skipped, callback query is not from an admin or creator",
				zap.Int64("from_id", fromID),
				zap.Int64("chat_id", chatID),
				zap.String("permission_check_result", err.Error()),
			)

			return nil, nil
		}
		if errors.Is(err, errOperationCanNotBeDone) {
			return nil, tgbot.
				NewMessageError(configureRecapGeneralInstructionMessage + "\n\n" + err.Error()).
				WithEdit(msg).
				WithParseModeHTML().
				WithReplyMarkup(tgbotapi.NewInlineKeyboardMarkup(msg.ReplyMarkup.InlineKeyboard...))
		}

		return nil, tgbot.
			NewExceptionError(err).
			WithMessage(generalErrorMessage).
			WithEdit(msg).
			WithReplyMarkup(tgbotapi.NewInlineKeyboardMarkup(msg.ReplyMarkup.InlineKeyboard...))
	}

	options, err := h.tgchats.FindOneRecapsOption(chatID)
	if err != nil {
		return nil, tgbot.
			NewExceptionError(err).
			WithMessage("暂时无法配置聊天记录回顾功能，请稍后再试！").
			WithEdit(c.Update.Message).
			WithReplyMarkup(tgbotapi.NewInlineKeyboardMarkup(msg.ReplyMarkup.InlineKeyboard...))
	}
	if options == nil {
		options = &ent.TelegramChatRecapsOptions{AutoRecapSendMode: int(tgchat.AutoRecapSendModePublicly)}
	}

	if actionData.Status {
		errMessage := configureRecapGeneralInstructionMessage + "\n\n" + "聊天记录回顾功能开启失败，请稍后再试！"

		err = h.tgchats.EnableChatHistoriesRecap(chatID, telegram.ChatType(chatType), chatTitle)
		if err != nil {
			return nil, tgbot.
				NewExceptionError(err).
				WithMessage(errMessage).
				WithEdit(msg).
				WithReplyMarkup(tgbotapi.NewInlineKeyboardMarkup(msg.ReplyMarkup.InlineKeyboard...))
		}

		err = h.tgchats.QueueOneSendChatHistoriesRecapTaskForChatID(chatID)
		if err != nil {
			return nil, tgbot.
				NewExceptionError(err).
				WithMessage(errMessage).
				WithEdit(msg).
				WithReplyMarkup(tgbotapi.NewInlineKeyboardMarkup(msg.ReplyMarkup.InlineKeyboard...))
		}
	} else {
		errMessage := configureRecapGeneralInstructionMessage + "\n\n" + "聊天记录回顾功能关闭失败，请稍后再试！"

		err = h.tgchats.DisableChatHistoriesRecap(chatID, telegram.ChatType(chatType), chatTitle)
		if err != nil {
			return nil, tgbot.
				NewExceptionError(err).
				WithMessage(errMessage).
				WithEdit(msg).
				WithReplyMarkup(tgbotapi.NewInlineKeyboardMarkup(msg.ReplyMarkup.InlineKeyboard...))
		}
	}

	markup, err := newRecapInlineKeyboardMarkup(c, chatID, fromID, actionData.Status, tgchat.AutoRecapSendMode(options.AutoRecapSendMode))
	if err != nil {
		return nil, tgbot.
			NewExceptionError(err).
			WithMessage("暂时无法配置聊天记录回顾功能，请稍后再试！").
			WithEdit(c.Update.Message).
			WithReplyMarkup(tgbotapi.NewInlineKeyboardMarkup(msg.ReplyMarkup.InlineKeyboard...))
	}

	return c.NewEditMessageTextAndReplyMarkup(messageID,
		lo.Ternary(
			actionData.Status,
			configureRecapGeneralInstructionMessage+"\n\n"+"聊天记录回顾功能已开启，开启后将会自动收集群组中的聊天记录并定时发送聊天回顾快报。",
			configureRecapGeneralInstructionMessage+"\n\n"+"聊天记录回顾功能已关闭，关闭后将不会再收集群组中的聊天记录了。",
		),
		markup,
	), nil
}

func (h *CallbackQueryHandler) handleCallbackQueryAssignMode(c *tgbot.Context) (tgbot.Response, error) {
	msg := c.Update.CallbackQuery.Message

	generalErrorMessage := configureRecapGeneralInstructionMessage + "\n\n" + "应用聊天记录回顾功能的配置时出现了问题，请稍后再试！"

	fromID := c.Update.CallbackQuery.From.ID
	chatID := msg.Chat.ID
	chatTitle := msg.Chat.Title
	messageID := msg.MessageID

	var actionData recap.ConfigureRecapAssignModeActionData

	err := c.BindFromCallbackQueryData(&actionData)
	if err != nil {
		return nil, tgbot.
			NewExceptionError(err).
			WithMessage(generalErrorMessage).
			WithEdit(msg).
			WithReplyMarkup(tgbotapi.NewInlineKeyboardMarkup(msg.ReplyMarkup.InlineKeyboard...))
	}

	shouldSkip := shouldSkipCallbackQueryHandlingByCheckingActionData(c, actionData, chatID, fromID)
	if shouldSkip {
		return nil, nil
	}

	// check whether the actor is admin or creator, and whether the bot is admin
	err = checkAssignMode(c, chatID, c.Update.CallbackQuery.From)
	if err != nil {
		if errors.Is(err, errAdministratorPermissionRequired) {
			h.logger.Debug("action skipped, callback query is not from an admin or creator",
				zap.Int64("from_id", fromID),
				zap.Int64("chat_id", chatID),
				zap.String("permission_check_result", err.Error()),
			)

			return nil, nil
		}
		if errors.Is(err, errOperationCanNotBeDone) || errors.Is(err, errCreatorPermissionRequired) {
			return nil, tgbot.
				NewMessageError(configureRecapGeneralInstructionMessage + "\n\n" + err.Error()).
				WithEdit(msg).
				WithParseModeHTML().
				WithReplyMarkup(tgbotapi.NewInlineKeyboardMarkup(msg.ReplyMarkup.InlineKeyboard...))
		}

		return nil, tgbot.
			NewExceptionError(err).
			WithMessage(generalErrorMessage).
			WithEdit(msg).
			WithReplyMarkup(tgbotapi.NewInlineKeyboardMarkup(msg.ReplyMarkup.InlineKeyboard...))
	}

	err = h.tgchats.SetRecapsRecapMode(chatID, actionData.Mode)
	if err != nil {
		return nil, tgbot.
			NewExceptionError(err).
			WithMessage(generalErrorMessage).
			WithEdit(msg).
			WithReplyMarkup(tgbotapi.NewInlineKeyboardMarkup(msg.ReplyMarkup.InlineKeyboard...))
	}

	h.logger.Info("assigned recap mode for chat", zap.String("recap_mode", actionData.Mode.String()))

	has, err := h.tgchats.HasChatHistoriesRecapEnabled(chatID, chatTitle)
	if err != nil {
		return nil, tgbot.
			NewExceptionError(err).
			WithMessage(configureRecapGeneralInstructionMessage + "\n\n" + "聊天记录回顾模式设定失败，请稍后再试！").
			WithEdit(msg).
			WithReplyMarkup(tgbotapi.NewInlineKeyboardMarkup(msg.ReplyMarkup.InlineKeyboard...))
	}

	markup, err := newRecapInlineKeyboardMarkup(c, chatID, fromID, has, actionData.Mode)
	if err != nil {
		return nil, tgbot.
			NewExceptionError(err).
			WithMessage("暂时无法配置聊天记录回顾功能，请稍后再试！").
			WithEdit(msg).
			WithReplyMarkup(tgbotapi.NewInlineKeyboardMarkup(msg.ReplyMarkup.InlineKeyboard...))
	}

	return c.NewEditMessageTextAndReplyMarkup(messageID,
		lo.Ternary(
			actionData.Mode == tgchat.AutoRecapSendModePublicly,
			configureRecapGeneralInstructionMessage+"\n\n"+"聊天记录回顾模式已切换为<b>"+tgchat.AutoRecapSendModePublicly.String()+"</b>，将会自动收集群组中的聊天记录并定时发送聊天回顾快报。",
			configureRecapGeneralInstructionMessage+"\n\n"+"聊天记录回顾模式已切换为<b>"+tgchat.AutoRecapSendModeOnlyPrivateSubscriptions.String()+"</b>，将会自动收集群组中的聊天记录并定时发送聊天回顾快报给通过 /subscribe_recap 命令订阅了本群组聊天回顾用户。",
		),
		markup,
	).WithParseModeHTML(), nil
}

func (h *CallbackQueryHandler) handleCallbackQueryComplete(c *tgbot.Context) (tgbot.Response, error) {
	msg := c.Update.CallbackQuery.Message

	generalErrorMessage := configureRecapGeneralInstructionMessage + "\n\n" + "应用聊天记录回顾功能的配置时出现了问题，请稍后再试！"

	fromID := c.Update.CallbackQuery.From.ID
	chatID := msg.Chat.ID
	messageID := msg.MessageID

	var actionData recap.ConfigureRecapCompleteActionData

	err := c.BindFromCallbackQueryData(&actionData)
	if err != nil {
		return nil, tgbot.
			NewExceptionError(err).
			WithMessage(generalErrorMessage).
			WithEdit(msg).
			WithReplyMarkup(tgbotapi.NewInlineKeyboardMarkup(msg.ReplyMarkup.InlineKeyboard...))
	}

	shouldSkip := shouldSkipCallbackQueryHandlingByCheckingActionData(c, actionData, chatID, fromID)
	if shouldSkip {
		return nil, nil
	}

	// check actor is admin or creator, bot is admin
	is, err := c.IsUserMemberStatus(fromID, []telegram.MemberStatus{telegram.MemberStatusCreator, telegram.MemberStatusAdministrator})
	if err != nil {
		return nil, tgbot.
			NewExceptionError(err).
			WithMessage("暂时无法配置聊天记录回顾功能，请稍后再试！").
			WithEdit(msg).
			WithReplyMarkup(tgbotapi.NewInlineKeyboardMarkup(msg.ReplyMarkup.InlineKeyboard...))
	}
	if !is && !c.Bot.IsGroupAnonymousBot(c.Update.CallbackQuery.From) {
		return nil, nil
	}

	_ = c.Bot.MayRequest(tgbotapi.NewDeleteMessage(chatID, messageID))
	if c.Update.CallbackQuery.Message.ReplyToMessage != nil {
		_ = c.Bot.MayRequest(tgbotapi.NewDeleteMessage(chatID, c.Update.CallbackQuery.Message.ReplyToMessage.MessageID))
	}

	return nil, nil
}

func (h *CallbackQueryHandler) handleCallbackQueryUnsubscribe(c *tgbot.Context) (tgbot.Response, error) {
	msg := c.Update.CallbackQuery.Message

	fromID := c.Update.CallbackQuery.From.ID
	chatID := msg.Chat.ID

	var actionData recap.UnsubscribeRecapActionData

	err := c.BindFromCallbackQueryData(&actionData)
	if err != nil {
		return nil, tgbot.
			NewExceptionError(err).
			WithMessage("取消订阅时出现了问题，请稍后再试！").
			WithEdit(msg).
			WithReplyMarkup(tgbotapi.NewInlineKeyboardMarkup(msg.ReplyMarkup.InlineKeyboard...))
	}
	if actionData.ChatID != chatID || actionData.FromID != fromID {
		return nil, nil
	}

	err = h.tgchats.UnsubscribeToAutoRecaps(chatID, fromID)
	if err != nil {
		return nil, tgbot.
			NewExceptionError(err).
			WithMessage("取消订阅时出现了问题，请稍后再试！").
			WithEdit(msg).
			WithReplyMarkup(tgbotapi.NewInlineKeyboardMarkup(msg.ReplyMarkup.InlineKeyboard...))
	}

	c.Bot.MayRequest(tgbotapi.NewEditMessageReplyMarkup(chatID, msg.MessageID, tgbotapi.InlineKeyboardMarkup{
		InlineKeyboard: make([][]tgbotapi.InlineKeyboardButton, 0),
	}))

	return c.NewMessage(fmt.Sprintf("已成功取消订阅群组 <b>%s</b> 的定时聊天回顾。", tgbot.EscapeHTMLSymbols(actionData.ChatTitle))).WithParseModeHTML(), nil
}
