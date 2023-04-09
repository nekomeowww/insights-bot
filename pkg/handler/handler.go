package handler

import (
	"errors"
	"net/url"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/gorilla/schema"
	"github.com/nekomeowww/insights-bot/pkg/logger"
	"github.com/nekomeowww/insights-bot/pkg/utils"
)

type Context struct {
	Bot    *Bot
	Update tgbotapi.Update
}

type Bot struct {
	*tgbotapi.BotAPI
	logger *logger.Logger
}

func (b *Bot) MustSend(chattable tgbotapi.Chattable) *tgbotapi.Message {
	message, err := b.Send(chattable)
	if err != nil {
		b.logger.Error("failed to send %v to telegram: %v", utils.SprintJSON(chattable), err)
		return nil
	}

	return &message
}

func NewContext(bot *tgbotapi.BotAPI, update tgbotapi.Update, logger *logger.Logger) *Context {
	return &Context{
		Bot:    &Bot{BotAPI: bot, logger: logger},
		Update: update,
	}
}

type HandleFunc func(ctx *Context)

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
