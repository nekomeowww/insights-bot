package lib

import (
	"github.com/nekomeowww/insights-bot/internal/configs"
	"github.com/nekomeowww/insights-bot/pkg/logger"
	"github.com/sirupsen/logrus"
	"go.uber.org/fx"
)

type NewLoggerParams struct {
	fx.In

	Configs *configs.Config
}

func NewLogger() func(NewLoggerParams) *logger.Logger {
	return func(params NewLoggerParams) *logger.Logger {
		logLevel, err := logrus.ParseLevel(params.Configs.LogLevel)
		if err != nil {
			logLevel = logrus.InfoLevel
		}

		var isFatalLevel bool
		if logLevel == logrus.FatalLevel {
			isFatalLevel = true
			logLevel = logrus.InfoLevel
		}

		logger := logger.NewLogger(logLevel, "insights-bot", "", make([]logrus.Hook, 0))
		if err != nil {
			logger.Errorf("failed to create logger: %v, fallbacks to info level", err)
		}
		if isFatalLevel {
			logger.Errorf("fatal log level is unacceptable, fallbacks to info level")
		}

		return logger
	}
}
