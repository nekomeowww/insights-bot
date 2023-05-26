package summarize

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/nekomeowww/insights-bot/internal/models/smr"
	smr2 "github.com/nekomeowww/insights-bot/internal/services/smr"
	"github.com/nekomeowww/insights-bot/internal/services/smr/types"
	"github.com/nekomeowww/insights-bot/pkg/bots/tgbot"
)

func (h *Handlers) Handle(c *tgbot.Context) (tgbot.Response, error) {
	urlString := c.Update.Message.CommandArguments()
	if urlString == "" && c.Update.Message.ReplyToMessage != nil && c.Update.Message.ReplyToMessage.Text != "" {
		urlString = c.Update.Message.ReplyToMessage.Text
	}

	err := smr2.CheckUrl(urlString)
	if err != nil {
		if smr2.IsUrlCheckError(err) {
			return nil, tgbot.NewMessageError(err.Error()).WithReply(c.Update.Message)
		}

		return nil, tgbot.NewMessageError("出现了一些问题，可以再试试？").WithReply(c.Update.Message)
	}

	message := tgbotapi.NewMessage(c.Update.Message.Chat.ID, "请稍等，量子速读中...")
	message.ReplyToMessageID = c.Update.Message.MessageID

	processingMessage, err := c.Bot.Send(message)
	if err != nil {
		return nil, tgbot.NewExceptionError(err)
	}

	err = h.smrService.AddTask(types.TaskInfo{
		Platform:  smr.FromPlatformTelegram,
		Url:       urlString,
		ChatID:    c.Update.Message.Chat.ID,
		MessageID: processingMessage.MessageID,
	})

	if err != nil {
		return nil, tgbot.NewMessageError("量子速读请求发送失败了，可以再试试？").WithEdit(&processingMessage)
	}

	return nil, nil
}
