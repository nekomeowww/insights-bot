package chathistories

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nekomeowww/insights-bot/ent"
	"github.com/nekomeowww/insights-bot/ent/chathistories"
	"github.com/nekomeowww/insights-bot/internal/configs"
	"github.com/nekomeowww/insights-bot/internal/datastore"
	"github.com/nekomeowww/insights-bot/internal/lib"
	"github.com/nekomeowww/insights-bot/pkg/openai"
	"github.com/nekomeowww/insights-bot/pkg/tutils"
	"github.com/nekomeowww/insights-bot/pkg/utils"
)

var model *Model

func TestMain(m *testing.M) {
	logger := lib.NewLogger()()

	ent, err := datastore.NewEnt()(datastore.NewEntParams{
		Lifecycle: tutils.NewEmtpyLifecycle(),
		Configs:   configs.NewTestConfig()(),
	})
	if err != nil {
		panic(err)
	}

	model, err = NewModel()(NewModelParams{
		Ent:    ent,
		Logger: logger,
	})
	if err != nil {
		panic(err)
	}

	os.Exit(m.Run())
}

func TestExtractTextFromMessage(t *testing.T) {
	t.Run("MixedUrlsAndTextLinks", func(t *testing.T) {
		assert := assert.New(t)
		message := &tgbotapi.Message{
			MessageID: 666,
			From:      &tgbotapi.User{ID: 23333333},
			Date:      1683386000,
			Chat:      &tgbotapi.Chat{ID: 0xc0001145e4},
			Text:      "看看这些链接：https://docs.swift.org/swift-book/documentation/the-swift-programming-language/stringsandcharacters/#Extended-Grapheme-Clusters 、https://www.youtube.com/watch?v=outcGtbnMuQ https://github.com/nekomeowww/insights-bot 还有 这个",
			Entities: []tgbotapi.MessageEntity{
				{Type: "url", Offset: 7, Length: 127, URL: "", Language: ""},
				{Type: "url", Offset: 136, Length: 43, URL: "", Language: ""},
				{Type: "url", Offset: 180, Length: 42, URL: "", Language: ""},
				{Type: "text_link", Offset: 226, Length: 2, URL: "https://matters.town/@1435Club/322889-%E8%BF%99%E5%87%A0%E5%A4%A9-web3%E5%9C%A8%E5%A4%A7%E7%90%86%E5%8F%91%E7%94%9F%E4%BA%86%E4%BB%80%E4%B9%88", Language: ""},
			},
			Photo: []tgbotapi.PhotoSize{},
		}
		expect := "看看这些链接：[Documentation](https://docs.swift.org/swift-book/documentation/the-swift-programming-language/stringsandcharacters/#Extended-Grapheme-Clusters) 、[GPT-4 Developer Livestream](https://www.youtube.com/watch?v=outcGtbnMuQ) [GitHub - nekomeowww/insights-bot: A bot works with OpenAI GPT models to provide insights for your info flows.](https://github.com/nekomeowww/insights-bot) 还有 [这个](https://matters.town/@1435Club/322889-这几天-web3在大理发生了什么)"
		assert.Equal(expect, model.ExtractTextFromMessage(message))
	})
}

func TestSaveOneTelegramChatHistory(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	message := &tgbotapi.Message{
		MessageID: int(utils.RandomInt64()),
		From: &tgbotapi.User{
			ID:        utils.RandomInt64(),
			FirstName: utils.RandomHashString(5),
			UserName:  utils.RandomHashString(10),
		},
		Chat: &tgbotapi.Chat{
			ID: utils.RandomInt64(),
		},
		Date: int(time.Now().Unix()),
		Text: utils.RandomHashString(10),
	}
	err := model.SaveOneTelegramChatHistory(message)
	require.NoError(err)

	chatHistory, err := model.ent.ChatHistories.
		Query().
		Where(
			chathistories.ChatID(message.Chat.ID),
			chathistories.MessageID(int64(message.MessageID)),
		).
		First(context.Background())
	require.NoError(err)
	require.NotNil(chatHistory)

	assert.Equal(message.Chat.ID, chatHistory.ChatID)
	assert.Equal(int64(message.MessageID), chatHistory.MessageID)
	assert.Equal(message.From.ID, chatHistory.UserID)
	assert.Equal(message.From.FirstName, chatHistory.FullName)
	assert.Equal(message.From.UserName, chatHistory.Username)
	assert.Equal(message.Text, chatHistory.Text)
	assert.Equal(time.Unix(int64(message.Date), 0).UnixMilli(), chatHistory.ChattedAt)
}

func TestFindLastOneHourChatHistories(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	chatID := utils.RandomInt64()

	message1 := &tgbotapi.Message{
		MessageID: 1,
		From: &tgbotapi.User{
			ID:        utils.RandomInt64(),
			FirstName: utils.RandomHashString(5),
			UserName:  utils.RandomHashString(10),
		},
		Chat: &tgbotapi.Chat{ID: chatID},
		Date: int(time.Now().Unix()),
		Text: utils.RandomHashString(10),
	}

	message2 := &tgbotapi.Message{
		MessageID: 2,
		From: &tgbotapi.User{
			ID:        utils.RandomInt64(),
			FirstName: utils.RandomHashString(5),
			UserName:  utils.RandomHashString(10),
		},
		Chat: &tgbotapi.Chat{ID: chatID},
		Date: int(time.Now().Unix()),
		Text: utils.RandomHashString(10),
	}

	message3 := &tgbotapi.Message{
		MessageID: 3,
		From: &tgbotapi.User{
			ID:        utils.RandomInt64(),
			FirstName: utils.RandomHashString(5),
			UserName:  utils.RandomHashString(10),
		},
		Chat: &tgbotapi.Chat{ID: chatID},
		Date: int(time.Now().Unix()),
		Text: utils.RandomHashString(10),
	}

	err := model.SaveOneTelegramChatHistory(message1)
	require.NoError(err)

	err = model.SaveOneTelegramChatHistory(message2)
	require.NoError(err)

	err = model.SaveOneTelegramChatHistory(message3)
	require.NoError(err)

	histories, err := model.FindLastOneHourChatHistories(chatID)
	require.NoError(err)
	require.Len(histories, 3)

	assert.Equal([]int{1, 2, 3}, lo.Map(histories, func(item *ent.ChatHistories, _ int) int64 {
		return item.MessageID
	}))
}

func TestRecapOutputTemplateExecute(t *testing.T) {
	sb := new(strings.Builder)
	err := RecapOutputTemplate.Execute(sb, RecapOutputTemplateInputs{
		ChatID: formatChatID(-100123456789),
		Recaps: []*openai.ChatHistorySummarizationOutputs{
			{
				TopicName:                        "Topic 1",
				SinceID:                          1,
				ParticipantsNamesWithoutUsername: []string{"User 1", "User 2"},
				Discussion: []*openai.ChatHistorySummarizationOutputsDiscussion{
					{
						Point:       "Point 1",
						CriticalIDs: []int64{1, 2},
					},
					{
						Point: "Point 2",
					},
				},
				Conclusion: "Conclusion 1",
			},
			{
				TopicName:                        "Topic 3",
				ParticipantsNamesWithoutUsername: []string{"User 1", "User 2"},
				Discussion: []*openai.ChatHistorySummarizationOutputsDiscussion{
					{
						Point: "Point 1",
					},
					{
						Point:       "Point 2",
						CriticalIDs: []int64{1, 2},
					},
				},
			},
			{
				TopicName:                        "Topic 1",
				SinceID:                          2,
				ParticipantsNamesWithoutUsername: []string{"User 1", "User 2"},
				Discussion: []*openai.ChatHistorySummarizationOutputsDiscussion{
					{
						Point:       "Point 1",
						CriticalIDs: []int64{1, 2},
					},
					{
						Point: "Point 2",
					},
				},
				Conclusion: "Conclusion 2",
			},
		},
	})
	require.NoError(t, err)
	expected := `## <a href="https://t.me/c/123456789/1">Topic 1</a>
参与人：User 1，User 2
讨论：
 - Point 1 <a href="https://t.me/c/123456789/1">[1]</a> <a href="https://t.me/c/123456789/2">[2]</a>
 - Point 2
结论：Conclusion 1

## Topic 3
参与人：User 1，User 2
讨论：
 - Point 1
 - Point 2 <a href="https://t.me/c/123456789/1">[1]</a> <a href="https://t.me/c/123456789/2">[2]</a>

## <a href="https://t.me/c/123456789/2">Topic 1</a>
参与人：User 1，User 2
讨论：
 - Point 1 <a href="https://t.me/c/123456789/1">[1]</a> <a href="https://t.me/c/123456789/2">[2]</a>
 - Point 2
结论：Conclusion 2`
	assert.Equal(t, expected, sb.String())
}

func TestFormatFullNameAndUsername(t *testing.T) {
	tests := []struct {
		name     string
		fullName string
		username string
		result   string
	}{
		{
			name:     `full name shorter than 10 chars`,
			fullName: "Full Name",
			username: "example_username",
			result:   "Full Name",
		},
		{
			name:     `full name longer than 10 chars`,
			fullName: "A Very Long Full Name",
			username: "example_username",
			result:   "example_username",
		},
		{
			name: `full name longer than 10 chars
			AND username is empty`,
			fullName: "A Very Long Full Name",
			username: "",
			result:   "A Very Long Full Name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatFullNameAndUsername(tt.fullName, tt.username)
			assert.Equal(t, tt.result, result)
		})
	}
}
