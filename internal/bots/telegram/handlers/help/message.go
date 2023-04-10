package help

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/nekomeowww/insights-bot/pkg/handler"
)

func (h *Handler) HandleHelpCommand(c *handler.Context) {
	message := tgbotapi.NewMessage(c.Update.Message.Chat.ID, `
你好，欢迎使用 Insights Bot！
量子速读 用法：/smr <链接>
聊天回顾 用法：/recap
开启聊天记录回顾（群组） 用法：/enable_recap
关闭聊天记录回顾（群组） 用法：/disable_recap
	`)
	message.ReplyToMessageID = c.Update.Message.MessageID
	c.Bot.MustSend(message)
}
