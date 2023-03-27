package summarize

import (
	"errors"
	"net/url"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/samber/lo"

	"github.com/nekomeowww/insights-bot/pkg/handler"
)

func (h *Handler) commandEmptyResponse(c *handler.Context) {
	message := tgbotapi.NewMessage(c.Update.Message.Chat.ID, "没有找到链接，可以发送一个有效的链接吗？用法：/smr <链接>")
	message.ReplyToMessageID = c.Update.Message.MessageID
	_, err := c.Bot.Request(message)
	if err != nil {
		h.Logger.Errorf("failed to send message to telegram... %v", err)
		return
	}
}

func (h *Handler) commandInvalidResponse(c *handler.Context) {
	message := tgbotapi.NewMessage(c.Update.Message.Chat.ID, "你发来的链接无法被理解，可以重新发一个试试。用法：/smr <链接>")
	message.ReplyToMessageID = c.Update.Message.MessageID
	_, err := c.Bot.Request(message)
	if err != nil {
		h.Logger.Errorf("failed to send message to telegram... %v", err)
		return
	}
}

func (h *Handler) HandleMessage(c *handler.Context) {
	if c.Update.Message.Command() != "smr" {
		return
	}

	urlString := c.Update.Message.CommandArguments()
	if urlString == "" {
		go h.commandEmptyResponse(c)
		return
	}

	parsedURL, err := url.Parse(urlString)
	if err != nil {
		go h.commandInvalidResponse(c)
		return
	}
	if parsedURL.Scheme == "" || !lo.Contains([]string{"http", "https"}, parsedURL.Scheme) {
		go h.commandInvalidResponse(c)
		return
	}

	message := tgbotapi.NewMessage(c.Update.Message.Chat.ID, "请稍等，量子速读中...")
	message.ReplyToMessageID = c.Update.Message.MessageID
	newMessage, err := c.Bot.Send(message)
	if err != nil {
		h.Logger.Errorf("failed to send message to telegram... %v", err)
		return
	}

	summarization, err := h.summarizeInputURL(urlString)
	if err != nil {
		h.Logger.Error(err)
		if errors.Is(err, ErrContentNotSupported) {
			_, err = c.Bot.Request(tgbotapi.EditMessageTextConfig{
				BaseEdit: tgbotapi.BaseEdit{
					ChatID:    c.Update.Message.Chat.ID,
					MessageID: newMessage.MessageID,
				},
				Text: "暂时不支持量子速读这样的内容呢，可以换个别的链接试试。",
			})
			if err != nil {
				h.Logger.Errorf("failed to send message to telegram... %v", err)
				return
			}
		}
		if errors.Is(err, ErrNetworkError) || errors.Is(err, ErrRequestFailed) {
			_, err = c.Bot.Request(tgbotapi.EditMessageTextConfig{
				BaseEdit: tgbotapi.BaseEdit{
					ChatID:    c.Update.Message.Chat.ID,
					MessageID: newMessage.MessageID,
				},
				Text: "量子速读的链接读取失败了哦。可以再试试？",
			})
			if err != nil {
				h.Logger.Errorf("failed to send message to telegram... %v", err)
				return
			}
		} else {
			_, err = c.Bot.Request(tgbotapi.EditMessageTextConfig{
				BaseEdit: tgbotapi.BaseEdit{
					ChatID:    c.Update.Message.Chat.ID,
					MessageID: newMessage.MessageID,
				},
				Text: "量子速读失败了。可以再试试？",
			})
			if err != nil {
				h.Logger.Errorf("failed to send message to telegram... %v", err)
				return
			}
		}

		return
	}

	message = tgbotapi.NewMessage(c.Update.Message.Chat.ID, summarization)
	message.ParseMode = "HTML"
	message.ReplyToMessageID = c.Update.Message.MessageID
	_, err = c.Bot.Request(message)
	if err != nil {
		h.Logger.Errorf("failed to send message to telegram... %v", err)
		return
	}
}
