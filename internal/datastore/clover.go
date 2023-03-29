package datastore

import (
	"os"
	"path/filepath"

	clover "github.com/ostafen/clover/v2"
	"go.uber.org/fx"

	"github.com/nekomeowww/insights-bot/internal/configs"
	"github.com/nekomeowww/insights-bot/pkg/utils"
)

type NewCloverParam struct {
	fx.In

	Config *configs.Config
}

type Clover struct {
	*clover.DB
}

func NewClover() func(NewCloverParam) (*Clover, error) {
	return func(param NewCloverParam) (*Clover, error) {
		db, err := clover.Open(param.Config.CloverDBPath)
		if err != nil {
			return nil, err
		}

		return &Clover{DB: db}, nil
	}
}

func NewTestClover() func() (*Clover, func()) {
	return func() (*Clover, func()) {
		testdataDir := utils.RelativePathOf("../../testdata")
		if _, err := os.Stat(testdataDir); os.IsNotExist(err) {
			err := os.Mkdir(testdataDir, 0755)
			if err != nil {
				panic(err)
			}
		}

		dbPath := filepath.Join(testdataDir, "insights_bot_clover_test_db.db")
		db, err := NewClover()(NewCloverParam{
			Config: &configs.Config{
				CloverDBPath: dbPath,
			},
		})
		if err != nil {
			panic(err)
		}

		return db, func() {
			err := os.RemoveAll(dbPath)
			if err != nil {
				panic(err)
			}
		}
	}
}
