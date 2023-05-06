package slack

import (
	"context"
	"errors"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/nekomeowww/insights-bot/ent"
	"github.com/nekomeowww/insights-bot/ent/slackoauthcredentials"
	"github.com/nekomeowww/insights-bot/pkg/bots/slackbot"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
)

func (s *SlackBot) postCommandInfo(ctx *gin.Context) {
	var body recivedCommandInfo
	if err := ctx.Bind(&body); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		s.logger.WithField("error", err.Error()).Warn("failed to bind request body, maybe slack request definition changed")

		return
	}

	s.logger.WithFields(logrus.Fields{
		"user_id":    body.UserID,
		"channel_id": body.ChannelID,
	}).Infof("slack: command received: /smr %s", body.Text)

	urlString := body.Text

	var err error

	if urlString == "" {
		err = errors.New("没有找到链接，可以发送一个有效的链接吗？用法：/smr <链接>")
	} else {
		var parsedURL *url.URL
		parsedURL, err = url.Parse(urlString)
		if err != nil {
			err = errors.New("你发来的链接无法被理解，可以重新发一个试试。用法：/smr <链接>")
		}
		if parsedURL.Scheme == "" || !lo.Contains([]string{"http", "https"}, parsedURL.Scheme) {
			err = errors.New("你发来的链接无法被理解，可以重新发一个试试。用法：/smr <链接>")
		}
	}

	if err != nil {
		ctx.JSON(http.StatusOK, slackbot.NewSlackWebhookMessage(err.Error()))
		return
	}

	// get access token
	token, err := s.ent.SlackOAuthCredentials.Query().Where(
		slackoauthcredentials.TeamID(body.TeamID),
	).First(context.Background())
	if err != nil {
		s.logger.WithField("error", err.Error()).Warn("slack: failed to get team's access token")
		if ent.IsNotFound(err) {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		ctx.AbortWithStatus(http.StatusInternalServerError)

		return
	}

	// send to channel
	s.processChan <- smrRequestInfo{
		accessToken: token.AccessToken,
		inputUrl:    body.Text,
		channelID:   body.ChannelID,
	}

	// response
	ctx.JSON(http.StatusOK, slackbot.NewSlackWebhookMessage("请稍等，量子速读中..."))
}

// Receive auth code and request for access token.
func (b *SlackBot) getInstallAuth(ctx *gin.Context) {
	code := ctx.Query("code")
	if code == "" {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	resp, err := slack.GetOAuthV2Response(&http.Client{}, b.config.Slack.ClientID, b.config.Slack.ClientSecret, code, "")
	if err != nil {
		b.logger.WithError(err).Error("slack: failed to get access token, interrupt")
		ctx.AbortWithStatus(http.StatusServiceUnavailable)

		return
	}

	err = b.createNewSlackCredential(resp.Team.ID, resp.AccessToken, resp.RefreshToken)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.Header("content-type", "text/html")
	_, _ = ctx.Writer.Write([]byte("<h1 style=\"text-align:center\">Success! Now you can close this page<h1>"))

	ctx.Status(http.StatusOK)
}
