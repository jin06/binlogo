package blog

import (
	"github.com/jin06/binlogo/configs"
	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetLevel(logrus.DebugLevel)
}

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

func Fatal(args ...interface{}) {
	logrus.Fatal(args...)
}
func Fatalln(args ...interface{}) {
	logrus.Fatalln(args...)
}
func Fatalf(format string, args ...interface{}) {
	logrus.Fatalf(format, args...)
}

func Info(args ...interface{}) {
	logrus.Info(args...)
}

func Infof(format string, args ...interface{}) {
	logrus.Infof(format, args...)
}

func Infoln(args ...interface{}) {
	logrus.Infoln(args...)
}

func Error(args ...interface{}) {
	logrus.Error(args...)
}

func Errorf(f string, args ...interface{}) {
	logrus.Errorf(f, args...)
}

func Errorln(args ...interface{}) {
	logrus.Errorln(args...)
}

func Debug(args ...interface{}) {
	logrus.Debug(args...)
}

func Debugf(f string, args ...interface{}) {
	logrus.Debugf(f, args...)
}
func Debugln(args ...interface{}) {
	logrus.Debugln(args...)
}

func Warn(args ...interface{}) {
	logrus.Warn(args...)
}

func Warnf(format string, args ...interface{}) {
	logrus.Warnf(format, args...)
}

func Warnln(args ...interface{}) {
	logrus.Warnln(args...)
}
