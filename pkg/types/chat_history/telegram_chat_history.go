package chat_history

type TelegramChatHistory struct {
	ID        string `clover:"_id"`
	ChatID    int64  `clover:"chat_id"`
	MessageID int    `clover:"message_id"`
	UserID    int64  `clover:"user_id"`
	Username  string `clover:"username"`
	FullName  string `clover:"full_name"`
	Text      string `clover:"text"`
	ChattedAt int64  `clover:"chatted_at"`
	Embedded  bool   `clover:"embedded"`
	CreatedAt int64  `clover:"created_at"`
	UpdatedAt int64  `clover:"updated_at"`
}

func (TelegramChatHistory) CollectionName() string {
	return "telegram_chat_histories"
}
