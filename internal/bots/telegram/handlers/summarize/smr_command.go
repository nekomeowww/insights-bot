package summarize

import (
	"fmt"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/nekomeowww/insights-bot/internal/services/smr"
	"github.com/nekomeowww/insights-bot/pkg/bots/tgbot"
	"github.com/nekomeowww/insights-bot/pkg/types/bot"
	types "github.com/nekomeowww/insights-bot/pkg/types/smr"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

func (h *Handlers) Handle(c *tgbot.Context) (tgbot.Response, error) {
	urlString := c.Update.Message.CommandArguments()
	if urlString == "" && c.Update.Message.ReplyToMessage != nil && c.Update.Message.ReplyToMessage.Text != "" {
		urlString = c.Update.Message.ReplyToMessage.Text
	}

	err, originErr := smr.CheckUrl(urlString)
	if err != nil {
		if smr.IsUrlCheckError(err) {
			return nil, tgbot.
				NewMessageError(smr.FormatUrlCheckError(err, bot.FromPlatformTelegram)).
				WithReply(c.Update.Message).
				WithParseModeHTML()
		}

		return nil, tgbot.NewExceptionError(originErr).WithReply(c.Update.Message)
	}

	message := tgbotapi.NewMessage(c.Update.Message.Chat.ID, "请稍等，量子速读中...")
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
			NewMessageError(fmt.Sprintf("很抱歉，您的操作触发了我们的限制机制，为了保证系统的可用性，本命令每最多 %d 分钟最多使用一次，请您耐心等待 %d 分钟后再试，感谢您的理解和支持。", perSeconds, lo.Ternary(ttl/time.Minute <= 1, 1, ttl/time.Minute))).
			WithReply(c.Update.Message)
	}

	err = h.smrQueue.AddTask(types.TaskInfo{
		Platform:  bot.FromPlatformTelegram,
		URL:       urlString,
		ChatID:    c.Update.Message.Chat.ID,
		MessageID: processingMessage.MessageID,
	})
	if err != nil {
		return nil, tgbot.NewExceptionError(err).WithMessage("量子速读失败了，可以再试试？").WithEdit(&processingMessage)
	}

	return nil, nil
}
