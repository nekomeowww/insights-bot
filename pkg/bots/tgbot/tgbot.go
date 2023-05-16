package tgbot

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/samber/lo"

	"github.com/nekomeowww/insights-bot/pkg/healthchecker"
	"github.com/nekomeowww/insights-bot/pkg/logger"
	"github.com/nekomeowww/insights-bot/pkg/utils"
)

type BotServiceOptions struct {
	webhookURL  string
	webhookPort string
	token       string
	dispatcher  *Dispatcher
	logger      *logger.Logger
}

type CallOption func(*BotServiceOptions)

func WithWebhookURL(url string) CallOption {
	return func(o *BotServiceOptions) {
		o.webhookURL = url
	}
}

func WithWebhookPort(port string) CallOption {
	return func(o *BotServiceOptions) {
		o.webhookPort = port
	}
}

func WithToken(token string) CallOption {
	return func(o *BotServiceOptions) {
		o.token = token
	}
}

func WithDispatcher(dispatcher *Dispatcher) CallOption {
	return func(o *BotServiceOptions) {
		o.dispatcher = dispatcher
	}
}

func WithLogger(logger *logger.Logger) CallOption {
	return func(o *BotServiceOptions) {
		o.logger = logger
	}
}

var _ healthchecker.HealthChecker = (*BotService)(nil)

type BotService struct {
	*tgbotapi.BotAPI

	opts       *BotServiceOptions
	logger     *logger.Logger
	dispatcher *Dispatcher

	webhookServer     *http.Server
	webhookUpdateChan chan tgbotapi.Update
	updateChan        tgbotapi.UpdatesChannel
	alreadyClose      bool
	ctxCancel         context.CancelFunc
	webhookStarted    bool
}

func NewBotService(callOpts ...CallOption) (*BotService, error) {
	opts := new(BotServiceOptions)
	for _, callOpt := range callOpts {
		callOpt(opts)
	}

	if opts.token == "" {
		return nil, errors.New("must supply a valid telegram bot token in configs or environment variable")
	}

	b, err := tgbotapi.NewBotAPI(opts.token)
	if err != nil {
		return nil, err
	}

	bot := &BotService{
		BotAPI:     b,
		opts:       opts,
		logger:     opts.logger,
		dispatcher: opts.dispatcher,
	}

	// init webhook server and set webhook
	if bot.opts.webhookURL != "" {
		parsed, err := url.Parse(bot.opts.webhookURL)
		if err != nil {
			return nil, err
		}

		bot.webhookUpdateChan = make(chan tgbotapi.Update, b.Buffer)
		bot.webhookServer = newWebhookServer(parsed.Path, bot.opts.webhookPort, bot.BotAPI, bot.webhookUpdateChan)

		err = setWebhook(bot.opts.webhookURL, bot.BotAPI)
		if err != nil {
			return nil, err
		}
	} else {
		u := tgbotapi.NewUpdate(0)
		u.Timeout = 60
		bot.updateChan = b.GetUpdatesChan(u)
	}

	// obtain webhook info
	webhookInfo, err := bot.GetWebhookInfo()
	if err != nil {
		return nil, err
	}
	if bot.opts.webhookURL != "" && webhookInfo.IsSet() && webhookInfo.LastErrorDate != 0 {
		bot.logger.Errorf("webhook callback failed: %s", webhookInfo.LastErrorMessage)
	}

	// cancel the previous set webhook
	if bot.opts.webhookURL == "" && webhookInfo.IsSet() {
		_, err := bot.Request(tgbotapi.DeleteWebhookConfig{DropPendingUpdates: true})
		if err != nil {
			return nil, err
		}
	}

	return bot, nil
}

func (b *BotService) getUpdateChan() tgbotapi.UpdatesChannel {
	if b.opts.webhookURL != "" {
		return b.webhookUpdateChan
	}

	return b.updateChan
}

func (b *BotService) Stop(ctx context.Context) error {
	if b.alreadyClose {
		return nil
	}

	b.alreadyClose = true

	if b.opts.webhookURL != "" {
		closeCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := b.webhookServer.Shutdown(closeCtx); err != nil && err != http.ErrServerClosed {
			return fmt.Errorf("failed to shutdown webhook server: %w", err)
		}

		close(b.webhookUpdateChan)
	} else {
		b.StopReceivingUpdates()
	}

	b.ctxCancel()

	return nil
}

func (b *BotService) startPullUpdates() {
	ctx, cancel := context.WithCancel(context.Background())
	b.ctxCancel = cancel

	for {
		if b.alreadyClose {
			b.logger.Info("stopped to receiving updates")

			return
		}

		select {
		case update := <-b.getUpdateChan():
			b.dispatcher.Dispatch(b.BotAPI, update)
		case <-ctx.Done():
			b.logger.Info("stopped to receiving updates")
			b.webhookStarted = false

			return
		}
	}
}

func (b *BotService) Start(ctx context.Context) error {
	return utils.Invoke0(ctx, func() error {
		if b.opts.webhookURL != "" && b.webhookServer != nil {
			l, err := net.Listen("tcp", b.webhookServer.Addr)
			if err != nil {
				return err
			}

			go func() {
				err := b.webhookServer.Serve(l)
				if err != nil && err != http.ErrServerClosed {
					b.logger.Fatal(err)
				}
			}()

			b.logger.Infof("Telegram Bot webhook server is listening on %s", b.webhookServer.Addr)
		}

		go b.startPullUpdates()
		b.webhookStarted = true
		return nil
	})
}

func (b *BotService) Check(ctx context.Context) error {
	// only check the webhookStarted field when running bot in webhook mode
	if b.opts.webhookURL != "" {
		return lo.Ternary(b.webhookStarted, nil, errors.New("bot service is not started yet"))
	}

	// otherwise return nil
	return nil
}

type Bot struct {
	*tgbotapi.BotAPI
	logger *logger.Logger
}

func (b *Bot) MustSend(chattable tgbotapi.Chattable) *tgbotapi.Message {
	message, err := b.Send(chattable)
	if err != nil {
		b.logger.Errorf("failed to send %v to telegram: %v", utils.SprintJSON(chattable), err)
		return nil
	}

	return &message
}

func (b *Bot) MustRequest(chattable tgbotapi.Chattable) *tgbotapi.APIResponse {
	resp, err := b.Request(chattable)
	if err != nil {
		b.logger.Errorf("failed to request %v to telegram: %v", utils.SprintJSON(chattable), err)
		return nil
	}

	return resp
}
