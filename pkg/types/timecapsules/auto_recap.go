package timecapsules

const (
	AutoRecapTimeCapsuleKey = "auto_recap_time_capsule"
)

type AutoRecapCapsule struct {
	ChatID int64 `json:"chat_id"`
}
