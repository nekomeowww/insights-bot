package slack

type recivedCommandInfo struct {
	Command     string `form:"command"`
	Text        string `form:"text"`
	ResponseUrl string `form:"response_url"`
	UserId      string `form:"user_id"`
	ChannelId   string `form:"channel_id"`
	TeamId      string `form:"team_id"`
}

type smrRequestInfo struct {
	inputUrl    string
	channelId   string
	accessToken string
}
