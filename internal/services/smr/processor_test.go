package smr

import (
	"testing"

	"github.com/nekomeowww/insights-bot/internal/models/smr"
	"github.com/stretchr/testify/assert"
)

func TestService_botExists(t *testing.T) {
	t.Run("BotNotExists", func(t *testing.T) {
		a := assert.New(t)
		a.False(testService.botExists(smr.FromPlatformDiscord))
		a.False(testService.botExists(smr.FromPlatformSlack))
		a.False(testService.botExists(smr.FromPlatformTelegram))
	})
}
