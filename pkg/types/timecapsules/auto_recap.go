package timecapsules

import "github.com/nekomeowww/insights-bot/pkg/types/telegram"

type AutoRecapCapsule struct {
	ChatID    int64             `json:"chat_id"`
	ChatType  telegram.ChatType `json:"chat_type"`
	ChatTitle string            `json:"chat_title"`
}
