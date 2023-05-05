package tgbot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/samber/lo"
)

type MessageResponse struct {
	messageConfig tgbotapi.MessageConfig
}

func NewMessage(chatID int64, message string) MessageResponse {
	return MessageResponse{
		messageConfig: tgbotapi.NewMessage(chatID, message),
	}
}

func NewMessageReplyTo(chatID int64, message string, replyToMessageID int) MessageResponse {
	m := MessageResponse{
		messageConfig: tgbotapi.NewMessage(chatID, message),
	}

	m.messageConfig.ReplyToMessageID = replyToMessageID

	return m
}

func (r MessageResponse) WithMessageConfig(config tgbotapi.MessageConfig) MessageResponse {
	r.messageConfig = config
	return r
}

func (r MessageResponse) WithParseModeHTML() MessageResponse {
	r.messageConfig.ParseMode = "HTML"
	return r
}

func (r MessageResponse) WithReplyMarkup(replyMarkup any) MessageResponse {
	r.messageConfig.ReplyMarkup = replyMarkup
	return r
}

type EditMessageResponse struct {
	textConfig         *tgbotapi.EditMessageTextConfig
	mediaConfig        *tgbotapi.EditMessageMediaConfig
	replyMarkupConfig  *tgbotapi.EditMessageReplyMarkupConfig
	captionConfig      *tgbotapi.EditMessageCaptionConfig
	liveLocationConfig *tgbotapi.EditMessageLiveLocationConfig
}

func NewEditMessageText(chatID int64, messageID int, text string) EditMessageResponse {
	return EditMessageResponse{
		textConfig: lo.ToPtr(tgbotapi.NewEditMessageText(chatID, messageID, text)),
	}
}

func NewEditMessageTextAndReplyMarkup(chatID int64, messageID int, text string, replyMarkup tgbotapi.InlineKeyboardMarkup) EditMessageResponse {
	return EditMessageResponse{
		textConfig: lo.ToPtr(tgbotapi.NewEditMessageTextAndMarkup(chatID, messageID, text, replyMarkup)),
	}
}

func (r EditMessageResponse) WithParseModeHTML() EditMessageResponse {
	r.textConfig.ParseMode = "HTML"
	return r
}

func (r EditMessageResponse) WithEditMessageTextConfig(config tgbotapi.EditMessageTextConfig) EditMessageResponse {
	r.textConfig = &config
	return r
}

func (r EditMessageResponse) WithEditMessageMediaConfig(config tgbotapi.EditMessageMediaConfig) EditMessageResponse {
	r.mediaConfig = &config
	return r
}

func (r EditMessageResponse) WithEditMessageReplyMarkupConfig(config tgbotapi.EditMessageReplyMarkupConfig) EditMessageResponse {
	r.replyMarkupConfig = &config
	return r
}

func (r EditMessageResponse) WithEditMessageCaptionConfig(config tgbotapi.EditMessageCaptionConfig) EditMessageResponse {
	r.captionConfig = &config
	return r
}

func (r EditMessageResponse) WithEditMessageLiveLocationConfig(config tgbotapi.EditMessageLiveLocationConfig) EditMessageResponse {
	r.liveLocationConfig = &config
	return r
}
