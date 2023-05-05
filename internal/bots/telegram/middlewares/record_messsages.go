package middlewares

import (
	"github.com/nekomeowww/insights-bot/internal/models/chathistories"
	"github.com/nekomeowww/insights-bot/internal/models/tgchats"
	"github.com/nekomeowww/insights-bot/pkg/bots/tgbot"
	"github.com/nekomeowww/insights-bot/pkg/types/telegram"
	"github.com/samber/lo"
)

func RecordMessage(chatHistories *chathistories.Model, tgchats *tgchats.Model) func(c *tgbot.Context, next func()) {
	return func(c *tgbot.Context, next func()) {
		if c.Update.Message == nil {
			return
		}

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
