package tgbot

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/nekomeowww/insights-bot/pkg/i18n"
	"github.com/nekomeowww/insights-bot/pkg/logger"
	"github.com/nekomeowww/insights-bot/pkg/types/telegram"
	"github.com/redis/rueidis"
	"go.uber.org/zap"
)

type UpdateType string

const (
	UpdateTypeUnknown            UpdateType = "unknown"
	UpdateTypeMessage            UpdateType = "message"
	UpdateTypeEditedMessage      UpdateType = "edited_message"
	UpdateTypeChannelPost        UpdateType = "channel_post"
	UpdateTypeEditedChannelPost  UpdateType = "edited_channel_post"
	UpdateTypeInlineQuery        UpdateType = "inline_query"
	UpdateTypeChosenInlineResult UpdateType = "chosen_inline_result"
	UpdateTypeCallbackQuery      UpdateType = "callback_query"
	UpdateTypeShippingQuery      UpdateType = "shipping_query"
	UpdateTypePreCheckoutQuery   UpdateType = "pre_checkout_query"
	UpdateTypePoll               UpdateType = "poll"
	UpdateTypePollAnswer         UpdateType = "poll_answer"
	UpdateTypeMyChatMember       UpdateType = "my_chat_member"
	UpdateTypeChatMember         UpdateType = "chat_member"
	UpdateTypeLeftChatMember     UpdateType = "left_chat_member"
	UpdateTypeNewChatMembers     UpdateType = "new_chat_members"
	UpdateTypeChatJoinRequest    UpdateType = "chat_join_request"
	UpdateTypeChatMigrationFrom  UpdateType = "chat_migration_from"
	UpdateTypeChatMigrationTo    UpdateType = "chat_migration_to"
)

type Context struct {
	Bot    *Bot
	Update tgbotapi.Update
	Logger *logger.Logger
	I18n   *i18n.I18n

	mutex         sync.Mutex
	rueidisClient rueidis.Client

	abort bool

	isCallbackQuery                       bool
	callbackQueryHandlerRoute             string
	callbackQueryHandlerRouteHash         string
	callbackQueryHandlerActionData        string
	callbackQueryHandlerActionDataHash    string
	callbackQueryHandlerActionDataIsEmpty bool
}

func NewContext(bot *tgbotapi.BotAPI, update tgbotapi.Update, logger *logger.Logger, i18n *i18n.I18n, rueidisClient rueidis.Client) *Context {
	return &Context{
		Bot:                                   &Bot{BotAPI: bot, logger: logger, rueidisClient: rueidisClient},
		Update:                                update,
		Logger:                                logger,
		I18n:                                  i18n,
		rueidisClient:                         rueidisClient,
		isCallbackQuery:                       false,
		callbackQueryHandlerRoute:             "",
		callbackQueryHandlerRouteHash:         "",
		callbackQueryHandlerActionData:        "",
		callbackQueryHandlerActionDataHash:    "",
		callbackQueryHandlerActionDataIsEmpty: false,
	}
}

func (c *Context) UpdateType() UpdateType {
	switch {
	case c.Update.Message != nil:
		switch {
		case c.Update.Message.NewChatMembers != nil:
			return UpdateTypeNewChatMembers
		case c.Update.Message.LeftChatMember != nil:
			return UpdateTypeLeftChatMember
		case c.Update.Message.MigrateFromChatID != 0:
			return UpdateTypeChatMigrationFrom
		case c.Update.Message.MigrateToChatID != 0:
			return UpdateTypeChatMigrationTo
		default:
			return UpdateTypeMessage
		}
	case c.Update.EditedMessage != nil:
		return UpdateTypeEditedMessage
	case c.Update.ChannelPost != nil:
		return UpdateTypeChannelPost
	case c.Update.EditedChannelPost != nil:
		return UpdateTypeEditedChannelPost
	case c.Update.InlineQuery != nil:
		return UpdateTypeInlineQuery
	case c.Update.ChosenInlineResult != nil:
		return UpdateTypeChosenInlineResult
	case c.Update.CallbackQuery != nil:
		return UpdateTypeCallbackQuery
	case c.Update.ShippingQuery != nil:
		return UpdateTypeShippingQuery
	case c.Update.PreCheckoutQuery != nil:
		return UpdateTypePreCheckoutQuery
	case c.Update.Poll != nil:
		return UpdateTypePoll
	case c.Update.PollAnswer != nil:
		return UpdateTypePollAnswer
	case c.Update.MyChatMember != nil:
		return UpdateTypeMyChatMember
	case c.Update.ChatMember != nil:
		return UpdateTypeChatMember
	case c.Update.ChatJoinRequest != nil:
		return UpdateTypeChatJoinRequest
	default:
		return UpdateTypeUnknown
	}
}

func (c *Context) Abort() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.abort = true
}

func (c *Context) IsAborted() bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return c.abort
}

func (c *Context) T(key string, args ...any) string {
	return c.I18n.TWithLanguage(c.Language(), key, args...)
}

func (c *Context) Language() string {
	if c.Update.SentFrom() == nil {
		c.Logger.Warn("update.SentFrom() is nil, fallback to 'en' language.")
		return "en"
	}

	languageCode := c.Update.SentFrom().LanguageCode
	if languageCode == "" {
		c.Logger.Warn("update.SentFrom().LanguageCode is empty, fallback to 'en' language.")
		return "en"
	}

	c.Logger.Debug("resolved language code", zap.String("languageCode", languageCode))

	return languageCode
}

func (c *Context) initForCallbackQuery(route, routeHash, actionDataHash string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.isCallbackQuery = true
	c.callbackQueryHandlerRoute = route
	c.callbackQueryHandlerRouteHash = routeHash
	c.callbackQueryHandlerActionDataHash = actionDataHash
}

func (c *Context) fetchActionDataForCallbackQueryHandler() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.callbackQueryHandlerActionDataIsEmpty = true

	if !c.isCallbackQuery {
		return fmt.Errorf("not a callback query")
	}
	if c.callbackQueryHandlerRouteHash == "" {
		return fmt.Errorf("callback query handler route hash is empty")
	}
	if c.callbackQueryHandlerActionDataHash == "" {
		return fmt.Errorf("callback query handler action data hash is empty")
	}
	if c.rueidisClient == nil {
		return fmt.Errorf("rueidis client is nil")
	}

	str, err := c.Bot.fetchCallbackQueryActionData(c.callbackQueryHandlerRoute, c.callbackQueryHandlerActionDataHash)
	if err != nil {
		return err
	}

	c.callbackQueryHandlerActionData = str
	c.callbackQueryHandlerActionDataIsEmpty = false

	return nil
}

func (c *Context) BindFromCallbackQueryData(dst any) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.callbackQueryHandlerActionDataIsEmpty {
		return errors.New("empty action data")
	}

	return json.Unmarshal([]byte(c.callbackQueryHandlerActionData), dst)
}

func (c *Context) IsBotAdministrator() (bool, error) {
	return c.Bot.IsBotAdministrator(c.Update.FromChat().ID)
}

func (c *Context) IsUserMemberStatus(userID int64, status []telegram.MemberStatus) (bool, error) {
	return c.Bot.IsUserMemberStatus(c.Update.FromChat().ID, userID, status)
}

func (c *Context) RateLimitForCommand(chatID int64, command string, rate int64, perDuration time.Duration) (int64, time.Duration, bool, error) {
	return c.Bot.RateLimitForCommand(chatID, command, rate, perDuration)
}

func (c *Context) NewMessage(message string) MessageResponse {
	return NewMessage(c.Update.FromChat().ID, message)
}

func (c *Context) NewMessageReplyTo(message string, replyToMessageID int) MessageResponse {
	return NewMessageReplyTo(c.Update.FromChat().ID, message, replyToMessageID)
}

func (c *Context) NewEditMessageText(messageID int, text string) EditMessageResponse {
	return NewEditMessageText(c.Update.FromChat().ID, messageID, text)
}

func (c *Context) NewEditMessageTextAndReplyMarkup(messageID int, text string, replyMarkup tgbotapi.InlineKeyboardMarkup) EditMessageResponse {
	return NewEditMessageTextAndReplyMarkup(c.Update.FromChat().ID, messageID, text, replyMarkup)
}

func (c *Context) NewEditMessageReplyMarkup(messageID int, replyMarkup tgbotapi.InlineKeyboardMarkup) EditMessageResponse {
	return NewEditMessageReplyMarkup(c.Update.FromChat().ID, messageID, replyMarkup)
}
