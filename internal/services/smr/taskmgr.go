package smr

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/nekomeowww/insights-bot/internal/services/smr/types"
	"github.com/redis/rueidis"
	"github.com/sourcegraph/conc/pool"
)

func (s *Service) AddTask(taskInfo types.TaskInfo) error {
	result, err := json.Marshal(&taskInfo)
	if err != nil {
		return err
	}

	err = s.redisClient.Do(context.Background(), s.redisClient.B().Lpush().Key("smr/task").Element(string(result)).Build()).Error()
	if err != nil {
		return err
	}

	s.logger.
		WithField("url", taskInfo.Url).
		WithField("platform", taskInfo.Platform).
		Info("smr service: task added")

	// TODO: #111 should reject ongoing smr request in the same chat
	return nil
}

func (s *Service) stop() {
	if s.alreadyClosed {
		return
	}
	s.closeChan <- struct{}{}
	close(s.closeChan)
}

func (s *Service) getTask() (types.TaskInfo, error) {
	var info types.TaskInfo

	res, err := s.redisClient.Do(context.Background(), s.redisClient.B().Brpop().Key("smr/task").Timeout(10).Build()).AsStrSlice()
	if err != nil {
		return info, err
	}

	err = json.Unmarshal([]byte(res[1]), &info)
	if err != nil {
		return info, err
	}

	return info, err
}

func (s *Service) run() {
	s.queue = types.NewTaskQueue()
	s.closeChan = make(chan struct{})
	s.started = true

	s.logger.Info("smr service started")

	needToClose := false

	taskRunner := pool.New().WithMaxGoroutines(10)

	for {
		select {
		case <-s.closeChan:
			s.logger.WithField("last tasks count", s.queue.Len()).Info("smr service: received stop signal, waiting for all tasks done")

			needToClose = true
		default:
		}

		// get task
		if s.queue.Len() >= 10 {
			continue
		}

		info, err := s.getTask()
		if err != nil {
			if errors.Is(err, rueidis.Nil) {
				continue
			}

			s.logger.WithError(err).Warn("smr service: failed to get task")

			continue
		}

		s.queue.AddTask(info)
		taskRunner.Go(func() {
			defer func() {
				err2 := recover()
				if err2 != nil {
					s.logger.
						WithField("err", err2).
						WithField("task", info).
						Error("smr service: task failed with panic")
				}
			}()

			s.processor(info)
			s.queue.RemoveTask()
		})

		if needToClose {
			break
		}
	}

	s.alreadyClosed = true

	taskRunner.Wait()
}
