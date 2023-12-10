package tgchat

type AutoRecapSendMode int

const (
	AutoRecapSendModePublicly                 AutoRecapSendMode = iota
	AutoRecapSendModeOnlyPrivateSubscriptions                   // Only users who subscribed to the recap will receive it
)

func (a AutoRecapSendMode) String() string {
	switch a {
	case AutoRecapSendModePublicly:
		return "公开"
	case AutoRecapSendModeOnlyPrivateSubscriptions:
		return "私聊"
	default:
		return "其他"
	}
}
