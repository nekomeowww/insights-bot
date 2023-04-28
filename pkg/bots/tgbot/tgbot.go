package tgbot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/nekomeowww/insights-bot/pkg/logger"
	"github.com/nekomeowww/insights-bot/pkg/utils"
)

type Bot struct {
	*tgbotapi.BotAPI
	logger *logger.Logger
}

func (b *Bot) MustSend(chattable tgbotapi.Chattable) *tgbotapi.Message {
	message, err := b.Send(chattable)
	if err != nil {
		b.logger.Errorf("failed to send %v to telegram: %v", utils.SprintJSON(chattable), err)
		return nil
	}

	return &message
}

func (b *Bot) MustRequest(chattable tgbotapi.Chattable) *tgbotapi.APIResponse {
	resp, err := b.Request(chattable)
	if err != nil {
		b.logger.Errorf("failed to request %v to telegram: %v", utils.SprintJSON(chattable), err)
		return nil
	}

	return resp
}
