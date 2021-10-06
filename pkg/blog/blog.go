package blog

import (
	"github.com/jin06/binlogo/config"
	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetLevel(logrus.DebugLevel)
}

func Env(env config.Env) {
	level := logrus.DebugLevel
	switch env {
	case config.ENV_PRO:
		{
			level = logrus.InfoLevel
		}
	case config.ENV_DEV:
		{
			level = logrus.DebugLevel
		}
	case config.ENV_TEST:
		{
			level = logrus.DebugLevel
		}
	}
	logrus.SetLevel(level)
	logrus.Info("Set log level to: ", level.String())
	return
}
