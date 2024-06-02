package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// TelegramChatRecapsOptions holds the schema definition for the TelegramChatRecapsOptions entity.
type TelegramChatRecapsOptions struct {
	ent.Schema
}

// Fields of the TelegramChatRecapsOptions.
func (TelegramChatRecapsOptions) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New).Unique().Immutable(),
		field.Int64("chat_id").Unique(),
		field.Int("auto_recap_send_mode").Default(0),
		field.Int64("manual_recap_rate_per_seconds").Default(0),
		field.Int("auto_recap_rates_per_day").Default(0),
		field.Bool("pin_auto_recap_message").Default(false),
		field.Int64("created_at").DefaultFunc(func() int64 { return time.Now().UnixMilli() }),
		field.Int64("updated_at").DefaultFunc(func() int64 { return time.Now().UnixMilli() }),
	}
}

// Edges of the TelegramChatRecapsOptions.
func (TelegramChatRecapsOptions) Edges() []ent.Edge {
	return nil
}
