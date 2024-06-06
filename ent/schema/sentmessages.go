package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
	"time"
)

// SentMessages holds the schema definition for the SentMessages entity.
type SentMessages struct {
	ent.Schema
}

// Fields of the SentMessages.
func (SentMessages) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New).Unique().Immutable(),
		field.Int64("chat_id").Default(0),
		field.Int("message_id").Default(0),
		field.Text("text").Default(""),
		field.Bool("is_pinned").Default(false),
		field.Int("from_platform").Default(0),
		field.Int("message_type").Default(0),
		field.Int64("created_at").DefaultFunc(func() int64 { return time.Now().UnixMilli() }),
		field.Int64("updated_at").DefaultFunc(func() int64 { return time.Now().UnixMilli() }),
	}
}

// Edges of the SentMessages.
func (SentMessages) Edges() []ent.Edge {
	return nil
}
