package middlewares

import (
	"github.com/nekomeowww/insights-bot/internal/models/chat_histories"
	"github.com/nekomeowww/insights-bot/internal/models/tgchats"
	"github.com/nekomeowww/insights-bot/pkg/bots/tgbot"
	"github.com/nekomeowww/insights-bot/pkg/types/telegram"
	"github.com/samber/lo"
)

func RecordMessage(chatHistories *chat_histories.Model, tgchats *tgchats.Model) func(c *tgbot.Context, next func()) {
	return func(c *tgbot.Context, next func()) {
		chatType := telegram.ChatType(c.Update.Message.Chat.Type)
		if !lo.Contains([]telegram.ChatType{telegram.ChatTypeGroup, telegram.ChatTypeSuperGroup}, chatType) {
			return
		}

		enabled, err := tgchats.HasChatHistoriesRecapEnabled(c.Update.Message.Chat.ID, chatType)
		if err != nil {
			c.Logger.Error(err)
			return
		}
		if !enabled {
			return
		}

		err = chatHistories.SaveOneTelegramChatHistory(c.Update.Message)
		if err != nil {
			c.Logger.Error(err)
			return
		}

		next()
	}
}
