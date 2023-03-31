package chat_history

type TelegramChatHistory struct {
	ID                 string `clover:"_id"`
	ChatID             int64  `clover:"chat_id"`
	MessageID          int    `clover:"message_id"`
	UserID             int64  `clover:"user_id"`
	Username           string `clover:"username"`
	FullName           string `clover:"full_name"`
	Text               string `clover:"text"`
	RepliedToMessageID int    `clover:"replied_to_message_id"`
	RepliedToUserID    int64  `clover:"replied_to_user_id"`
	RepliedToFullName  string `clover:"replied_to_full_name"`
	RepliedToUsername  string `clover:"replied_to_username"`
	RepliedToText      string `clover:"replied_to_text"`
	ChattedAt          int64  `clover:"chatted_at"`
	Embedded           bool   `clover:"embedded"`
	CreatedAt          int64  `clover:"created_at"`
	UpdatedAt          int64  `clover:"updated_at"`
}

func (TelegramChatHistory) CollectionName() string {
	return "telegram_chat_histories"
}
