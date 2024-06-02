package chathistories

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/nekomeowww/insights-bot/ent/sentmessages"
	"github.com/nekomeowww/xo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestSaveOneTelegramSentMessage(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	message := &tgbotapi.Message{
		MessageID: int(xo.RandomInt64()),
		From: &tgbotapi.User{
			ID:        xo.RandomInt64(),
			FirstName: xo.RandomHashString(5),
			UserName:  xo.RandomHashString(10),
		},
		Chat: &tgbotapi.Chat{
			ID: xo.RandomInt64(),
		},
		Date: int(time.Now().Unix()),
		Text: xo.RandomHashString(10),
	}
	isPinned := false

	err := model.SaveOneTelegramSentMessage(message, isPinned)
	require.NoError(err)

	sentMessage, err := model.ent.SentMessages.
		Query().
		Where(
			sentmessages.ChatID(message.Chat.ID),
			sentmessages.MessageID(message.MessageID),
		).
		First(context.Background())
	require.NoError(err)
	require.NotNil(sentMessage)

	assert.Equal(message.Chat.ID, sentMessage.ChatID)
	assert.Equal(message.MessageID, sentMessage.MessageID)
	assert.Equal(message.Text, sentMessage.Text)
	assert.Equal(isPinned, sentMessage.IsPinned)
	assert.Equal(int(autoRecapMessage), sentMessage.MessageType)
}

func TestUpdatePinnedMessage(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	chatID := xo.RandomInt64()
	messageID := int(xo.RandomInt64())
	isPinned := true

	// Save a message first
	message := &tgbotapi.Message{
		MessageID: messageID,
		Chat: &tgbotapi.Chat{
			ID: chatID,
		},
	}
	err := model.SaveOneTelegramSentMessage(message, isPinned)
	require.NoError(err)

	// Update the pinned status of the message
	isPinned = false
	err = model.UpdatePinnedMessage(chatID, messageID, isPinned)
	require.NoError(err)

	// Retrieve the message and check the pinned status
	sentMessage, err := model.ent.SentMessages.
		Query().
		Where(
			sentmessages.ChatID(chatID),
			sentmessages.MessageID(messageID),
		).
		First(context.Background())
	require.NoError(err)
	require.NotNil(sentMessage)

	assert.Equal(isPinned, sentMessage.IsPinned)
}
