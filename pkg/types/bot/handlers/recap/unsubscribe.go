package recap

type UnsubscribeRecapActionData struct {
	ChatID    int64  `json:"chatId"`
	ChatTitle string `json:"chatTitle"`
	FromID    int64  `json:"fromId"`
}
