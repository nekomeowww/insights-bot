package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// TelegramChatAutoRecapsSubscribers holds the schema definition for the TelegramChatAutoRecapsSubscribers entity.
type TelegramChatAutoRecapsSubscribers struct {
	ent.Schema
}

// Fields of the TelegramChatAutoRecapsSubscribers.
func (TelegramChatAutoRecapsSubscribers) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New).Unique().Immutable(),
		field.Int64("chat_id").Default(0),
		field.Int64("user_id").Default(0),
		field.Int64("created_at").DefaultFunc(func() int64 { return time.Now().UnixMilli() }),
		field.Int64("updated_at").DefaultFunc(func() int64 { return time.Now().UnixMilli() }),
	}
}

// Edges of the TelegramChatAutoRecapsSubscribers.
func (TelegramChatAutoRecapsSubscribers) Edges() []ent.Edge {
	return nil
}
