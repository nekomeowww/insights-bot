package slack

import (
	"errors"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
)

type recivedCommandInfo struct {
	Command     string `form:"command"`
	Text        string `form:"text"`
	ResponseUrl string `form:"response_url"`
	UserId      string `form:"user_id"`
	ChannelId   string `form:"channel_id"`
}

func (s *SlackBot) postCommandInfo(ctx *gin.Context) {
	var body recivedCommandInfo
	if err := ctx.Bind(&body); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		s.Logger.WithField("error", err.Error()).Warn("failed to bind request body, maybe slack request definition changed")
		return
	}

	s.Logger.WithFields(logrus.Fields{
		"user_id":    body.UserId,
		"channel_id": body.ChannelId,
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
		ctx.JSON(http.StatusOK, &slack.WebhookMessage{
			Parse:        slack.MarkdownType,
			Text:         err.Error(),
			ResponseType: slack.ResponseTypeInChannel,
		})
		return
	}

	// send to channel
	s.processChan <- body

	// response
	ctx.JSON(http.StatusOK, &slack.WebhookMessage{
		Parse:        slack.MarkdownType,
		Text:         "请稍等，量子速读中...",
		ResponseType: slack.ResponseTypeInChannel,
	})
}
