package lib

import (
	"fmt"

	"github.com/nekomeowww/insights-bot/internal/configs"
	"github.com/nekomeowww/insights-bot/pkg/logger"
	"github.com/sirupsen/logrus"
	"go.uber.org/fx"
	"go.uber.org/zap/zapcore"
)

type NewLoggerParams struct {
	fx.In

	Configs *configs.Config
}

func NewLogger() func(NewLoggerParams) (*logger.Logger, error) {
	return func(params NewLoggerParams) (*logger.Logger, error) {
		logLevel, err := zapcore.ParseLevel(params.Configs.LogLevel)
		if err != nil {
			logLevel = zapcore.InfoLevel
		}

		var isFatalLevel bool
		if logLevel == zapcore.FatalLevel {
			isFatalLevel = true
			logLevel = zapcore.InfoLevel
		}

		logger, err := logger.NewLogger(logLevel, "insights-bot", "", make([]logrus.Hook, 0))
		if err != nil {
			return nil, fmt.Errorf("failed to create logger: %w", err)
		}
		if isFatalLevel {
			logger.Error("fatal log level is unacceptable, fallbacks to info level")
		}

		return logger, nil
	}
}
