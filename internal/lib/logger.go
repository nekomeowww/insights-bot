package lib

import (
	"github.com/nekomeowww/insights-bot/pkg/logger"
	"github.com/sirupsen/logrus"
)

func NewLogger() func() *logger.Logger {
	return func() *logger.Logger {
		return logger.NewLogger(logrus.InfoLevel, "insights-bot", "", make([]logrus.Hook, 0))
	}
}
