package slackbot

import (
	"net/http"

	"github.com/slack-go/slack"
)

func NewSlackWebhookMessage(msg string) *slack.WebhookMessage {
	return &slack.WebhookMessage{
		Parse:        slack.MarkdownType,
		Text:         msg,
		ResponseType: slack.ResponseTypeInChannel,
	}
}

type HttpClientForSlack interface {
	Do(req *http.Request) (*http.Response, error)
}

type SlackCli struct {
	*slack.Client

	httpclient   HttpClientForSlack
	clientID     string
	clientSecret string
	refreshToken string
}

func newOriginSlackCli(httpCli HttpClientForSlack, accessToken string) *slack.Client {
	var opt []slack.Option
	if httpCli != nil {
		opt = append(opt, slack.OptionHTTPClient(httpCli))
	}

	return slack.New(accessToken, opt...)
}

func NewSlackCli(httpCli HttpClientForSlack, clientID, clientSecret, refreshToken, accessToken string) *SlackCli {
	return &SlackCli{
		Client:       newOriginSlackCli(httpCli, accessToken),
		clientID:     clientID,
		clientSecret: clientSecret,
		refreshToken: refreshToken,
		httpclient:   httpCli,
	}
}

type StoreNewTokenFunc func(accessToken string, refreshToken string) error

// SendMessageWithTokenExpirationCheck will checks if the error is "token_expired" error,
// if so, will get new token and try again.
func (cli *SlackCli) SendMessageWithTokenExpirationCheck(channel string, storeFn StoreNewTokenFunc, options ...slack.MsgOption) (channelID string, msgTimestamp string, respText string, err error) {
	channelID, msgTimestamp, respText, err = cli.SendMessage(channel, options...)
	if err == nil || err.Error() != "token_expired" {
		return
	}

	resp, err := slack.RefreshOAuthV2Token(cli.httpclient, cli.clientID, cli.clientSecret, cli.refreshToken)
	if err != nil {
		return
	}

	err = storeFn(resp.AccessToken, resp.RefreshToken)
	if err != nil {
		return
	}
	// create new slack client
	cli.Client = newOriginSlackCli(cli.httpclient, resp.AccessToken)

	return cli.SendMessageWithTokenExpirationCheck(channel, storeFn, options...)
}
