package chathistories

import (
	"context"
	"testing"

	"github.com/nekomeowww/insights-bot/pkg/types/redis"
	"github.com/nekomeowww/xo"
	"github.com/redis/rueidis"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHasOngoingRecapForwardedFromPrivateMessages(t *testing.T) {
	userID := xo.RandomInt64()

	// assign batch for user ID
	setCmd := model.redis.B().
		Set().
		Key(redis.RecapReplayFromPrivateMessageControl1.Format(userID)).
		Value("1").
		ExSeconds(60 * 60 * 2).
		Build()

	err := model.redis.Do(context.Background(), setCmd).Error()
	require.NoError(t, err)

	has, err := model.HasOngoingRecapForwardedFromPrivateMessages(userID)
	require.NoError(t, err)
	assert.True(t, has)
}

func TestEnabledRecapForwardedFromPrivateMessages(t *testing.T) {
	userID := xo.RandomInt64()

	err := model.EnabledRecapForwardedFromPrivateMessages(userID)
	require.NoError(t, err)

	getCmd := model.redis.B().
		Get().
		Key(redis.RecapReplayFromPrivateMessageControl1.Format(userID)).
		Build()

	res := model.redis.Do(context.Background(), getCmd)
	require.NoError(t, res.Error())

	str, err := res.ToString()
	require.NoError(t, err)
	assert.Equal(t, "1", str)
}

func TestDisableRecapForwardedFromPrivateMessages(t *testing.T) {
	userID := xo.RandomInt64()

	err := model.EnabledRecapForwardedFromPrivateMessages(userID)
	require.NoError(t, err)

	err = model.DisableRecapForwardedFromPrivateMessages(userID)
	require.NoError(t, err)

	getCmd := model.redis.B().
		Get().
		Key(redis.RecapReplayFromPrivateMessageControl1.Format(userID)).
		Build()

	res := model.redis.Do(context.Background(), getCmd)
	require.Error(t, res.Error())
	assert.True(t, rueidis.IsRedisNil(res.Error()))
}
