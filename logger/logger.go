package logger

import (
	"github.com/sirupsen/logrus"
)

var App *logrus.Logger

func New(logLevel string) *logrus.Logger {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})

	switch logLevel {
	case "debug", "5":
		logger.SetLevel(logrus.DebugLevel)
	case "info", "4":
		logger.SetLevel(logrus.InfoLevel)
	case "warn", "3":
		logger.SetLevel(logrus.WarnLevel)
	case "error", "2":
		logger.SetLevel(logrus.ErrorLevel)
	case "fatal", "1":
		logger.SetLevel(logrus.FatalLevel)
	default:
		logger.SetLevel(logrus.InfoLevel)
	}
	return logger
}
