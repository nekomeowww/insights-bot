package smr

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/nekomeowww/insights-bot/internal/models/smr"
	"github.com/nekomeowww/insights-bot/internal/services/smr/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestService_AddTask(t *testing.T) {
	a := assert.New(t)
	r := require.New(t)
	taskInfo := types.TaskInfo{
		Platform:  smr.FromPlatformDiscord,
		Url:       "https://an.example.url/article",
		ChatID:    114514,
		MessageID: 1919810,
		ChannelID: "CHANNEL",
		TeamID:    "A_TEAM",
	}
	err := testService.AddTask(taskInfo)
	a.Empty(err)

	// clean up
	defer func() {
		err = testService.redisClient.Do(context.Background(), testService.redisClient.B().Del().Key("smr/task").Build()).Error()
		r.Empty(err)
	}()

	// try get task
	var taskResult []string
	taskResult, err = testService.redisClient.Do(context.Background(), testService.redisClient.B().Brpop().Key("smr/task").Timeout(0).Build()).AsStrSlice()
	r.Empty(err)
	a.Equal("smr/task", taskResult[0])

	expect, err := json.Marshal(&taskInfo)
	r.Empty(err)
	a.JSONEq(string(expect), taskResult[1])
}

func TestService_getTask(t *testing.T) {
	a := assert.New(t)
	r := require.New(t)
	expect := types.TaskInfo{
		Platform:  smr.FromPlatformDiscord,
		Url:       "https://an.example.url/article",
		ChatID:    114514,
		MessageID: 1919810,
		ChannelID: "CHANNEL",
		TeamID:    "A_TEAM",
	}

	expectJson, err := json.Marshal(&expect)
	r.Empty(err)

	err = testService.redisClient.Do(context.Background(), testService.redisClient.B().Lpush().Key("smr/task").Element(string(expectJson)).Build()).Error()
	r.Empty(err)

	// clean up
	defer func() {
		err = testService.redisClient.Do(context.Background(), testService.redisClient.B().Del().Key("smr/task").Build()).Error()
		r.Empty(err)
	}()

	// try get task
	actual, err := testService.getTask()
	r.Empty(err)
	a.Equal(expect, actual)
}
