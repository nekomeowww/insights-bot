package chathistories

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"text/template"
	"time"
	"unicode/utf16"
	"unicode/utf8"

	"entgo.io/ent/dialect/sql"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/samber/lo"
	lop "github.com/samber/lo/parallel"
	"github.com/sirupsen/logrus"
	"go.uber.org/fx"

	"github.com/nekomeowww/insights-bot/ent"
	"github.com/nekomeowww/insights-bot/ent/chathistories"
	"github.com/nekomeowww/insights-bot/internal/datastore"
	"github.com/nekomeowww/insights-bot/pkg/bots/tgbot"
	"github.com/nekomeowww/insights-bot/pkg/linkprev"
	"github.com/nekomeowww/insights-bot/pkg/logger"
	"github.com/nekomeowww/insights-bot/pkg/openai"
	"github.com/nekomeowww/insights-bot/pkg/utils"
)

type FromPlatform int

const (
	FromPlatformTelegram FromPlatform = iota
)

type NewModelParams struct {
	fx.In

	Lifecycle fx.Lifecycle

	Logger *logger.Logger
	Ent    *datastore.Ent
	OpenAI openai.Client
}

type Model struct {
	logger   *logger.Logger
	ent      *datastore.Ent
	openAI   openai.Client
	linkprev *linkprev.Client
}

func NewModel() func(NewModelParams) (*Model, error) {
	return func(param NewModelParams) (*Model, error) {
		return &Model{
			logger:   param.Logger,
			ent:      param.Ent,
			openAI:   param.OpenAI,
			linkprev: linkprev.NewClient(),
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
				m.logger.Errorf("🔗Failed to generate link preview for %s, error %+v", href, err)
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
				m.logger.Errorf("🔗Failed to summarize title for %s, error %+v", href, err)
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
		resp, err := m.openAI.SummarizeWithOneChatHistory(context.Background(), text)
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

func (m *Model) SaveOneTelegramChatHistory(message *tgbotapi.Message) error {
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

	telegramChatHistoryCreate := m.ent.ChatHistories.
		Create().
		SetChatID(message.Chat.ID).
		SetChatTitle(message.Chat.Title).
		SetMessageID(int64(message.MessageID)).
		SetUserID(message.From.ID).
		SetUsername(message.From.UserName).
		SetFullName(tgbot.FullNameFromFirstAndLastName(message.From.FirstName, message.From.LastName)).
		SetChattedAt(time.Unix(int64(message.Date), 0).UnixMilli()).
		SetFromPlatform(int(FromPlatformTelegram))

	if message.ForwardFrom != nil {
		telegramChatHistoryCreate.SetText(fmt.Sprintf("[forwarded from %s]: %s", tgbot.FullNameFromFirstAndLastName(message.ForwardFrom.FirstName, message.ForwardFrom.LastName), text))
	} else if message.ForwardFromChat != nil {
		telegramChatHistoryCreate.SetText(fmt.Sprintf("[forwarded from %s]: %s", message.ForwardFromChat.Title, text))
	} else {
		telegramChatHistoryCreate.SetText(text)
	}
	if message.ReplyToMessage != nil {
		repliedToText, err := m.extractTextWithSummarization(message.ReplyToMessage)
		if err != nil {
			return err
		}
		if repliedToText != "" {
			telegramChatHistoryCreate.SetRepliedToMessageID(int64(message.ReplyToMessage.MessageID))
			telegramChatHistoryCreate.SetRepliedToUserID(message.ReplyToMessage.From.ID)
			telegramChatHistoryCreate.SetRepliedToFullName(tgbot.FullNameFromFirstAndLastName(message.ReplyToMessage.From.FirstName, message.ReplyToMessage.From.LastName))
			telegramChatHistoryCreate.SetRepliedToUsername(message.ReplyToMessage.From.UserName)
			telegramChatHistoryCreate.SetRepliedToText(repliedToText)
		}
	}

	telegramChatHistory, err := telegramChatHistoryCreate.Save(context.TODO())
	if err != nil {
		return err
	}

	m.logger.WithFields(logrus.Fields{
		"id":         telegramChatHistory.ID,
		"chat_id":    telegramChatHistory.ChatID,
		"message_id": telegramChatHistory.MessageID,
		"text":       strings.ReplaceAll(telegramChatHistory.Text, "\n", " "),
	}).Debug("saved one telegram chat history")

	return nil
}

func (m *Model) FindLastOneHourChatHistories(chatID int64) ([]*ent.ChatHistories, error) {
	return m.FindChatHistoriesByTimeBefore(chatID, time.Hour)
}

func (m *Model) FindLastSixHourChatHistories(chatID int64) ([]*ent.ChatHistories, error) {
	return m.FindChatHistoriesByTimeBefore(chatID, 6*time.Hour)
}

func (m *Model) FindChatHistoriesByTimeBefore(chatID int64, before time.Duration) ([]*ent.ChatHistories, error) {
	m.logger.Infof("querying chat histories for %d", chatID)

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

type RecapOutputTemplateInputs struct {
	ChatID string
	Recap  *openai.ChatHistorySummarizationOutputs
}

func formatChatID(chatID int64) string {
	chatIDStr := strconv.FormatInt(chatID, 10)
	if strings.HasPrefix(chatIDStr, "-100") {
		return strings.TrimPrefix(chatIDStr, "-100")
	}

	return chatIDStr
}

var RecapOutputTemplate = lo.Must(template.
	New("recap output markdown template").
	Funcs(template.FuncMap{
		"join":   strings.Join,
		"sub":    func(a, b int) int { return a - b },
		"add":    func(a, b int) int { return a + b },
		"escape": tgbot.EscapeHTMLSymbols,
	}).
	Parse(`{{ $chatID := .ChatID }}{{ if .Recap.SinceID }}## <a href="https://t.me/c/{{ $chatID }}/{{ .Recap.SinceID }}">{{ escape .Recap.TopicName }}</a>{{ else }}## {{ escape .Recap.TopicName }}{{ end }}
参与人：{{ join .Recap.ParticipantsNamesWithoutUsername "，" }}
讨论：{{ range $di, $d := .Recap.Discussion }}
 - {{ escape $d.Point }}{{ if len $d.KeyIDs }} {{ range $cIndex, $c := $d.KeyIDs }}<a href="https://t.me/c/{{ $chatID }}/{{ $c }}">[{{ add $cIndex 1 }}]</a>{{ if not (eq $cIndex (sub (len $d.KeyIDs) 1)) }} {{ end }}{{ end }}{{ end }}{{ end }}{{ if .Recap.Conclusion }}
结论：{{ escape .Recap.Conclusion }}{{ end }}`))

func (m *Model) summarizeChatHistoriesSlice(s string) ([]*openai.ChatHistorySummarizationOutputs, error) {
	m.logger.Infof("✍️ summarizing last one hour chat histories:\n%s", s)

	resp, err := m.openAI.SummarizeWithChatHistories(context.Background(), s)
	if err != nil {
		return nil, err
	}
	if len(resp.Choices) == 0 {
		return nil, nil
	}

	m.logger.WithFields(logrus.Fields{
		"prompt_token_usage":     resp.Usage.PromptTokens,
		"completion_token_usage": resp.Usage.CompletionTokens,
		"total_token_usage":      resp.Usage.TotalTokens,
	}).Info("✅ summarized last one hour chat histories")
	if resp.Choices[0].Message.Content == "" {
		return nil, nil
	}

	var outputs []*openai.ChatHistorySummarizationOutputs

	err = json.Unmarshal([]byte(resp.Choices[0].Message.Content), &outputs)
	if err != nil {
		m.logger.Errorf("failed to unmarshal chat history summarization output: %s", resp.Choices[0].Message.Content)
		return nil, err
	}

	m.logger.Infof("✅ unmarshaled chat history summarization output: %s", utils.SprintJSON(outputs))

	return outputs, nil
}

func (m *Model) SummarizeChatHistories(chatID int64, histories []*ent.ChatHistories) ([]string, error) {
	historiesLLMFriendly := make([]string, 0, len(histories))

	for _, message := range histories {
		if message.RepliedToMessageID == 0 {
			historiesLLMFriendly = append(historiesLLMFriendly, fmt.Sprintf(
				"msgId:%d: %s sent: %s",
				message.MessageID,
				formatFullNameAndUsername(message.FullName, message.Username),
				message.Text,
			))
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
		}
	}

	chatHistories := strings.Join(historiesLLMFriendly, "\n")
	chatHistoriesSlices := m.openAI.SplitContentBasedByTokenLimitations(chatHistories, 2800)
	chatHistoriesSummarizations := make([]*openai.ChatHistorySummarizationOutputs, 0, len(chatHistoriesSlices))

	for _, s := range chatHistoriesSlices {
		var outputs []*openai.ChatHistorySummarizationOutputs

		_, _, err := lo.AttemptWithDelay(3, time.Second, func(tried int, delay time.Duration) error {
			o, err := m.summarizeChatHistoriesSlice(s)
			if err != nil {
				m.logger.Errorf("failed to summarize chat histories slice: %s, tried %d...", s, tried)
				return err
			}

			outputs = o
			return nil
		})
		if err != nil {
			return make([]string, 0), err
		}
		if outputs == nil {
			continue
		}

		// filter out empty outputs
		outputs = lo.Filter(outputs, func(item *openai.ChatHistorySummarizationOutputs, _ int) bool {
			return item != nil &&
				item.TopicName != "" && // filter out empty topic name
				item.SinceID != 0 && // filter out empty since id
				len(item.ParticipantsNamesWithoutUsername) > 0 && // filter out empty participants
				len(item.Discussion) > 0 // filter out empty discussion
		})

		// limit key ids to 5
		for _, o := range outputs {
			for _, d := range o.Discussion {
				d.KeyIDs = lo.UniqBy(d.KeyIDs, func(item int64) int64 {
					return item
				})
				d.KeyIDs = lo.Filter(d.KeyIDs, func(item int64, _ int) bool {
					return item != 0
				})

				if len(d.KeyIDs) > 5 {
					d.KeyIDs = d.KeyIDs[:5]
				}
			}
		}

		chatHistoriesSummarizations = append(chatHistoriesSummarizations, outputs...)
	}

	ss := make([]string, 0)

	for _, r := range chatHistoriesSummarizations {
		sb := new(strings.Builder)

		err := RecapOutputTemplate.Execute(sb, RecapOutputTemplateInputs{
			ChatID: formatChatID(chatID),
			Recap:  r,
		})
		if err != nil {
			return make([]string, 0), err
		}

		ss = append(ss, sb.String())
	}

	err := m.ent.LogChatHistoriesRecap.
		Create().
		SetChatID(chatID).
		SetRecapInputs(chatHistories).
		SetRecapOutputs(strings.Join(ss, "\n")).
		Exec(context.Background())
	if err != nil {
		return make([]string, 0), err
	}

	return ss, nil
}
