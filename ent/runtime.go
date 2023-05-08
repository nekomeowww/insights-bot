// Code generated by ent, DO NOT EDIT.

package ent

import (
	"github.com/google/uuid"
	"github.com/nekomeowww/insights-bot/ent/chathistories"
	"github.com/nekomeowww/insights-bot/ent/schema"
	"github.com/nekomeowww/insights-bot/ent/slackoauthcredentials"
	"github.com/nekomeowww/insights-bot/ent/telegramchatfeatureflags"
)

// The init function reads all schema descriptors with runtime code
// (default values, validators, hooks and policies) and stitches it
// to their package variables.
func init() {
	chathistoriesFields := schema.ChatHistories{}.Fields()
	_ = chathistoriesFields
	// chathistoriesDescChatID is the schema descriptor for chat_id field.
	chathistoriesDescChatID := chathistoriesFields[1].Descriptor()
	// chathistories.DefaultChatID holds the default value on creation for the chat_id field.
	chathistories.DefaultChatID = chathistoriesDescChatID.Default.(int64)
	// chathistoriesDescMessageID is the schema descriptor for message_id field.
	chathistoriesDescMessageID := chathistoriesFields[2].Descriptor()
	// chathistories.DefaultMessageID holds the default value on creation for the message_id field.
	chathistories.DefaultMessageID = chathistoriesDescMessageID.Default.(int64)
	// chathistoriesDescUserID is the schema descriptor for user_id field.
	chathistoriesDescUserID := chathistoriesFields[3].Descriptor()
	// chathistories.DefaultUserID holds the default value on creation for the user_id field.
	chathistories.DefaultUserID = chathistoriesDescUserID.Default.(int64)
	// chathistoriesDescUsername is the schema descriptor for username field.
	chathistoriesDescUsername := chathistoriesFields[4].Descriptor()
	// chathistories.DefaultUsername holds the default value on creation for the username field.
	chathistories.DefaultUsername = chathistoriesDescUsername.Default.(string)
	// chathistoriesDescFullName is the schema descriptor for full_name field.
	chathistoriesDescFullName := chathistoriesFields[5].Descriptor()
	// chathistories.DefaultFullName holds the default value on creation for the full_name field.
	chathistories.DefaultFullName = chathistoriesDescFullName.Default.(string)
	// chathistoriesDescText is the schema descriptor for text field.
	chathistoriesDescText := chathistoriesFields[6].Descriptor()
	// chathistories.DefaultText holds the default value on creation for the text field.
	chathistories.DefaultText = chathistoriesDescText.Default.(string)
	// chathistoriesDescRepliedToMessageID is the schema descriptor for replied_to_message_id field.
	chathistoriesDescRepliedToMessageID := chathistoriesFields[7].Descriptor()
	// chathistories.DefaultRepliedToMessageID holds the default value on creation for the replied_to_message_id field.
	chathistories.DefaultRepliedToMessageID = chathistoriesDescRepliedToMessageID.Default.(int64)
	// chathistoriesDescRepliedToUserID is the schema descriptor for replied_to_user_id field.
	chathistoriesDescRepliedToUserID := chathistoriesFields[8].Descriptor()
	// chathistories.DefaultRepliedToUserID holds the default value on creation for the replied_to_user_id field.
	chathistories.DefaultRepliedToUserID = chathistoriesDescRepliedToUserID.Default.(int64)
	// chathistoriesDescRepliedToFullName is the schema descriptor for replied_to_full_name field.
	chathistoriesDescRepliedToFullName := chathistoriesFields[9].Descriptor()
	// chathistories.DefaultRepliedToFullName holds the default value on creation for the replied_to_full_name field.
	chathistories.DefaultRepliedToFullName = chathistoriesDescRepliedToFullName.Default.(string)
	// chathistoriesDescRepliedToUsername is the schema descriptor for replied_to_username field.
	chathistoriesDescRepliedToUsername := chathistoriesFields[10].Descriptor()
	// chathistories.DefaultRepliedToUsername holds the default value on creation for the replied_to_username field.
	chathistories.DefaultRepliedToUsername = chathistoriesDescRepliedToUsername.Default.(string)
	// chathistoriesDescRepliedToText is the schema descriptor for replied_to_text field.
	chathistoriesDescRepliedToText := chathistoriesFields[11].Descriptor()
	// chathistories.DefaultRepliedToText holds the default value on creation for the replied_to_text field.
	chathistories.DefaultRepliedToText = chathistoriesDescRepliedToText.Default.(string)
	// chathistoriesDescChattedAt is the schema descriptor for chatted_at field.
	chathistoriesDescChattedAt := chathistoriesFields[12].Descriptor()
	// chathistories.DefaultChattedAt holds the default value on creation for the chatted_at field.
	chathistories.DefaultChattedAt = chathistoriesDescChattedAt.Default.(func() int64)
	// chathistoriesDescEmbedded is the schema descriptor for embedded field.
	chathistoriesDescEmbedded := chathistoriesFields[13].Descriptor()
	// chathistories.DefaultEmbedded holds the default value on creation for the embedded field.
	chathistories.DefaultEmbedded = chathistoriesDescEmbedded.Default.(bool)
	// chathistoriesDescCreatedAt is the schema descriptor for created_at field.
	chathistoriesDescCreatedAt := chathistoriesFields[14].Descriptor()
	// chathistories.DefaultCreatedAt holds the default value on creation for the created_at field.
	chathistories.DefaultCreatedAt = chathistoriesDescCreatedAt.Default.(func() int64)
	// chathistoriesDescUpdatedAt is the schema descriptor for updated_at field.
	chathistoriesDescUpdatedAt := chathistoriesFields[15].Descriptor()
	// chathistories.DefaultUpdatedAt holds the default value on creation for the updated_at field.
	chathistories.DefaultUpdatedAt = chathistoriesDescUpdatedAt.Default.(func() int64)
	// chathistoriesDescID is the schema descriptor for id field.
	chathistoriesDescID := chathistoriesFields[0].Descriptor()
	// chathistories.DefaultID holds the default value on creation for the id field.
	chathistories.DefaultID = chathistoriesDescID.Default.(func() uuid.UUID)
	slackoauthcredentialsFields := schema.SlackOAuthCredentials{}.Fields()
	_ = slackoauthcredentialsFields
	// slackoauthcredentialsDescTeamID is the schema descriptor for team_id field.
	slackoauthcredentialsDescTeamID := slackoauthcredentialsFields[1].Descriptor()
	// slackoauthcredentials.TeamIDValidator is a validator for the "team_id" field. It is called by the builders before save.
	slackoauthcredentials.TeamIDValidator = slackoauthcredentialsDescTeamID.Validators[0].(func(string) error)
	// slackoauthcredentialsDescRefreshToken is the schema descriptor for refresh_token field.
	slackoauthcredentialsDescRefreshToken := slackoauthcredentialsFields[2].Descriptor()
	// slackoauthcredentials.DefaultRefreshToken holds the default value on creation for the refresh_token field.
	slackoauthcredentials.DefaultRefreshToken = slackoauthcredentialsDescRefreshToken.Default.(string)
	// slackoauthcredentialsDescAccessToken is the schema descriptor for access_token field.
	slackoauthcredentialsDescAccessToken := slackoauthcredentialsFields[3].Descriptor()
	// slackoauthcredentials.AccessTokenValidator is a validator for the "access_token" field. It is called by the builders before save.
	slackoauthcredentials.AccessTokenValidator = slackoauthcredentialsDescAccessToken.Validators[0].(func(string) error)
	// slackoauthcredentialsDescCreatedAt is the schema descriptor for created_at field.
	slackoauthcredentialsDescCreatedAt := slackoauthcredentialsFields[4].Descriptor()
	// slackoauthcredentials.DefaultCreatedAt holds the default value on creation for the created_at field.
	slackoauthcredentials.DefaultCreatedAt = slackoauthcredentialsDescCreatedAt.Default.(func() int64)
	// slackoauthcredentialsDescUpdatedAt is the schema descriptor for updated_at field.
	slackoauthcredentialsDescUpdatedAt := slackoauthcredentialsFields[5].Descriptor()
	// slackoauthcredentials.DefaultUpdatedAt holds the default value on creation for the updated_at field.
	slackoauthcredentials.DefaultUpdatedAt = slackoauthcredentialsDescUpdatedAt.Default.(func() int64)
	// slackoauthcredentialsDescID is the schema descriptor for id field.
	slackoauthcredentialsDescID := slackoauthcredentialsFields[0].Descriptor()
	// slackoauthcredentials.DefaultID holds the default value on creation for the id field.
	slackoauthcredentials.DefaultID = slackoauthcredentialsDescID.Default.(func() uuid.UUID)
	telegramchatfeatureflagsFields := schema.TelegramChatFeatureFlags{}.Fields()
	_ = telegramchatfeatureflagsFields
	// telegramchatfeatureflagsDescCreatedAt is the schema descriptor for created_at field.
	telegramchatfeatureflagsDescCreatedAt := telegramchatfeatureflagsFields[4].Descriptor()
	// telegramchatfeatureflags.DefaultCreatedAt holds the default value on creation for the created_at field.
	telegramchatfeatureflags.DefaultCreatedAt = telegramchatfeatureflagsDescCreatedAt.Default.(func() int64)
	// telegramchatfeatureflagsDescUpdatedAt is the schema descriptor for updated_at field.
	telegramchatfeatureflagsDescUpdatedAt := telegramchatfeatureflagsFields[5].Descriptor()
	// telegramchatfeatureflags.DefaultUpdatedAt holds the default value on creation for the updated_at field.
	telegramchatfeatureflags.DefaultUpdatedAt = telegramchatfeatureflagsDescUpdatedAt.Default.(func() int64)
	// telegramchatfeatureflagsDescID is the schema descriptor for id field.
	telegramchatfeatureflagsDescID := telegramchatfeatureflagsFields[0].Descriptor()
	// telegramchatfeatureflags.DefaultID holds the default value on creation for the id field.
	telegramchatfeatureflags.DefaultID = telegramchatfeatureflagsDescID.Default.(func() uuid.UUID)
}
