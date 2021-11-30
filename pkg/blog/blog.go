package blog

import (
	"github.com/jin06/binlogo/configs"
	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetLevel(logrus.DebugLevel)
}

// Env sets log level
func Env(env configs.Env) {
	level := logrus.DebugLevel
	switch env {
	case configs.ENV_PRO:
		{
			level = logrus.InfoLevel
		}
	case configs.ENV_DEV:
		{
			level = logrus.DebugLevel
			logrus.SetReportCaller(true)
			logrus.SetFormatter(&logrus.TextFormatter{
				//DisableColors:true,
				TimestampFormat: "2006-01-02 15:03:04",
			})
		}
	case configs.ENV_TEST:
		{
			level = logrus.DebugLevel
		}
	}
	logrus.SetLevel(level)
	logrus.Info("Set log level to: ", level.String())
	return
}
