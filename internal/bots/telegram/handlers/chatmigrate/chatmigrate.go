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
}

type Handlers struct {
	tgchats       *tgchats.Model
	chatHistories *chathistories.Model
	logs          *logs.Model
	logger        *logger.Logger
}

func NewHandlers() func(param NewHandlersParams) *Handlers {
	return func(param NewHandlersParams) *Handlers {
		return &Handlers{
			tgchats:       param.TgChats,
			chatHistories: param.ChatHistories,
			logs:          param.Logs,
			logger:        param.Logger,
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
			h.logger.Error(fmt.Sprint(a...))
		}))

	may.Invoke(func() error {
		err := h.tgchats.MigrateFeatureFlagsOfChatFromChatIDToChatID(fromChatID, toChatID)
		if err != nil {
			return fmt.Errorf("failed to migrate all feature flags of chat from chat id %d to chat id %d: %w", fromChatID, toChatID, err)
		}

		return nil
	}())
	may.Invoke(func() error {
		err := h.tgchats.MigrateOptionOfChatFromChatIDToChatID(fromChatID, toChatID)
		if err != nil {
			return fmt.Errorf("failed to migrate all options of chat from chat id %d to chat id %d: %w", fromChatID, toChatID, err)
		}

		return nil
	}())
	may.Invoke(func() error {
		err := h.tgchats.MigrateSubscribersOfChatFromChatIDToChatID(fromChatID, toChatID)
		if err != nil {
			return fmt.Errorf("failed to migrate all subscribers of chat from chat id %d to chat id %d: %w", fromChatID, toChatID, err)
		}

		return nil
	}())
	may.Invoke(func() error {
		err := h.chatHistories.MigrateChatHistoriesOfChatFromChatIDToChatID(fromChatID, toChatID)
		if err != nil {
			return fmt.Errorf("failed to migrate all chat histories of chat from chat id %d to chat id %d: %w", fromChatID, toChatID, err)
		}

		return nil
	}())
	may.Invoke(func() error {
		err := h.logs.MigrateLogsOfChatFromChatIDToChatID(fromChatID, toChatID)
		if err != nil {
			return fmt.Errorf("failed to migrate all logs of chat from chat id %d to chat id %d: %w", fromChatID, toChatID, err)
		}

		return nil
	}())

	h.logger.Info(fmt.Sprintf("成功将群组 %d 的所有数据迁移到群组 %d", fromChatID, toChatID), zap.Int64("from_chat_id", fromChatID), zap.Int64("to_chat_id", toChatID))

	message := tgbotapi.NewMessage(toChatID,
		fmt.Sprintf(""+
			"%s @%s 监测到您的群组已从 <b>群组（group）</b> 升级为了 <b>超级群组（supergroup）</b>，届时"+
			"，群组的 ID 将会发生变更，<b>现已自动将过去的历史记录和数据留存自动迁移到了新的群组 ID 名下</b>，"+
			"之前的设置将会保留并继续沿用，不过需要注意的是，由于 Telegram 官方的限制，迁移事件前的消息 ID 将无"+
			"法与今后发送的消息 ID 相兼容，所以当下一次总结消息时将不会包含在迁移事件发生前所发送的消息，由此带来"+
			"的不便敬请谅解。",
			tgbot.FullNameFromFirstAndLastName(c.Bot.Self.FirstName, c.Bot.Self.LastName),
			c.Bot.Self.UserName,
		))
	message.ParseMode = tgbotapi.ModeHTML

	c.Bot.MaySend(message)

	return nil, nil
}
