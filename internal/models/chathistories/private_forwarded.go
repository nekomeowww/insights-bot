package chathistories

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/nekomeowww/insights-bot/internal/thirdparty/openai"
	"github.com/nekomeowww/insights-bot/pkg/bots/tgbot"
	"github.com/nekomeowww/insights-bot/pkg/types/redis"
	"github.com/nekomeowww/insights-bot/pkg/types/telegram"
	"github.com/redis/rueidis"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

type telegramPrivateForwardedReplayChatHistory struct {
	ChatID           int64  `json:"chat_id"`
	ChatType         string `json:"chat_type"`
	ChatTitle        string `json:"chat_title"`
	MessageID        int    `json:"message_id"`
	ActorID          int64  `json:"actor_id"`
	ActorUsername    string `json:"actor_username"`
	ActorDisplayName string `json:"actor_display_name"`
	Text             string `json:"text"`
	ChattedAt        int64  `json:"chatted_at"`
}

func (m *Model) HasOngoingRecapForwardedFromPrivateMessages(userID int64) (bool, error) {
	getCmd := m.redis.B().
		Get().
		Key(redis.RecapReplayFromPrivateMessageControl1.Format(userID)).
		Build()

	str, err := m.redis.Do(context.Background(), getCmd).ToString()
	if err != nil {
		if rueidis.IsRedisNil(err) {
			return false, nil
		}

		return false, err
	}

	return str == "1", nil
}

func (m *Model) EnabledRecapForwardedFromPrivateMessages(userID int64) error {
	setCmd := m.redis.B().
		Set().
		Key(redis.RecapReplayFromPrivateMessageControl1.Format(userID)).
		Value("1").
		ExSeconds(60 * 60 * 2).
		Build()

	err := m.redis.Do(context.Background(), setCmd).Error()
	if err != nil {
		return err
	}

	return nil
}

func (m *Model) DisableRecapForwardedFromPrivateMessages(userID int64) error {
	delCmd := m.redis.B().
		Del().
		Key(redis.RecapReplayFromPrivateMessageBatch1.Format(userID)).
		Key(redis.RecapReplayFromPrivateMessageControl1.Format(userID)).
		Build()

	err := m.redis.Do(context.Background(), delCmd).Error()
	if err != nil {
		return err
	}

	return nil
}

func (m *Model) SaveOneTelegramPrivateForwardedReplayChatHistory(message *tgbotapi.Message) error {
	has, err := m.HasOngoingRecapForwardedFromPrivateMessages(message.From.ID)
	if err != nil {
		return err
	}
	if !has {
		return nil
	}

	text, err := m.extractTextFromMessage(message)
	if err != nil {
		return err
	}
	if text == "" {
		return nil
	}

	chatActorFullname := tgbot.FullNameFromFirstAndLastName(message.From.FirstName, message.From.LastName)

	telegramChatHistory := telegramPrivateForwardedReplayChatHistory{
		ChatID:           message.Chat.ID,
		ChatType:         message.Chat.Type,
		ChatTitle:        chatActorFullname,
		MessageID:        message.MessageID,
		ActorUsername:    message.From.UserName,
		ActorDisplayName: chatActorFullname,
		ChattedAt:        time.Unix(int64(message.Date), 0).UnixMilli(),
	}

	if message.ForwardFrom != nil {
		telegramChatHistory.ActorID = message.ForwardFrom.ID
		telegramChatHistory.ActorUsername = message.ForwardFrom.UserName
		telegramChatHistory.ActorDisplayName = tgbot.FullNameFromFirstAndLastName(message.ForwardFrom.FirstName, message.ForwardFrom.LastName)
	}
	if message.ForwardFrom == nil && message.ForwardSenderName != "" {
		telegramChatHistory.ActorID = 0
		telegramChatHistory.ActorUsername = message.ForwardSenderName
		telegramChatHistory.ActorDisplayName = message.ForwardSenderName
	}
	if message.ForwardFromChat != nil {
		telegramChatHistory.Text = fmt.Sprintf("[forwarded from %s]: %s", message.ForwardFromChat.Title, text)
	} else {
		telegramChatHistory.Text = text
	}

	zaddCmd := m.redis.B().
		Zadd().
		Key(redis.RecapReplayFromPrivateMessageBatch1.Format(message.From.ID)).
		ScoreMember().
		ScoreMember(
			float64(telegramChatHistory.ChattedAt),
			string(lo.Must(json.Marshal(telegramChatHistory))),
		).
		Build()

	err = m.redis.Do(context.Background(), zaddCmd).Error()
	if err != nil {
		return err
	}

	expireCmd := m.redis.B().
		Expire().
		Key(redis.RecapReplayFromPrivateMessageControl1.Format(message.From.ID)).
		Seconds(60 * 60 * 2).
		Build()

	err = m.redis.Do(context.Background(), expireCmd).Error()
	if err != nil {
		return err
	}

	expireCmd = m.redis.B().
		Expire().
		Key(redis.RecapReplayFromPrivateMessageBatch1.Format(message.From.ID)).
		Seconds(60 * 60 * 2).
		Build()

	err = m.redis.Do(context.Background(), expireCmd).Error()
	if err != nil {
		return err
	}

	m.logger.Debug("saved one telegram private forwarded replay chat history",
		zap.Int64("chat_id", telegramChatHistory.ChatID),
		zap.Int("message_id", telegramChatHistory.MessageID),
		zap.String("text", strings.ReplaceAll(telegramChatHistory.Text, "\n", " ")),
	)

	return nil
}

func (m *Model) FindPrivateForwardedChatHistories(userID int64) ([]telegramPrivateForwardedReplayChatHistory, error) {
	zrevrangeCmd := m.redis.B().
		Zrevrange().
		Key(redis.RecapReplayFromPrivateMessageBatch1.Format(userID)).
		Start(0).
		Stop(-1).
		Build()

	replayChatHistories, err := m.redis.Client.Do(context.Background(), zrevrangeCmd).AsStrSlice()
	if err != nil {
		if rueidis.IsRedisNil(err) {
			return make([]telegramPrivateForwardedReplayChatHistory, 0), nil
		}

		return make([]telegramPrivateForwardedReplayChatHistory, 0), err
	}
	if len(replayChatHistories) == 0 {
		return make([]telegramPrivateForwardedReplayChatHistory, 0), nil
	}

	replayChatHistories = lo.Reverse(replayChatHistories)

	telegramPrivateForwardedReplayChatHistories := make([]telegramPrivateForwardedReplayChatHistory, 0, len(replayChatHistories))

	for _, h := range replayChatHistories {
		var v telegramPrivateForwardedReplayChatHistory

		err = json.Unmarshal([]byte(h), &v)
		if err != nil {
			return make([]telegramPrivateForwardedReplayChatHistory, 0), err
		}

		telegramPrivateForwardedReplayChatHistories = append(telegramPrivateForwardedReplayChatHistories, v)
	}

	return telegramPrivateForwardedReplayChatHistories, nil
}

func (m *Model) SummarizePrivateForwardedChatHistories(userID int64, histories []telegramPrivateForwardedReplayChatHistory) ([]string, error) {
	historiesLLMFriendly := make([]string, 0, len(histories))
	historiesIncludedMessageIDs := make([]int64, 0)

	for _, message := range histories {
		historiesLLMFriendly = append(historiesLLMFriendly, fmt.Sprintf(
			"msgId:%d: %s sent: %s",
			message.MessageID,
			formatFullNameAndUsername(message.ActorDisplayName, message.ActorUsername),
			message.Text,
		))

		historiesIncludedMessageIDs = append(historiesIncludedMessageIDs, int64(message.MessageID))
	}

	chatHistories := strings.Join(historiesLLMFriendly, "\n")

	summarizations, statusUsage, err := m.summarizeChatHistories(userID, historiesIncludedMessageIDs, chatHistories)
	if err != nil {
		return make([]string, 0), err
	}

	summarizations = lo.Map(summarizations, func(item *openai.ChatHistorySummarizationOutputs, _ int) *openai.ChatHistorySummarizationOutputs {
		item.SinceID = 0
		item.Discussion = lo.Map(item.Discussion, func(item *openai.ChatHistorySummarizationOutputsDiscussion, _ int) *openai.ChatHistorySummarizationOutputsDiscussion {
			item.KeyIDs = nil
			return item
		})

		return item
	})

	ss, err := m.renderRecapTemplates(0, telegram.ChatTypePrivate, summarizations)
	if err != nil {
		return make([]string, 0), err
	}

	err = m.ent.LogChatHistoriesRecap.
		Create().
		SetChatID(userID).
		SetRecapInputs(chatHistories).
		SetRecapOutputs(strings.Join(ss, "\n")).
		SetCompletionTokenUsage(statusUsage.CompletionTokens).
		SetPromptTokenUsage(statusUsage.PromptTokens).
		SetTotalTokenUsage(statusUsage.TotalTokens).
		SetFromPlatform(int(FromPlatformTelegram)).
		SetRecapType(int(RecapTypeForPrivateForwarded)).
		Exec(context.Background())
	if err != nil {
		return make([]string, 0), err
	}

	return ss, nil
}
