// Code generated by ent, DO NOT EDIT.

package telegramchatrecapsoptions

import (
	"entgo.io/ent/dialect/sql"
	"github.com/google/uuid"
)

const (
	// Label holds the string label denoting the telegramchatrecapsoptions type in the database.
	Label = "telegram_chat_recaps_options"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldChatID holds the string denoting the chat_id field in the database.
	FieldChatID = "chat_id"
	// FieldAutoRecapSendMode holds the string denoting the auto_recap_send_mode field in the database.
	FieldAutoRecapSendMode = "auto_recap_send_mode"
	// FieldManualRecapRatePerSeconds holds the string denoting the manual_recap_rate_per_seconds field in the database.
	FieldManualRecapRatePerSeconds = "manual_recap_rate_per_seconds"
	// FieldAutoRecapRatesPerDay holds the string denoting the auto_recap_rates_per_day field in the database.
	FieldAutoRecapRatesPerDay = "auto_recap_rates_per_day"
	// FieldPinAutoRecapMessage holds the string denoting the pin_auto_recap_message field in the database.
	FieldPinAutoRecapMessage = "pin_auto_recap_message"
	// FieldCreatedAt holds the string denoting the created_at field in the database.
	FieldCreatedAt = "created_at"
	// FieldUpdatedAt holds the string denoting the updated_at field in the database.
	FieldUpdatedAt = "updated_at"
	// Table holds the table name of the telegramchatrecapsoptions in the database.
	Table = "telegram_chat_recaps_options"
)

// Columns holds all SQL columns for telegramchatrecapsoptions fields.
var Columns = []string{
	FieldID,
	FieldChatID,
	FieldAutoRecapSendMode,
	FieldManualRecapRatePerSeconds,
	FieldAutoRecapRatesPerDay,
	FieldPinAutoRecapMessage,
	FieldCreatedAt,
	FieldUpdatedAt,
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	return false
}

var (
	// DefaultAutoRecapSendMode holds the default value on creation for the "auto_recap_send_mode" field.
	DefaultAutoRecapSendMode int
	// DefaultManualRecapRatePerSeconds holds the default value on creation for the "manual_recap_rate_per_seconds" field.
	DefaultManualRecapRatePerSeconds int64
	// DefaultAutoRecapRatesPerDay holds the default value on creation for the "auto_recap_rates_per_day" field.
	DefaultAutoRecapRatesPerDay int
	// DefaultPinAutoRecapMessage holds the default value on creation for the "pin_auto_recap_message" field.
	DefaultPinAutoRecapMessage bool
	// DefaultCreatedAt holds the default value on creation for the "created_at" field.
	DefaultCreatedAt func() int64
	// DefaultUpdatedAt holds the default value on creation for the "updated_at" field.
	DefaultUpdatedAt func() int64
	// DefaultID holds the default value on creation for the "id" field.
	DefaultID func() uuid.UUID
)

// OrderOption defines the ordering options for the TelegramChatRecapsOptions queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByChatID orders the results by the chat_id field.
func ByChatID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldChatID, opts...).ToFunc()
}

// ByAutoRecapSendMode orders the results by the auto_recap_send_mode field.
func ByAutoRecapSendMode(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldAutoRecapSendMode, opts...).ToFunc()
}

// ByManualRecapRatePerSeconds orders the results by the manual_recap_rate_per_seconds field.
func ByManualRecapRatePerSeconds(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldManualRecapRatePerSeconds, opts...).ToFunc()
}

// ByAutoRecapRatesPerDay orders the results by the auto_recap_rates_per_day field.
func ByAutoRecapRatesPerDay(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldAutoRecapRatesPerDay, opts...).ToFunc()
}

// ByPinAutoRecapMessage orders the results by the pin_auto_recap_message field.
func ByPinAutoRecapMessage(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldPinAutoRecapMessage, opts...).ToFunc()
}

// ByCreatedAt orders the results by the created_at field.
func ByCreatedAt(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldCreatedAt, opts...).ToFunc()
}

// ByUpdatedAt orders the results by the updated_at field.
func ByUpdatedAt(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldUpdatedAt, opts...).ToFunc()
}
