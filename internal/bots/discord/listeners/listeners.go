package listeners

import (
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/nekomeowww/insights-bot/internal/models/smr"
	"github.com/nekomeowww/insights-bot/internal/services/smr/smrqueue"
	"github.com/nekomeowww/insights-bot/internal/services/smr/smrutils"
	"github.com/nekomeowww/insights-bot/internal/services/smr/types"
	"github.com/nekomeowww/insights-bot/pkg/logger"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func NewModules() fx.Option {
	return fx.Options(
		fx.Provide(NewListeners()),
	)
}

type NewListenersParam struct {
	fx.In

	Logger   *logger.Logger
	SmrQueue *smrqueue.Queue
}

type Listeners struct {
	logger   *logger.Logger
	smrQueue *smrqueue.Queue
}

func NewListeners() func(param NewListenersParam) *Listeners {
	return func(param NewListenersParam) *Listeners {
		return &Listeners{
			logger:   param.Logger,
			smrQueue: param.SmrQueue,
		}
	}
}

func (b *Listeners) smrCmd(event *events.ApplicationCommandInteractionCreate, data discord.SlashCommandInteractionData) {
	urlString := data.String("link")

	b.logger.Info(fmt.Sprintf("discord: command received: /smr %s", urlString))

	// url check
	err, originErr := smrutils.CheckUrl(urlString)
	if err != nil {
		if smrutils.IsUrlCheckError(err) {
			err = event.CreateMessage(
				discord.NewMessageCreateBuilder().
					SetContent(smrutils.FormatUrlCheckError(err, smr.FromPlatformDiscord)).
					Build(),
			)
			if err != nil {
				b.logger.Warn("discord: failed to send error message", zap.Error(err))
			}

			return
		}

		err = event.CreateMessage(discord.NewMessageCreateBuilder().SetContent("出现了一些问题，可以再试试？").Build())
		if err != nil {
			b.logger.Warn("discord: failed to send error message", zap.Error(err), zap.NamedError("original_error", originErr))
		}

		return
	}

	// must reply the interaction as soon as possible
	err = event.CreateMessage(discord.NewMessageCreateBuilder().
		SetContent("请稍等，量子速读中...").
		Build())
	if err != nil {
		b.logger.Warn("discord: failed to send response message", zap.Error(err))
		return
	}

	err = b.smrQueue.AddTask(types.TaskInfo{
		Platform:  smr.FromPlatformDiscord,
		URL:       urlString,
		ChannelID: event.Channel().ID.String(),
	})
	if err != nil {
		b.logger.Warn("discord: failed to add task", zap.Error(err))

		err = event.CreateMessage(discord.NewMessageCreateBuilder().SetContent("出现了一些问题，可以再试试？").Build())
		if err != nil {
			b.logger.Warn("discord: failed to send error message", zap.Error(err))
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
