package chat_with_chat_history

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/nekomeowww/insights-bot/pkg/handler"
	"github.com/nekomeowww/insights-bot/pkg/types/telegram"
	"github.com/samber/lo"
)

func (h *Handler) checkTogglingRecapPermission(bot *handler.Bot, chatID, userID int64) (bool, error) {
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

func (h *Handler) HandleEnableRecapCommand(c *handler.Context) {
	chatType := telegram.ChatType(c.Update.Message.Chat.Type)
	if !lo.Contains([]telegram.ChatType{telegram.ChatTypeGroup, telegram.ChatTypeSuperGroup}, chatType) {
		message := tgbotapi.NewMessage(c.Update.Message.Chat.ID, "聊天记录回顾功能只有群组和超级群组可以配置开启哦！")
		message.ReplyToMessageID = c.Update.Message.MessageID
		c.Bot.MustSend(message)
		return
	}

	hasTogglingRecapPermission, err := h.checkTogglingRecapPermission(c.Bot, c.Update.Message.Chat.ID, c.Update.Message.From.ID)
	if err != nil {
		h.Logger.Errorf("failed to check toggling recap permission: %v", err)

		message := tgbotapi.NewMessage(c.Update.Message.Chat.ID, "聊天记录回顾功能开启失败，请稍后再试！")
		message.ReplyToMessageID = c.Update.Message.MessageID
		c.Bot.MustSend(message)
		return
	}
	if !hasTogglingRecapPermission {
		message := tgbotapi.NewMessage(c.Update.Message.Chat.ID, "你没有权限开启聊天记录回顾功能哦！")
		message.ReplyToMessageID = c.Update.Message.MessageID
		c.Bot.MustSend(message)
		return
	}

	err = h.TelegramChatFeatureFlags.EnableChatHistoriesRecap(c.Update.Message.Chat.ID, chatType)
	if err != nil {
		h.Logger.Errorf("failed to enable chat histories recap: %v", err)

		message := tgbotapi.NewMessage(c.Update.Message.Chat.ID, "聊天记录回顾功能开启失败，请稍后再试！")
		message.ReplyToMessageID = c.Update.Message.MessageID
		c.Bot.MustSend(message)
		return
	}

	message := tgbotapi.NewMessage(c.Update.Message.Chat.ID, "聊天记录回顾功能已开启，开启后将会自动收集群组中的聊天记录并定时发送聊天回顾快报！")
	message.ReplyToMessageID = c.Update.Message.MessageID
	c.Bot.MustSend(message)
}

func (h *Handler) HandleDisableRecapCommand(c *handler.Context) {
	chatType := telegram.ChatType(c.Update.Message.Chat.Type)
	if !lo.Contains([]telegram.ChatType{telegram.ChatTypeGroup, telegram.ChatTypeSuperGroup}, chatType) {
		message := tgbotapi.NewMessage(c.Update.Message.Chat.ID, "聊天记录回顾功能只有群组和超级群组可以配置关闭哦！")
		message.ReplyToMessageID = c.Update.Message.MessageID
		c.Bot.MustSend(message)
		return
	}

	hasTogglingRecapPermission, err := h.checkTogglingRecapPermission(c.Bot, c.Update.Message.Chat.ID, c.Update.Message.From.ID)
	if err != nil {
		h.Logger.Errorf("failed to check toggling recap permission: %v", err)

		message := tgbotapi.NewMessage(c.Update.Message.Chat.ID, "聊天记录回顾功能开启失败，请稍后再试！")
		message.ReplyToMessageID = c.Update.Message.MessageID
		c.Bot.MustSend(message)
		return
	}
	if !hasTogglingRecapPermission {
		message := tgbotapi.NewMessage(c.Update.Message.Chat.ID, "你没有权限关闭聊天记录回顾功能哦！")
		message.ReplyToMessageID = c.Update.Message.MessageID
		c.Bot.MustSend(message)
		return
	}

	err = h.TelegramChatFeatureFlags.DisableChatHistoriesRecap(c.Update.Message.Chat.ID, chatType)
	if err != nil {
		h.Logger.Errorf("failed to enable chat histories recap: %v", err)

		message := tgbotapi.NewMessage(c.Update.Message.Chat.ID, "聊天记录回顾功能关闭失败，请稍后再试！")
		message.ReplyToMessageID = c.Update.Message.MessageID
		c.Bot.MustSend(message)
		return
	}

	message := tgbotapi.NewMessage(c.Update.Message.Chat.ID, "聊天记录回顾功能已关闭，关闭后将不会自动收集群组中的聊天记录并定时发送聊天回顾快报了。")
	message.ReplyToMessageID = c.Update.Message.MessageID
	c.Bot.MustSend(message)
}
