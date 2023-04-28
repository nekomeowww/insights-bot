package tgchats

import (
	"testing"

	"github.com/nekomeowww/insights-bot/pkg/types/telegram"
	"github.com/nekomeowww/insights-bot/pkg/types/telegram/tgchat"
	"github.com/nekomeowww/insights-bot/pkg/utils"
	"github.com/ostafen/clover/v2"
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

		query := clover.
			NewQuery(tgchat.FeatureFlag{}.CollectionName()).
			Where(clover.Field("chat_id").Eq(chatID))

		doc, err := model.Clover.FindFirst(query)
		require.NoError(err)
		require.NotNil(doc)

		var featureFlag tgchat.FeatureFlag
		err = doc.Unmarshal(&featureFlag)
		require.NoError(err)

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

		query := clover.
			NewQuery(tgchat.FeatureFlag{}.CollectionName()).
			Where(clover.Field("chat_id").Eq(chatID))

		doc, err := model.Clover.FindFirst(query)
		require.NoError(err)
		require.NotNil(doc)

		var featureFlag tgchat.FeatureFlag
		err = doc.Unmarshal(&featureFlag)
		require.NoError(err)

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
