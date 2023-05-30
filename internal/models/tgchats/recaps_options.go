package tgchats

import (
	"context"

	"github.com/nekomeowww/insights-bot/ent"
	"github.com/nekomeowww/insights-bot/ent/telegramchatrecapsoptions"
	"github.com/nekomeowww/insights-bot/pkg/types/tgchat"
)

func (m *Model) FindOneRecapsOption(chatID int64) (*ent.TelegramChatRecapsOptions, error) {
	option, err := m.ent.TelegramChatRecapsOptions.
		Query().
		Where(telegramchatrecapsoptions.ChatID(chatID)).
		First(context.Background())
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		}

		return nil, err
	}

	return option, nil
}

func (m *Model) SetRecapsRecapMode(chatID int64, recapMode tgchat.AutoRecapSendMode) error {
	option, err := m.FindOneRecapsOption(chatID)
	if err != nil {
		return err
	}
	if option == nil {
		err := m.ent.TelegramChatRecapsOptions.
			Create().
			SetChatID(chatID).
			SetAutoRecapSendMode(int(recapMode)).
			Exec(context.Background())
		if err != nil {
			return err
		}

		return nil
	}
	if option.AutoRecapSendMode == int(recapMode) {
		return nil
	}

	return m.ent.TelegramChatRecapsOptions.
		UpdateOne(option).
		SetAutoRecapSendMode(int(recapMode)).
		Exec(context.Background())
}
