package telegram

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/samber/lo"
	"go.uber.org/fx"

	"github.com/nekomeowww/insights-bot/internal/bots/telegram/handlers"
	"github.com/nekomeowww/insights-bot/internal/bots/telegram/middlewares"
	"github.com/nekomeowww/insights-bot/internal/configs"
	"github.com/nekomeowww/insights-bot/internal/models/chathistories"
	"github.com/nekomeowww/insights-bot/internal/models/tgchats"
	"github.com/nekomeowww/insights-bot/pkg/bots/tgbot"
	"github.com/nekomeowww/insights-bot/pkg/logger"
)

func NewModules() fx.Option {
	return fx.Options(
		fx.Provide(NewBot()),
		fx.Provide(tgbot.NewDispatcher()),
		fx.Options(handlers.NewModules()),
	)
}

func newWebhookServer(patternPath, port string, bot *tgbotapi.BotAPI, updateChan chan<- tgbotapi.Update) *http.Server {
	srv := http.NewServeMux()
	srv.HandleFunc(patternPath+"/"+bot.Token, func(w http.ResponseWriter, r *http.Request) {
		update, err := bot.HandleUpdate(r)
		if err != nil {
			errMsg, _ := json.Marshal(map[string]string{"error": err.Error()})
			w.WriteHeader(http.StatusBadRequest)
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write(errMsg)
			return
		}

		updateChan <- *update
	})

	return &http.Server{
		Addr:              net.JoinHostPort("", lo.Ternary(port == "", "7071", port)),
		ReadTimeout:       time.Second * 15,
		ReadHeaderTimeout: time.Second * 15,
		Handler:           srv,
	}
}

func setWebhook(webhookURL string, bot *tgbotapi.BotAPI) error {
	webhookConfig, err := tgbotapi.NewWebhook(webhookURL + "/" + bot.Token)
	if err != nil {
		return fmt.Errorf("failed to create webhook config: %w", err)
	}

	_, err = bot.Request(webhookConfig)
	if err != nil {
		return fmt.Errorf("failed to set webhook: %w", err)
	}

	return nil
}

type NewBotParam struct {
	fx.In

	Lifecycle fx.Lifecycle

	Config     *configs.Config
	Logger     *logger.Logger
	Dispatcher *tgbot.Dispatcher
	Handlers   *handlers.Handlers

	ChatHistories *chathistories.Model
	TgChats       *tgchats.Model
}

type Bot struct {
	*tgbotapi.BotAPI

	config        *configs.Config
	logger        *logger.Logger
	dispatcher    *tgbot.Dispatcher
	chatHistories *chathistories.Model

	webhookServer     *http.Server
	webhookUpdateChan chan tgbotapi.Update
	updateChan        tgbotapi.UpdatesChannel
	alreadyClose      bool
	ctxCancel         context.CancelFunc
}

func NewBot() func(param NewBotParam) (*Bot, error) {
	return func(param NewBotParam) (*Bot, error) {
		if param.Config.Telegram.BotToken == "" {
			param.Logger.Fatal("must supply a valid telegram bot token in configs or environment variable")
		}

		b, err := tgbotapi.NewBotAPI(param.Config.Telegram.BotToken)
		if err != nil {
			return nil, err
		}

		bot := &Bot{
			BotAPI:        b,
			config:        param.Config,
			logger:        param.Logger,
			dispatcher:    param.Dispatcher,
			chatHistories: param.ChatHistories,
		}

		// init webhook server and set webhook
		if bot.config.Telegram.BotWebhookURL != "" {
			parsed, err := url.Parse(bot.config.Telegram.BotWebhookURL)
			if err != nil {
				return nil, err
			}

			bot.webhookUpdateChan = make(chan tgbotapi.Update, b.Buffer)
			bot.webhookServer = newWebhookServer(parsed.Path, bot.config.Telegram.BotWebhookPort, bot.BotAPI, bot.webhookUpdateChan)

			err = setWebhook(bot.config.Telegram.BotWebhookURL, bot.BotAPI)
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
		if bot.config.Telegram.BotWebhookURL != "" && webhookInfo.IsSet() && webhookInfo.LastErrorDate != 0 {
			param.Logger.Errorf("webhook callback failed: %s", webhookInfo.LastErrorMessage)
		}

		// cancel the previous set webhook
		if bot.config.Telegram.BotWebhookURL == "" && webhookInfo.IsSet() {
			_, err := bot.Request(tgbotapi.DeleteWebhookConfig{DropPendingUpdates: true})
			if err != nil {
				return nil, err
			}
		}

		param.Lifecycle.Append(fx.Hook{
			OnStop: func(ctx context.Context) error {
				return bot.stop(ctx)
			},
		})

		param.Dispatcher.Use(middlewares.RecordMessage(param.ChatHistories, param.TgChats))
		param.Handlers.InstallAll()
		param.Logger.Infof("Authorized as bot @%s", bot.Self.UserName)

		return bot, nil
	}
}

func (b *Bot) getUpdateChan() tgbotapi.UpdatesChannel {
	if b.config.Telegram.BotWebhookURL != "" {
		return b.webhookUpdateChan
	}

	return b.updateChan
}

func (b *Bot) stop(ctx context.Context) error {
	if b.alreadyClose {
		return nil
	}

	b.alreadyClose = true

	if b.config.Telegram.BotWebhookURL != "" {
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

func (b *Bot) startPullUpdates() {
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
			return
		}
	}
}

func (b *Bot) start() error {
	if b.config.Telegram.BotWebhookURL != "" && b.webhookServer != nil {
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

	return nil
}

func Run() func(bot *Bot) error {
	return func(bot *Bot) error {
		return bot.start()
	}
}
