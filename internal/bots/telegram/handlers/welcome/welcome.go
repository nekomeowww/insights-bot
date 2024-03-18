package welcome

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/nekomeowww/fo"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/nekomeowww/insights-bot/internal/models/chathistories"
	"github.com/nekomeowww/insights-bot/internal/models/logs"
	"github.com/nekomeowww/insights-bot/internal/models/tgchats"
	"github.com/nekomeowww/insights-bot/pkg/bots/tgbot"
	"github.com/nekomeowww/insights-bot/pkg/i18n"
	"github.com/nekomeowww/insights-bot/pkg/logger"
	"github.com/nekomeowww/insights-bot/pkg/types/telegram"
)

func NewModules() fx.Option {
	return fx.Options(
		fx.Provide(NewHandlers()),
	)
}

type NewHandlersParams struct {
	fx.In

	TgChats       *tgchats.Model
	ChatHistories *chathistories.Model
	Logs          *logs.Model
	Logger        *logger.Logger
	I18n          *i18n.I18n
}

type Handlers struct {
	tgchats       *tgchats.Model
	chatHistories *chathistories.Model
	logs          *logs.Model
	logger        *logger.Logger
	i18n          *i18n.I18n
}

func NewHandlers() func(param NewHandlersParams) *Handlers {
	return func(param NewHandlersParams) *Handlers {
		return &Handlers{
			tgchats:       param.TgChats,
			chatHistories: param.ChatHistories,
			logs:          param.Logs,
			logger:        param.Logger,
			i18n:          param.I18n,
		}
	}
}

func (h *Handlers) Install(dispatcher *tgbot.Dispatcher) {
	dispatcher.OnMyChatMember(tgbot.NewHandler(h.handleChatMember))
}

func (h *Handlers) handleChatMember(c *tgbot.Context) (tgbot.Response, error) {
	if c.Update.MyChatMember == nil {
		return nil, nil
	}
	if telegram.MemberStatus(c.Update.MyChatMember.NewChatMember.Status) == telegram.MemberStatusLeft {
		h.handleBotLeftChat(c.Update.MyChatMember.Chat.ID)
		return nil, nil
	}
	if telegram.MemberStatus(c.Update.MyChatMember.NewChatMember.Status) == telegram.MemberStatusMember {
		h.handleBotJoinChat(c)
		return nil, nil
	}

	return nil, nil
}

func (h *Handlers) handleBotLeftChat(chatID int64) {
	may := fo.
		NewMay0().
		Use(fo.WithLogFuncHandler(func(a ...any) {
			h.logger.Error(fmt.Sprint(a...), zap.Int64("chat_id", chatID))
		}))

	may.Invoke(func() error {
		err := h.tgchats.DeleteAllSubscribersByChatID(chatID)
		if err != nil {
			return fmt.Errorf("failed to delete all subscribers by chat id %d: %w", chatID, err)
		}

		h.logger.Info("deleted all subscribers by chat id", zap.Int64("chat_id", chatID))

		return nil
	}())
	may.Invoke(func() error {
		err := h.tgchats.DeleteOneFeatureFlagByChatID(chatID)
		if err != nil {
			return fmt.Errorf("failed to delete one feature flag by chat id %d: %w", chatID, err)
		}

		h.logger.Info("deleted all feature flags by chat id", zap.Int64("chat_id", chatID))

		return nil
	}())
	may.Invoke(func() error {
		err := h.tgchats.DeleteOneOptionByChatID(chatID)
		if err != nil {
			return fmt.Errorf("failed to delete one option by chat id %d: %w", chatID, err)
		}

		h.logger.Info("deleted all chat options by chat id", zap.Int64("chat_id", chatID))

		return nil
	}())
	may.Invoke(func() error {
		err := h.chatHistories.DeleteAllChatHistoriesByChatID(chatID)
		if err != nil {
			return fmt.Errorf("failed to delete all chat histories by chat id %d: %w", chatID, err)
		}

		h.logger.Info("deleted all chat histories by chat id", zap.Int64("chat_id", chatID))

		return nil
	}())
	may.Invoke(func() error {
		err := h.logs.PruneAllLogsContentForChatID(chatID)
		if err != nil {
			return fmt.Errorf("failed to prune all related content for chat id %d: %w", chatID, err)
		}

		h.logger.Info("pruned all related metrics logs for chat id", zap.Int64("chat_id", chatID))

		return nil
	}())

	h.logger.Info("pruned all related content for chat id", zap.Int64("chat_id", chatID))
}

func (h *Handlers) handleBotJoinChat(c *tgbot.Context) {
	chatID := c.Update.MyChatMember.Chat.ID
	chatType := telegram.ChatType(c.Update.MyChatMember.Chat.Type)
	chatTitle := c.Update.MyChatMember.Chat.Title
	language := c.Update.MyChatMember.From.LanguageCode

	hasJoinedBefore, err := h.tgchats.HasJoinedGroupsBefore(chatID, chatTitle)
	if err != nil {
		h.logger.Error("failed to check if bot has joined groups before",
			zap.Error(err),
			zap.Int64("chat_id", chatID),
			zap.String("chat_title", chatTitle),
			zap.String("chat_type", string(chatType)),
			zap.String("language", language),
		)

		return
	}
	if hasJoinedBefore {
		return
	}

	err = h.tgchats.SetLanguageForGroups(chatID, chatType, chatTitle, language)
	if err != nil {
		h.logger.Error("failed to set language for groups",
			zap.Error(err),
			zap.Int64("chat_id", chatID),
			zap.String("chat_title", chatTitle),
			zap.String("chat_type", string(chatType)),
			zap.String("language", language),
		)
	}

	msg := tgbotapi.NewMessage(
		chatID,
		h.i18n.TWithLanguage(
			language,
			"modules.telegram.welcome.messageNormalGroup",
			i18n.M{
				"Username": c.Bot.Self.UserName,
			},
		),
	)

	msg.ParseMode = tgbotapi.ModeHTML

	c.Bot.MaySend(msg)
}
