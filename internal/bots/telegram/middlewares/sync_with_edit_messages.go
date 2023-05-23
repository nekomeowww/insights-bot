package middlewares

import (
	"github.com/nekomeowww/insights-bot/internal/models/chathistories"
	"github.com/nekomeowww/insights-bot/pkg/bots/tgbot"
)

func SyncWithEditedMessage(chatHistories *chathistories.Model) func(c *tgbot.Context, next func()) {
	return func(c *tgbot.Context, next func()) {
		if c.Update.EditedMessage == nil {
			return
		}

		err := chatHistories.UpdateOneMessage(c.Update.EditedMessage)
		if err != nil {
			c.Logger.Error(err)
		}

		next()
	}
}
