package smr

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	"github.com/samber/lo"

	"github.com/nekomeowww/insights-bot/ent"
	"github.com/nekomeowww/insights-bot/ent/feedbacksummarizationsreactions"
	"github.com/nekomeowww/insights-bot/pkg/bots/tgbot"
	"github.com/nekomeowww/insights-bot/pkg/types/bot/handlers/recap"
)

type FeedbackSummarizationsReactionsCounts struct {
	UpVotes   int
	DownVotes int
	Lmao      int
}

func (m *Model) FindFeedbackSummarizationsReactionCountsForChatIDAndLogID(chatID int64, logID uuid.UUID) (FeedbackSummarizationsReactionsCounts, error) {
	votes, err := m.ent.FeedbackSummarizationsReactions.
		Query().
		Where(
			feedbacksummarizationsreactions.ChatIDEQ(chatID),
			feedbacksummarizationsreactions.LogIDEQ(logID),
		).
		All(context.TODO())
	if err != nil {
		return FeedbackSummarizationsReactionsCounts{}, err
	}

	upVotes := len(lo.Filter(votes, func(item *ent.FeedbackSummarizationsReactions, _ int) bool {
		return item.Type == feedbacksummarizationsreactions.TypeUpVote
	}))
	downVotes := len(lo.Filter(votes, func(item *ent.FeedbackSummarizationsReactions, _ int) bool {
		return item.Type == feedbacksummarizationsreactions.TypeDownVote
	}))
	lmaos := len(lo.Filter(votes, func(item *ent.FeedbackSummarizationsReactions, _ int) bool {
		return item.Type == feedbacksummarizationsreactions.TypeLmao
	}))

	return FeedbackSummarizationsReactionsCounts{
		UpVotes:   upVotes,
		DownVotes: downVotes,
		Lmao:      lmaos,
	}, nil
}

func (m *Model) HasFeedbackReactSummarizationsToChatIDAndLogID(chatID int64, logID uuid.UUID, userID int64, reactionType feedbacksummarizationsreactions.Type) (bool, error) {
	existing, err := m.ent.FeedbackSummarizationsReactions.
		Query().
		Where(
			feedbacksummarizationsreactions.ChatIDEQ(chatID),
			feedbacksummarizationsreactions.LogIDEQ(logID),
			feedbacksummarizationsreactions.UserIDEQ(userID),
			feedbacksummarizationsreactions.TypeEQ(reactionType),
		).
		First(context.Background())
	if err != nil {
		if ent.IsNotFound(err) {
			return false, nil
		}

		return false, err
	}

	return existing != nil, nil
}

func (m *Model) FeedbackReactSummarizationsToChatIDAndLogID(chatID int64, logID uuid.UUID, userID int64, reactionType feedbacksummarizationsreactions.Type) error {
	affectedRows, err := m.ent.FeedbackSummarizationsReactions.
		Delete().
		Where(
			feedbacksummarizationsreactions.ChatIDEQ(chatID),
			feedbacksummarizationsreactions.LogIDEQ(logID),
			feedbacksummarizationsreactions.UserIDEQ(userID),
			feedbacksummarizationsreactions.TypeEQ(reactionType),
		).
		Exec(context.Background())
	if err != nil {
		return err
	}

	if affectedRows > 0 {
		return nil
	}

	_, err = m.ent.FeedbackSummarizationsReactions.
		Delete().
		Where(
			feedbacksummarizationsreactions.ChatIDEQ(chatID),
			feedbacksummarizationsreactions.LogIDEQ(logID),
			feedbacksummarizationsreactions.UserIDEQ(userID),
		).
		Exec(context.Background())
	if err != nil {
		return err
	}

	err = m.ent.FeedbackSummarizationsReactions.
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

func (m *Model) NewFeedbackSummarizationsUpVoteReactionButton(bot *tgbot.Bot, chatID int64, logID uuid.UUID, upVoteCount int) (tgbotapi.InlineKeyboardButton, error) {
	upVoteData, err := bot.AssignOneCallbackQueryData("smr/summarization/feedback/react", recap.FeedbackSummarizationReactionActionData{ChatID: chatID, Type: feedbacksummarizationsreactions.TypeUpVote, LogID: logID.String()})
	if err != nil {
		return tgbotapi.InlineKeyboardButton{}, err
	}

	upVoteButtonText := fmt.Sprintf("👍 %d", upVoteCount)
	if upVoteCount <= 0 {
		upVoteButtonText = "👍"
	}

	return tgbotapi.NewInlineKeyboardButtonData(upVoteButtonText, upVoteData), nil
}

func (m *Model) NewFeedbackSummarizationsDownVoteReactionButton(bot *tgbot.Bot, chatID int64, logID uuid.UUID, downVoteCount int) (tgbotapi.InlineKeyboardButton, error) {
	downVoteData, err := bot.AssignOneCallbackQueryData("smr/summarization/feedback/react", recap.FeedbackSummarizationReactionActionData{ChatID: chatID, Type: feedbacksummarizationsreactions.TypeDownVote, LogID: logID.String()})
	if err != nil {
		return tgbotapi.InlineKeyboardButton{}, err
	}

	downVoteButtonText := fmt.Sprintf("👎 %d", downVoteCount)
	if downVoteCount <= 0 {
		downVoteButtonText = "👎"
	}

	return tgbotapi.NewInlineKeyboardButtonData(downVoteButtonText, downVoteData), nil
}

func (m *Model) NewFeedbackSummarizationsLmaoReactionButton(bot *tgbot.Bot, chatID int64, logID uuid.UUID, downVoteCount int) (tgbotapi.InlineKeyboardButton, error) {
	lmaoData, err := bot.AssignOneCallbackQueryData("smr/summarization/feedback/react", recap.FeedbackSummarizationReactionActionData{ChatID: chatID, Type: feedbacksummarizationsreactions.TypeLmao, LogID: logID.String()})
	if err != nil {
		return tgbotapi.InlineKeyboardButton{}, err
	}

	lmaoButtonText := fmt.Sprintf("🤣 %d", downVoteCount)
	if downVoteCount <= 0 {
		lmaoButtonText = "🤣"
	}

	return tgbotapi.NewInlineKeyboardButtonData(lmaoButtonText, lmaoData), nil
}

func (m *Model) NewVoteSummarizationsReactionsInlineKeyboardMarkup(bot *tgbot.Bot, chatID int64, logID uuid.UUID, upVoteCount int, downVoteCount int, lmaoCount int) (tgbotapi.InlineKeyboardMarkup, error) {
	upVoteButton, err := m.NewFeedbackSummarizationsUpVoteReactionButton(bot, chatID, logID, upVoteCount)
	if err != nil {
		return tgbotapi.InlineKeyboardMarkup{}, err
	}

	downVoteButton, err := m.NewFeedbackSummarizationsDownVoteReactionButton(bot, chatID, logID, downVoteCount)
	if err != nil {
		return tgbotapi.InlineKeyboardMarkup{}, err
	}

	lmaoButton, err := m.NewFeedbackSummarizationsLmaoReactionButton(bot, chatID, logID, lmaoCount)
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
