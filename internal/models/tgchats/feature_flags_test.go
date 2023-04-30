package tgchats

import (
	"context"
	"testing"

	"github.com/nekomeowww/insights-bot/ent/telegramchatfeatureflags"
	"github.com/nekomeowww/insights-bot/pkg/types/telegram"
	"github.com/nekomeowww/insights-bot/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChatHistoriesRecap(t *testing.T) {
	chatID := utils.RandomInt64()

	t.Run("Enable", func(t *testing.T) {
		assert := assert.New(t)
		require := require.New(t)

		err := model.EnableChatHistoriesRecap(chatID, telegram.ChatTypeGroup)
		require.NoError(err)

		featureFlag, err := model.ent.TelegramChatFeatureFlags.
			Query().
			Where(
				telegramchatfeatureflags.ChatID(chatID),
			).
			First(context.Background())
		require.NoError(err)
		require.NotNil(featureFlag)

		assert.True(featureFlag.FeatureChatHistoriesRecap)

		enabled, err := model.HasChatHistoriesRecapEnabled(chatID, telegram.ChatTypeGroup)
		require.NoError(err)
		assert.True(enabled)
	})

	t.Run("Disable", func(t *testing.T) {
		assert := assert.New(t)
		require := require.New(t)

		err := model.DisableChatHistoriesRecap(chatID, telegram.ChatTypeGroup)
		require.NoError(err)

		featureFlag, err := model.ent.TelegramChatFeatureFlags.
			Query().
			Where(
				telegramchatfeatureflags.ChatID(chatID),
			).
			First(context.Background())
		require.NoError(err)
		require.NotNil(featureFlag)
		assert.False(featureFlag.FeatureChatHistoriesRecap)

		enabled, err := model.HasChatHistoriesRecapEnabled(chatID, telegram.ChatTypeGroup)
		require.NoError(err)
		assert.False(enabled)
	})
}

func TestListChatHistoriesRecapEnabledChats(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	chatID1 := utils.RandomInt64()
	chatID2 := utils.RandomInt64()
	chatID3 := utils.RandomInt64()

	err := model.EnableChatHistoriesRecap(chatID1, telegram.ChatTypeGroup)
	require.NoError(err)

	err = model.EnableChatHistoriesRecap(chatID2, telegram.ChatTypeGroup)
	require.NoError(err)

	err = model.EnableChatHistoriesRecap(chatID3, telegram.ChatTypeGroup)
	require.NoError(err)

	chatIDs, err := model.ListChatHistoriesRecapEnabledChats()
	require.NoError(err)
	require.Len(chatIDs, 3)
	assert.ElementsMatch([]int64{chatID1, chatID2, chatID3}, chatIDs)
}
