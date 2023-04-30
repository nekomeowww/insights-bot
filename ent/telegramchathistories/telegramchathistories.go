// Code generated by ent, DO NOT EDIT.

package telegramchathistories

import (
	"entgo.io/ent/dialect/sql"
	"github.com/google/uuid"
)

const (
	// Label holds the string label denoting the telegramchathistories type in the database.
	Label = "telegram_chat_histories"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldChatID holds the string denoting the chat_id field in the database.
	FieldChatID = "chat_id"
	// FieldMessageID holds the string denoting the message_id field in the database.
	FieldMessageID = "message_id"
	// FieldUserID holds the string denoting the user_id field in the database.
	FieldUserID = "user_id"
	// FieldUsername holds the string denoting the username field in the database.
	FieldUsername = "username"
	// FieldFullName holds the string denoting the full_name field in the database.
	FieldFullName = "full_name"
	// FieldText holds the string denoting the text field in the database.
	FieldText = "text"
	// FieldRepliedToMessageID holds the string denoting the replied_to_message_id field in the database.
	FieldRepliedToMessageID = "replied_to_message_id"
	// FieldRepliedToUserID holds the string denoting the replied_to_user_id field in the database.
	FieldRepliedToUserID = "replied_to_user_id"
	// FieldRepliedToFullName holds the string denoting the replied_to_full_name field in the database.
	FieldRepliedToFullName = "replied_to_full_name"
	// FieldRepliedToUsername holds the string denoting the replied_to_username field in the database.
	FieldRepliedToUsername = "replied_to_username"
	// FieldRepliedToText holds the string denoting the replied_to_text field in the database.
	FieldRepliedToText = "replied_to_text"
	// FieldChattedAt holds the string denoting the chatted_at field in the database.
	FieldChattedAt = "chatted_at"
	// FieldEmbedded holds the string denoting the embedded field in the database.
	FieldEmbedded = "embedded"
	// FieldCreatedAt holds the string denoting the created_at field in the database.
	FieldCreatedAt = "created_at"
	// FieldUpdatedAt holds the string denoting the updated_at field in the database.
	FieldUpdatedAt = "updated_at"
	// Table holds the table name of the telegramchathistories in the database.
	Table = "telegram_chat_histories"
)

// Columns holds all SQL columns for telegramchathistories fields.
var Columns = []string{
	FieldID,
	FieldChatID,
	FieldMessageID,
	FieldUserID,
	FieldUsername,
	FieldFullName,
	FieldText,
	FieldRepliedToMessageID,
	FieldRepliedToUserID,
	FieldRepliedToFullName,
	FieldRepliedToUsername,
	FieldRepliedToText,
	FieldChattedAt,
	FieldEmbedded,
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
	// DefaultChatID holds the default value on creation for the "chat_id" field.
	DefaultChatID int64
	// DefaultMessageID holds the default value on creation for the "message_id" field.
	DefaultMessageID int64
	// DefaultUserID holds the default value on creation for the "user_id" field.
	DefaultUserID int64
	// DefaultUsername holds the default value on creation for the "username" field.
	DefaultUsername string
	// DefaultFullName holds the default value on creation for the "full_name" field.
	DefaultFullName string
	// DefaultText holds the default value on creation for the "text" field.
	DefaultText string
	// DefaultRepliedToMessageID holds the default value on creation for the "replied_to_message_id" field.
	DefaultRepliedToMessageID int64
	// DefaultRepliedToUserID holds the default value on creation for the "replied_to_user_id" field.
	DefaultRepliedToUserID int64
	// DefaultRepliedToFullName holds the default value on creation for the "replied_to_full_name" field.
	DefaultRepliedToFullName string
	// DefaultRepliedToUsername holds the default value on creation for the "replied_to_username" field.
	DefaultRepliedToUsername string
	// DefaultRepliedToText holds the default value on creation for the "replied_to_text" field.
	DefaultRepliedToText string
	// DefaultChattedAt holds the default value on creation for the "chatted_at" field.
	DefaultChattedAt int64
	// DefaultEmbedded holds the default value on creation for the "embedded" field.
	DefaultEmbedded bool
	// DefaultCreatedAt holds the default value on creation for the "created_at" field.
	DefaultCreatedAt func() int64
	// DefaultUpdatedAt holds the default value on creation for the "updated_at" field.
	DefaultUpdatedAt func() int64
	// DefaultID holds the default value on creation for the "id" field.
	DefaultID func() uuid.UUID
)

// OrderOption defines the ordering options for the TelegramChatHistories queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByChatID orders the results by the chat_id field.
func ByChatID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldChatID, opts...).ToFunc()
}

// ByMessageID orders the results by the message_id field.
func ByMessageID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldMessageID, opts...).ToFunc()
}

// ByUserID orders the results by the user_id field.
func ByUserID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldUserID, opts...).ToFunc()
}

// ByUsername orders the results by the username field.
func ByUsername(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldUsername, opts...).ToFunc()
}

// ByFullName orders the results by the full_name field.
func ByFullName(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldFullName, opts...).ToFunc()
}

// ByText orders the results by the text field.
func ByText(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldText, opts...).ToFunc()
}

// ByRepliedToMessageID orders the results by the replied_to_message_id field.
func ByRepliedToMessageID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldRepliedToMessageID, opts...).ToFunc()
}

// ByRepliedToUserID orders the results by the replied_to_user_id field.
func ByRepliedToUserID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldRepliedToUserID, opts...).ToFunc()
}

// ByRepliedToFullName orders the results by the replied_to_full_name field.
func ByRepliedToFullName(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldRepliedToFullName, opts...).ToFunc()
}

// ByRepliedToUsername orders the results by the replied_to_username field.
func ByRepliedToUsername(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldRepliedToUsername, opts...).ToFunc()
}

// ByRepliedToText orders the results by the replied_to_text field.
func ByRepliedToText(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldRepliedToText, opts...).ToFunc()
}

// ByChattedAt orders the results by the chatted_at field.
func ByChattedAt(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldChattedAt, opts...).ToFunc()
}

// ByEmbedded orders the results by the embedded field.
func ByEmbedded(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldEmbedded, opts...).ToFunc()
}

// ByCreatedAt orders the results by the created_at field.
func ByCreatedAt(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldCreatedAt, opts...).ToFunc()
}

// ByUpdatedAt orders the results by the updated_at field.
func ByUpdatedAt(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldUpdatedAt, opts...).ToFunc()
}
