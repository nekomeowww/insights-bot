package tgchats

import (
	"testing"

	"github.com/nekomeowww/insights-bot/internal/configs"
	"github.com/nekomeowww/insights-bot/internal/datastore"
	"github.com/nekomeowww/insights-bot/pkg/tutils"
)

var model *Model

func TestMain(m *testing.M) {
	ent, err := datastore.NewEnt()(datastore.NewEntParams{
		Lifecycle: tutils.NewEmtpyLifecycle(),
		Configs:   configs.NewTestConfig()(),
	})
	if err != nil {
		panic(err)
	}

	model, err = NewModel()(NewModelParams{
		Ent: ent,
	})
	if err != nil {
		panic(err)
	}

	m.Run()
}
