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

		enabled, err = m.tgchats.HasChatHistoriesRecapEnabledForGroups(capsule.Payload.ChatID, "")
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
		queueErr := m.tgchats.QueueOneSendChatHistoriesRecapTaskForChatID(capsule.Payload.ChatID, options)
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
	err := m.tgchats.QueueOneSendChatHistoriesRecapTaskForChatID(capsule.Payload.ChatID, options)
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
	m.logger.Info("generating chat histories recap for chat",
		zap.Int64("chat_id", chatID),
		zap.String("module", "autorecap"),
		zap.Int("auto_recap_rates", options.AutoRecapRatesPerDay),
	)

	chat, err := m.botService.GetChat(tgbotapi.ChatInfoConfig{
		ChatConfig: tgbotapi.ChatConfig{
			ChatID: chatID,
		},
	})
	if err != nil {
		m.logger.Error("failed to get chat",
			zap.Int64("chat_id", chatID),
			zap.String("module", "autorecap"),
			zap.Int("auto_recap_rates", options.AutoRecapRatesPerDay),
			zap.Error(err),
		)

		return
	}

	chatType := telegram.ChatType(chat.Type)

	mAutoRecapRatesPerDayHours := map[int]int{
		4: 6,
		3: 8,
		2: 12,
	}

	hours, ok := mAutoRecapRatesPerDayHours[options.AutoRecapRatesPerDay]
	if !ok {
		hours = 6
	}

	mFindChatHistoriesHoursBefore := map[int]func(chatID int64) ([]*ent.ChatHistories, error){
		6:  m.chathistories.FindLast6HourChatHistories,
		8:  m.chathistories.FindLast8HourChatHistories,
		12: m.chathistories.FindLast12HourChatHistories,
	}

	findChatHistories, ok := mFindChatHistoriesHoursBefore[hours]
	if !ok {
		findChatHistories = m.chathistories.FindLast6HourChatHistories
	}

	histories, err := findChatHistories(chatID)
	if err != nil {
		m.logger.Error(fmt.Sprintf("failed to find last %d hour chat histories", hours),
			zap.Int64("chat_id", chatID),
			zap.String("module", "autorecap"),
			zap.Int("auto_recap_rates", options.AutoRecapRatesPerDay),
			zap.Error(err),
		)

		return
	}
	if len(histories) <= 5 {
		m.logger.Warn("no enough chat histories")
		return
	}

	chatTitle := histories[len(histories)-1].ChatTitle

	logID, summarizations, err := m.chathistories.SummarizeChatHistories(chatID, chatType, histories)
	if err != nil {
		m.logger.Error(fmt.Sprintf("failed to summarize last %d hour chat histories", hours),
			zap.Int64("chat_id", chatID),
			zap.String("module", "autorecap"),
			zap.Int("auto_recap_rates", options.AutoRecapRatesPerDay),
			zap.Error(err),
		)

		return
	}

	counts, err := m.chathistories.FindFeedbackRecapsReactionCountsForChatIDAndLogID(chatID, logID)
	if err != nil {
		m.logger.Error("failed to find feedback recaps votes for chat",
			zap.Int64("chat_id", chatID),
			zap.String("module", "autorecap"),
			zap.Int("auto_recap_rates", options.AutoRecapRatesPerDay),
			zap.Error(err),
		)

		return
	}

	inlineKeyboardMarkup, err := m.chathistories.NewVoteRecapInlineKeyboardMarkup(m.botService.Bot(), chatID, logID, counts.UpVotes, counts.DownVotes, counts.Lmao)
	if err != nil {
		m.logger.Error("failed to create vote recap inline keyboard markup",
			zap.Int64("chat_id", chatID),
			zap.String("log_id", logID.String()),
			zap.String("module", "autorecap"),
			zap.Int("auto_recap_rates", options.AutoRecapRatesPerDay),
			zap.Error(err),
		)

		return
	}

	summarizations = lo.Filter(summarizations, func(item string, _ int) bool { return item != "" })
	if len(summarizations) == 0 {
		m.logger.Warn("summarization is empty",
			zap.Int64("chat_id", chatID),
			zap.String("module", "autorecap"),
			zap.Int("auto_recap_rates", options.AutoRecapRatesPerDay),
		)

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
				zap.String("module", "autorecap"),
				zap.Int("auto_recap_rates", options.AutoRecapRatesPerDay),
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
						zap.String("module", "autorecap"),
						zap.Int("auto_recap_rates", options.AutoRecapRatesPerDay),
					)

					return err
				}

				return nil
			})
			if err != nil {
				m.logger.Error("failed to unsubscribe to auto recaps",
					zap.Int64("chat_id", chatID),
					zap.Int("auto_recap_rates", options.AutoRecapRatesPerDay),
					zap.Error(err),
				)
			}

			msg := tgbotapi.NewMessage(subscriber.UserID, fmt.Sprintf("由于您已不再是 <b>%s</b> 的成员，因此已自动帮您取消了您所订阅的聊天记录回顾。", tgbot.EscapeHTMLSymbols(chatTitle)))
			msg.ParseMode = tgbotapi.ModeHTML

			_, err = m.botService.Send(msg)
			if err != nil {
				m.logger.Error("failed to send the auto un-subscription message",
					zap.Int64("user_id", subscriber.UserID),
					zap.Int64("chat_id", chatID),
					zap.Int("auto_recap_rates", options.AutoRecapRatesPerDay),
					zap.Error(err),
				)
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
			content = fmt.Sprintf("%s\n\n(%d/%d)\n%s#recap #recap_auto\n<em>🤖️ Generated by chatGPT</em>",
				strings.Join(b, "\n\n"),
				i+1,
				len(summarizationBatches),
				lo.Ternary(chatType == telegram.ChatTypeGroup, "\n<b>Tips: </b>由于群组不是超级群组（supergroup），因此消息链接引用暂时被禁用了，如果希望使用该功能，请通过短时间内将群组开放为公共群组并还原回私有群组，或通过其他操作将本群组升级为超级群组后，该功能方可恢复正常运作。\n\n", ""),
			)
		} else {
			content = fmt.Sprintf("%s\n\n%s#recap #recap_auto\n<em>🤖️ Generated by chatGPT</em>",
				strings.Join(b, "\n\n"),
				lo.Ternary(chatType == telegram.ChatTypeGroup, "<b>Tips: </b>由于群组不是超级群组（supergroup），因此消息链接引用暂时被禁用了，如果希望使用该功能，请通过短时间内将群组开放为公共群组并还原回私有群组，或通过其他操作将本群组升级为超级群组后，该功能方可恢复正常运作。\n\n", ""),
			)
		}

		for _, targetChat := range targetChats {
			limiter.Take()
			m.logger.Info("sending chat histories recap for chat", zap.Int64("summarized_for_chat_id", chatID), zap.Int64("sending_target_chat_id", targetChat.chatID))

			msg := tgbotapi.NewMessage(targetChat.chatID, "")
			msg.ParseMode = tgbotapi.ModeHTML

			if targetChat.isPrivateSubscriber {
				msg.Text = fmt.Sprintf("您好，这是您订阅的 <b>%s</b> 群组的定时聊天回顾。\n\n%s", tgbot.EscapeHTMLSymbols(chatTitle), content)

				inlineKeyboardMarkup, err := m.chathistories.NewVoteRecapWithUnsubscribeInlineKeyboardMarkup(m.botService.Bot(), chatID, chatTitle, targetChat.chatID, logID, counts.UpVotes, counts.DownVotes, counts.Lmao)
				if err != nil {
					m.logger.Error("failed to assign callback query data",
						zap.Int64("chat_id", chatID),
						zap.Int("auto_recap_rates", options.AutoRecapRatesPerDay),
						zap.Error(err),
					)

					continue
				}

				msg.ReplyMarkup = inlineKeyboardMarkup
			} else {
				msg.Text = content
				msg.ReplyMarkup = inlineKeyboardMarkup
			}

			sentMsg, err := m.botService.Send(msg)
			if err != nil {
				m.logger.Error("failed to send chat histories recap",
					zap.Int64("chat_id", chatID),
					zap.Int("auto_recap_rates", options.AutoRecapRatesPerDay),
					zap.Error(err),
				)
			}

			// Check whether the first message of the batch needs to be pinned, if not, skip the pinning process
			if i != 0 || !options.PinAutoRecapMessage {
				err = m.chathistories.SaveOneTelegramSentMessage(&sentMsg, false)
				if err != nil {
					m.logger.Error("failed to save one telegram sent message",
						zap.Int64("chat_id", chatID),
						zap.Error(err))
				}

				return
			}

			may := fo.NewMay0().Use(func(err error, messageArgs ...any) {
				if len(messageArgs) == 0 {
					m.logger.Error(err.Error())
					return
				}
				prefix, _ := messageArgs[0].(string)

				if len(messageArgs) == 1 {
					m.logger.Error(prefix, zap.Error(err))
					return
				}
				fields := make([]zap.Field, 0)
				fields = append(fields, zap.Error(err))

				for i, v := range messageArgs[1:] {
					field, ok := v.(zap.Field)
					if !ok {
						fields = append(fields, zap.Any(fmt.Sprintf("error_field_%d", i), field))
					} else {
						fields = append(fields, field)
					}
				}

				m.logger.Error(prefix, fields...)
			})

			// Unpin the last pinned message
			lastPinnedMessage, err := m.chathistories.FindLastTelegramPinnedMessage(chatID)
			if err != nil {
				m.logger.Error("failed to find last pinned message",
					zap.Int64("chat_id", chatID),
					zap.Error(err),
				)
			}

			may.Invoke(m.botService.UnpinChatMessage(tgbot.NewUnpinChatMessageConfig(chatID, lastPinnedMessage.MessageID)), "failed to unpin chat message", zap.Int64("chat_id", chatID), zap.Int("message_id", lastPinnedMessage.MessageID))
			may.Invoke(m.chathistories.UpdatePinnedMessage(lastPinnedMessage.ChatID, lastPinnedMessage.MessageID, false), "failed to save one telegram sent message", zap.Int64("chat_id", lastPinnedMessage.ChatID), zap.Int("message_id", lastPinnedMessage.MessageID))
			may.Invoke(m.botService.PinChatMessage(tgbot.NewPinChatMessageConfig(chatID, sentMsg.MessageID)), "failed to pin chat message", zap.Int64("chat_id", chatID), zap.Int("message_id", sentMsg.MessageID))
			may.Invoke(m.chathistories.SaveOneTelegramSentMessage(&sentMsg, true), "failed to save one telegram sent message")
		}
	}
}
