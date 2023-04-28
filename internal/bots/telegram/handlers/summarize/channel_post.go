package summarize

import (
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/nekomeowww/insights-bot/pkg/bots/tgbot"
)

func (h *Handlers) HandleChannelPost(c *tgbot.Context) error {
	// 转发的消息不处理
	if c.Update.ChannelPost.ForwardFrom != nil {
		return nil
	}
	// 转发的消息不处理
	if c.Update.ChannelPost.ForwardFromChat != nil {
		return nil
	}
	// 若无 /s 命令则不处理
	if !strings.HasPrefix(c.Update.ChannelPost.Text, "/smr ") {
		return nil
	}

	urlString := strings.TrimSpace(strings.TrimPrefix(c.Update.ChannelPost.Text, "/smr "))
	summarization, err := h.smr.SummarizeInputURL(urlString)
	if err != nil {
		return tgbot.NewExceptionError(err)
	}

	_, err = c.Bot.Request(tgbotapi.EditMessageTextConfig{
		BaseEdit: tgbotapi.BaseEdit{
			ChatID:    c.Update.ChannelPost.Chat.ID,
			MessageID: c.Update.ChannelPost.MessageID,
		},
		ParseMode: "HTML", // 这个结构体里面需要给 ParseMode 字段赋一个字面量为 HTML 的枚举值
		Text:      summarization,
	})
	if err != nil {
		return tgbot.NewExceptionError(err)
	}

	return nil
}
