package datastore

import (
	"context"
	"time"

	"github.com/nekomeowww/timecapsule/v2"
	"go.uber.org/fx"

	"github.com/nekomeowww/insights-bot/pkg/logger"
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
}

func NewAutoRecapTimeCapsuleDigger() func(NewAutoRecapTimeCapsuleDiggerParams) (*AutoRecapTimeCapsuleDigger, error) {
	return func(params NewAutoRecapTimeCapsuleDiggerParams) (*AutoRecapTimeCapsuleDigger, error) {
		dataloader := timecapsule.NewRueidisDataloader[timecapsules.AutoRecapCapsule](timecapsules.AutoRecapTimeCapsuleKey, params.Redis)
		digger := timecapsule.NewDigger[timecapsules.AutoRecapCapsule](dataloader, time.Second, timecapsule.TimeCapsuleDiggerOption{
			Logger: params.Logger,
		})

		params.Lifecycle.Append(fx.Hook{
			OnStart: func(ctx context.Context) error {
				go digger.Start()
				return nil
			},
			OnStop: func(ctx context.Context) error {
				digger.Stop()
				return nil
			},
		})

		return &AutoRecapTimeCapsuleDigger{TimeCapsuleDigger: digger}, nil
	}
}
