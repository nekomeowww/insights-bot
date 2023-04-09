package summarize

import (
	"os"
	"testing"

	"github.com/nekomeowww/insights-bot/internal/configs"
	"github.com/nekomeowww/insights-bot/internal/thirdparty"
	"github.com/nekomeowww/insights-bot/pkg/logger"
	"github.com/sirupsen/logrus"
)

var h *Handler

func TestMain(m *testing.M) {
	logger := logger.NewLogger(logrus.DebugLevel, "insights-bot", "", nil)
	config := configs.NewConfig()()
	openaiClient, err := thirdparty.NewOpenAIClient()(thirdparty.NewOpenAIClientParam{
		Config: config,
	})
	if err != nil {
		os.Exit(1)
	}

	h = NewHandler()(NewHandlerParam{
		Logger: logger,
		OpenAI: openaiClient,
	})

	os.Exit(m.Run())
}
