package chat_histories

import (
	"os"
	"strings"
	"testing"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/ostafen/clover/v2"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nekomeowww/insights-bot/internal/datastore"
	"github.com/nekomeowww/insights-bot/internal/lib"
	"github.com/nekomeowww/insights-bot/pkg/openai"
	"github.com/nekomeowww/insights-bot/pkg/types/chat_history"
	"github.com/nekomeowww/insights-bot/pkg/utils"
)

var model *ChatHistoriesModel

func TestMain(m *testing.M) {
	logger := lib.NewLogger()()

	db, cancel := datastore.NewTestClover()()
	defer cancel()

	var err error
	model, err = NewChatHistoriesModel()(NewChatHistoriesModelParam{
		Clover: db,
		Logger: logger,
	})
	if err != nil {
		panic(err)
	}

	os.Exit(m.Run())
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

	query := clover.
		NewQuery(chat_history.TelegramChatHistory{}.CollectionName()).
		Where(clover.Field("chat_id").Eq(message.Chat.ID)).
		Where(clover.Field("message_id").Eq(message.MessageID))

	doc, err := model.Clover.FindFirst(query)
	require.NoError(err)
	require.NotNil(doc)

	var chatHistory chat_history.TelegramChatHistory
	err = doc.Unmarshal(&chatHistory)
	require.NoError(err)

	assert.Equal(message.Chat.ID, chatHistory.ChatID)
	assert.Equal(message.MessageID, chatHistory.MessageID)
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

	assert.Equal([]int{1, 2, 3}, lo.Map(histories, func(item *chat_history.TelegramChatHistory, _ int) int {
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
				SinceMsgID:                       1,
				ParticipantsNamesWithoutUsername: []string{"User 1", "User 2"},
				Discussion: []*openai.ChatHistorySummarizationOutputsDiscussion{
					{
						Point:              "Point 1",
						CriticalMessageIDs: []int64{1, 2},
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
						Point:              "Point 2",
						CriticalMessageIDs: []int64{1, 2},
					},
				},
			},
			{
				TopicName:                        "Topic 1",
				SinceMsgID:                       2,
				ParticipantsNamesWithoutUsername: []string{"User 1", "User 2"},
				Discussion: []*openai.ChatHistorySummarizationOutputsDiscussion{
					{
						Point:              "Point 1",
						CriticalMessageIDs: []int64{1, 2},
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
