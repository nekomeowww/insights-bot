package recap

import (
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/nekomeowww/insights-bot/pkg/bots/tgbot"
	"github.com/nekomeowww/insights-bot/pkg/types/telegram"
	"github.com/nekomeowww/insights-bot/pkg/utils"
	"github.com/samber/lo"
)

func (h *CommandHandler) handleSubscribeRecapCommand(c *tgbot.Context) (tgbot.Response, error) {
	chatType := telegram.ChatType(c.Update.Message.Chat.Type)
	if !lo.Contains([]telegram.ChatType{telegram.ChatTypeGroup, telegram.ChatTypeSuperGroup}, chatType) {
		return nil, tgbot.NewMessageError("只有在群组和超级群组内才可以订阅定时的聊天记录回顾哦！").WithReply(c.Update.Message)
	}

	fromID := c.Update.Message.From.ID
	chatID := c.Update.Message.Chat.ID
	chatTitle := c.Update.Message.Chat.Title

	has, err := h.tgchats.HasChatHistoriesRecapEnabled(chatID, chatTitle)
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

	msg := tgbotapi.NewMessage(fromID, fmt.Sprintf("您已成功订阅群组 <b>%s</b> 的定时聊天回顾！", c.Update.Message.Chat.Title))
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
		h.logger.Errorf("failed to send private message %s to user %d: %v", utils.SprintJSON(msg), c.Update.Message.From.ID, err)
	}

	return nil, nil
}

func (h *CommandHandler) handleStartCommandWithRecapSubscription(c *tgbot.Context) (tgbot.Response, error) {
	args := strings.Split(c.Update.Message.CommandArguments(), " ")
	if len(args) != 1 {
		return nil, nil
	}

	context, err := h.getSubscribeStartCommandContext(args[0])
	if err != nil {
		h.logger.Errorf("failed to get private subscription recap start command context: %s", err)
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
		h.logger.Errorf("failed to delete all delete later messages: %v", err)
	}

	return c.
		NewMessage(fmt.Sprintf("您已成功订阅群组 <b>%s</b> 的定时聊天回顾！", context.ChatTitle)).
		WithParseModeHTML(), nil
}

func (h *CommandHandler) handleUnsubscribeRecapCommand(c *tgbot.Context) (tgbot.Response, error) {
	chatType := telegram.ChatType(c.Update.Message.Chat.Type)
	if !lo.Contains([]telegram.ChatType{telegram.ChatTypeGroup, telegram.ChatTypeSuperGroup}, chatType) {
		return nil, tgbot.NewMessageError("只有在群组和超级群组内才可以取消订阅定时的聊天记录回顾哦！").WithReply(c.Update.Message)
	}

	fromID := c.Update.Message.From.ID
	chatID := c.Update.Message.Chat.ID
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

	msg := tgbotapi.NewMessage(fromID, fmt.Sprintf("您已成功取消订阅群组 <b>%s</b> 的定时聊天回顾！", chatTitle))
	msg.ParseMode = tgbotapi.ModeHTML

	_, err = c.Bot.Send(msg)
	if err != nil {
		if c.Bot.IsCannotInitiateChatWithUserErr(err) || c.Bot.IsBotWasBlockedByTheUserErr(err) {
			c.Bot.MayRequest(tgbotapi.NewDeleteMessage(chatID, c.Update.Message.MessageID))
		} else {
			h.logger.Errorf("failed to send private message %s to user %d: %v", utils.SprintJSON(msg), c.Update.Message.From.ID, err)
		}
	}

	return nil, nil
}
