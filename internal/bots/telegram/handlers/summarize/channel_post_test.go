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
	openaiClient := thirdparty.NewOpenAIClient()(thirdparty.NewOpenAIClientParam{
		Config: config,
	})

	h = NewHandler()(NewHandlerParam{
		Logger: logger,
		OpenAI: openaiClient,
	})

	os.Exit(m.Run())
}
