package tgchats

import (
	"testing"

	"github.com/nekomeowww/insights-bot/internal/datastore"
)

var model *Model

func TestMain(m *testing.M) {
	clover, cancel := datastore.NewTestClover()()
	defer cancel()

	var err error
	model, err = NewModel()(NewModelParams{
		Clover: clover,
	})
	if err != nil {
		panic(err)
	}

	m.Run()
}
