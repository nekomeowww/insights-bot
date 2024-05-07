package summarize

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/nekomeowww/insights-bot/pkg/bots/tgbot"
	"github.com/nekomeowww/insights-bot/pkg/types/smr"
	"go.uber.org/zap"
)

func (h *Handlers) handleCallbackQueryRetry(c *tgbot.Context) (tgbot.Response, error) {
	messageID := c.Update.CallbackQuery.Message.MessageID
	var data smr.TaskInfo

	err := c.BindFromCallbackQueryData(&data)
	if err != nil {
		h.logger.Error("failed to bind callback query data when retry smr",
			zap.Error(err),
			zap.Int("message_id", messageID),
			zap.Int64("chat_id", c.Update.CallbackQuery.Message.Chat.ID),
			zap.Int64("from_id", c.Update.CallbackQuery.From.ID),
			zap.String("data", c.Update.CallbackQuery.Data),
		)

		return nil, nil
	}

	err = h.smrQueue.AddTask(data)

	if err != nil {
		h.logger.Error("failed to move task back to queue", zap.Error(err))
		return nil, nil
	}

	// remove the retry button
	c.Bot.MayRequest(tgbotapi.NewEditMessageTextAndMarkup(
		data.ChatID,
		messageID,
		h.i18n.TWithLanguage(data.Language, "commands.groups.summarization.commands.smr.reading"),
		tgbotapi.InlineKeyboardMarkup{
			InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{},
		}))

	return nil, nil
}
