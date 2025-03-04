package blog

import (
	"github.com/jin06/binlogo/v2/configs"
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

func SetLevel(log *logrus.Logger, level string) {
	switch level {
	case DEBUG:
		log.SetLevel(logrus.DebugLevel)
		log.SetReportCaller(true)
	case INFO:
		log.SetLevel(logrus.InfoLevel)
	case WARN:
		log.SetLevel(logrus.WarnLevel)
	case ERROR:
		log.SetLevel(logrus.ErrorLevel)
	case FATAL:
		log.SetLevel(logrus.FatalLevel)
	case PANIC:
		log.SetLevel(logrus.PanicLevel)
	case TRACE:
		log.SetLevel(logrus.TraceLevel)
	default:
		log.SetLevel(logrus.InfoLevel)
	}
}

func NewLog() *logrus.Logger {
	log := logrus.New()
	log.SetFormatter(&formatter{})
	SetLevel(log, configs.Default.LogLevel)
	return log
}
