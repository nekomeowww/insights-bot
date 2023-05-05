package slack

import (
	"context"
	"log"
	"testing"

	"github.com/nekomeowww/insights-bot/internal/configs"
	"github.com/nekomeowww/insights-bot/internal/datastore"
	"github.com/nekomeowww/insights-bot/pkg/logger"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
)

type lifeCycleMock struct{}

func (lifeCycleMock) Append(hook fx.Hook) {}

func newTestSlackBot() *SlackBot {
	config := configs.NewTestConfig()()

	ent, err := datastore.NewEnt()(datastore.NewEntParams{
		Configs:   config,
		Lifecycle: lifeCycleMock{},
	})
	if err != nil {
		log.Fatal("datastore init failed")
	}

	return &SlackBot{
		ent:    ent,
		logger: logger.NewLogger(logrus.InfoLevel, "insights-bot", "", make([]logrus.Hook, 0)),
	}
}

func cleanSlackCredential(bot *SlackBot, r *require.Assertions) {
	_, err := bot.ent.SlackOAuthCredentials.Delete().Exec(context.Background())
	r.Empty(err)
}

func TestSlackBot_createNewSlackCredential(t *testing.T) {
	bot := newTestSlackBot()

	t.Run("no record", func(t *testing.T) {
		a := assert.New(t)
		r := require.New(t)

		defer cleanSlackCredential(bot, r)

		expectTeamId := "TEAM_ID"
		expectAccessToken := "ACCESS_TOKEN"

		r.Empty(bot.createNewSlackCredential(expectTeamId, expectAccessToken))

		// query
		cre, err := bot.ent.SlackOAuthCredentials.Query().First(context.Background())
		r.Empty(err)
		a.Equal(expectTeamId, cre.TeamID)
		a.Equal(expectAccessToken, cre.AccessToken)
	})

	t.Run("exists", func(t *testing.T) {
		a := assert.New(t)
		r := require.New(t)

		defer cleanSlackCredential(bot, r)

		expectTeamId := "TEAM_ID"
		expectAccessToken := "ACCESS_TOKEN"

		r.Empty(bot.createNewSlackCredential(expectTeamId, "ANOTHER_ACCESS_TOKEN"))
		r.Empty(bot.createNewSlackCredential(expectTeamId, expectAccessToken))

		// query
		cre, err := bot.ent.SlackOAuthCredentials.Query().First(context.Background())
		r.Empty(err)
		a.Equal(expectTeamId, cre.TeamID)
		a.Equal(expectAccessToken, cre.AccessToken)
	})
}
