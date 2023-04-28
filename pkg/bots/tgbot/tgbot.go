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

func (b *Bot) EditMessageText(chatID int64, messageID int, text string) error {
	_, err := b.Request(tgbotapi.EditMessageTextConfig{
		BaseEdit: tgbotapi.BaseEdit{ChatID: chatID, MessageID: messageID},
		Text:     text,
	})
	if err != nil {
		return err
	}

	return nil
}

func (b *Bot) MustEditMessageText(chatID int64, messageID int, text string) {
	err := b.EditMessageText(chatID, messageID, text)
	if err != nil {
		b.logger.Errorf("failed to edit message text: %v", err)
	}
}
