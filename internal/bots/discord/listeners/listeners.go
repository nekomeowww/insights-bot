package listeners

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	smr2 "github.com/nekomeowww/insights-bot/internal/models/smr"
	"github.com/nekomeowww/insights-bot/internal/services/smr"
	"github.com/nekomeowww/insights-bot/internal/services/smr/types"
	"github.com/nekomeowww/insights-bot/pkg/logger"
	"go.uber.org/fx"
)

func NewModules() fx.Option {
	return fx.Options(
		fx.Provide(NewListeners()),
	)
}

type NewListenersParam struct {
	fx.In

	Logger *logger.Logger
	Smr    *smr.Service
}

type Listeners struct {
	logger *logger.Logger
	smr    *smr.Service
}

func NewListeners() func(param NewListenersParam) *Listeners {
	return func(param NewListenersParam) *Listeners {
		return &Listeners{
			logger: param.Logger,
			smr:    param.Smr,
		}
	}
}

func (b *Listeners) smrCmd(event *events.ApplicationCommandInteractionCreate, data discord.SlashCommandInteractionData) {
	urlString := data.String("link")

	b.logger.Infof("discordbot: command received: /smr %s", urlString)

	// url check
	err := smr.CheckUrl(urlString)
	if err != nil {
		if smr.IsUrlCheckError(err) {
			err = event.CreateMessage(discord.NewMessageCreateBuilder().SetContent(err.Error()).Build())
			if err != nil {
				b.logger.WithField("error", err.Error()).Warn("discordbot: failed to send error message")
			}

			return
		}

		err = event.CreateMessage(discord.NewMessageCreateBuilder().SetContent("出现了一些问题，可以再试试？").Build())
		if err != nil {
			b.logger.WithField("error", err.Error()).Warn("discordbot: failed to send error message")
		}

		return
	}

	// must reply the interaction as soon as possible
	err = event.CreateMessage(discord.NewMessageCreateBuilder().
		SetContent("请稍等，量子速读中...").
		Build())
	if err != nil {
		b.logger.WithField("error", err.Error()).Warn("discordbot: failed to send response message")
		return
	}

	err = b.smr.AddTask(types.TaskInfo{
		Platform:  smr2.FromPlatformDiscord,
		Url:       urlString,
		ChannelID: event.Channel().ID.String(),
	})
	if err != nil {
		b.logger.WithField("error", err.Error()).Warn("discordbot: failed to add task")

		err = event.CreateMessage(discord.NewMessageCreateBuilder().SetContent("出现了一些问题，可以再试试？").Build())
		if err != nil {
			b.logger.WithField("error", err.Error()).Warn("discordbot: failed to send error message")
		}

		return
	}
}

func (b *Listeners) CommandListener(event *events.ApplicationCommandInteractionCreate) {
	data := event.SlashCommandInteractionData()
	switch data.CommandName() {
	case "smr":
		b.smrCmd(event, data)
	}
}
