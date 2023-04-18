package help

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/nekomeowww/insights-bot/pkg/handler"
)

func (h *Handler) HandleHelpCommand(c *handler.Context) {
	message := tgbotapi.NewMessage(c.Update.Message.Chat.ID, ""+
		"你好，欢迎使用 Insights Bot！\n"+
		"我当前支持这些命令：\n"+
		"量子速读 用法：/smr <code>&lt;链接&gt;</code>\n"+
		"聊天回顾 用法：/recap\n"+
		"<em>以下命令仅在群聊中可用：</em>\n"+
		"开启聊天记录回顾（需要管理权限） 用法：/enable_recap\n"+
		"关闭聊天记录回顾（需要管理权限） 用法：/disable_recap"+
		"")
	message.ReplyToMessageID = c.Update.Message.MessageID
	message.ParseMode = "HTML"
	c.Bot.MustSend(message)
}
