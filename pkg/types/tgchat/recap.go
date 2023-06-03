package tgchat

type AutoRecapSendMode int

const (
	AutoRecapSendModePublicly                 AutoRecapSendMode = iota
	AutoRecapSendModeOnlyPrivateSubscriptions                   // Only users who subscribed to the recap will receive it
)

func (a AutoRecapSendMode) String() string {
	switch a {
	case AutoRecapSendModePublicly:
		return "公开模式"
	case AutoRecapSendModeOnlyPrivateSubscriptions:
		return "私聊订阅模式"
	default:
		return "其他模式"
	}
}
