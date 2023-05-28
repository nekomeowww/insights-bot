package smr

import (
	"testing"

	"github.com/nekomeowww/insights-bot/internal/configs"
	"github.com/nekomeowww/insights-bot/internal/datastore"
	"github.com/nekomeowww/insights-bot/internal/lib"
	"github.com/nekomeowww/insights-bot/internal/models/smr"
	"github.com/nekomeowww/insights-bot/pkg/tutils"
)

var testService *Service

func TestMain(m *testing.M) {
	config := configs.NewTestConfig()()
	redis, _ := datastore.NewRedis()(datastore.NewRedisParams{
		Configs: config,
	})
	testService, _ = NewService()(NewServiceParam{
		Config: config,
		Model: smr.NewModel()(smr.NewModelParams{
			Logger: lib.NewLogger()(lib.NewLoggerParams{
				Configs: config,
			}),
		}),
		Logger: lib.NewLogger()(lib.NewLoggerParams{
			Configs: config,
		}),
		RedisClient: redis,
		LifeCycle:   tutils.NewEmtpyLifecycle(),
	})

	m.Run()
}
