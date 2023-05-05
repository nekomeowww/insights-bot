package main

import (
	"context"
	"log"
	"time"

	"go.uber.org/fx"

	"github.com/nekomeowww/insights-bot/internal/bots/slack"
	"github.com/nekomeowww/insights-bot/internal/bots/telegram"
	"github.com/nekomeowww/insights-bot/internal/configs"
	"github.com/nekomeowww/insights-bot/internal/datastore"
	"github.com/nekomeowww/insights-bot/internal/lib"
	"github.com/nekomeowww/insights-bot/internal/models"
	"github.com/nekomeowww/insights-bot/internal/services"
	"github.com/nekomeowww/insights-bot/internal/services/autorecap"
	"github.com/nekomeowww/insights-bot/internal/services/pprof"
	"github.com/nekomeowww/insights-bot/internal/thirdparty"
)

func main() {
	app := fx.New(fx.Options(
		fx.Provide(configs.NewConfig()),
		fx.Options(lib.NewModules()),
		fx.Options(datastore.NewModules()),
		fx.Options(models.NewModules()),
		fx.Options(thirdparty.NewModules()),
		fx.Options(services.NewModules()),
		fx.Options(telegram.NewModules()),
		fx.Options(slack.NewModules()),
		fx.Invoke(telegram.Run()),
		fx.Invoke(autorecap.Run()),
		fx.Invoke(slack.Run()),
		fx.Invoke(pprof.Run()),
	))

	app.Run()

	stopCtx, stopCtxCancel := context.WithTimeout(context.Background(), time.Second*15)
	defer stopCtxCancel()

	if err := app.Stop(stopCtx); err != nil {
		log.Fatal(err)
	}
}
