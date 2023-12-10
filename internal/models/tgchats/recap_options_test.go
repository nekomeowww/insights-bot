package tgchats

import (
	"testing"

	"github.com/nekomeowww/xo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFindOneOrCreateRecapsOption(t *testing.T) {
	chatID := xo.RandomInt64()

	option, err := model.FindOneOrCreateRecapsOption(chatID)
	require.NoError(t, err)
	require.NotNil(t, option)

	assert.Equal(t, chatID, option.ChatID)
	assert.Empty(t, option.AutoRecapSendMode)
	assert.Equal(t, 4, option.AutoRecapRatesPerDay)

	option2, err := model.FindOneRecapsOption(chatID)
	require.NoError(t, err)
	require.NotNil(t, option2)

	assert.Equal(t, option.ID, option2.ID)

	option3, err := model.FindOneOrCreateRecapsOption(chatID)
	require.NoError(t, err)
	require.NotNil(t, option3)

	assert.Equal(t, option.ID, option3.ID)
}

func TestSetAutoRecapRatesPerDay(t *testing.T) {
	chatID := xo.RandomInt64()

	err := model.SetAutoRecapRatesPerDay(chatID, 4)
	require.NoError(t, err)

	option, err := model.FindOneRecapsOption(chatID)
	require.NoError(t, err)
	require.NotNil(t, option)

	assert.Equal(t, chatID, option.ChatID)
	assert.Equal(t, 4, option.AutoRecapRatesPerDay)

	err = model.SetAutoRecapRatesPerDay(chatID, 10)
	require.NoError(t, err)

	option2, err := model.FindOneRecapsOption(chatID)
	require.NoError(t, err)
	require.NotNil(t, option2)

	assert.Equal(t, chatID, option2.ChatID)
	assert.Equal(t, 10, option2.AutoRecapRatesPerDay)
}
