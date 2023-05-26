package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// LogChatHistoriesRecap holds the schema definition for the LogChatHistoriesRecap entity.
type LogChatHistoriesRecap struct {
	ent.Schema
}

// Fields of the LogChatHistoriesRecap.
func (LogChatHistoriesRecap) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New).Unique().Immutable(),
		field.Int64("chat_id").Default(0),
		field.Text("recap_inputs").Default(""),
		field.Text("recap_outputs").Default(""),
		field.Int("from_platform").Default(0),
		field.Int("prompt_token_usage").Default(0),
		field.Int("completion_token_usage").Default(0),
		field.Int("total_token_usage").Default(0),
		field.Int("recap_type").Default(0),
		field.Int64("created_at").DefaultFunc(func() int64 { return time.Now().UnixMilli() }),
		field.Int64("updated_at").DefaultFunc(func() int64 { return time.Now().UnixMilli() }),
	}
}

// Edges of the LogChatHistoriesRecap.
func (LogChatHistoriesRecap) Edges() []ent.Edge {
	return nil
}
