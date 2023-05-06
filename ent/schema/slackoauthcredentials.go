package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// SlackOAuthCredentials holds the schema definition for the SlackOAuthCredentials entity.
type SlackOAuthCredentials struct {
	ent.Schema
}

// Fields of the SavedSlackToken.
func (SlackOAuthCredentials) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New).Unique().Immutable(),
		field.Text("team_id").NotEmpty().Unique().Immutable(),
		field.Text("refresh_token").NotEmpty(),
		field.Text("access_token").NotEmpty(),
		field.Int64("created_at").DefaultFunc(func() int64 { return time.Now().UnixMilli() }),
		field.Int64("updated_at").DefaultFunc(func() int64 { return time.Now().UnixMilli() }),
	}
}

// Edges of the SavedSlackToken.
func (SlackOAuthCredentials) Edges() []ent.Edge {
	return nil
}
