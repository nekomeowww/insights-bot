package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// SavedSlackToken holds the schema definition for the SavedSlackToken entity.
type SavedSlackToken struct {
	ent.Schema
}

// Fields of the SavedSlackToken.
func (SavedSlackToken) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New).Unique().Immutable(),
		field.Text("team_id").NotEmpty().Unique().Immutable(),
		field.Text("access_token").NotEmpty(),
		field.Int64("created_at").DefaultFunc(func() int64 { return time.Now().UnixMilli() }),
		field.Int64("updated_at").DefaultFunc(func() int64 { return time.Now().UnixMilli() }),
	}
}

// Edges of the SavedSlackToken.
func (SavedSlackToken) Edges() []ent.Edge {
	return nil
}
