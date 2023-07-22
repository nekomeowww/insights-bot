package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// MetricOpenAIChatCompletionTokenUsage holds the schema definition for the MetricOpenAIChatCompletionTokenUsage entity.
type MetricOpenAIChatCompletionTokenUsage struct {
	ent.Schema
}

// Fields of the MetricOpenAIChatCompletionTokenUsage.
func (MetricOpenAIChatCompletionTokenUsage) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New).Unique().Immutable(),
		field.Text("prompt_operation").Default(""),
		field.Int("prompt_character_length").Default(0),
		field.Int("prompt_token_usage").Default(0),
		field.Int("completion_character_length").Default(0),
		field.Int("completion_token_usage").Default(0),
		field.Int("total_token_usage").Default(0),
		field.String("model_name").Default(""),
		field.Int64("created_at").DefaultFunc(func() int64 { return time.Now().UnixMilli() }),
	}
}

// Edges of the MetricOpenAIChatCompletionTokenUsage.
func (MetricOpenAIChatCompletionTokenUsage) Edges() []ent.Edge {
	return nil
}
