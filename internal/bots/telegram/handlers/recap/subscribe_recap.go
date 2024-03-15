package recap

import (
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/nekomeowww/insights-bot/pkg/bots/tgbot"
	"github.com/nekomeowww/insights-bot/pkg/types/telegram"
	"github.com/nekomeowww/xo"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

func (h *CommandHandler) handleSubscribeRecapCommand(c *tgbot.Context) (tgbot.Response, error) {
	fromID := c.Update.Message.From.ID
	chatID := c.Update.Message.Chat.ID

	chatType := telegram.ChatType(c.Update.Message.Chat.Type)
	if !lo.Contains([]telegram.ChatType{telegram.ChatTypeGroup, telegram.ChatTypeSuperGroup}, chatType) {
		return nil, tgbot.NewMessageError("只有在群组和超级群组内才可以订阅定时的聊天记录回顾哦！").WithReply(c.Update.Message)
	}
	if c.Bot.IsGroupAnonymousBot(c.Update.Message.From) {
		return nil, tgbot.
			NewMessageError("匿名管理员无法订阅定时的聊天记录回顾哦！如果需要订阅定时的聊天记录回顾，必须先将发送角色切换为普通用户然后再试哦。").
			WithReply(c.Update.Message).
			WithDeleteLater(fromID, chatID)
	}

	chatTitle := c.Update.Message.Chat.Title

	has, err := h.tgchats.HasChatHistoriesRecapEnabledForGroups(chatID, chatTitle)
	if err != nil {
		return nil, tgbot.
			NewExceptionError(err).
			WithMessage("订阅群组定时聊天回顾时出现问题，请稍后再试！").
			WithReply(c.Update.Message).
			WithDeleteLater(fromID, chatID)
	}
	if !has {
		return nil, tgbot.
			NewMessageError("聊天记录回顾功能在当前群组尚未启用，需要在群组管理员通过 /configure_recap 命令配置功能启用后才可以订阅聊天回顾哦。").
			WithReply(c.Update.Message).
			WithDeleteLater(fromID, chatID)
	}

	msg := tgbotapi.NewMessage(fromID, fmt.Sprintf("您已成功订阅群组 <b>%s</b> 的定时聊天回顾！", tgbot.EscapeHTMLSymbols(c.Update.Message.Chat.Title)))
	msg.ParseMode = tgbotapi.ModeHTML

	_, err = c.Bot.Send(msg)
	if err == nil {
		err = h.tgchats.SubscribeToAutoRecaps(chatID, fromID)
		if err != nil {
			return nil, tgbot.
				NewExceptionError(err).
				WithMessage("订阅群组定时聊天回顾时出现问题，请稍后再试！").
				WithReply(c.Update.Message).
				WithDeleteLater(fromID, chatID)
		}

		c.Bot.MayRequest(tgbotapi.NewDeleteMessage(chatID, c.Update.Message.MessageID))

		err = c.Bot.DeleteAllDeleteLaterMessages(fromID)
		if err != nil {
			h.logger.Error("failed to delete all delete later messages", zap.Error(err))
		}

		return nil, nil
	}

	hashKey, hashKeyErr := h.setSubscribeStartCommandContext(chatID, chatTitle)
	if hashKeyErr != nil {
		return nil, tgbot.
			NewExceptionError(hashKeyErr).
			WithMessage("订阅群组定时聊天回顾时出现问题，请稍后再试！").
			WithReply(c.Update.Message).
			WithDeleteLater(fromID, chatID)
	}

	if c.Bot.IsCannotInitiateChatWithUserErr(err) {
		return h.handleUserNeverStartedChatOrBlockedErr(c, chatID, chatTitle, newSubscribeRecapCommandWhenUserNeverStartedChat(c.Bot, hashKey))
	} else if c.Bot.IsBotWasBlockedByTheUserErr(err) {
		return h.handleUserNeverStartedChatOrBlockedErr(c, chatID, chatTitle, newSubscribeRecapCommandWhenUserBlockedMessage(c.Bot, hashKey))
	} else {
		h.logger.Error("failed to send private message to user",
			zap.String("message", xo.SprintJSON(msg)),
			zap.Int64("chat_id", c.Update.Message.From.ID),
			zap.Error(err),
		)
	}

	return nil, nil
}

func (h *CommandHandler) handleStartCommandWithRecapSubscription(c *tgbot.Context) (tgbot.Response, error) {
	if c.Bot.IsGroupAnonymousBot(c.Update.Message.From) {
		return nil, nil
	}

	args := strings.Split(c.Update.Message.CommandArguments(), " ")
	if len(args) != 1 {
		return nil, nil
	}

	context, err := h.getSubscribeStartCommandContext(args[0])
	if err != nil {
		h.logger.Error("failed to get private subscription recap start command context", zap.Error(err))
		return nil, nil
	}
	if context == nil {
		return nil, nil
	}

	err = h.tgchats.SubscribeToAutoRecaps(context.ChatID, c.Update.Message.From.ID)
	if err != nil {
		return nil, tgbot.
			NewExceptionError(err).
			WithMessage("订阅群组定时聊天回顾时出现问题，请稍后再试！").
			WithReply(c.Update.Message)
	}

	err = c.Bot.DeleteAllDeleteLaterMessages(c.Update.Message.From.ID)
	if err != nil {
		h.logger.Error("failed to delete all delete later messages", zap.Error(err))
	}

	return c.
		NewMessage(fmt.Sprintf("您已成功订阅群组 <b>%s</b> 的定时聊天回顾！", tgbot.EscapeHTMLSymbols(context.ChatTitle))).
		WithParseModeHTML(), nil
}

func (h *CommandHandler) handleUnsubscribeRecapCommand(c *tgbot.Context) (tgbot.Response, error) {
	chatType := telegram.ChatType(c.Update.Message.Chat.Type)
	if !lo.Contains([]telegram.ChatType{telegram.ChatTypeGroup, telegram.ChatTypeSuperGroup}, chatType) {
		return nil, tgbot.NewMessageError("只有在群组和超级群组内才可以取消订阅定时的聊天记录回顾哦！").WithReply(c.Update.Message)
	}

	chatID := c.Update.Message.Chat.ID

	if c.Bot.IsGroupAnonymousBot(c.Update.Message.From) {
		c.Bot.MayRequest(tgbotapi.NewDeleteMessage(chatID, c.Update.Message.MessageID))
		return nil, nil
	}

	fromID := c.Update.Message.From.ID
	chatTitle := c.Update.Message.Chat.Title

	err := h.tgchats.UnsubscribeToAutoRecaps(chatID, fromID)
	if err != nil {
		return nil, tgbot.
			NewExceptionError(err).
			WithMessage("订阅群组定时聊天回顾时出现问题，请稍后再试！").
			WithReply(c.Update.Message).
			WithDeleteLater(fromID, chatID)
	}

	c.Bot.MayRequest(tgbotapi.NewDeleteMessage(chatID, c.Update.Message.MessageID))

	msg := tgbotapi.NewMessage(fromID, fmt.Sprintf("您已成功取消订阅群组 <b>%s</b> 的定时聊天回顾！", tgbot.EscapeHTMLSymbols(chatTitle)))
	msg.ParseMode = tgbotapi.ModeHTML

	_, err = c.Bot.Send(msg)
	if err != nil {
		if c.Bot.IsCannotInitiateChatWithUserErr(err) || c.Bot.IsBotWasBlockedByTheUserErr(err) {
			c.Bot.MayRequest(tgbotapi.NewDeleteMessage(chatID, c.Update.Message.MessageID))
		} else {
			h.logger.Error("failed to send private message to user",
				zap.String("message", xo.SprintJSON(msg)),
				zap.Int64("chat_id", c.Update.Message.From.ID),
				zap.Error(err),
			)
		}
	}

	return nil, nil
}
