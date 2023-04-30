package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// TelegramChatHistories holds the schema definition for the TelegramChatHistories entity.
type TelegramChatHistories struct {
	ent.Schema
}

// Fields of the TelegramChatHistories.
func (TelegramChatHistories) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New).Unique().Immutable(),
		field.Int64("chat_id").Default(0),
		field.Int64("message_id").Default(0),
		field.Int64("user_id").Default(0),
		field.Text("username").Default(""),
		field.Text("full_name").Default(""),
		field.Text("text").Default(""),
		field.Int64("replied_to_message_id").Default(0),
		field.Int64("replied_to_user_id").Default(0),
		field.Text("replied_to_full_name").Default(""),
		field.Text("replied_to_username").Default(""),
		field.Text("replied_to_text").Default(""),
		field.Int64("chatted_at").Default(time.Now().Unix()),
		field.Bool("embedded").Default(false),
		field.Int64("created_at").DefaultFunc(func() int64 { return time.Now().UnixMilli() }),
		field.Int64("updated_at").DefaultFunc(func() int64 { return time.Now().UnixMilli() }),
	}
}

// Edges of the TelegramChatHistories.
func (TelegramChatHistories) Edges() []ent.Edge {
	return nil
}
