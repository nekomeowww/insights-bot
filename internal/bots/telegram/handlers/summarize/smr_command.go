package summarize

import (
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/samber/lo"
	"go.uber.org/zap"

	"github.com/nekomeowww/insights-bot/internal/services/smr"
	"github.com/nekomeowww/insights-bot/pkg/bots/tgbot"
	"github.com/nekomeowww/insights-bot/pkg/i18n"
	"github.com/nekomeowww/insights-bot/pkg/types/bot"
	types "github.com/nekomeowww/insights-bot/pkg/types/smr"
)

func (h *Handlers) Handle(c *tgbot.Context) (tgbot.Response, error) {
	urlString := c.Update.Message.CommandArguments()
	if urlString == "" && c.Update.Message.ReplyToMessage != nil && c.Update.Message.ReplyToMessage.Text != "" {
		urlString = c.Update.Message.ReplyToMessage.Text
	}
	if urlString == "" {
		return nil, tgbot.
			NewMessageError(c.T("commands.groups.summarization.commands.smr.noLinksFound.telegram")).
			WithReply(c.Update.Message)
	}

	urlString = strings.TrimSpace(urlString)
	if !strings.HasPrefix(urlString, "http://") && !strings.HasPrefix(urlString, "https://") {
		urlString = "https://" + urlString
	}

	err, originErr := smr.CheckUrl(urlString)
	if err != nil {
		if smr.IsUrlCheckError(err) {
			return nil, tgbot.
				NewMessageError(smr.FormatUrlCheckError(err, bot.FromPlatformTelegram, c.Language(), h.i18n)).
				WithReply(c.Update.Message).
				WithParseModeHTML()
		}

		return nil, tgbot.NewExceptionError(originErr).WithReply(c.Update.Message)
	}

	message := tgbotapi.NewMessage(c.Update.Message.Chat.ID, c.T("commands.groups.summarization.commands.smr.reading"))
	message.ReplyToMessageID = c.Update.Message.MessageID

	processingMessage, err := c.Bot.Send(message)
	if err != nil {
		return nil, tgbot.NewExceptionError(err)
	}

	chatID := c.Update.Message.Chat.ID
	perSeconds := h.smr.SummarizeWebpageRatePerSeconds()

	_, ttl, ok, err := c.RateLimitForCommand(chatID, "/smr", 1, perSeconds)
	if err != nil {
		h.logger.Error("failed to check rate limit for command /smr", zap.Error(err))
	}
	if !ok {
		return nil, tgbot.
			NewMessageError(c.T("", i18n.M{
				"Seconds":           perSeconds,
				"SecondsToBeWaited": lo.Ternary(ttl/time.Minute <= 1, 1, ttl/time.Minute),
			})).
			WithReply(c.Update.Message)
	}

	err = h.smrQueue.AddTask(types.TaskInfo{
		Platform:  bot.FromPlatformTelegram,
		URL:       urlString,
		ChatID:    c.Update.Message.Chat.ID,
		MessageID: processingMessage.MessageID,
		Language:  c.Language(),
	})
	if err != nil {
		return nil, tgbot.
			NewExceptionError(err).
			WithMessage(c.T("commands.groups.summarization.commands.smr.failedToRead")).
			WithEdit(&processingMessage)
	}

	return nil, nil
}
