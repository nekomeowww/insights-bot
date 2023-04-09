package chat_with_chat_history

import (
	"github.com/nekomeowww/insights-bot/pkg/handler"
	"github.com/nekomeowww/insights-bot/pkg/types/telegram"
	"github.com/samber/lo"
)

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
