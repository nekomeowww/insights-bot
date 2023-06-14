package bot

type FromPlatform int

func (f FromPlatform) String() string {
	switch f {
	case FromPlatformTelegram:
		return "Telegram"
	case FromPlatformSlack:
		return "Slack"
	case FromPlatformDiscord:
		return "Discord"
	default:
		return "Unknown"
	}
}

const (
	FromPlatformTelegram FromPlatform = iota
	FromPlatformSlack
	FromPlatformDiscord
)
