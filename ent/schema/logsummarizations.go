package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// LogSummarizations holds the schema definition for the LogSummarizations entity.
type LogSummarizations struct {
	ent.Schema
}

// Fields of the LogSummarizations.
func (LogSummarizations) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New).Unique().Immutable(),
		field.Text("content_url").Default(""),
		field.Text("content_title").Default(""),
		field.Text("content_author").Default(""),
		field.Text("content_text").Default(""),
		field.Text("content_summarized_outputs").Default(""),
		field.Int("from_platform").Default(0),
		field.Int("prompt_token_usage").Default(0),
		field.Int("completion_token_usage").Default(0),
		field.Int("total_token_usage").Default(0),
		field.Int64("created_at").DefaultFunc(func() int64 { return time.Now().UnixMilli() }),
		field.Int64("updated_at").DefaultFunc(func() int64 { return time.Now().UnixMilli() }),
	}
}

// Edges of the LogSummarizations.
func (LogSummarizations) Edges() []ent.Edge {
	return nil
}
