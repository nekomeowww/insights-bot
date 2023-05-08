package summarize

import (
	"context"
	"errors"
	"net/url"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/samber/lo"

	"github.com/nekomeowww/insights-bot/internal/models/smr"
	"github.com/nekomeowww/insights-bot/pkg/bots/tgbot"
)

func (h Handlers) Command() string {
	return "smr"
}

func (h Handlers) CommandHelp() string {
	return "量子速读网页文章（也支持在频道中使用） 用法：/smr <code>&lt;链接&gt;</code>"
}

func (h *Handlers) Handle(c *tgbot.Context) (tgbot.Response, error) {
	urlString := c.Update.Message.CommandArguments()
	if urlString == "" {
		return nil, tgbot.NewMessageError("没有找到链接，可以发送一个有效的链接吗？用法：/smr <链接>").WithReply(c.Update.Message)
	}

	parsedURL, err := url.Parse(urlString)
	if err != nil {
		return nil, tgbot.NewMessageError("你发来的链接无法被理解，可以重新发一个试试。用法：/smr <链接>").WithReply(c.Update.Message)
	}
	if parsedURL.Scheme == "" || !lo.Contains([]string{"http", "https"}, parsedURL.Scheme) {
		return nil, tgbot.NewMessageError("你发来的链接无法被理解，可以重新发一个试试。用法：/smr <链接>").WithReply(c.Update.Message)
	}

	message := tgbotapi.NewMessage(c.Update.Message.Chat.ID, "请稍等，量子速读中...")
	message.ReplyToMessageID = c.Update.Message.MessageID

	processingMessage, err := c.Bot.Send(message)
	if err != nil {
		return nil, tgbot.NewExceptionError(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*2)
	defer cancel()

	summarization, err := h.smr.SummarizeInputURL(ctx, urlString)
	if err != nil {
		if errors.Is(err, smr.ErrContentNotSupported) {
			return nil, tgbot.NewMessageError("暂时不支持量子速读这样的内容呢，可以换个别的链接试试。").WithEdit(&processingMessage)
		}
		if errors.Is(err, smr.ErrNetworkError) || errors.Is(err, smr.ErrRequestFailed) {
			return nil, tgbot.NewMessageError("量子速读的链接读取失败了哦。可以再试试？").WithEdit(&processingMessage)
		}

		return nil, tgbot.NewMessageError("量子速读失败了。可以再试试？").WithEdit(&processingMessage)
	}

	return c.NewMessageReplyTo(summarization.FormatSummarizationAsHTML(), c.Update.Message.MessageID).WithParseModeHTML(), nil
}
