package slackbot

import "github.com/slack-go/slack"

func NewSlackWebhookMessage(msg string) *slack.WebhookMessage {
	return &slack.WebhookMessage{
		Parse:        slack.MarkdownType,
		Text:         msg,
		ResponseType: slack.ResponseTypeInChannel,
	}
}
