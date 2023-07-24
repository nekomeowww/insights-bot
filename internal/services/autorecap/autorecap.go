package autorecap

import (
	"context"
	"fmt"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/nekomeowww/fo"
	"github.com/nekomeowww/timecapsule/v2"
	"github.com/samber/lo"
	"github.com/sourcegraph/conc/pool"
	"go.uber.org/fx"
	"go.uber.org/multierr"
	"go.uber.org/ratelimit"
	"go.uber.org/zap"

	"github.com/nekomeowww/insights-bot/ent"
	"github.com/nekomeowww/insights-bot/internal/datastore"
	"github.com/nekomeowww/insights-bot/internal/models/chathistories"
	"github.com/nekomeowww/insights-bot/internal/models/tgchats"
	"github.com/nekomeowww/insights-bot/pkg/bots/tgbot"
	"github.com/nekomeowww/insights-bot/pkg/logger"
	"github.com/nekomeowww/insights-bot/pkg/types/telegram"
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

		service.digger.SetHandler(service.sendChatHistoriesRecapTimeCapsuleHandler)
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

func (m *AutoRecapService) sendChatHistoriesRecapTimeCapsuleHandler(
	digger *timecapsule.TimeCapsuleDigger[timecapsules.AutoRecapCapsule],
	capsule *timecapsule.TimeCapsule[timecapsules.AutoRecapCapsule],
) {
	m.logger.Debug("send chat histories recap time capsule handler invoked", zap.Int64("chat_id", capsule.Payload.ChatID))

	var enabled bool
	var options *ent.TelegramChatRecapsOptions
	var subscribers []*ent.TelegramChatAutoRecapsSubscribers

	may := fo.NewMay[int]()

	_ = may.Invoke(lo.Attempt(10, func(index int) error {
		var err error

		enabled, err = m.tgchats.HasChatHistoriesRecapEnabled(capsule.Payload.ChatID, "")
		if err != nil {
			m.logger.Error("failed to check chat histories recap enabled", zap.Error(err))
		}

		return err
	}))
	_ = may.Invoke(lo.Attempt(10, func(index int) error {
		var err error

		options, err = m.tgchats.FindOneRecapsOption(capsule.Payload.ChatID)
		if err != nil {
			m.logger.Error("failed to find chat recap options", zap.Error(err))
		}

		return err
	}))
	_ = may.Invoke(lo.Attempt(10, func(index int) error {
		var err error

		subscribers, err = m.tgchats.FindAutoRecapsSubscribers(capsule.Payload.ChatID)
		if err != nil {
			m.logger.Error("failed to find chat recap subscribers", zap.Error(err))
		}

		return err
	}))

	may.HandleErrors(func(errs []error) {
		// requeue if failed
		queueErr := m.tgchats.QueueOneSendChatHistoriesRecapTaskForChatID(capsule.Payload.ChatID)
		if queueErr != nil {
			m.logger.Error("failed to queue one send chat histories recap task for chat", zap.Int64("chat_id", capsule.Payload.ChatID), zap.Error(queueErr))
		}

		m.logger.Error("failed to check chat histories recap enabled, options or subscribers", zap.Error(multierr.Combine(errs...)))
	})
	if !enabled {
		m.logger.Debug("chat histories recap disabled, skipping...", zap.Int64("chat_id", capsule.Payload.ChatID))

		return
	}

	// always requeue
	err := m.tgchats.QueueOneSendChatHistoriesRecapTaskForChatID(capsule.Payload.ChatID)
	if err != nil {
		m.logger.Error("failed to queue one send chat histories recap task for chat", zap.Int64("chat_id", capsule.Payload.ChatID), zap.Error(err))
	}
	if options != nil && tgchat.AutoRecapSendMode(options.AutoRecapSendMode) == tgchat.AutoRecapSendModeOnlyPrivateSubscriptions && len(subscribers) == 0 {
		m.logger.Debug("chat histories recap send mode is only private subscriptions, but no subscribers, skipping...", zap.Int64("chat_id", capsule.Payload.ChatID))

		return
	}

	pool := pool.New().WithMaxGoroutines(20)
	pool.Go(func() {
		m.summarize(capsule.Payload.ChatID, options, subscribers)
	})
}

func (m *AutoRecapService) summarize(chatID int64, options *ent.TelegramChatRecapsOptions, subscribers []*ent.TelegramChatAutoRecapsSubscribers) {
	m.logger.Info("generating chat histories recap for chat", zap.Int64("chat_id", chatID))

	histories, err := m.chathistories.FindLastSixHourChatHistories(chatID)
	if err != nil {
		m.logger.Error("failed to find last six hour chat histories", zap.Error(err), zap.Int64("chat_id", chatID))
		return
	}
	if len(histories) <= 5 {
		m.logger.Warn("no enough chat histories")
		return
	}

	chatTitle := histories[len(histories)-1].ChatTitle

	logID, summarizations, err := m.chathistories.SummarizeChatHistories(chatID, histories)
	if err != nil {
		m.logger.Error("failed to summarize last six hour chat histories", zap.Error(err), zap.Int64("chat_id", chatID))
		return
	}

	counts, err := m.chathistories.FindFeedbackRecapsReactionCountsForChatIDAndLogID(chatID, logID)
	if err != nil {
		m.logger.Error("failed to find feedback recaps votes for chat", zap.Error(err), zap.Int64("chat_id", chatID))
		return
	}

	inlineKeyboardMarkup, err := m.chathistories.NewVoteRecapInlineKeyboardMarkup(m.botService.Bot(), chatID, logID, counts.UpVotes, counts.DownVotes, counts.Lmao)
	if err != nil {
		m.logger.Error("failed to create vote recap inline keyboard markup", zap.Error(err), zap.Int64("chat_id", chatID), zap.String("log_id", logID.String()))
		return
	}

	summarizations = lo.Filter(summarizations, func(item string, _ int) bool { return item != "" })
	if len(summarizations) == 0 {
		m.logger.Warn("summarization is empty", zap.Int64("chat_id", chatID))
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
		member, err := m.botService.GetChatMember(tgbotapi.GetChatMemberConfig{
			ChatConfigWithUser: tgbotapi.ChatConfigWithUser{
				ChatID: chatID,
				UserID: subscriber.UserID,
			},
		})
		if err != nil {
			m.logger.Error("failed to get chat member", zap.Error(err), zap.Int64("chat_id", chatID))
			continue
		}
		if !lo.Contains([]telegram.MemberStatus{
			telegram.MemberStatusAdministrator,
			telegram.MemberStatusCreator,
			telegram.MemberStatusMember,
			telegram.MemberStatusRestricted,
		}, telegram.MemberStatus(member.Status)) {
			m.logger.Warn("subscriber is not a member, auto unsubscribing...",
				zap.String("status", member.Status),
				zap.Int64("chat_id", chatID),
				zap.Int64("user_id", subscriber.UserID),
			)

			_, _, err := lo.AttemptWithDelay(1000, time.Minute, func(iter int, _ time.Duration) error {
				err := m.tgchats.UnsubscribeToAutoRecaps(chatID, subscriber.UserID)
				if err != nil {
					m.logger.Error("failed to auto unsubscribe to auto recaps",
						zap.Error(err),
						zap.String("status", member.Status),
						zap.Int64("chat_id", chatID),
						zap.Int64("user_id", subscriber.UserID),
						zap.Int("iter", iter),
						zap.Int("max_iter", 100),
					)

					return err
				}

				return nil
			})
			if err != nil {
				m.logger.Error("failed to unsubscribe to auto recaps", zap.Error(err), zap.Int64("chat_id", chatID))
			}

			msg := tgbotapi.NewMessage(subscriber.UserID, fmt.Sprintf("由于您已不再是 <b>%s</b> 的成员，因此已自动帮您取消了您所订阅的聊天记录回顾。", tgbot.EscapeHTMLSymbols(chatTitle)))
			msg.ParseMode = tgbotapi.ModeHTML

			_, err = m.botService.Send(msg)
			if err != nil {
				m.logger.Error("failed to send the auto un-subscription message", zap.Error(err), zap.Int64("user_id", subscriber.UserID), zap.Int64("chat_id", chatID))
			}

			continue
		}

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
			m.logger.Info("sending chat histories recap for chat", zap.Int64("summarized_for_chat_id", chatID), zap.Int64("sending_target_chat_id", targetChat.chatID))

			msg := tgbotapi.NewMessage(targetChat.chatID, "")
			msg.ParseMode = tgbotapi.ModeHTML

			if targetChat.isPrivateSubscriber {
				msg.Text = fmt.Sprintf("您好，这是您订阅的 <b>%s</b> 群组的定时聊天回顾。\n\n%s", tgbot.EscapeHTMLSymbols(chatTitle), content)

				inlineKeyboardMarkup, err := m.chathistories.NewVoteRecapWithUnsubscribeInlineKeyboardMarkup(m.botService.Bot(), targetChat.chatID, chatTitle, targetChat.chatID, logID, counts.UpVotes, counts.DownVotes, counts.Lmao)
				if err != nil {
					m.logger.Error("failed to assign callback query data", zap.Error(err), zap.Int64("chat_id", chatID))
					continue
				}

				msg.ReplyMarkup = inlineKeyboardMarkup
			} else {
				msg.Text = content
				msg.ReplyMarkup = inlineKeyboardMarkup
			}

			_, err = m.botService.Send(msg)
			if err != nil {
				m.logger.Error("failed to send chat histories recap", zap.Error(err), zap.Int64("chat_id", chatID))
			}
		}
	}
}
