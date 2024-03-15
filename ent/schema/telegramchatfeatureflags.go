package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// TelegramChatFeatureFlags holds the schema definition for the TelegramChatFeatureFlags entity.
type TelegramChatFeatureFlags struct {
	ent.Schema
}

// Fields of the TelegramChatFeatureFlags.
func (TelegramChatFeatureFlags) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New).Unique().Immutable(),
		field.Int64("chat_id").Unique(),
		field.Text("chat_type"),
		field.Text("chat_title").Default(""),
		field.Bool("feature_chat_histories_recap").Default(false),
		field.Text("feature_language").Default(""),
		field.Int64("created_at").DefaultFunc(func() int64 { return time.Now().UnixMilli() }),
		field.Int64("updated_at").DefaultFunc(func() int64 { return time.Now().UnixMilli() }),
	}
}

// Edges of the TelegramChatFeatureFlags.
func (TelegramChatFeatureFlags) Edges() []ent.Edge {
	return nil
}
