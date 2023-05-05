package models

import (
	"go.uber.org/fx"

	"github.com/nekomeowww/insights-bot/internal/models/chathistories"
	"github.com/nekomeowww/insights-bot/internal/models/smr"
	"github.com/nekomeowww/insights-bot/internal/models/tgchats"
)

func NewModules() fx.Option {
	return fx.Options(
		fx.Provide(chathistories.NewModel()),
		fx.Provide(tgchats.NewModel()),
		fx.Provide(smr.NewModel()),
	)
}
