package recap

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	"github.com/nekomeowww/insights-bot/ent/feedbackchathistoriesrecapsreactions"
	"github.com/nekomeowww/insights-bot/pkg/bots/tgbot"
	"github.com/nekomeowww/insights-bot/pkg/types/bot/handlers/recap"
	"go.uber.org/zap"
)

func (h *CallbackQueryHandler) handleCallbackQueryReact(c *tgbot.Context) (tgbot.Response, error) {
	messageID := c.Update.CallbackQuery.Message.MessageID

	var data recap.FeedbackRecapReactionActionData

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
	case feedbackchathistoriesrecapsreactions.TypeNone:
		return nil, nil
	case feedbackchathistoriesrecapsreactions.TypeUpVote, feedbackchathistoriesrecapsreactions.TypeDownVote, feedbackchathistoriesrecapsreactions.TypeLmao:
		err = h.chatHistories.FeedbackRecapsReactToChatIDAndLogID(data.ChatID, logID, c.Update.CallbackQuery.From.ID, data.Type)
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

	counts, err := h.chatHistories.FindFeedbackRecapsReactionCountsForChatIDAndLogID(data.ChatID, logID)
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

	upVoteButton, err := h.chatHistories.NewFeedbackRecapsUpVoteButton(c.Bot, data.ChatID, logID, counts.UpVotes)
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

	downVoteButton, err := h.chatHistories.NewFeedbackRecapsDownVoteButton(c.Bot, data.ChatID, logID, counts.DownVotes)
	if err != nil {
		h.logger.Error("failed to new down vote recap inline keyboard markup",
			zap.Error(err),
			zap.Int("message_id", messageID),
			zap.Int64("from_id", c.Update.CallbackQuery.From.ID),
			zap.Int64("chat_id", data.ChatID),
			zap.String("log_id", data.LogID),
		)

		return nil, nil
	}

	lmaoButton, err := h.chatHistories.NewFeedbackRecapsLmaoButton(c.Bot, data.ChatID, logID, counts.Lmao)
	if err != nil {
		h.logger.Error("failed to new lmao recap inline keyboard markup",
			zap.Error(err),
			zap.Int("message_id", messageID),
			zap.Int64("from_id", c.Update.CallbackQuery.From.ID),
			zap.Int64("chat_id", data.ChatID),
			zap.String("log_id", data.LogID),
		)

		return nil, nil
	}

	inlineKeyboardMarkup := c.Update.CallbackQuery.Message.ReplyMarkup
	if inlineKeyboardMarkup == nil || len(inlineKeyboardMarkup.InlineKeyboard) == 0 {
		return nil, nil
	}

	for i := range inlineKeyboardMarkup.InlineKeyboard {
		for j := range inlineKeyboardMarkup.InlineKeyboard[i] {
			if inlineKeyboardMarkup.InlineKeyboard[i][j].CallbackData == nil {
				continue
			}

			if *inlineKeyboardMarkup.InlineKeyboard[i][j].CallbackData != c.Update.CallbackQuery.Data {
				continue
			}

			inlineKeyboardMarkup.InlineKeyboard[i] = tgbotapi.NewInlineKeyboardRow(upVoteButton, downVoteButton, lmaoButton)
		}
	}

	c.Bot.MayRequest(tgbotapi.NewEditMessageReplyMarkup(data.ChatID, messageID, *inlineKeyboardMarkup))

	return nil, nil
}
