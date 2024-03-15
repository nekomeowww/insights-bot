package tgchats

import (
	"context"
	"testing"

	"github.com/nekomeowww/insights-bot/ent"
	"github.com/nekomeowww/insights-bot/ent/telegramchatfeatureflags"
	"github.com/nekomeowww/insights-bot/pkg/types/telegram"
	"github.com/nekomeowww/xo"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChatHistoriesRecap(t *testing.T) {
	chatID := xo.RandomInt64()

	t.Run("Enable", func(t *testing.T) {
		assert := assert.New(t)
		require := require.New(t)

		err := model.EnableChatHistoriesRecapForGroups(chatID, telegram.ChatTypeGroup, xo.RandomHashString(6))
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

		enabled, err := model.HasChatHistoriesRecapEnabledForGroups(chatID, "")
		require.NoError(err)
		assert.True(enabled)
	})

	t.Run("Disable", func(t *testing.T) {
		assert := assert.New(t)
		require := require.New(t)

		err := model.DisableChatHistoriesRecapForGroups(chatID, telegram.ChatTypeGroup, xo.RandomHashString(6))
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

		enabled, err := model.HasChatHistoriesRecapEnabledForGroups(chatID, "")
		require.NoError(err)
		assert.False(enabled)
	})
}

func TestListChatHistoriesRecapEnabledChats(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	chatID1 := xo.RandomInt64()
	chatTitle1 := xo.RandomHashString(6)
	chatID2 := xo.RandomInt64()
	chatTitle2 := xo.RandomHashString(6)
	chatID3 := xo.RandomInt64()
	chatTitle3 := xo.RandomHashString(6)

	err := model.EnableChatHistoriesRecapForGroups(chatID1, telegram.ChatTypeGroup, chatTitle1)
	require.NoError(err)

	err = model.EnableChatHistoriesRecapForGroups(chatID2, telegram.ChatTypeGroup, chatTitle2)
	require.NoError(err)

	err = model.EnableChatHistoriesRecapForGroups(chatID3, telegram.ChatTypeGroup, chatTitle3)
	require.NoError(err)

	defer func() {
		_, err := model.ent.TelegramChatFeatureFlags.Delete().Exec(context.Background())
		assert.NoError(err)
	}()

	chats, err := model.ListChatHistoriesRecapEnabledChatsForGroups()
	require.NoError(err)
	require.Len(chats, 3)
	assert.ElementsMatch([]int64{chatID1, chatID2, chatID3}, lo.Map(chats, func(item *ent.TelegramChatFeatureFlags, _ int) int64 { return item.ChatID }))
}
