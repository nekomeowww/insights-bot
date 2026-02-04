package summarize

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	"github.com/nekomeowww/insights-bot/ent/feedbacksummarizationsreactions"
	"github.com/nekomeowww/insights-bot/pkg/bots/tgbot"
	"github.com/nekomeowww/insights-bot/pkg/types/bot/handlers/recap"
	"go.uber.org/zap"
)

func (h *Handlers) handleCallbackQueryReact(c *tgbot.Context) (tgbot.Response, error) {
	messageID := c.Update.CallbackQuery.Message.MessageID

	var data recap.FeedbackSummarizationReactionActionData

	err := c.BindFromCallbackQueryData(&data)
	if err != nil {
		h.logger.Error("failed to bind callback query data",
			zap.Error(err),
			zap.Int("message_id", messageID),
			zap.Int64("chat_id", c.Update.CallbackQuery.Message.Chat.ID),
			zap.Int64("from_id", c.Update.CallbackQuery.From.ID),
			zap.String("data", c.Update.CallbackQuery.Data),
		)

		return nil, nil
	}

	logID, err := uuid.Parse(data.LogID)
	if err != nil {
		h.logger.Error("failed to parse log id",
			zap.Error(err),
			zap.Int("message_id", messageID),
			zap.Int64("chat_id", data.ChatID),
			zap.Int64("from_id", c.Update.CallbackQuery.From.ID),
			zap.String("log_id", data.LogID),
		)

		return nil, nil
	}

	switch data.Type {
	case feedbacksummarizationsreactions.TypeNone:
		return nil, nil
	case feedbacksummarizationsreactions.TypeUpVote, feedbacksummarizationsreactions.TypeDownVote, feedbacksummarizationsreactions.TypeLmao:
		err = h.smr.FeedbackReactSummarizationsToChatIDAndLogID(data.ChatID, logID, c.Update.CallbackQuery.From.ID, data.Type)
	default:
		return nil, nil
	}

	if err != nil {
		h.logger.Error("failed to react to recap of chat id and log id",
			zap.Error(err),
			zap.Int("message_id", messageID),
			zap.Int64("from_id", c.Update.CallbackQuery.From.ID),
			zap.Int64("chat_id", data.ChatID),
			zap.String("log_id", data.LogID),
			zap.String("action", "up vote"),
		)

		return nil, nil
	}

	counts, err := h.smr.FindFeedbackSummarizationsReactionCountsForChatIDAndLogID(data.ChatID, logID)
	if err != nil {
		h.logger.Error("failed to find feedback recaps reactions for chat id and log id",
			zap.Error(err),
			zap.Int("message_id", messageID),
			zap.Int64("from_id", c.Update.CallbackQuery.From.ID),
			zap.Int64("chat_id", data.ChatID),
			zap.String("log_id", data.LogID),
		)

		return nil, nil
	}

	inlineKeyboardMarkup, err := h.smr.NewVoteSummarizationsReactionsInlineKeyboardMarkup(c.Bot, data.ChatID, logID, counts.UpVotes, counts.DownVotes, counts.Lmao)
	if err != nil {
		h.logger.Error("failed to new up vote recap inline keyboard markup",
			zap.Error(err),
			zap.Int("message_id", messageID),
			zap.Int64("from_id", c.Update.CallbackQuery.From.ID),
			zap.Int64("chat_id", data.ChatID),
			zap.String("log_id", data.LogID),
		)

		return nil, nil
	}

	c.Bot.MayRequest(tgbotapi.NewEditMessageReplyMarkup(data.ChatID, messageID, inlineKeyboardMarkup))

	return nil, nil
}
