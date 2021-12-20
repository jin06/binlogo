package promeths

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	MessageTotalCounter   *prometheus.CounterVec
	MessageSendCounter    *prometheus.CounterVec
	MessageSendErrCounter *prometheus.CounterVec
	MessageFilterCounter  *prometheus.CounterVec
	MessageSendHistogram  *prometheus.HistogramVec
)

func Init() {
	logrus.Info("init prometheus")
	pipelineLabels := []string{"pipeline", "node"}
	nameSpace := "binlogo"
	subSystem := viper.GetString("cluster.name")

	MessageTotalCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: nameSpace,
			Subsystem: subSystem,
			Name:      "message_total",
		},
		pipelineLabels,
	)
	prometheus.Register(MessageTotalCounter)
	MessageSendCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: nameSpace,
			Subsystem: subSystem,
			Name:      "message_send",
		},
		pipelineLabels,
	)
	prometheus.Register(MessageSendCounter)
	MessageSendErrCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: nameSpace,
			Subsystem: subSystem,
			Name:      "message_send_error",
		},
		pipelineLabels,
	)
	prometheus.Register(MessageSendErrCounter)
	MessageFilterCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: nameSpace,
			Subsystem: subSystem,
			Name:      "message_filter",
		},
		pipelineLabels,
	)
	prometheus.Register(MessageFilterCounter)
	MessageSendHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: nameSpace,
			Subsystem: subSystem,
			Name:      "message_send_time",
			Buckets:   prometheus.LinearBuckets(0, 1, 10),
		},
		pipelineLabels,
	)
	prometheus.Register(MessageSendHistogram)

}
