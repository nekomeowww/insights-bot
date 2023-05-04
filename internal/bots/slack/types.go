package slack

type recivedCommandInfo struct {
	Command     string `form:"command"`
	Text        string `form:"text"`
	ResponseUrl string `form:"response_url"`
	UserID      string `form:"user_id"`
	ChannelID   string `form:"channel_id"`
	TeamID      string `form:"team_id"`
}

type smrRequestInfo struct {
	inputUrl    string
	channelID   string
	accessToken string
}
