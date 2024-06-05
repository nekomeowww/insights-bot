package tgbot

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type PinChatMessageConfig struct {
	ChatID              int64
	ChannelUsername     string
	MessageID           int
	DisableNotification bool
}

func (config PinChatMessageConfig) method() string {
	return "pinChatMessage"
}

func (config PinChatMessageConfig) params() (tgbotapi.Params, error) {
	params := make(tgbotapi.Params)

	err := params.AddFirstValid("chat_id", config.ChatID, config.ChannelUsername)
	if err != nil {
		return nil, err
	}

	params.AddNonZero("message_id", config.MessageID)
	params.AddBool("disable_notification", true)

	return params, err
}

func NewPinChatMessageConfig(chatID int64, messageID int) PinChatMessageConfig {
	return PinChatMessageConfig{
		ChatID:              chatID,
		MessageID:           messageID,
		DisableNotification: true,
	}
}

type UnpinChatMessageConfig struct {
	ChatID          int64
	ChannelUsername string
	MessageID       int
}

func (config UnpinChatMessageConfig) method() string {
	return "unpinChatMessage"
}

func (config UnpinChatMessageConfig) params() (tgbotapi.Params, error) {
	params := make(tgbotapi.Params)

	err := params.AddFirstValid("chat_id", config.ChatID, config.ChannelUsername)
	if err != nil {
		return nil, err
	}

	params.AddNonZero("message_id", config.MessageID)

	return params, err
}

func NewUnpinChatMessageConfig(chatID int64, messageID int) UnpinChatMessageConfig {
	return UnpinChatMessageConfig{
		ChatID:    chatID,
		MessageID: messageID,
	}
}
