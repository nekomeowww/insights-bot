package chathistories

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	"github.com/samber/lo"

	"github.com/nekomeowww/insights-bot/ent"
	"github.com/nekomeowww/insights-bot/ent/feedbackchathistoriesrecapsreactions"
	"github.com/nekomeowww/insights-bot/pkg/bots/tgbot"
	"github.com/nekomeowww/insights-bot/pkg/types/bot/handlers/recap"
)

type FeedbackChatHistoriesRecapsReactionsCounts struct {
	UpVotes   int
	DownVotes int
	Lmao      int
}

func (m *Model) FindFeedbackRecapsReactionCountsForChatIDAndLogID(chatID int64, logID uuid.UUID) (FeedbackChatHistoriesRecapsReactionsCounts, error) {
	votes, err := m.ent.FeedbackChatHistoriesRecapsReactions.
		Query().
		Where(
			feedbackchathistoriesrecapsreactions.ChatIDEQ(chatID),
			feedbackchathistoriesrecapsreactions.LogIDEQ(logID),
		).
		All(context.TODO())
	if err != nil {
		return FeedbackChatHistoriesRecapsReactionsCounts{}, err
	}

	upVotes := len(lo.Filter(votes, func(item *ent.FeedbackChatHistoriesRecapsReactions, _ int) bool {
		return item.Type == feedbackchathistoriesrecapsreactions.TypeUpVote
	}))
	downVotes := len(lo.Filter(votes, func(item *ent.FeedbackChatHistoriesRecapsReactions, _ int) bool {
		return item.Type == feedbackchathistoriesrecapsreactions.TypeDownVote
	}))
	lmaos := len(lo.Filter(votes, func(item *ent.FeedbackChatHistoriesRecapsReactions, _ int) bool {
		return item.Type == feedbackchathistoriesrecapsreactions.TypeLmao
	}))

	return FeedbackChatHistoriesRecapsReactionsCounts{
		UpVotes:   upVotes,
		DownVotes: downVotes,
		Lmao:      lmaos,
	}, nil
}

func (m *Model) FeedbackRecapsReactToChatIDAndLogID(chatID int64, logID uuid.UUID, userID int64, reactionType feedbackchathistoriesrecapsreactions.Type) error {
	affectedRows, err := m.ent.FeedbackChatHistoriesRecapsReactions.
		Delete().
		Where(
			feedbackchathistoriesrecapsreactions.ChatIDEQ(chatID),
			feedbackchathistoriesrecapsreactions.LogIDEQ(logID),
			feedbackchathistoriesrecapsreactions.UserIDEQ(userID),
			feedbackchathistoriesrecapsreactions.TypeEQ(reactionType),
		).
		Exec(context.Background())
	if err != nil {
		return err
	}
	if affectedRows > 0 {
		return nil
	}

	_, err = m.ent.FeedbackChatHistoriesRecapsReactions.
		Delete().
		Where(
			feedbackchathistoriesrecapsreactions.ChatIDEQ(chatID),
			feedbackchathistoriesrecapsreactions.LogIDEQ(logID),
			feedbackchathistoriesrecapsreactions.UserIDEQ(userID),
		).
		Exec(context.Background())
	if err != nil {
		return err
	}

	err = m.ent.FeedbackChatHistoriesRecapsReactions.
		Create().
		SetChatID(chatID).
		SetLogID(logID).
		SetUserID(userID).
		SetType(reactionType).
		Exec(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func (m *Model) NewFeedbackRecapsUpVoteButton(bot *tgbot.Bot, chatID int64, logID uuid.UUID, upVoteCount int) (tgbotapi.InlineKeyboardButton, error) {
	upVoteData, err := bot.AssignOneCallbackQueryData("smr/summarization/feedback/react", recap.FeedbackRecapReactionActionData{ChatID: chatID, Type: feedbackchathistoriesrecapsreactions.TypeUpVote, LogID: logID.String()})
	if err != nil {
		return tgbotapi.InlineKeyboardButton{}, err
	}

	upVoteButtonText := fmt.Sprintf("👍 %d", upVoteCount)
	if upVoteCount <= 0 {
		upVoteButtonText = "👍"
	}

	return tgbotapi.NewInlineKeyboardButtonData(upVoteButtonText, upVoteData), nil
}

func (m *Model) NewFeedbackRecapsDownVoteButton(bot *tgbot.Bot, chatID int64, logID uuid.UUID, downVoteCount int) (tgbotapi.InlineKeyboardButton, error) {
	downVoteData, err := bot.AssignOneCallbackQueryData("smr/summarization/feedback/react", recap.FeedbackRecapReactionActionData{ChatID: chatID, Type: feedbackchathistoriesrecapsreactions.TypeDownVote, LogID: logID.String()})
	if err != nil {
		return tgbotapi.InlineKeyboardButton{}, err
	}

	downVoteButtonText := fmt.Sprintf("👎 %d", downVoteCount)
	if downVoteCount <= 0 {
		downVoteButtonText = "👎"
	}

	return tgbotapi.NewInlineKeyboardButtonData(downVoteButtonText, downVoteData), nil
}

func (m *Model) NewFeedbackRecapsLmaoButton(bot *tgbot.Bot, chatID int64, logID uuid.UUID, downVoteCount int) (tgbotapi.InlineKeyboardButton, error) {
	lmaoData, err := bot.AssignOneCallbackQueryData("smr/summarization/feedback/react", recap.FeedbackRecapReactionActionData{ChatID: chatID, Type: feedbackchathistoriesrecapsreactions.TypeLmao, LogID: logID.String()})
	if err != nil {
		return tgbotapi.InlineKeyboardButton{}, err
	}

	lmaoButtonText := fmt.Sprintf("🤣 %d", downVoteCount)
	if downVoteCount <= 0 {
		lmaoButtonText = "🤣"
	}

	return tgbotapi.NewInlineKeyboardButtonData(lmaoButtonText, lmaoData), nil
}

func (m *Model) NewVoteRecapInlineKeyboardMarkup(bot *tgbot.Bot, chatID int64, logID uuid.UUID, upVoteCount int, downVoteCount int, lmaoCount int) (tgbotapi.InlineKeyboardMarkup, error) {
	upVoteButton, err := m.NewFeedbackRecapsUpVoteButton(bot, chatID, logID, upVoteCount)
	if err != nil {
		return tgbotapi.InlineKeyboardMarkup{}, err
	}

	downVoteButton, err := m.NewFeedbackRecapsDownVoteButton(bot, chatID, logID, downVoteCount)
	if err != nil {
		return tgbotapi.InlineKeyboardMarkup{}, err
	}

	lmaoButton, err := m.NewFeedbackRecapsLmaoButton(bot, chatID, logID, lmaoCount)
	if err != nil {
		return tgbotapi.InlineKeyboardMarkup{}, err
	}

	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			upVoteButton,
			downVoteButton,
			lmaoButton,
		),
	), nil
}

func (m *Model) NewVoteRecapWithUnsubscribeInlineKeyboardMarkup(bot *tgbot.Bot, chatID int64, chatTitle string, fromID int64, logID uuid.UUID, upVoteCount int, downVoteCount int, lmaoCount int) (tgbotapi.InlineKeyboardMarkup, error) {
	upVoteButton, err := m.NewFeedbackRecapsUpVoteButton(bot, chatID, logID, upVoteCount)
	if err != nil {
		return tgbotapi.InlineKeyboardMarkup{}, err
	}

	downVoteButton, err := m.NewFeedbackRecapsDownVoteButton(bot, chatID, logID, downVoteCount)
	if err != nil {
		return tgbotapi.InlineKeyboardMarkup{}, err
	}

	lmaoButton, err := m.NewFeedbackRecapsLmaoButton(bot, chatID, logID, lmaoCount)
	if err != nil {
		return tgbotapi.InlineKeyboardMarkup{}, err
	}

	buttonData, err := bot.AssignOneCallbackQueryData("recap/unsubscribe_recap", recap.UnsubscribeRecapActionData{
		ChatID:    chatID,
		ChatTitle: chatTitle,
		FromID:    fromID,
	})
	if err != nil {
		return tgbotapi.InlineKeyboardMarkup{}, err
	}

	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			upVoteButton,
			downVoteButton,
			lmaoButton,
		),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("取消订阅", buttonData)),
	), nil
}
