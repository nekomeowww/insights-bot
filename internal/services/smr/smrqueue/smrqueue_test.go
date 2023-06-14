package smrqueue

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nekomeowww/insights-bot/internal/configs"
	"github.com/nekomeowww/insights-bot/internal/datastore"
	"github.com/nekomeowww/insights-bot/internal/lib"
	"github.com/nekomeowww/insights-bot/pkg/types/bot"
	types "github.com/nekomeowww/insights-bot/pkg/types/smr"
)

var testQueue *Queue

func TestMain(m *testing.M) {
	config := configs.NewTestConfig()()

	logger, err := lib.NewLogger()(lib.NewLoggerParams{
		Configs: config,
	})
	if err != nil {
		panic(err)
	}

	redis, _ := datastore.NewRedis()(datastore.NewRedisParams{
		Configs: config,
	})
	testQueue = NewQueue()(NewQueueParams{
		Logger:      logger,
		RedisClient: redis,
	})

	os.Exit(m.Run())
}

func TestQueue_AddTask(t *testing.T) {
	a := assert.New(t)
	r := require.New(t)
	taskInfo := types.TaskInfo{
		Platform:  bot.FromPlatformDiscord,
		URL:       "https://an.example.url/article",
		ChatID:    114514,
		MessageID: 1919810,
		ChannelID: "CHANNEL",
		TeamID:    "A_TEAM",
	}
	err := testQueue.AddTask(taskInfo)
	a.Empty(err)

	// clean up
	defer func() {
		err = testQueue.redisClient.Do(context.Background(), testQueue.redisClient.B().Del().Key("smr/task").Build()).Error()
		r.Empty(err)
	}()

	// try get task
	var taskResult []string
	taskResult, err = testQueue.redisClient.Do(context.Background(), testQueue.redisClient.B().Brpop().Key("smr/task").Timeout(0).Build()).AsStrSlice()
	r.Len(taskResult, 2)
	r.Empty(err)
	a.Equal("smr/task", taskResult[0])

	expect, err := json.Marshal(&taskInfo)
	r.Empty(err)
	a.JSONEq(string(expect), taskResult[1])
}

func TestService_GetTask(t *testing.T) {
	a := assert.New(t)
	r := require.New(t)
	expect := types.TaskInfo{
		Platform:  bot.FromPlatformDiscord,
		URL:       "https://an.example.url/article",
		ChatID:    114514,
		MessageID: 1919810,
		ChannelID: "CHANNEL",
		TeamID:    "A_TEAM",
	}

	expectJson, err := json.Marshal(&expect)
	r.Empty(err)

	// Add 11 tasks
	for i := 0; i < 11; i++ {
		err = testQueue.redisClient.Do(context.Background(), testQueue.redisClient.B().Lpush().Key("smr/task").Element(string(expectJson)).Build()).Error()
		r.Empty(err)
	}

	// clean up
	defer func() {
		err = testQueue.redisClient.Do(context.Background(), testQueue.redisClient.B().Del().Key("smr/task").Build()).Error()
		r.Empty(err)
	}()

	var actual types.TaskInfo
	// try to get task 10 times
	for i := 0; i < 10; i++ {
		actual, err = testQueue.GetTask()
		a.Empty(err)
		a.Equal(expect, actual)
	}

	// try to get last task
	actual, err = testQueue.GetTask()
	a.ErrorIs(err, ErrQueueFull)
	a.Zero(actual)

	// finish a task
	testQueue.FinishTask()

	// try to get last task again
	actual, err = testQueue.GetTask()
	a.Empty(err)
	a.Equal(expect, actual)
}
