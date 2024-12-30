package blog

import (
	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&formatter{})
}

const (
	DEBUG = "debug"
	INFO  = "info"
	WARN  = "warn"
	ERROR = "error"
	FATAL = "fatal"
	PANIC = "panic"
	TRACE = "trace"
)

func SetLevel(level string) {
	switch level {
	case DEBUG:
		logrus.SetLevel(logrus.DebugLevel)
	case INFO:
		logrus.SetLevel(logrus.InfoLevel)
	case WARN:
		logrus.SetLevel(logrus.WarnLevel)
	case ERROR:
		logrus.SetLevel(logrus.ErrorLevel)
	case FATAL:
		logrus.SetLevel(logrus.FatalLevel)
	case PANIC:
		logrus.SetLevel(logrus.PanicLevel)
	case TRACE:
		logrus.SetLevel(logrus.TraceLevel)
	default:
		logrus.SetLevel(logrus.InfoLevel)
	}
}
