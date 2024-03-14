package smr

import (
	"sync"

	"github.com/nekomeowww/insights-bot/pkg/types/bot"
)

type TaskInfo struct {
	Platform bot.FromPlatform `json:"platform"`
	URL      string           `json:"url"` // url to summarize
	Language string           `json:"language"`

	ChatID    int64 `json:"chatID"`    // only for telegram
	MessageID int   `json:"messageID"` // used to edit the reply message of request, not work in slack or discordbot currently

	ChannelID string `json:"channelID"` // for slack and discordbot

	TeamID string `json:"teamID"` // only for slack, used to query access token and refresh token
}

type OngoingTaskPool struct {
	tasks []TaskInfo
	mu    *sync.RWMutex
}

func NewOngoingTaskPool() *OngoingTaskPool {
	return &OngoingTaskPool{
		mu: &sync.RWMutex{},
	}
}

func (t *OngoingTaskPool) Add(info TaskInfo) {
	t.mu.Lock()
	t.tasks = append(t.tasks, info)
	t.mu.Unlock()
}

func (t *OngoingTaskPool) Remove() {
	t.mu.Lock()
	if len(t.tasks) == 0 {
		t.mu.Unlock()
		return
	}

	t.tasks = t.tasks[1:]
	t.mu.Unlock()
}

func (t *OngoingTaskPool) Len() int {
	t.mu.RLock()
	l := len(t.tasks)
	t.mu.RUnlock()

	return l
}
