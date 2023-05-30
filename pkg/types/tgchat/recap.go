package tgchat

type AutoRecapSendMode int

const (
	AutoRecapSendModePublicly                 AutoRecapSendMode = iota
	AutoRecapSendModeOnlyPrivateSubscriptions                   // Only users who subscribed to the recap will receive it
)
