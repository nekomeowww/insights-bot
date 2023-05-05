package tgbot

import (
	"errors"
	"net/url"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/gorilla/schema"
	"github.com/nekomeowww/insights-bot/pkg/logger"
	"github.com/samber/lo"
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
	UpdateTypeChatJoinRequest    UpdateType = "chat_join_request"
)

type Context struct {
	Bot    *Bot
	Update tgbotapi.Update
	Logger *logger.Logger
}

func NewContext(bot *tgbotapi.BotAPI, update tgbotapi.Update, logger *logger.Logger) *Context {
	return &Context{
		Bot:    &Bot{BotAPI: bot, logger: logger},
		Update: update,
		Logger: logger,
	}
}

func (c Context) UpdateType() UpdateType {
	switch {
	case c.Update.Message != nil:
		return UpdateTypeMessage
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

func (c *Context) CallbackQueryDataBindQuery(dst interface{}) error {
	if c.Update.CallbackQuery == nil {
		return errors.New("callback query is nil")
	}
	if c.Update.CallbackQuery.Data == "" {
		return nil
	}

	parsedURL, err := url.Parse(c.Update.CallbackQuery.Data)
	if err != nil {
		return err
	}

	decoder := schema.NewDecoder()

	err = decoder.Decode(dst, parsedURL.Query())
	if err != nil {
		return err
	}

	return nil
}

func (c *Context) NewMessage(message string) MessageResponse {
	return NewMessage(c.Update.FromChat().ID, message)
}

func (c *Context) NewMessageReplyTo(message string, replyToMessageID int) MessageResponse {
	return NewMessageReplyTo(c.Update.FromChat().ID, message, replyToMessageID)
}

func (c *Context) NewEditMessageText(messageID int, text string) EditMessageResponse {
	return EditMessageResponse{
		textConfig: lo.ToPtr(tgbotapi.NewEditMessageText(c.Update.FromChat().ID, messageID, text)),
	}
}

func (c *Context) NewEditMessageTextAndReplyMarkup(messageID int, text string, replyMarkup tgbotapi.InlineKeyboardMarkup) EditMessageResponse {
	return EditMessageResponse{
		textConfig: lo.ToPtr(tgbotapi.NewEditMessageTextAndMarkup(c.Update.FromChat().ID, messageID, text, replyMarkup)),
	}
}
