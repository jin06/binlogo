package promeths

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/viper"
)

var (
	MessageTotalCounter  *prometheus.CounterVec
	MessageSendCounter   *prometheus.CounterVec
	MessageSendHistogram *prometheus.HistogramVec
)

func Init() {
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
	MessageSendCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: nameSpace,
			Subsystem: subSystem,
			Name:      "message_send",
		},
		pipelineLabels,
	)
	MessageSendHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: nameSpace,
			Subsystem: subSystem,
			Name:      "message_send_time",
			Buckets:   prometheus.LinearBuckets(0, 1, 10),
		},
		pipelineLabels,
	)
	prometheus.Register(MessageTotalCounter)
	prometheus.Register(MessageSendCounter)
	prometheus.Register(MessageSendHistogram)
}
