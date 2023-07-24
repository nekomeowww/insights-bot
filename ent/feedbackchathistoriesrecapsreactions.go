// Code generated by ent, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/google/uuid"
	"github.com/nekomeowww/insights-bot/ent/feedbackchathistoriesrecapsreactions"
)

// FeedbackChatHistoriesRecapsReactions is the model entity for the FeedbackChatHistoriesRecapsReactions schema.
type FeedbackChatHistoriesRecapsReactions struct {
	config `json:"-"`
	// ID of the ent.
	ID uuid.UUID `json:"id,omitempty"`
	// ChatID holds the value of the "chat_id" field.
	ChatID int64 `json:"chat_id,omitempty"`
	// LogID holds the value of the "log_id" field.
	LogID uuid.UUID `json:"log_id,omitempty"`
	// UserID holds the value of the "user_id" field.
	UserID int64 `json:"user_id,omitempty"`
	// Type holds the value of the "type" field.
	Type feedbackchathistoriesrecapsreactions.Type `json:"type,omitempty"`
	// CreatedAt holds the value of the "created_at" field.
	CreatedAt int64 `json:"created_at,omitempty"`
	// UpdatedAt holds the value of the "updated_at" field.
	UpdatedAt    int64 `json:"updated_at,omitempty"`
	selectValues sql.SelectValues
}

// scanValues returns the types for scanning values from sql.Rows.
func (*FeedbackChatHistoriesRecapsReactions) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case feedbackchathistoriesrecapsreactions.FieldChatID, feedbackchathistoriesrecapsreactions.FieldUserID, feedbackchathistoriesrecapsreactions.FieldCreatedAt, feedbackchathistoriesrecapsreactions.FieldUpdatedAt:
			values[i] = new(sql.NullInt64)
		case feedbackchathistoriesrecapsreactions.FieldType:
			values[i] = new(sql.NullString)
		case feedbackchathistoriesrecapsreactions.FieldID, feedbackchathistoriesrecapsreactions.FieldLogID:
			values[i] = new(uuid.UUID)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the FeedbackChatHistoriesRecapsReactions fields.
func (fchrr *FeedbackChatHistoriesRecapsReactions) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case feedbackchathistoriesrecapsreactions.FieldID:
			if value, ok := values[i].(*uuid.UUID); !ok {
				return fmt.Errorf("unexpected type %T for field id", values[i])
			} else if value != nil {
				fchrr.ID = *value
			}
		case feedbackchathistoriesrecapsreactions.FieldChatID:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field chat_id", values[i])
			} else if value.Valid {
				fchrr.ChatID = value.Int64
			}
		case feedbackchathistoriesrecapsreactions.FieldLogID:
			if value, ok := values[i].(*uuid.UUID); !ok {
				return fmt.Errorf("unexpected type %T for field log_id", values[i])
			} else if value != nil {
				fchrr.LogID = *value
			}
		case feedbackchathistoriesrecapsreactions.FieldUserID:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field user_id", values[i])
			} else if value.Valid {
				fchrr.UserID = value.Int64
			}
		case feedbackchathistoriesrecapsreactions.FieldType:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field type", values[i])
			} else if value.Valid {
				fchrr.Type = feedbackchathistoriesrecapsreactions.Type(value.String)
			}
		case feedbackchathistoriesrecapsreactions.FieldCreatedAt:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field created_at", values[i])
			} else if value.Valid {
				fchrr.CreatedAt = value.Int64
			}
		case feedbackchathistoriesrecapsreactions.FieldUpdatedAt:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field updated_at", values[i])
			} else if value.Valid {
				fchrr.UpdatedAt = value.Int64
			}
		default:
			fchrr.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the FeedbackChatHistoriesRecapsReactions.
// This includes values selected through modifiers, order, etc.
func (fchrr *FeedbackChatHistoriesRecapsReactions) Value(name string) (ent.Value, error) {
	return fchrr.selectValues.Get(name)
}

// Update returns a builder for updating this FeedbackChatHistoriesRecapsReactions.
// Note that you need to call FeedbackChatHistoriesRecapsReactions.Unwrap() before calling this method if this FeedbackChatHistoriesRecapsReactions
// was returned from a transaction, and the transaction was committed or rolled back.
func (fchrr *FeedbackChatHistoriesRecapsReactions) Update() *FeedbackChatHistoriesRecapsReactionsUpdateOne {
	return NewFeedbackChatHistoriesRecapsReactionsClient(fchrr.config).UpdateOne(fchrr)
}

// Unwrap unwraps the FeedbackChatHistoriesRecapsReactions entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (fchrr *FeedbackChatHistoriesRecapsReactions) Unwrap() *FeedbackChatHistoriesRecapsReactions {
	_tx, ok := fchrr.config.driver.(*txDriver)
	if !ok {
		panic("ent: FeedbackChatHistoriesRecapsReactions is not a transactional entity")
	}
	fchrr.config.driver = _tx.drv
	return fchrr
}

// String implements the fmt.Stringer.
func (fchrr *FeedbackChatHistoriesRecapsReactions) String() string {
	var builder strings.Builder
	builder.WriteString("FeedbackChatHistoriesRecapsReactions(")
	builder.WriteString(fmt.Sprintf("id=%v, ", fchrr.ID))
	builder.WriteString("chat_id=")
	builder.WriteString(fmt.Sprintf("%v", fchrr.ChatID))
	builder.WriteString(", ")
	builder.WriteString("log_id=")
	builder.WriteString(fmt.Sprintf("%v", fchrr.LogID))
	builder.WriteString(", ")
	builder.WriteString("user_id=")
	builder.WriteString(fmt.Sprintf("%v", fchrr.UserID))
	builder.WriteString(", ")
	builder.WriteString("type=")
	builder.WriteString(fmt.Sprintf("%v", fchrr.Type))
	builder.WriteString(", ")
	builder.WriteString("created_at=")
	builder.WriteString(fmt.Sprintf("%v", fchrr.CreatedAt))
	builder.WriteString(", ")
	builder.WriteString("updated_at=")
	builder.WriteString(fmt.Sprintf("%v", fchrr.UpdatedAt))
	builder.WriteByte(')')
	return builder.String()
}

// FeedbackChatHistoriesRecapsReactionsSlice is a parsable slice of FeedbackChatHistoriesRecapsReactions.
type FeedbackChatHistoriesRecapsReactionsSlice []*FeedbackChatHistoriesRecapsReactions