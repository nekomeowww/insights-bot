package types

import (
	"sync"

	"github.com/nekomeowww/insights-bot/internal/models/smr"
)

type TaskInfo struct {
	Platform smr.FromPlatform `json:"platform"`
	Url      string           `json:"url"` // url to summarize

	ChatID    int64 `json:"chatID"`    // only for telegram
	MessageID int   `json:"messageID"` // used to edit the reply message of request, not work in slack or discordbot currently

	ChannelID string `json:"channelID"` // for slack and discordbot

	TeamID string `json:"teamID"` // only for slack, used to query access token and refresh token
}

type TaskQueue struct {
	queue []TaskInfo
	mu    *sync.RWMutex
}

func NewTaskQueue() *TaskQueue {
	return &TaskQueue{
		mu: &sync.RWMutex{},
	}
}

func (t *TaskQueue) AddTask(info TaskInfo) {
	t.mu.Lock()
	if len(t.queue) >= 10 {
		t.mu.Unlock()
		return
	}

	t.queue = append(t.queue, info)
	t.mu.Unlock()
}

func (t *TaskQueue) RemoveTask() {
	t.mu.Lock()
	if len(t.queue) == 0 {
		t.mu.Unlock()
		return
	}

	t.queue = t.queue[1:]
}

func (t *TaskQueue) Len() int {
	t.mu.RLock()
	l := len(t.queue)
	t.mu.RUnlock()

	return l
}
