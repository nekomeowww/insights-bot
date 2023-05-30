package recap

import (
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/nekomeowww/insights-bot/pkg/bots/tgbot"
	"github.com/nekomeowww/insights-bot/pkg/types/bots/handlers/recap"
	"github.com/nekomeowww/insights-bot/pkg/types/telegram"
	"github.com/nekomeowww/insights-bot/pkg/types/tgchat"
	"github.com/nekomeowww/insights-bot/pkg/utils"
	"github.com/samber/lo"
)

func newRecapSelectHoursInlineKeyboardButtons(ctx *tgbot.Context, chatID int64, chatTitle string, recapMode tgchat.AutoRecapSendMode) (tgbotapi.InlineKeyboardMarkup, error) {
	buttons := make([]tgbotapi.InlineKeyboardButton, 0, len(RecapSelectHourAvailable))

	for _, v := range RecapSelectHourAvailable {
		data, err := ctx.Bot.AssignOneCallbackQueryData("recap/recap/select_hours", recap.SelectHourCallbackQueryData{
			Hour:      v,
			ChatID:    chatID,
			ChatTitle: chatTitle,
			RecapMode: recapMode,
		})
		if err != nil {
			return tgbotapi.InlineKeyboardMarkup{}, err
		}

		buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData(
			RecapSelectHourAvailableText[v],
			data,
		))
	}

	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			buttons...,
		),
	), nil
}

func (h *CommandHandler) handleRecapCommand(c *tgbot.Context) (tgbot.Response, error) {
	chatType := telegram.ChatType(c.Update.Message.Chat.Type)
	if !lo.Contains([]telegram.ChatType{telegram.ChatTypeGroup, telegram.ChatTypeSuperGroup}, chatType) {
		return nil, tgbot.NewMessageError("只有在群组和超级群组内才可以创建聊天记录回顾哦！").WithReply(c.Update.Message)
	}

	chatID := c.Update.Message.Chat.ID
	chatTitle := c.Update.Message.Chat.Title

	has, err := h.tgchats.HasChatHistoriesRecapEnabled(chatID, chatTitle)
	if err != nil {
		return nil, tgbot.
			NewExceptionError(err).
			WithMessage("聊天记录回顾生成失败，请稍后再试！").
			WithReply(c.Update.Message)
	}
	if !has {
		return nil, tgbot.
			NewMessageError("聊天记录回顾功能在当前群组尚未启用，需要在群组管理员通过 /configure_recap 命令配置功能启用后才可以创建聊天回顾哦。").
			WithReply(c.Update.Message)
	}

	options, err := h.tgchats.FindOneRecapsOption(chatID)
	if err != nil {
		return nil, tgbot.
			NewExceptionError(err).
			WithMessage("聊天记录回顾生成失败，请稍后再试！").
			WithReply(c.Update.Message)
	}
	if options != nil && tgchat.AutoRecapSendMode(options.AutoRecapSendMode) == tgchat.AutoRecapSendModeOnlyPrivateSubscriptions {
		return h.handleRecapCommandForPrivateSubscriptionsMode(c)
	}

	inlineKeyboardButtons, err := newRecapSelectHoursInlineKeyboardButtons(c, chatID, chatTitle, tgchat.AutoRecapSendModePublicly)
	if err != nil {
		return nil, tgbot.
			NewExceptionError(err).
			WithMessage("聊天记录回顾生成失败，请稍后再试！").
			WithReply(c.Update.Message)
	}

	return c.
		NewMessageReplyTo("请问您要为过去几个小时内的聊天创建回顾呢？", c.Update.Message.MessageID).
		WithReplyMarkup(inlineKeyboardButtons), nil
}

func (h *CommandHandler) handleRecapCommandForPrivateSubscriptionsMode(c *tgbot.Context) (tgbot.Response, error) {
	chatID := c.Update.Message.Chat.ID
	chatTitle := c.Update.Message.Chat.Title
	fromID := c.Update.Message.From.ID
	msg := tgbotapi.NewMessage(fromID, fmt.Sprintf("您正在请求为群组 <b>%s</b> 创建聊天回顾。\n请问您要为过去几个小时内的聊天创建回顾呢？", c.Update.Message.Chat.Title))
	msg.ParseMode = tgbotapi.ModeHTML

	inlineKeyboardButtons, err := newRecapSelectHoursInlineKeyboardButtons(c, chatID, chatTitle, tgchat.AutoRecapSendModeOnlyPrivateSubscriptions)
	if err != nil {
		return nil, tgbot.
			NewExceptionError(err).
			WithMessage("聊天记录回顾生成失败，请稍后再试！").
			WithReply(c.Update.Message)
	}

	msg.ReplyMarkup = &inlineKeyboardButtons

	_, err = c.Bot.Send(msg)
	if err == nil {
		c.Bot.MayRequest(tgbotapi.NewDeleteMessage(chatID, c.Update.Message.MessageID))
		return nil, nil
	}

	hashKey, hashKeyErr := h.setRecapForPrivateSubscriptionModeStartCommandContext(chatID, chatTitle)
	if hashKeyErr != nil {
		return nil, tgbot.
			NewExceptionError(hashKeyErr).
			WithMessage("聊天记录回顾生成失败，请稍后再试！").
			WithReply(c.Update.Message)
	}

	if c.Bot.IsCannotInitiateChatWithUserErr(err) {
		return h.handleUserNeverStartedChatOrBlockedErr(c, chatID, chatTitle, newRecapCommandWhenUserNeverStartedChat(c.Bot, hashKey))
	} else if c.Bot.IsBotWasBlockedByTheUserErr(err) {
		return h.handleUserNeverStartedChatOrBlockedErr(c, chatID, chatTitle, newRecapCommandWhenUserBlockedMessage(c.Bot, hashKey))
	} else {
		h.logger.Errorf("failed to send private message %s to user %d: %v", utils.SprintJSON(msg), c.Update.Message.From.ID, err)
	}

	return nil, nil
}

func (h *CommandHandler) handleStartCommandWithPrivateSubscriptionsRecap(c *tgbot.Context) (tgbot.Response, error) {
	args := strings.Split(c.Update.Message.CommandArguments(), " ")
	if len(args) != 1 {
		return nil, nil
	}

	context, err := h.getRecapForPrivateSubscriptionModeStartCommandContext(args[0])
	if err != nil {
		h.logger.Errorf("failed to get private subscription recap start command context: %s", err)
		return nil, nil
	}
	if context == nil {
		return nil, nil
	}

	inlineKeyboardButtons, err := newRecapSelectHoursInlineKeyboardButtons(c, context.ChatID, context.ChatTitle, tgchat.AutoRecapSendModeOnlyPrivateSubscriptions)
	if err != nil {
		return nil, tgbot.
			NewExceptionError(err).
			WithMessage("聊天记录回顾生成失败，请稍后再试！").
			WithReply(c.Update.Message)
	}

	err = c.Bot.DeleteAllDeleteLaterMessages(c.Update.Message.From.ID)
	if err != nil {
		h.logger.Errorf("failed to delete all delete later messages: %v", err)
	}

	return c.
		NewMessageReplyTo(fmt.Sprintf("您正在请求为群组 <b>%s</b> 创建聊天回顾。\n请问您要为过去几个小时内的聊天创建回顾呢？", context.ChatTitle), c.Update.Message.MessageID).
		WithReplyMarkup(inlineKeyboardButtons).
		WithParseModeHTML(), nil
}
