package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// FeedbackChatHistoriesRecapsReactions holds the schema definition for the FeedbackChatHistoriesRecapsReactions entity.
type FeedbackChatHistoriesRecapsReactions struct {
	ent.Schema
}

// Fields of the FeedbackChatHistoriesRecapsReactions.
func (FeedbackChatHistoriesRecapsReactions) Fields() []ent.Field {
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

// Edges of the FeedbackChatHistoriesRecapsReactions.
func (FeedbackChatHistoriesRecapsReactions) Edges() []ent.Edge {
	return nil
}
