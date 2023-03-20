package handler

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Context struct {
	Bot    *tgbotapi.BotAPI
	Update tgbotapi.Update
}

func NewContext(bot *tgbotapi.BotAPI, update tgbotapi.Update) *Context {
	return &Context{
		Bot:    bot,
		Update: update,
	}
}

type HandleFunc func(ctx *Context)
