package logging

import (
	"github.com/sirupsen/logrus"
)

func New() *logrus.Logger {
	logger := logrus.New()

	// logger.SetReportCaller(true)
	logger.SetLevel(logrus.DebugLevel)

	logger.SetFormatter(&logrus.TextFormatter{
		ForceColors:   true,
		DisableColors: false,
	})

	return logger
}
