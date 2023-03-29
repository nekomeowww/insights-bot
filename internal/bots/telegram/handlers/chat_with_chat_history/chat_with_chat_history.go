package chat_with_chat_history

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/nekomeowww/insights-bot/internal/models/chat_histories"
	"github.com/nekomeowww/insights-bot/internal/models/telegram_chat_feature_flags"
	"github.com/nekomeowww/insights-bot/pkg/handler"
	"github.com/nekomeowww/insights-bot/pkg/logger"
	"github.com/nekomeowww/insights-bot/pkg/types/telegram"
	"github.com/samber/lo"
	"go.uber.org/fx"
)

type NewHandlerParam struct {
	fx.In

	Logger *logger.Logger

	ChatHistories            *chat_histories.ChatHistoriesModel
	TelegramChatFeatureFlags *telegram_chat_feature_flags.TelegramChatFeatureFlagsModel
}

type Handler struct {
	Logger                   *logger.Logger
	ChatHistories            *chat_histories.ChatHistoriesModel
	TelegramChatFeatureFlags *telegram_chat_feature_flags.TelegramChatFeatureFlagsModel
}

func NewHandler() func(NewHandlerParam) *Handler {
	return func(param NewHandlerParam) *Handler {
		return &Handler{
			Logger:                   param.Logger,
			ChatHistories:            param.ChatHistories,
			TelegramChatFeatureFlags: param.TelegramChatFeatureFlags,
		}
	}
}

func (h *Handler) HandleRecordMessage(c *handler.Context) {
	chatType := telegram.ChatType(c.Update.Message.Chat.Type)
	if !lo.Contains([]telegram.ChatType{telegram.ChatTypeGroup, telegram.ChatTypeSuperGroup}, chatType) {
		return
	}

	enabled, err := h.TelegramChatFeatureFlags.HasChatHistoriesRecapEnabled(c.Update.Message.Chat.ID, chatType)
	if err != nil {
		h.Logger.Errorf("failed to check if chat history recap is enabled: %v", err)
		return
	}
	if !enabled {
		return
	}

	err = h.ChatHistories.SaveOneTelegramChatHistory(c.Update.Message)
	if err != nil {
		h.Logger.Errorf("failed to save telegram chat history: %v", err)
		return
	}
}

func (h *Handler) HandleRecapCommand(c *handler.Context) {
	chatID := c.Update.Message.Chat.ID
	h.Logger.Infof("generating chat histories recap for chat %d", chatID)

	message := tgbotapi.NewMessage(chatID, "稍等我翻一下聊天记录看看你们都聊了什么哦～")
	message.ReplyToMessageID = c.Update.Message.MessageID
	_, err := c.Bot.Send(message)
	if err != nil {
		h.Logger.Errorf("failed to send chat histories recap: %v", err)
		return
	}

	summarization, err := h.ChatHistories.SummarizeLastOneHourChatHistories(chatID)
	if err != nil {
		h.Logger.Errorf("failed to summarize last one hour chat histories: %v", err)

		errMessage := tgbotapi.NewMessage(chatID, "聊天记录回顾生成失败，请稍后再试！")
		errMessage.ReplyToMessageID = c.Update.Message.MessageID
		_, err = c.Bot.Send(errMessage)
		if err != nil {
			h.Logger.Errorf("failed to send chat histories recap: %v", err)
			return
		}

		return
	}
	if summarization == "" {
		h.Logger.Warn("summarization is empty")

		errMessage := tgbotapi.NewMessage(chatID, "暂时没有聊天记录可以生成聊天回顾哦，要再多聊点之后再试试吗？")
		errMessage.ReplyToMessageID = c.Update.Message.MessageID
		_, err = c.Bot.Send(errMessage)
		if err != nil {
			h.Logger.Errorf("failed to send chat histories recap: %v", err)
			return
		}

		return
	}

	h.Logger.Infof("sending chat histories recap for chat %d", chatID)
	message = tgbotapi.NewMessage(chatID, summarization)
	message.ReplyToMessageID = c.Update.Message.MessageID
	if err != nil {
		h.Logger.Errorf("failed to send chat histories recap: %v", err)
		return
	}
}

func (h *Handler) HandleEnableRecapCommand(c *handler.Context) {
	chatType := telegram.ChatType(c.Update.Message.Chat.Type)
	if !lo.Contains([]telegram.ChatType{telegram.ChatTypeGroup, telegram.ChatTypeSuperGroup}, chatType) {
		message := tgbotapi.NewMessage(c.Update.Message.Chat.ID, "聊天记录回顾功能只有群组和超级群组可以配置开启哦！")
		message.ReplyToMessageID = c.Update.Message.MessageID
		_, err := c.Bot.Send(message)
		if err != nil {
			h.Logger.Errorf("failed to send message to telegram: %v", err)
			return
		}

		return
	}

	err := h.TelegramChatFeatureFlags.EnableChatHistoriesRecap(c.Update.Message.Chat.ID, chatType)
	if err != nil {
		message := tgbotapi.NewMessage(c.Update.Message.Chat.ID, "聊天记录回顾功能开启失败，请稍后再试！")
		message.ReplyToMessageID = c.Update.Message.MessageID
		_, err := c.Bot.Send(message)
		if err != nil {
			h.Logger.Errorf("failed to send message to telegram: %v", err)
			return
		}

		h.Logger.Errorf("failed to enable chat histories recap: %v", err)
		return
	}

	message := tgbotapi.NewMessage(c.Update.Message.Chat.ID, "聊天记录回顾功能已开启，开启后将会自动收集群组中的聊天记录并定时发送聊天回顾快报！")
	message.ReplyToMessageID = c.Update.Message.MessageID
	_, err = c.Bot.Send(message)
	if err != nil {
		h.Logger.Errorf("failed to send message to telegram: %v", err)
		return
	}
}

func (h *Handler) HandleDisableRecapCommand(c *handler.Context) {
	chatType := telegram.ChatType(c.Update.Message.Chat.Type)
	if !lo.Contains([]telegram.ChatType{telegram.ChatTypeGroup, telegram.ChatTypeSuperGroup}, chatType) {
		message := tgbotapi.NewMessage(c.Update.Message.Chat.ID, "聊天记录回顾功能只有群组和超级群组可以配置关闭哦！")
		message.ReplyToMessageID = c.Update.Message.MessageID
		_, err := c.Bot.Send(message)
		if err != nil {
			h.Logger.Errorf("failed to send message to telegram: %v", err)
			return
		}

		return
	}

	err := h.TelegramChatFeatureFlags.DisableChatHistoriesRecap(c.Update.Message.Chat.ID, chatType)
	if err != nil {
		message := tgbotapi.NewMessage(c.Update.Message.Chat.ID, "聊天记录回顾功能关闭失败，请稍后再试！")
		message.ReplyToMessageID = c.Update.Message.MessageID
		_, err := c.Bot.Send(message)
		if err != nil {
			h.Logger.Errorf("failed to send message to telegram: %v", err)
			return
		}

		h.Logger.Errorf("failed to enable chat histories recap: %v", err)
		return
	}

	message := tgbotapi.NewMessage(c.Update.Message.Chat.ID, "聊天记录回顾功能已关闭，关闭后将不会自动收集群组中的聊天记录并定时发送聊天回顾快报了。")
	message.ReplyToMessageID = c.Update.Message.MessageID
	_, err = c.Bot.Send(message)
	if err != nil {
		h.Logger.Errorf("failed to send message to telegram: %v", err)
		return
	}
}
