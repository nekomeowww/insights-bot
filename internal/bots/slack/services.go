package slack

import (
	"context"

	"github.com/nekomeowww/insights-bot/ent/slackoauthcredentials"
)

func (b *SlackBot) newStoreFuncForRefresh(teamID string) func(accessToken, refreshToken string) error {
	return func(accessToken, refreshToken string) error {
		return b.createOrUpdateSlackCredential(teamID, accessToken, refreshToken)
	}
}

func (b *SlackBot) createOrUpdateSlackCredential(teamID, accessToken, refreshToken string) error {
	affectRows, err := b.ent.SlackOAuthCredentials.Update().
		Where(slackoauthcredentials.TeamID(teamID)).
		SetAccessToken(accessToken).
		SetRefreshToken(refreshToken).
		Save(context.Background())
	if err != nil {
		b.logger.WithError(err).Warn("slack: failed to update access token")
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
			b.logger.WithError(err).Warn("slack: failed to save access token")
			return err
		}
	}

	return nil
}
