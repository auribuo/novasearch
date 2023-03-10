package log

import (
	"github.com/sirupsen/logrus"
)

// Setup reads the config and configures log level and log output of all loggers
func Setup(logLevel string) error {
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:  true,
		PadLevelText: true,
	})

	lvl, err := logrus.ParseLevel(logLevel)
	if err != nil {
		return err
	}

	logrus.SetLevel(lvl)

	return nil
}
