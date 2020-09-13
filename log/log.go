package log

import (
	"time"

	"github.com/4lie/nats-static-monitoring/config"

	"github.com/sirupsen/logrus"
)

func SetupLogger(logger config.Logger) {
	logLevel, err := logrus.ParseLevel(logger.Level)
	if err != nil {
		logLevel = logrus.ErrorLevel
	}

	logrus.SetLevel(logLevel)

	if logLevel == logrus.DebugLevel {
		logrus.SetReportCaller(true)
		logrus.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: time.RFC3339,
		})
	} else {
		logrus.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: time.RFC3339,
		})
	}
}
