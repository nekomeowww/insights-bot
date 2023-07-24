package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// FeedbackSummarizationsReactions holds the schema definition for the FeedbackSummarizationsReactions entity.
type FeedbackSummarizationsReactions struct {
	ent.Schema
}

// Fields of the FeedbackSummarizationsReactions.
func (FeedbackSummarizationsReactions) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New).Unique().Immutable(),
		field.Int64("chat_id").Default(0),
		field.UUID("log_id", uuid.UUID{}).Default(uuid.New).Immutable(),
		field.Int64("user_id").Default(0),
		field.Enum("type").Values("none", "up_vote", "down_vote", "lmao").Default("none"),
		field.Int64("created_at").DefaultFunc(func() int64 { return time.Now().UnixMilli() }),
		field.Int64("updated_at").DefaultFunc(func() int64 { return time.Now().UnixMilli() }),
	}
}

// Edges of the FeedbackSummarizationsReactions.
func (FeedbackSummarizationsReactions) Edges() []ent.Edge {
	return nil
}
