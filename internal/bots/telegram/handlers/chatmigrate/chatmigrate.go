package chatmigrate

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
	dispatcher.OnChatMigrationFrom(tgbot.NewHandler(h.handleChatMigrationFrom))
}

func (h *Handlers) handleChatMigrationFrom(c *tgbot.Context) (tgbot.Response, error) {
	if c.Update.Message == nil {
		return nil, nil
	}
	if c.Update.Message.MigrateFromChatID == 0 {
		return nil, nil
	}

	fromChatID := c.Update.Message.MigrateFromChatID
	toChatID := c.Update.Message.Chat.ID

	may := fo.
		NewMay0().
		Use(fo.WithLogFuncHandler(func(a ...any) {
			h.logger.Error(fmt.Sprint(a...), zap.Int64("from_chat_id", fromChatID), zap.Int64("to_chat_id", toChatID))
		}))

	may.Invoke(func() error {
		err := h.tgchats.MigrateFeatureFlagsOfChatFromChatIDToChatID(fromChatID, toChatID)
		if err != nil {
			return fmt.Errorf("failed to migrate all feature flags of chat from chat id %d to chat id %d: %w", fromChatID, toChatID, err)
		}

		h.logger.Info("successfully migrated feature flags of chat",
			zap.Int64("from_chat_id", fromChatID),
			zap.Int64("to_chat_id", toChatID),
		)

		return nil
	}())
	may.Invoke(func() error {
		err := h.tgchats.MigrateOptionOfChatFromChatIDToChatID(fromChatID, toChatID)
		if err != nil {
			return fmt.Errorf("failed to migrate all options of chat from chat id %d to chat id %d: %w", fromChatID, toChatID, err)
		}

		h.logger.Info("successfully migrated options of chat",
			zap.Int64("from_chat_id", fromChatID),
			zap.Int64("to_chat_id", toChatID),
		)

		return nil
	}())
	may.Invoke(func() error {
		err := h.tgchats.MigrateSubscribersOfChatFromChatIDToChatID(fromChatID, toChatID)
		if err != nil {
			return fmt.Errorf("failed to migrate all subscribers of chat from chat id %d to chat id %d: %w", fromChatID, toChatID, err)
		}

		h.logger.Info("successfully migrated subscribers of chat",
			zap.Int64("from_chat_id", fromChatID),
			zap.Int64("to_chat_id", toChatID),
		)

		return nil
	}())
	may.Invoke(func() error {
		err := h.chatHistories.MigrateChatHistoriesOfChatFromChatIDToChatID(fromChatID, toChatID)
		if err != nil {
			return fmt.Errorf("failed to migrate all chat histories of chat from chat id %d to chat id %d: %w", fromChatID, toChatID, err)
		}

		h.logger.Info("successfully migrated chat histories of chat",
			zap.Int64("from_chat_id", fromChatID),
			zap.Int64("to_chat_id", toChatID),
		)

		return nil
	}())
	may.Invoke(func() error {
		err := h.logs.MigrateLogsOfChatFromChatIDToChatID(fromChatID, toChatID)
		if err != nil {
			return fmt.Errorf("failed to migrate all logs of chat from chat id %d to chat id %d: %w", fromChatID, toChatID, err)
		}

		h.logger.Info("successfully migrated logs of chat",
			zap.Int64("from_chat_id", fromChatID),
			zap.Int64("to_chat_id", toChatID),
		)

		return nil
	}())

	h.logger.Info(
		fmt.Sprintf("Successfully migrated all data with chat id %d to supergroup chat id %d", fromChatID, toChatID),
		zap.Int64("from_chat_id", fromChatID),
		zap.Int64("to_chat_id", toChatID),
	)

	language, err := h.tgchats.FindLanguageForGroups(toChatID, "")
	if err != nil {
		h.logger.Error("failed to find language for groups",
			zap.Int64("chat_id", toChatID),
			zap.Error(err),
		)
	}

	message := tgbotapi.NewMessage(
		toChatID,
		h.i18n.TWithLanguage(language, "modules.chatMigration.notification", i18n.M{
			"Name":     tgbot.FullNameFromFirstAndLastName(c.Bot.Self.FirstName, c.Bot.Self.LastName),
			"Username": c.Bot.Self.UserName,
		}),
	)

	message.ParseMode = tgbotapi.ModeHTML

	c.Bot.MaySend(message)

	return nil, nil
}
