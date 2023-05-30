package redis

import "fmt"

// Key key.
type Key string

// Format format.
func (k Key) Format(params ...interface{}) string {
	return fmt.Sprintf(string(k), params...)
}

// Timecapsule keys.

const (
	// TimeCapsuleAutoRecapSortedSetKey is the key for auto recap used timecapsule queue.
	TimeCapsuleAutoRecapSortedSetKey Key = "time_capsule/auto_recap_capsules" //  SortedSet
)

// Recap keys.

const (
	// RecapReplayFromPrivateMessageControl1 is the key for recapturing the replayed message from private messages.
	// params: enabled for user id
	RecapReplayFromPrivateMessageControl1 Key = "recap/replay_from_private_message/%d"

	// RecapReplayFromPrivateMessageBatch1 is the key for recapturing the replayed message from private messages.
	// params: enabled for user id
	RecapReplayFromPrivateMessageBatch1 Key = "recap/replay_from_private_message/%d/batch"

	// RecapPrivateSubscriptionStartCommandContext1 is the key for storing the recap private subscription start command context.
	// params: hash key
	RecapPrivateSubscriptionStartCommandContext1 Key = "recap/private_subscription/start_command_context/%s"

	// RecapSubscribeRecapStartCommandContext1 is the key for storing the recap subscribe recap start command context.
	// params: hash key
	RecapSubscribeRecapStartCommandContext1 Key = "recap/subscribe_recap/start_command_context/%s"
)

// Common keys.

const (
	// SessionRecapReplayFromPrivateMessage0 is the key for cancelling the recap replayed message from private messages.
	SessionRecapReplayFromPrivateMessage0 Key = "session/cancellable/recap_replay_from_private_message"

	// SessionDeleteLaterMessagesForActor1 is the key for deleting later messages for actor.
	// params: actor id
	SessionDeleteLaterMessagesForActor1 Key = "session/delete_later_messages_for_actor/%d" // List
)

// CallbackQueryData keys.
const (
	// CallbackQueryData2 is the key for storing callback query data.
	// params: handler route, action hash
	CallbackQueryData2 Key = "callback_query/button_data/%s/%s"
)
