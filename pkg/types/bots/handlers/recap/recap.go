package recap

import "github.com/nekomeowww/insights-bot/pkg/types/tgchat"

type SelectHourCallbackQueryData struct {
	Hour      int64                    `json:"hour"`
	ChatID    int64                    `json:"chat_id"`
	ChatTitle string                   `json:"chat_title"`
	RecapMode tgchat.AutoRecapSendMode `json:"recap_mode"`
}
