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
)

// Session keys.

const (
	// SessionRecapReplayFromPrivateMessage0 is the key for cancelling the recap replayed message from private messages.
	SessionRecapReplayFromPrivateMessage0 Key = "session/cancellable/recap_replay_from_private_message"
)
