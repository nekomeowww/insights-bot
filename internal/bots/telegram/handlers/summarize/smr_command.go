package summarize

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/nekomeowww/insights-bot/internal/models/smr"
	"github.com/nekomeowww/insights-bot/internal/services/smr/smrutils"
	"github.com/nekomeowww/insights-bot/internal/services/smr/types"
	"github.com/nekomeowww/insights-bot/pkg/bots/tgbot"
)

func (h *Handlers) Handle(c *tgbot.Context) (tgbot.Response, error) {
	urlString := c.Update.Message.CommandArguments()
	if urlString == "" && c.Update.Message.ReplyToMessage != nil && c.Update.Message.ReplyToMessage.Text != "" {
		urlString = c.Update.Message.ReplyToMessage.Text
	}

	err, originErr := smrutils.CheckUrl(urlString)
	if err != nil {
		if smrutils.IsUrlCheckError(err) {
			return nil, tgbot.NewMessageError(smrutils.FormatUrlCheckError(err, smr.FromPlatformTelegram)).WithParseModeHTML().WithReply(c.Update.Message)
		}

		return nil, tgbot.NewExceptionError(originErr).WithReply(c.Update.Message)
	}

	message := tgbotapi.NewMessage(c.Update.Message.Chat.ID, "请稍等，量子速读中...")
	message.ReplyToMessageID = c.Update.Message.MessageID

	processingMessage, err := c.Bot.Send(message)
	if err != nil {
		return nil, tgbot.NewExceptionError(err)
	}

	err = h.smrQueue.AddTask(types.TaskInfo{
		Platform:  smr.FromPlatformTelegram,
		URL:       urlString,
		ChatID:    c.Update.Message.Chat.ID,
		MessageID: processingMessage.MessageID,
	})

	if err != nil {
		return nil, tgbot.NewExceptionError(err).WithMessage("量子速读失败了，可以再试试？").WithEdit(&processingMessage)
	}

	return nil, nil
}
