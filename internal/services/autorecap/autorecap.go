package autorecap

import (
	"context"
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/nekomeowww/fo"
	"github.com/nekomeowww/timecapsule/v2"
	"github.com/samber/lo"
	"github.com/sourcegraph/conc/pool"
	"go.uber.org/fx"
	"go.uber.org/multierr"
	"go.uber.org/ratelimit"

	"github.com/nekomeowww/insights-bot/ent"
	"github.com/nekomeowww/insights-bot/internal/datastore"
	"github.com/nekomeowww/insights-bot/internal/models/chathistories"
	"github.com/nekomeowww/insights-bot/internal/models/tgchats"
	"github.com/nekomeowww/insights-bot/pkg/bots/tgbot"
	"github.com/nekomeowww/insights-bot/pkg/logger"
	"github.com/nekomeowww/insights-bot/pkg/types/bots/handlers/recap"
	"github.com/nekomeowww/insights-bot/pkg/types/tgchat"
	"github.com/nekomeowww/insights-bot/pkg/types/timecapsules"
)

type NewAutoRecapParams struct {
	fx.In

	Lifecycle fx.Lifecycle

	Logger        *logger.Logger
	Bot           *tgbot.BotService
	ChatHistories *chathistories.Model
	TgChats       *tgchats.Model
	Digger        *datastore.AutoRecapTimeCapsuleDigger
}

type AutoRecapService struct {
	logger        *logger.Logger
	botService    *tgbot.BotService
	chathistories *chathistories.Model
	tgchats       *tgchats.Model

	digger  *datastore.AutoRecapTimeCapsuleDigger
	started bool
}

func NewAutoRecapService() func(NewAutoRecapParams) (*AutoRecapService, error) {
	return func(params NewAutoRecapParams) (*AutoRecapService, error) {
		service := &AutoRecapService{
			logger:        params.Logger,
			botService:    params.Bot,
			chathistories: params.ChatHistories,
			tgchats:       params.TgChats,
			digger:        params.Digger,
		}

		service.digger.SetHandler(service.sendChatHistoriesRecap)
		service.tgchats.QueueSendChatHistoriesRecapTask()

		return service, nil
	}
}

func (s *AutoRecapService) Check(ctx context.Context) error {
	return lo.Ternary(s.started, nil, fmt.Errorf("auto recap not started yet"))
}

func Run() func(service *AutoRecapService) {
	return func(service *AutoRecapService) {
		service.started = true
	}
}

func (m *AutoRecapService) sendChatHistoriesRecap(
	digger *timecapsule.TimeCapsuleDigger[timecapsules.AutoRecapCapsule],
	capsule *timecapsule.TimeCapsule[timecapsules.AutoRecapCapsule],
) {
	var enabled bool
	var options *ent.TelegramChatRecapsOptions
	var subscribers []*ent.TelegramChatAutoRecapsSubscribers

	may := fo.NewMay[int]()

	_ = may.Invoke(lo.Attempt(10, func(index int) error {
		var err error

		enabled, err = m.tgchats.HasChatHistoriesRecapEnabled(capsule.Payload.ChatID, "")
		if err != nil {
			m.logger.Errorf("failed to check chat histories recap enabled: %v", err)
		}

		return err
	}))
	_ = may.Invoke(lo.Attempt(10, func(index int) error {
		var err error

		options, err = m.tgchats.FindOneRecapsOption(capsule.Payload.ChatID)
		if err != nil {
			m.logger.Errorf("failed to find chat recap options: %v", err)
		}

		return err
	}))
	_ = may.Invoke(lo.Attempt(10, func(index int) error {
		var err error

		subscribers, err = m.tgchats.FindAutoRecapsSubscribers(capsule.Payload.ChatID)
		if err != nil {
			m.logger.Errorf("failed to find chat recap subscribers: %v", err)
		}

		return err
	}))

	may.HandleErrors(func(errs []error) {
		// requeue if failed
		queueErr := m.tgchats.QueueOneSendChatHistoriesRecapTaskForChatID(capsule.Payload.ChatID)
		if queueErr != nil {
			m.logger.Errorf("failed to queue one send chat histories recap task for chat %d: %v", capsule.Payload.ChatID, queueErr)
		}

		m.logger.Errorf("failed to check chat histories recap enabled, options or subscribers: %v", multierr.Combine(errs...))
	})
	if !enabled {
		return
	}

	// always requeue
	err := m.tgchats.QueueOneSendChatHistoriesRecapTaskForChatID(capsule.Payload.ChatID)
	if err != nil {
		m.logger.Errorf("failed to queue one send chat histories recap task for chat %d: %v", capsule.Payload.ChatID, err)
	}
	if options != nil && tgchat.AutoRecapSendMode(options.AutoRecapSendMode) == tgchat.AutoRecapSendModeOnlyPrivateSubscriptions && len(subscribers) == 0 {
		return
	}

	pool := pool.New().WithMaxGoroutines(20)
	pool.Go(func() {
		m.summarize(capsule.Payload.ChatID, options, subscribers)
	})
}

func (m *AutoRecapService) summarize(chatID int64, options *ent.TelegramChatRecapsOptions, subscribers []*ent.TelegramChatAutoRecapsSubscribers) {
	m.logger.Infof("generating chat histories recap for chat %d", chatID)

	histories, err := m.chathistories.FindLastSixHourChatHistories(chatID)
	if err != nil {
		m.logger.Errorf("failed to find last six hour chat histories: %v", err)
		return
	}
	if len(histories) <= 5 {
		m.logger.Warn("no enough chat histories")
		return
	}

	chatTitle := histories[len(histories)-1].ChatTitle

	summarizations, err := m.chathistories.SummarizeChatHistories(chatID, histories)
	if err != nil {
		m.logger.Errorf("failed to summarize last six hour chat histories: %v", err)
		return
	}

	summarizations = lo.Filter(summarizations, func(item string, _ int) bool { return item != "" })
	if len(summarizations) == 0 {
		m.logger.Warn("summarization is empty")
		return
	}

	for i, s := range summarizations {
		summarizations[i] = tgbot.ReplaceMarkdownTitlesToTelegramBoldElement(s)
	}

	summarizationBatches := tgbot.SplitMessagesAgainstLengthLimitIntoMessageGroups(summarizations)

	limiter := ratelimit.New(5)

	type targetChat struct {
		chatID              int64
		isPrivateSubscriber bool
	}

	targetChats := make([]targetChat, 0)

	if options == nil || tgchat.AutoRecapSendMode(options.AutoRecapSendMode) == tgchat.AutoRecapSendModePublicly {
		targetChats = append(targetChats, targetChat{
			chatID:              chatID,
			isPrivateSubscriber: false,
		})
	}

	for _, subscriber := range subscribers {
		targetChats = append(targetChats, targetChat{
			chatID:              subscriber.UserID,
			isPrivateSubscriber: true,
		})
	}

	for i, b := range summarizationBatches {
		var content string
		if len(summarizationBatches) > 1 {
			content = fmt.Sprintf("%s\n\n(%d/%d)\n#recap #recap_auto\n<em>🤖️ Generated by chatGPT</em>", strings.Join(b, "\n\n"), i+1, len(summarizationBatches))
		} else {
			content = fmt.Sprintf("%s\n\n#recap #recap_auto\n<em>🤖️ Generated by chatGPT</em>", strings.Join(b, "\n\n"))
		}

		for _, targetChat := range targetChats {
			limiter.Take()
			m.logger.Infof("sending chat histories recap for chat %d", targetChat.chatID)

			msg := tgbotapi.NewMessage(targetChat.chatID, "")
			msg.ParseMode = tgbotapi.ModeHTML

			if targetChat.isPrivateSubscriber {
				msg.Text = fmt.Sprintf("您好，这是您订阅的 <b>%s</b> 群组的定时聊天回顾。\n\n%s", chatTitle, content)

				buttonData, err := m.botService.Bot().AssignOneCallbackQueryData("recap/unsubscribe_recap", recap.UnsubscribeRecapActionData{
					ChatID:    chatID,
					ChatTitle: chatTitle,
					FromID:    targetChat.chatID,
				})
				if err != nil {
					m.logger.Errorf("failed to assign callback query data: %v", err)
					continue
				}

				msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("取消订阅", buttonData)))
			} else {
				msg.Text = content
			}

			_, err = m.botService.Send(msg)
			if err != nil {
				m.logger.Errorf("failed to send chat histories recap: %v", err)
			}
		}
	}
}
