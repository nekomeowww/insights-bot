package recap

import "github.com/nekomeowww/insights-bot/pkg/types/tgchat"

type ConfigureRecapToggleActionData struct {
	Status bool  `json:"status"`
	ChatID int64 `json:"chatId"`
	FromID int64 `json:"fromId"`
}

type ConfigureRecapAssignModeActionData struct {
	Mode   tgchat.AutoRecapSendMode `json:"mode"`
	ChatID int64                    `json:"chatId"`
	FromID int64                    `json:"fromId"`
}

type ConfigureRecapCompleteActionData struct {
	ChatID int64 `json:"chatId"`
	FromID int64 `json:"fromId"`
}
