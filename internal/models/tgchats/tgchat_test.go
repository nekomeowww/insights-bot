package tgchats

import (
	"os"
	"testing"

	"github.com/nekomeowww/insights-bot/internal/configs"
	"github.com/nekomeowww/insights-bot/internal/datastore"
	"github.com/nekomeowww/insights-bot/internal/lib"
	"github.com/nekomeowww/insights-bot/pkg/tutils"
)

var model *Model

func TestMain(m *testing.M) {
	logger, err := lib.NewLogger()(lib.NewLoggerParams{
		Configs: configs.NewTestConfig()(),
	})
	if err != nil {
		panic(err)
	}

	ent, err := datastore.NewEnt()(datastore.NewEntParams{
		Lifecycle: tutils.NewEmtpyLifecycle(),
		Configs:   configs.NewTestConfig()(),
	})
	if err != nil {
		panic(err)
	}

	model, err = NewModel()(NewModelParams{
		Ent:    ent,
		Logger: logger,
	})
	if err != nil {
		panic(err)
	}

	os.Exit(m.Run())
}
