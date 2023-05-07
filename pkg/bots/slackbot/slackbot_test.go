package slackbot

import (
	"errors"
	"net/http"
	"testing"

	"github.com/nekomeowww/insights-bot/pkg/tutils"
	"github.com/slack-go/slack"
	"github.com/stretchr/testify/assert"
)

func TestSlackCli_SendMessageWithTokenExpirationCheck(t *testing.T) {
	t.Run("Default", func(t *testing.T) {
		a := assert.New(t)

		httpCli := &tutils.SlackHttpClientForTokenExpiresTest{
			TokenExpired: true,
		}
		slackCli := NewSlackCli(httpCli, "ID", "SECRET", "REFRESH_TOKEN_0", "ACCESS_TOKEN_0")

		_, _, _, err := slackCli.SendMessageWithTokenExpirationCheck("CHANNEL", func(accessToken, refreshToken string) error {
			a.Equal("ACCESS_TOKEN", accessToken)
			a.Equal("REFRESH_TOKEN", refreshToken)

			return nil
		})

		a.Empty(err)
		a.Equal(3, httpCli.GetInvokeNums())
	})

	t.Run("StatusError", func(t *testing.T) {
		a := assert.New(t)

		httpCli := &tutils.SlackHttpClientForTokenExpiresTest{
			OtherError: true,
		}
		slackCli := NewSlackCli(httpCli, "ID", "SECRET", "REFRESH_TOKEN_0", "ACCESS_TOKEN_0")

		_, _, _, err := slackCli.SendMessageWithTokenExpirationCheck("CHANNEL", func(accessToken, refreshToken string) error {
			return nil
		})

		statusCodeError, ok := err.(slack.StatusCodeError)
		a.True(ok)

		a.Equal(http.StatusInternalServerError, statusCodeError.HTTPStatusCode())
		a.Equal(1, httpCli.GetInvokeNums())
	})

	t.Run("StoreFuncError", func(t *testing.T) {
		a := assert.New(t)

		httpCli := &tutils.SlackHttpClientForTokenExpiresTest{
			TokenExpired: true,
		}
		slackCli := NewSlackCli(httpCli, "ID", "SECRET", "REFRESH_TOKEN_0", "ACCESS_TOKEN_0")

		storeFuncErr := errors.New("expected error")
		_, _, _, err := slackCli.SendMessageWithTokenExpirationCheck("CHANNEL", func(accessToken, refreshToken string) error {
			return storeFuncErr
		})

		a.ErrorIs(err, storeFuncErr)
		a.Equal(2, httpCli.GetInvokeNums())
	})

	t.Run("GetOAuthTokenError", func(t *testing.T) {
		a := assert.New(t)

		httpCli := &tutils.SlackHttpClientForTokenExpiresTest{
			TokenExpired:       true,
			GetOAuthTokenError: true,
		}
		slackCli := NewSlackCli(httpCli, "ID", "SECRET", "REFRESH_TOKEN_0", "ACCESS_TOKEN_0")

		_, _, _, err := slackCli.SendMessageWithTokenExpirationCheck("CHANNEL", func(accessToken, refreshToken string) error {
			return nil
		})

		statusCodeError, ok := err.(slack.StatusCodeError)
		a.True(ok)

		a.Equal(http.StatusInternalServerError, statusCodeError.HTTPStatusCode())
		a.Equal(2, httpCli.GetInvokeNums())
	})
}
