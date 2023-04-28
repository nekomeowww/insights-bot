package recap

import (
	"github.com/nekomeowww/insights-bot/internal/models/chat_histories"
	"github.com/nekomeowww/insights-bot/internal/models/tgchats"
	"github.com/nekomeowww/insights-bot/pkg/bots/tgbot"
	"github.com/nekomeowww/insights-bot/pkg/types/telegram"
	"github.com/samber/lo"
	"go.uber.org/fx"
)

type NewMessageHandlerParams struct {
	fx.In

	ChatHistories *chat_histories.Model
	TgChats       *tgchats.Model
}

type MessageHandler struct {
	chatHistories *chat_histories.Model
	tgchats       *tgchats.Model
}

func NewMessageHandler() func(NewMessageHandlerParams) *MessageHandler {
	return func(params NewMessageHandlerParams) *MessageHandler {
		return &MessageHandler{
			chatHistories: params.ChatHistories,
			tgchats:       params.TgChats,
		}
	}
}

func (h *MessageHandler) HandleRecordMessage(c *tgbot.Context) error {
	chatType := telegram.ChatType(c.Update.Message.Chat.Type)
	if !lo.Contains([]telegram.ChatType{telegram.ChatTypeGroup, telegram.ChatTypeSuperGroup}, chatType) {
		return nil
	}

	enabled, err := h.tgchats.HasChatHistoriesRecapEnabled(c.Update.Message.Chat.ID, chatType)
	if err != nil {
		return tgbot.NewExceptionError(err)
	}
	if !enabled {
		return nil
	}

	err = h.chatHistories.SaveOneTelegramChatHistory(c.Update.Message)
	if err != nil {
		return tgbot.NewExceptionError(err)
	}

	return nil
}
