package services

import (
	"context"

	"github.com/nekomeowww/insights-bot/internal/datastore"
	"github.com/nekomeowww/insights-bot/pkg/logger"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/nekomeowww/insights-bot/ent/slackoauthcredentials"
)

type NewServicesParam struct {
	fx.In

	Logger *logger.Logger
	Ent    *datastore.Ent
}

type Services struct {
	logger *logger.Logger
	ent    *datastore.Ent
}

func NewServices() func(param NewServicesParam) *Services {
	return func(param NewServicesParam) *Services {
		return &Services{
			logger: param.Logger,
			ent:    param.Ent,
		}
	}
}

func (b *Services) NewStoreFuncForRefresh(teamID string) func(accessToken, refreshToken string) error {
	return func(accessToken, refreshToken string) error {
		return b.CreateOrUpdateSlackCredential(teamID, accessToken, refreshToken)
	}
}

func (b *Services) CreateOrUpdateSlackCredential(teamID, accessToken, refreshToken string) error {
	affectRows, err := b.ent.SlackOAuthCredentials.Update().
		Where(slackoauthcredentials.TeamID(teamID)).
		SetAccessToken(accessToken).
		SetRefreshToken(refreshToken).
		Save(context.Background())
	if err != nil {
		b.logger.Warn("slack: failed to update access token", zap.Error(err))
		return err
	}

	if affectRows == 0 {
		// create
		err = b.ent.SlackOAuthCredentials.Create().
			SetTeamID(teamID).
			SetAccessToken(accessToken).
			SetRefreshToken(refreshToken).
			Exec(context.Background())
		if err != nil {
			b.logger.Warn("slack: failed to save access token", zap.Error(err))
			return err
		}
	}

	return nil
}
