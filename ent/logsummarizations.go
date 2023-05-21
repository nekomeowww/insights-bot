// Code generated by ent, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/google/uuid"
	"github.com/nekomeowww/insights-bot/ent/logsummarizations"
)

// LogSummarizations is the model entity for the LogSummarizations schema.
type LogSummarizations struct {
	config `json:"-"`
	// ID of the ent.
	ID uuid.UUID `json:"id,omitempty"`
	// ContentURL holds the value of the "content_url" field.
	ContentURL string `json:"content_url,omitempty"`
	// ContentTitle holds the value of the "content_title" field.
	ContentTitle string `json:"content_title,omitempty"`
	// ContentAuthor holds the value of the "content_author" field.
	ContentAuthor string `json:"content_author,omitempty"`
	// ContentText holds the value of the "content_text" field.
	ContentText string `json:"content_text,omitempty"`
	// ContentSummarizedOutputs holds the value of the "content_summarized_outputs" field.
	ContentSummarizedOutputs string `json:"content_summarized_outputs,omitempty"`
	// FromPlatform holds the value of the "from_platform" field.
	FromPlatform int `json:"from_platform,omitempty"`
	// CreatedAt holds the value of the "created_at" field.
	CreatedAt int64 `json:"created_at,omitempty"`
	// UpdatedAt holds the value of the "updated_at" field.
	UpdatedAt    int64 `json:"updated_at,omitempty"`
	selectValues sql.SelectValues
}

// scanValues returns the types for scanning values from sql.Rows.
func (*LogSummarizations) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case logsummarizations.FieldFromPlatform, logsummarizations.FieldCreatedAt, logsummarizations.FieldUpdatedAt:
			values[i] = new(sql.NullInt64)
		case logsummarizations.FieldContentURL, logsummarizations.FieldContentTitle, logsummarizations.FieldContentAuthor, logsummarizations.FieldContentText, logsummarizations.FieldContentSummarizedOutputs:
			values[i] = new(sql.NullString)
		case logsummarizations.FieldID:
			values[i] = new(uuid.UUID)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the LogSummarizations fields.
func (ls *LogSummarizations) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case logsummarizations.FieldID:
			if value, ok := values[i].(*uuid.UUID); !ok {
				return fmt.Errorf("unexpected type %T for field id", values[i])
			} else if value != nil {
				ls.ID = *value
			}
		case logsummarizations.FieldContentURL:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field content_url", values[i])
			} else if value.Valid {
				ls.ContentURL = value.String
			}
		case logsummarizations.FieldContentTitle:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field content_title", values[i])
			} else if value.Valid {
				ls.ContentTitle = value.String
			}
		case logsummarizations.FieldContentAuthor:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field content_author", values[i])
			} else if value.Valid {
				ls.ContentAuthor = value.String
			}
		case logsummarizations.FieldContentText:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field content_text", values[i])
			} else if value.Valid {
				ls.ContentText = value.String
			}
		case logsummarizations.FieldContentSummarizedOutputs:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field content_summarized_outputs", values[i])
			} else if value.Valid {
				ls.ContentSummarizedOutputs = value.String
			}
		case logsummarizations.FieldFromPlatform:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field from_platform", values[i])
			} else if value.Valid {
				ls.FromPlatform = int(value.Int64)
			}
		case logsummarizations.FieldCreatedAt:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field created_at", values[i])
			} else if value.Valid {
				ls.CreatedAt = value.Int64
			}
		case logsummarizations.FieldUpdatedAt:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field updated_at", values[i])
			} else if value.Valid {
				ls.UpdatedAt = value.Int64
			}
		default:
			ls.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the LogSummarizations.
// This includes values selected through modifiers, order, etc.
func (ls *LogSummarizations) Value(name string) (ent.Value, error) {
	return ls.selectValues.Get(name)
}

// Update returns a builder for updating this LogSummarizations.
// Note that you need to call LogSummarizations.Unwrap() before calling this method if this LogSummarizations
// was returned from a transaction, and the transaction was committed or rolled back.
func (ls *LogSummarizations) Update() *LogSummarizationsUpdateOne {
	return NewLogSummarizationsClient(ls.config).UpdateOne(ls)
}

// Unwrap unwraps the LogSummarizations entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (ls *LogSummarizations) Unwrap() *LogSummarizations {
	_tx, ok := ls.config.driver.(*txDriver)
	if !ok {
		panic("ent: LogSummarizations is not a transactional entity")
	}
	ls.config.driver = _tx.drv
	return ls
}

// String implements the fmt.Stringer.
func (ls *LogSummarizations) String() string {
	var builder strings.Builder
	builder.WriteString("LogSummarizations(")
	builder.WriteString(fmt.Sprintf("id=%v, ", ls.ID))
	builder.WriteString("content_url=")
	builder.WriteString(ls.ContentURL)
	builder.WriteString(", ")
	builder.WriteString("content_title=")
	builder.WriteString(ls.ContentTitle)
	builder.WriteString(", ")
	builder.WriteString("content_author=")
	builder.WriteString(ls.ContentAuthor)
	builder.WriteString(", ")
	builder.WriteString("content_text=")
	builder.WriteString(ls.ContentText)
	builder.WriteString(", ")
	builder.WriteString("content_summarized_outputs=")
	builder.WriteString(ls.ContentSummarizedOutputs)
	builder.WriteString(", ")
	builder.WriteString("from_platform=")
	builder.WriteString(fmt.Sprintf("%v", ls.FromPlatform))
	builder.WriteString(", ")
	builder.WriteString("created_at=")
	builder.WriteString(fmt.Sprintf("%v", ls.CreatedAt))
	builder.WriteString(", ")
	builder.WriteString("updated_at=")
	builder.WriteString(fmt.Sprintf("%v", ls.UpdatedAt))
	builder.WriteByte(')')
	return builder.String()
}

// LogSummarizationsSlice is a parsable slice of LogSummarizations.
type LogSummarizationsSlice []*LogSummarizations
