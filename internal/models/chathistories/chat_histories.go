package chathistories

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"
	"unicode/utf16"
	"unicode/utf8"

	"entgo.io/ent/dialect/sql"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	"github.com/samber/lo"
	lop "github.com/samber/lo/parallel"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/nekomeowww/insights-bot/ent"
	"github.com/nekomeowww/insights-bot/ent/chathistories"
	"github.com/nekomeowww/insights-bot/internal/configs"
	"github.com/nekomeowww/insights-bot/internal/datastore"
	"github.com/nekomeowww/insights-bot/internal/thirdparty/openai"
	"github.com/nekomeowww/insights-bot/pkg/bots/tgbot"
	"github.com/nekomeowww/insights-bot/pkg/linkprev"
	"github.com/nekomeowww/insights-bot/pkg/logger"
	"github.com/nekomeowww/insights-bot/pkg/types/telegram"
)

type FromPlatform int

const (
	FromPlatformTelegram FromPlatform = iota
)

type RecapType int

const (
	RecapTypeForGroup RecapType = iota
	RecapTypeForPrivateForwarded
)

type NewModelParams struct {
	fx.In

	Lifecycle fx.Lifecycle

	Config *configs.Config
	Logger *logger.Logger
	Ent    *datastore.Ent
	OpenAI openai.Client
	Redis  *datastore.Redis
}

type Model struct {
	config   *configs.Config
	logger   *logger.Logger
	ent      *datastore.Ent
	openAI   openai.Client
	linkprev *linkprev.Client
	redis    *datastore.Redis
}

func NewModel() func(NewModelParams) (*Model, error) {
	return func(param NewModelParams) (*Model, error) {
		return &Model{
			config:   param.Config,
			logger:   param.Logger,
			ent:      param.Ent,
			openAI:   param.OpenAI,
			linkprev: linkprev.NewClient(),
			redis:    param.Redis,
		}, nil
	}
}

func (m *Model) ExtractTextFromMessage(message *tgbotapi.Message) string {
	text := lo.Ternary(message.Caption != "", message.Caption, message.Text)

	type MarkdownLink struct {
		Markdown []uint16
		Start    int
		End      int
	}

	textUTF16 := utf16.Encode([]rune(text))
	links := lop.Map(message.Entities, func(entity tgbotapi.MessageEntity, i int) MarkdownLink {
		startIndex := entity.Offset
		endIndex := startIndex + entity.Length
		var title string
		var href string

		switch entity.Type {
		case "url":
			href = string(utf16.Decode(textUTF16[startIndex:endIndex]))

			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			meta, err := m.linkprev.Preview(ctx, href)
			if err != nil {
				m.logger.Error("🔗Failed to generate link preview", zap.String("url", href), zap.Error(err))
				return MarkdownLink{[]uint16{}, -1, -1}
			}

			title = lo.Ternary(meta.Title != "", meta.Title, meta.OpenGraph.Title)
		case "text_link":
			href = entity.URL
			title = string(utf16.Decode(textUTF16[startIndex:endIndex]))
		default:
			return MarkdownLink{[]uint16{}, -1, -1}
		}

		if utf8.RuneCountInString(title) > 200 {
			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
			defer cancel()

			resp, err := m.openAI.SummarizeAny(ctx, title)
			if err != nil {
				m.logger.Error("🔗Failed to summarize title", zap.String("url", href), zap.Error(err), zap.String("title", title))
				return MarkdownLink{[]uint16{}, -1, -1}
			}
			if len(resp.Choices) != 0 {
				title = resp.Choices[0].Message.Content
			}
		}

		unescaped, err := url.QueryUnescape(href)
		if err == nil {
			href = unescaped
		}

		md := fmt.Sprintf("[%s](%s)", title, href)
		mdUTF16 := utf16.Encode([]rune(md))

		return MarkdownLink{mdUTF16, startIndex, endIndex}
	})

	for i := len(links) - 1; i >= 0; i-- {
		if links[i].Start == -1 {
			continue
		}

		temp := append(links[i].Markdown, textUTF16[links[i].End:]...)
		textUTF16 = append(textUTF16[:links[i].Start], temp...)
	}

	return string(utf16.Decode(textUTF16))
}

func (m *Model) extractTextWithSummarization(message *tgbotapi.Message) (string, error) {
	text := m.ExtractTextFromMessage(message)
	if text == "" {
		return "", nil
	}
	if utf8.RuneCountInString(text) >= 300 {
		resp, err := m.openAI.SummarizeOneChatHistory(context.Background(), text)
		if err != nil {
			return "", err
		}
		if len(resp.Choices) == 0 {
			return "", nil
		}

		return resp.Choices[0].Message.Content, nil
	}

	return text, nil
}

func (m *Model) extractTextFromMessage(message *tgbotapi.Message) (string, error) {
	if message == nil {
		return "", nil
	}
	if message.Text == "" && message.Caption == "" {
		m.logger.Warn("message text is empty")
		return "", nil
	}

	text, err := m.extractTextWithSummarization(message)
	if err != nil {
		return "", err
	}
	if text == "" {
		m.logger.Warn("message text is empty")
		return "", nil
	}

	return text, nil
}

func (m *Model) assignReplyMessageDataForChatHistory(entity *ent.ChatHistoriesCreate, message *tgbotapi.Message) error {
	if message.ReplyToMessage == nil {
		return nil
	}

	repliedToText, err := m.extractTextWithSummarization(message.ReplyToMessage)
	if err != nil {
		return err
	}
	if repliedToText != "" {
		entity.
			SetRepliedToMessageID(int64(message.ReplyToMessage.MessageID)).
			SetRepliedToUserID(message.ReplyToMessage.From.ID).
			SetRepliedToFullName(tgbot.FullNameFromFirstAndLastName(message.ReplyToMessage.From.FirstName, message.ReplyToMessage.From.LastName)).
			SetRepliedToUsername(message.ReplyToMessage.From.UserName).
			SetRepliedToText(repliedToText).
			SetRepliedToChatType(message.ReplyToMessage.Chat.Type)
	}

	return nil
}

func (m *Model) SaveOneTelegramChatHistory(message *tgbotapi.Message) error {
	text, err := m.extractTextFromMessage(message)
	if err != nil {
		return err
	}
	if text == "" {
		return nil
	}

	telegramChatHistoryCreate := m.ent.ChatHistories.
		Create().
		SetChatID(message.Chat.ID).
		SetChatType(message.Chat.Type).
		SetChatTitle(message.Chat.Title).
		SetMessageID(int64(message.MessageID)).
		SetUserID(message.From.ID).
		SetUsername(message.From.UserName).
		SetFullName(tgbot.FullNameFromFirstAndLastName(message.From.FirstName, message.From.LastName)).
		SetFromPlatform(int(FromPlatformTelegram)).
		SetChattedAt(time.Unix(int64(message.Date), 0).UnixMilli())

	if message.ForwardFrom != nil {
		telegramChatHistoryCreate.SetText(fmt.Sprintf("[forwarded from %s]: %s", tgbot.FullNameFromFirstAndLastName(message.ForwardFrom.FirstName, message.ForwardFrom.LastName), text))
	} else if message.ForwardFromChat != nil {
		telegramChatHistoryCreate.SetText(fmt.Sprintf("[forwarded from %s]: %s", message.ForwardFromChat.Title, text))
	} else {
		telegramChatHistoryCreate.SetText(text)
	}

	err = m.assignReplyMessageDataForChatHistory(telegramChatHistoryCreate, message)
	if err != nil {
		return err
	}

	telegramChatHistory, err := telegramChatHistoryCreate.Save(context.TODO())
	if err != nil {
		return err
	}

	m.logger.Debug("saved one telegram chat history",
		zap.String("id", telegramChatHistory.ID.String()),
		zap.Int64("chat_id", telegramChatHistory.ChatID),
		zap.Int64("message_id", telegramChatHistory.MessageID),
		zap.String("text", strings.ReplaceAll(telegramChatHistory.Text, "\n", " ")),
	)

	return nil
}

func (m *Model) UpdateOneTelegramChatHistory(message *tgbotapi.Message) error {
	if message == nil {
		return nil
	}
	if message.Text == "" && message.Caption == "" {
		m.logger.Warn("message text is empty")
		return nil
	}

	text, err := m.extractTextWithSummarization(message)
	if err != nil {
		return err
	}
	if text == "" {
		m.logger.Warn("message text is empty")
		return nil
	}

	err = m.ent.ChatHistories.
		Update().
		Where(
			chathistories.ChatID(message.Chat.ID),
			chathistories.MessageID(int64(message.MessageID)),
		).
		SetText(text).
		Exec(context.Background())
	if err != nil {
		return err
	}

	m.logger.Debug("updated one message",
		zap.Int64("chat_id", message.Chat.ID),
		zap.Int("message_id", message.MessageID),
		zap.String("text", strings.ReplaceAll(text, "\n", " ")),
	)

	return nil
}

func (m *Model) DeleteAllChatHistoriesByChatID(chatID int64) error {
	_, err := m.ent.ChatHistories.
		Delete().
		Where(chathistories.ChatID(chatID)).
		Exec(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func (m *Model) MigrateChatHistoriesOfChatFromChatIDToChatID(fromChatID, toChatID int64) error {
	affectedRows, err := m.ent.ChatHistories.
		Update().
		Where(chathistories.ChatID(fromChatID)).
		SetChatID(toChatID).
		SetChatType(string(telegram.ChatTypeSuperGroup)).
		Save(context.Background())
	if err != nil {
		return err
	}

	m.logger.Info("successfully migrated options of chat",
		zap.Int64("from_chat_id", fromChatID),
		zap.Int64("to_chat_id", toChatID),
		zap.Int("affected_rows", affectedRows),
	)

	return nil
}

func (m *Model) FindLastOneHourChatHistories(chatID int64) ([]*ent.ChatHistories, error) {
	return m.FindChatHistoriesByTimeBefore(chatID, time.Hour)
}

func (m *Model) FindLast6HourChatHistories(chatID int64) ([]*ent.ChatHistories, error) {
	return m.FindChatHistoriesByTimeBefore(chatID, 6*time.Hour)
}

func (m *Model) FindLast8HourChatHistories(chatID int64) ([]*ent.ChatHistories, error) {
	return m.FindChatHistoriesByTimeBefore(chatID, 8*time.Hour)
}

func (m *Model) FindLast12HourChatHistories(chatID int64) ([]*ent.ChatHistories, error) {
	return m.FindChatHistoriesByTimeBefore(chatID, 12*time.Hour)
}

func (m *Model) FindChatHistoriesByTimeBefore(chatID int64, before time.Duration) ([]*ent.ChatHistories, error) {
	m.logger.Info("querying chat histories", zap.Int64("chat_id", chatID))

	telegramChatHistories, err := m.ent.ChatHistories.
		Query().
		Where(
			chathistories.ChatID(chatID),
			chathistories.ChattedAtGT(time.Now().Add(-before).UnixMilli()),
		).
		Order(
			chathistories.ByMessageID(sql.OrderAsc()),
		).
		All(context.TODO())
	if err != nil {
		return make([]*ent.ChatHistories, 0), err
	}

	return telegramChatHistories, nil
}

func formatFullNameAndUsername(fullName, username string) string {
	if utf8.RuneCountInString(fullName) >= 10 && username != "" {
		return username
	}

	return strings.ReplaceAll(fullName, "#", "")
}

func (m *Model) encodeMessageIDIntoVirtualMessageID(histories []*ent.ChatHistories) map[int64]int64 {
	virtualMessageID := int64(1)
	mMessageIDToVirtualMessageID := make(map[int64]int64)

	for _, message := range histories {
		mMessageIDToVirtualMessageID[virtualMessageID] = message.MessageID
		message.MessageID = virtualMessageID
		virtualMessageID++

		if message.RepliedToMessageID != 0 {
			mMessageIDToVirtualMessageID[virtualMessageID] = message.RepliedToMessageID
			message.RepliedToMessageID = virtualMessageID
			virtualMessageID++
		}
	}

	return mMessageIDToVirtualMessageID
}

func (m *Model) decodeMessageIDFromVirtualMessageID(mMessageIDToVirtualMessageID map[int64]int64, outputs []*openai.ChatHistorySummarizationOutputs) {
	for _, o := range outputs {
		for _, d := range o.Discussion {
			d.KeyIDs = lo.Map(d.KeyIDs, func(virtualMessageID int64, i int) int64 {
				return mMessageIDToVirtualMessageID[virtualMessageID]
			})
		}

		o.SinceID = mMessageIDToVirtualMessageID[o.SinceID]
	}
}

func (m *Model) SummarizeChatHistories(chatID int64, chatType telegram.ChatType, histories []*ent.ChatHistories) (uuid.UUID, []string, error) {
	historiesLLMFriendly := make([]string, 0, len(histories))
	historiesIncludedMessageIDs := make([]int64, 0)

	mMessageIDToVirtualMessageID := m.encodeMessageIDIntoVirtualMessageID(histories)

	for _, message := range histories {
		if message.RepliedToMessageID == 0 {
			historiesLLMFriendly = append(historiesLLMFriendly, fmt.Sprintf(
				"msgId:%d: %s sent: %s",
				message.MessageID,
				formatFullNameAndUsername(message.FullName, message.Username),
				message.Text,
			))

			historiesIncludedMessageIDs = append(historiesIncludedMessageIDs, message.MessageID)
		} else {
			repliedToPartialContextMessage := fmt.Sprintf(
				"%s sent msgId:%d",
				formatFullNameAndUsername(message.RepliedToFullName, message.RepliedToUsername),
				message.RepliedToMessageID,
			)
			historiesLLMFriendly = append(historiesLLMFriendly, fmt.Sprintf(
				"msgId:%d: %s replying to [%s]: %s",
				message.MessageID,
				formatFullNameAndUsername(message.FullName, message.Username),
				repliedToPartialContextMessage,
				message.Text,
			))

			historiesIncludedMessageIDs = append(historiesIncludedMessageIDs, message.MessageID)
		}
	}

	chatHistories := strings.Join(historiesLLMFriendly, "\n")

	summarizations, statusUsage, err := m.summarizeChatHistories(chatID, historiesIncludedMessageIDs, chatHistories)
	if err != nil {
		return uuid.Nil, make([]string, 0), err
	}

	// reverse virtual message id to real message id
	m.decodeMessageIDFromVirtualMessageID(mMessageIDToVirtualMessageID, summarizations)

	ss, err := m.renderRecapTemplates(chatID, chatType, summarizations)
	if err != nil {
		return uuid.Nil, make([]string, 0), err
	}

	saved, err := m.ent.LogChatHistoriesRecap.
		Create().
		SetChatID(chatID).
		SetRecapInputs(chatHistories).
		SetRecapOutputs(strings.Join(ss, "\n")).
		SetCompletionTokenUsage(statusUsage.CompletionTokens).
		SetPromptTokenUsage(statusUsage.PromptTokens).
		SetTotalTokenUsage(statusUsage.TotalTokens).
		SetFromPlatform(int(FromPlatformTelegram)).
		SetRecapType(int(RecapTypeForGroup)).
		SetModelName(m.openAI.GetModelName()).
		Save(context.Background())
	if err != nil {
		return uuid.Nil, make([]string, 0), err
	}

	return saved.ID, ss, nil
}
