package services

import (
	"context"
	"log"
	"testing"

	"github.com/nekomeowww/insights-bot/internal/configs"
	"github.com/nekomeowww/insights-bot/internal/datastore"
	"github.com/nekomeowww/insights-bot/pkg/logger"
	"github.com/nekomeowww/insights-bot/pkg/tutils"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

func newTestServices() *Services {
	config := configs.NewTestConfig()()

	ent, err := datastore.NewEnt()(datastore.NewEntParams{
		Configs:   config,
		Lifecycle: tutils.NewEmtpyLifecycle(),
	})
	if err != nil {
		log.Fatal("datastore init failed")
	}

	logger, err := logger.NewLogger(zapcore.DebugLevel, "insights-bot", "", make([]logrus.Hook, 0))
	if err != nil {
		log.Fatal("logger init failed")
	}

	return &Services{
		Ent:    ent,
		logger: logger,
	}
}

func cleanSlackCredential(s *Services, r *require.Assertions) {
	_, err := s.Ent.SlackOAuthCredentials.Delete().Exec(context.Background())
	r.Empty(err)
}

func TestSlackBot_createNewSlackCredential(t *testing.T) {
	s := newTestServices()

	t.Run("no record", func(t *testing.T) {
		a := assert.New(t)
		r := require.New(t)

		defer cleanSlackCredential(s, r)

		expectTeamID := "TEAM_ID"
		expectAccessToken := "ACCESS_TOKEN"
		expectRefreshToken := "REFRESH_TOKEN"

		r.Empty(s.CreateOrUpdateSlackCredential(expectTeamID, expectAccessToken, expectRefreshToken))

		// query
		cre, err := s.Ent.SlackOAuthCredentials.Query().First(context.Background())
		r.Empty(err)
		a.Equal(expectTeamID, cre.TeamID)
		a.Equal(expectAccessToken, cre.AccessToken)
		a.Equal(expectRefreshToken, cre.RefreshToken)
	})

	t.Run("exists", func(t *testing.T) {
		a := assert.New(t)
		r := require.New(t)

		defer cleanSlackCredential(s, r)

		expectTeamID := "TEAM_ID"
		expectAccessToken := "ACCESS_TOKEN"
		expectRefreshToken := "REFRESH_TOKEN"

		r.Empty(s.CreateOrUpdateSlackCredential(expectTeamID, "ANOTHER_ACCESS_TOKEN", expectRefreshToken))
		r.Empty(s.CreateOrUpdateSlackCredential(expectTeamID, expectAccessToken, expectRefreshToken))

		// query
		cre, err := s.Ent.SlackOAuthCredentials.Query().First(context.Background())
		r.Empty(err)
		a.Equal(expectTeamID, cre.TeamID)
		a.Equal(expectAccessToken, cre.AccessToken)
		a.Equal(expectRefreshToken, cre.RefreshToken)
	})
}
