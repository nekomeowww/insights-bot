package tgbot

import (
	"encoding/json"
	"testing"

	"github.com/nekomeowww/insights-bot/pkg/logger"
	"github.com/redis/rueidis"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAssignOneCallbackQueryData(t *testing.T) {
	c, err := rueidis.NewClient(rueidis.ClientOption{
		InitAddress:  []string{"localhost:6379"},
		DisableCache: true,
	})
	require.NoError(t, err)

	data := struct {
		Hello string `json:"hello"`
	}{
		Hello: "world",
	}

	bot := Bot{logger: logger.NewLogger(logrus.DebugLevel, "insights-bot", "", nil), rueidisClient: c}

	callbackQueryData, err := bot.AssignOneCallbackQueryData("test", data)
	require.NoError(t, err)

	routeHash, dataHash := bot.routeHashAndActionHashFromData(callbackQueryData)
	require.NotEmpty(t, routeHash)
	require.NotEmpty(t, dataHash)

	dataStr, err := bot.fetchCallbackQueryActionData("test", dataHash)
	require.NoError(t, err)

	assert.Equal(t, string(lo.Must(json.Marshal(data))), dataStr)
}
