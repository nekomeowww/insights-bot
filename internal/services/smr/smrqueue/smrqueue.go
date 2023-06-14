package smrqueue

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/nekomeowww/insights-bot/internal/datastore"
	"github.com/nekomeowww/insights-bot/pkg/logger"
	types "github.com/nekomeowww/insights-bot/pkg/types/smr"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var (
	ErrQueueFull = errors.New("task queue full")
)

func NewModules() fx.Option {
	return fx.Options(
		fx.Provide(NewQueue()),
	)
}

func NewQueue() func(NewQueueParams) *Queue {
	return func(params NewQueueParams) *Queue {
		return &Queue{
			logger:       params.Logger,
			redisClient:  params.RedisClient,
			ongoingTasks: types.NewOngoingTaskPool(),
		}
	}
}

type NewQueueParams struct {
	fx.In

	Logger      *logger.Logger
	RedisClient *datastore.Redis
}

type Queue struct {
	logger       *logger.Logger
	redisClient  *datastore.Redis
	ongoingTasks *types.OngoingTaskPool
}

func (q *Queue) AddTask(taskInfo types.TaskInfo) error {
	result, err := json.Marshal(&taskInfo)
	if err != nil {
		return err
	}

	err = q.redisClient.Do(context.Background(), q.redisClient.B().Lpush().Key("smr/task").Element(string(result)).Build()).Error()
	if err != nil {
		return err
	}

	q.logger.Info("smr service: task added",
		zap.String("url", taskInfo.URL),
		zap.String("platform", taskInfo.Platform.String()),
	)

	// TODO: #111 should reject ongoing smr request in the same chat
	return nil
}

func (q *Queue) GetTask() (types.TaskInfo, error) {
	var info types.TaskInfo

	if q.Count() >= 10 {
		return info, ErrQueueFull
	}

	res, err := q.redisClient.Do(context.Background(), q.redisClient.B().Brpop().Key("smr/task").Timeout(10).Build()).AsStrSlice()
	if err != nil {
		return info, err
	}

	err = json.Unmarshal([]byte(res[1]), &info)
	if err != nil {
		return info, err
	}

	q.ongoingTasks.Add(info)

	return info, err
}

func (q *Queue) Count() int {
	return q.ongoingTasks.Len()
}

func (q *Queue) FinishTask() {
	q.ongoingTasks.Remove()
}
