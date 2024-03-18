package recap

import (
	"fmt"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/nekomeowww/insights-bot/ent"
	"github.com/nekomeowww/insights-bot/pkg/bots/tgbot"
	"github.com/nekomeowww/insights-bot/pkg/types/bot/handlers/recap"
	"github.com/nekomeowww/insights-bot/pkg/types/telegram"
	"github.com/nekomeowww/insights-bot/pkg/types/tgchat"
	"github.com/nekomeowww/xo"
	"github.com/samber/lo"
	"go.uber.org/zap"
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

	has, err := h.tgchats.HasChatHistoriesRecapEnabledForGroups(chatID, chatTitle)
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

	perSeconds := h.tgchats.ManualRecapRatePerSeconds(options)

	_, ttl, ok, err := c.RateLimitForCommand(chatID, "/recap", 1, perSeconds)
	if err != nil {
		h.logger.Error("failed to check rate limit for command /recap", zap.Error(err))
	}
	if !ok {
		return nil, tgbot.
			NewMessageError(fmt.Sprintf("很抱歉，您的操作触发了我们的限制机制，为了保证系统的可用性，本命令每最多 %d 分钟最多使用一次，请您耐心等待 %d 分钟后再试，感谢您的理解和支持。", perSeconds, lo.Ternary(ttl/time.Minute <= 1, 1, ttl/time.Minute))).
			WithReply(c.Update.Message)
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
	fromID := c.Update.Message.From.ID

	if c.Bot.IsGroupAnonymousBot(c.Update.Message.From) {
		return nil, tgbot.
			NewMessageError("匿名管理员无法在设定为私聊回顾模式的群组内请求创建聊天记录回顾哦！如果需要创建聊天记录回顾，必须先将发送角色切换为普通用户然后再试哦。").
			WithReply(c.Update.Message).
			WithDeleteLater(fromID, chatID)
	}

	chatTitle := c.Update.Message.Chat.Title
	msg := tgbotapi.NewMessage(fromID, fmt.Sprintf("您正在请求为群组 <b>%s</b> 创建聊天回顾。\n请问您要为过去几个小时内的聊天创建回顾呢？", tgbot.EscapeHTMLSymbols(c.Update.Message.Chat.Title)))
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

		err = c.Bot.DeleteAllDeleteLaterMessages(fromID)
		if err != nil {
			h.logger.Error("failed to delete all delete later messages", zap.Error(err))
		}

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
		h.logger.Error("failed to send private message to user",
			zap.String("message", xo.SprintJSON(msg)),
			zap.Int64("chat_id", c.Update.Message.From.ID),
			zap.Error(err),
		)
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
		h.logger.Error("failed to get private subscription recap start command context", zap.Error(err))
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
		h.logger.Error("failed to delete all delete later messages", zap.Error(err))
	}

	return c.
		NewMessageReplyTo(fmt.Sprintf("您正在请求为群组 <b>%s</b> 创建聊天回顾。\n请问您要为过去几个小时内的聊天创建回顾呢？", tgbot.EscapeHTMLSymbols(context.ChatTitle)), c.Update.Message.MessageID).
		WithReplyMarkup(inlineKeyboardButtons).
		WithParseModeHTML(), nil
}

func (h *CommandHandler) handleChatMemberLeft(c *tgbot.Context) (tgbot.Response, error) {
	if c.Update.Message.LeftChatMember == nil {
		return nil, nil
	}

	chatID := c.Update.Message.Chat.ID
	userID := c.Update.Message.LeftChatMember.ID

	var err error
	var subscriber *ent.TelegramChatAutoRecapsSubscribers

	_, _, err = lo.AttemptWithDelay(1000, time.Minute, func(iter int, _ time.Duration) error {
		subscriber, err = h.tgchats.FindOneAutoRecapsSubscriber(chatID, userID)
		if err != nil {
			h.logger.Error("failed to query subscriber of auto recaps",
				zap.Error(err),
				zap.Int64("chat_id", chatID),
				zap.Int64("user_id", userID),
				zap.Int("iter", iter),
				zap.Int("max_iter", 100),
			)

			return err
		}

		return nil
	})
	if err != nil {
		h.logger.Error("failed to query subscriber of auto recaps",
			zap.Error(err),
			zap.Int64("chat_id", chatID),
			zap.Int64("user_id", userID),
		)

		return nil, nil
	}
	if subscriber == nil {
		return nil, nil
	}

	h.logger.Warn("subscriber is no longer a member, auto unsubscribing...",
		zap.Int64("chat_id", chatID),
		zap.Int64("user_id", userID),
	)

	_, _, err = lo.AttemptWithDelay(1000, time.Minute, func(iter int, _ time.Duration) error {
		err := h.tgchats.UnsubscribeToAutoRecaps(chatID, userID)
		if err != nil {
			h.logger.Error("failed to auto unsubscribe to auto recaps",
				zap.Error(err),
				zap.Int64("chat_id", chatID),
				zap.Int64("user_id", userID),
				zap.Int("iter", iter),
				zap.Int("max_iter", 100),
			)

			return err
		}

		return nil
	})
	if err != nil {
		h.logger.Error("failed to unsubscribe to auto recaps", zap.Error(err))
	}

	msg := tgbotapi.NewMessage(subscriber.UserID, fmt.Sprintf("由于您已不再是 <b>%s</b> 的成员，因此已自动帮您取消了您所订阅的聊天记录回顾。", tgbot.EscapeHTMLSymbols(c.Update.Message.Chat.Title)))
	msg.ParseMode = tgbotapi.ModeHTML
	c.Bot.MaySend(msg)

	return nil, nil
}
