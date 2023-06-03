package datastore

import (
	"context"
	"errors"
	"time"

	"github.com/nekomeowww/timecapsule/v2"
	"github.com/samber/lo"
	"go.uber.org/fx"

	"github.com/nekomeowww/insights-bot/pkg/logger"
	"github.com/nekomeowww/insights-bot/pkg/types/redis"
	"github.com/nekomeowww/insights-bot/pkg/types/timecapsules"
)

type NewAutoRecapTimeCapsuleDiggerParams struct {
	fx.In

	Lifecycle fx.Lifecycle

	Logger *logger.Logger
	Redis  *Redis
}

type AutoRecapTimeCapsuleDigger struct {
	*timecapsule.TimeCapsuleDigger[timecapsules.AutoRecapCapsule]
	started bool
}

func (d *AutoRecapTimeCapsuleDigger) Check(ctx context.Context) error {
	return lo.Ternary(d.started, nil, errors.New("digger not started"))
}

func NewAutoRecapTimeCapsuleDigger() func(NewAutoRecapTimeCapsuleDiggerParams) (*AutoRecapTimeCapsuleDigger, error) {
	return func(params NewAutoRecapTimeCapsuleDiggerParams) (*AutoRecapTimeCapsuleDigger, error) {
		dataloader := timecapsule.NewRueidisDataloader[timecapsules.AutoRecapCapsule](redis.TimeCapsuleAutoRecapSortedSetKey.Format(), params.Redis)

		digger := &AutoRecapTimeCapsuleDigger{TimeCapsuleDigger: timecapsule.NewDigger[timecapsules.AutoRecapCapsule](
			dataloader,
			time.Second,
			timecapsule.TimeCapsuleDiggerOption{Logger: params.Logger.LogrusLogger},
		)}

		params.Lifecycle.Append(fx.Hook{
			OnStart: func(ctx context.Context) error {
				go digger.Start()
				digger.started = true
				return nil
			},
			OnStop: func(ctx context.Context) error {
				digger.Stop()
				return nil
			},
		})

		return digger, nil
	}
}
