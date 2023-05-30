package tgbot

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/redis/rueidis"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"

	"github.com/nekomeowww/fo"
	"github.com/nekomeowww/insights-bot/pkg/healthchecker"
	"github.com/nekomeowww/insights-bot/pkg/logger"
	"github.com/nekomeowww/insights-bot/pkg/types/redis"
	"github.com/nekomeowww/insights-bot/pkg/types/telegram"
	"github.com/nekomeowww/insights-bot/pkg/utils"
)

type BotServiceOptions struct {
	webhookURL    string
	webhookPort   string
	token         string
	dispatcher    *Dispatcher
	logger        *logger.Logger
	rueidisClient rueidis.Client
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

func WithRueidisClient(client rueidis.Client) CallOption {
	return func(o *BotServiceOptions) {
		o.rueidisClient = client
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
			b.dispatcher.Dispatch(b.BotAPI, update, b.opts.rueidisClient)
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

func (b *BotService) Bot() *Bot {
	return &Bot{
		BotAPI:        b.BotAPI,
		logger:        b.logger,
		rueidisClient: b.opts.rueidisClient,
	}
}

type Bot struct {
	*tgbotapi.BotAPI
	logger        *logger.Logger
	rueidisClient rueidis.Client
}

func (b *Bot) MaySend(chattable tgbotapi.Chattable) *tgbotapi.Message {
	may := fo.NewMay[tgbotapi.Message]().Use(func(err error, messageArgs ...any) {
		b.logger.Errorf("failed to send %v to telegram: %v", utils.SprintJSON(chattable), err)
	})

	return lo.ToPtr(may.Invoke(b.Send(chattable)))
}

func (b *Bot) MayRequest(chattable tgbotapi.Chattable) *tgbotapi.APIResponse {
	may := fo.NewMay[*tgbotapi.APIResponse]().Use(func(err error, messageArgs ...any) {
		b.logger.Errorf("failed to request %v to telegram: %v", utils.SprintJSON(chattable), err)
	})

	return may.Invoke(b.Request(chattable))
}

func (b *Bot) IsCannotInitiateChatWithUserErr(err error) bool {
	if err == nil {
		return false
	}

	tgbotapiErr, ok := err.(*tgbotapi.Error)
	if !ok {
		return false
	}

	return tgbotapiErr.Code == 403 && tgbotapiErr.Message == "Forbidden: bot can't initiate conversation with a user"
}

func (b *Bot) IsBotWasBlockedByTheUserErr(err error) bool {
	if err == nil {
		return false
	}

	tgbotapiErr, ok := err.(*tgbotapi.Error)
	if !ok {
		return false
	}

	return tgbotapiErr.Code == 403 && tgbotapiErr.Message == "Forbidden: bot was blocked by the user"
}

func (b *Bot) IsBotAdministrator(chatID int64) (bool, error) {
	botMember, err := b.GetChatMember(tgbotapi.GetChatMemberConfig{ChatConfigWithUser: tgbotapi.ChatConfigWithUser{ChatID: chatID, UserID: b.Self.ID}})
	if err != nil {
		return false, err
	}
	if botMember.Status == string(telegram.MemberStatusAdministrator) {
		return true, err
	}

	return false, err
}

func (b *Bot) IsUserMemberStatus(chatID int64, userID int64, status []telegram.MemberStatus) (bool, error) {
	member, err := b.GetChatMember(tgbotapi.GetChatMemberConfig{ChatConfigWithUser: tgbotapi.ChatConfigWithUser{ChatID: chatID, UserID: userID}})
	if err != nil {
		return false, err
	}
	if lo.Contains(status, telegram.MemberStatus(member.Status)) {
		return true, nil
	}

	return false, nil
}

func (b *Bot) PushOneDeleteLaterMessage(forUserID int64, chatID int64, messageID int) error {
	if forUserID == 0 || chatID == 0 || messageID == 0 {
		return nil
	}

	lpushCmd := b.rueidisClient.B().
		Lpush().
		Key(redis.SessionDeleteLaterMessagesForActor1.Format(forUserID)).
		Element(fmt.Sprintf("%d;%d", chatID, messageID)).
		Build()

	exCmd := b.rueidisClient.B().
		Expire().
		Key(redis.SessionDeleteLaterMessagesForActor1.Format(forUserID)).
		Seconds(24 * 60 * 60).
		Build()

	res := b.rueidisClient.DoMulti(context.Background(), lpushCmd, exCmd)
	for _, v := range res {
		if v.Error() != nil {
			return v.Error()
		}
	}

	b.logger.WithFields(logrus.Fields{
		"from_id":    forUserID,
		"chat_id":    chatID,
		"message_id": messageID,
	}).Trace("pushed one delete later message for user")

	return nil
}

func (b *Bot) DeleteAllDeleteLaterMessages(forUserID int64) error {
	if forUserID == 0 {
		return nil
	}

	lrangeCmd := b.rueidisClient.B().
		Lrange().
		Key(redis.SessionDeleteLaterMessagesForActor1.Format(forUserID)).
		Start(0).
		Stop(-1).
		Build()

	elems, err := b.rueidisClient.Do(context.Background(), lrangeCmd).AsStrSlice()
	if err != nil {
		return err
	}
	if len(elems) == 0 {
		return nil
	}

	delCmd := b.rueidisClient.B().
		Del().
		Key(redis.SessionDeleteLaterMessagesForActor1.Format(forUserID)).
		Build()

	res := b.rueidisClient.Do(context.Background(), delCmd)

	for _, v := range elems {
		pairs := strings.Split(v, ";")
		if len(pairs) != 2 {
			continue
		}

		chatID, err := strconv.ParseInt(pairs[0], 10, 64)
		if err != nil {
			continue
		}

		messageID, err := strconv.Atoi(pairs[1])
		if err != nil {
			continue
		}
		if chatID == 0 || messageID == 0 {
			continue
		}

		b.MayRequest(tgbotapi.NewDeleteMessage(chatID, messageID))
		b.logger.WithFields(logrus.Fields{
			"from_id":    forUserID,
			"chat_id":    chatID,
			"message_id": messageID,
		}).Trace("deleted one delete later message for user")
	}

	return res.Error()
}

func (b *Bot) AssignOneCallbackQueryData(route string, data any) (string, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	routeHash := fmt.Sprintf("%x", sha256.Sum256([]byte(route)))[0:16]
	actionHash := fmt.Sprintf("%x", sha256.Sum256(jsonData))[0:16]

	setCmd := b.rueidisClient.B().
		Set().
		Key(redis.CallbackQueryData2.Format(route, actionHash)).
		Value(string(jsonData)).
		ExSeconds(24 * 60 * 60).
		Build()

	err = b.rueidisClient.Do(context.Background(), setCmd).Error()
	if err != nil {
		return fmt.Sprintf("%s;%s", routeHash, actionHash), err
	}

	b.logger.Tracef("assigned callback query data: %s;%s for route %s with data %s", route, actionHash, route, string(jsonData))

	return fmt.Sprintf("%s;%s", routeHash, actionHash), nil
}

func (b *Bot) routeHashAndActionHashFromData(callbackQueryData string) (string, string) {
	handlerIdentifierPairs := strings.Split(callbackQueryData, ";")
	if len(handlerIdentifierPairs) != 2 {
		return "", ""
	}

	return handlerIdentifierPairs[0], handlerIdentifierPairs[1]
}

func (b *Bot) fetchCallbackQueryActionData(route string, dataHash string) (string, error) {
	getCmd := b.rueidisClient.B().
		Get().
		Key(redis.CallbackQueryData2.Format(route, dataHash)).
		Build()

	str, err := b.rueidisClient.Do(context.Background(), getCmd).ToString()
	if err != nil {
		if rueidis.IsRedisNil(err) {
			return "", nil
		}

		return "", err
	}

	return str, nil
}
