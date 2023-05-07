package slack

import (
	"context"

	"github.com/nekomeowww/insights-bot/ent/slackoauthcredentials"
)

func (b *SlackBot) createNewSlackCredential(teamID, accessToken string) error {
	affectRows, err := b.ent.SlackOAuthCredentials.Update().
		Where(slackoauthcredentials.TeamID(teamID)).
		SetAccessToken(accessToken).
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
			Exec(context.Background())
		if err != nil {
			b.logger.WithError(err).Warn("slack: failed to save access token")
			return err
		}
	}

	return nil
}
