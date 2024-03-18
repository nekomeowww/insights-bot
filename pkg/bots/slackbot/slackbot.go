package slackbot

import (
	"context"
	"errors"
	"net"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nekomeowww/insights-bot/internal/configs"
	"github.com/nekomeowww/insights-bot/pkg/bots/slackbot/services"
	"github.com/nekomeowww/insights-bot/pkg/healthchecker"
	"github.com/nekomeowww/insights-bot/pkg/logger"
	"github.com/samber/lo"
	"go.uber.org/zap"

	"github.com/slack-go/slack"
)

func NewSlackWebhookMessage(msg string) *slack.WebhookMessage {
	return &slack.WebhookMessage{
		Parse:        slack.MarkdownType,
		Text:         msg,
		ResponseType: slack.ResponseTypeInChannel,
	}
}

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Client struct {
	*slack.Client

	httpClient   HttpClient
	clientID     string
	clientSecret string
	refreshToken string
}

func newOriginSlackCli(httpCli HttpClient, accessToken string) *slack.Client {
	var opt []slack.Option
	if httpCli != nil {
		opt = append(opt, slack.OptionHTTPClient(httpCli))
	}

	return slack.New(accessToken, opt...)
}

func NewSlackCli(httpCli HttpClient, clientID, clientSecret, refreshToken, accessToken string) *Client {
	return &Client{
		Client:       newOriginSlackCli(httpCli, accessToken),
		clientID:     clientID,
		clientSecret: clientSecret,
		refreshToken: refreshToken,
		httpClient:   httpCli,
	}
}

type StoreNewTokenFunc func(accessToken string, refreshToken string) error

// SendMessageWithTokenExpirationCheck will checks if the error is "token_expired" error,
// if so, will get new token and try again.
func (cli *Client) SendMessageWithTokenExpirationCheck(channel string, storeFn StoreNewTokenFunc, options ...slack.MsgOption) (channelID string, msgTimestamp string, respText string, err error) {
	channelID, msgTimestamp, respText, err = cli.SendMessage(channel, options...)
	if err == nil || err.Error() != "token_expired" {
		return
	}

	resp, err := slack.RefreshOAuthV2Token(cli.httpClient, cli.clientID, cli.clientSecret, cli.refreshToken)
	if err != nil {
		return
	}

	err = storeFn(resp.AccessToken, resp.RefreshToken)
	if err != nil {
		return
	}
	// create new slack client
	cli.Client = newOriginSlackCli(cli.httpClient, resp.AccessToken)

	return cli.SendMessageWithTokenExpirationCheck(channel, storeFn, options...)
}

// GetUserInfoWithTokenExpirationCheck will checks if the error is "token_expired" error,
// if so, will get new token and try again.
func (cli *Client) GetUserInfoWithTokenExpirationCheck(channel string, storeFn StoreNewTokenFunc, options ...slack.MsgOption) (slackUser *slack.User, err error) {
	slackUser, err = cli.GetUserInfo(channel)
	if err == nil || err.Error() != "token_expired" {
		return
	}

	resp, err := slack.RefreshOAuthV2Token(cli.httpClient, cli.clientID, cli.clientSecret, cli.refreshToken)
	if err != nil {
		return
	}

	err = storeFn(resp.AccessToken, resp.RefreshToken)
	if err != nil {
		return
	}
	// create new slack client
	cli.Client = newOriginSlackCli(cli.httpClient, resp.AccessToken)

	return cli.GetUserInfoWithTokenExpirationCheck(channel, storeFn, options...)
}

var _ healthchecker.HealthChecker = (*BotService)(nil)

type BotService struct {
	logger *logger.Logger

	services     *services.Services
	serverEngine *gin.Engine
	server       *http.Server

	started bool
}

func (s *BotService) SetService(services *services.Services) {
	s.services = services
}

func (s *BotService) GetService() *services.Services {
	return s.services
}

func NewBotService(slackConfig configs.SectionSlack) *BotService {
	engine := gin.Default()
	server := &http.Server{
		Addr:              lo.Ternary(slackConfig.Port == "", ":7070", net.JoinHostPort("", slackConfig.Port)),
		Handler:           engine,
		ReadHeaderTimeout: time.Second * 10,
	}

	return &BotService{
		serverEngine: engine,
		server:       server,
	}
}

func (s *BotService) Handle(method, path string, handler gin.HandlerFunc) {
	s.serverEngine.Handle(method, path, handler)
}

func (s *BotService) SetLogger(logger *logger.Logger) {
	s.logger = logger
}

func (s *BotService) Check(ctx context.Context) error {
	return lo.Ternary(s.started, nil, errors.New("slack bot not started yet"))
}

func (s *BotService) Run() error {
	listener, err := net.Listen("tcp", s.server.Addr)
	if err != nil {
		return err
	}

	go func() {
		err = s.server.Serve(listener)
		if err != nil && err != http.ErrServerClosed {
			s.logger.Fatal("slack bot server error", zap.Error(err))
		}
	}()

	s.logger.Info("Slack Bot/App webhook server is listening", zap.String("addr", s.server.Addr))
	s.started = true

	return nil
}

func (s *BotService) Stop(ctx context.Context) error {
	err := s.server.Shutdown(ctx)
	if err != nil {
		s.logger.Error("slack bot server shutdown failed", zap.Error(err))
		return err
	}

	s.logger.Info("stopped to receiving new requests")

	return nil
}
