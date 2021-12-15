package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
)

// RegisterPipelineCounter register a new pipeline counter
func RegisterPipelineCounter(pName string, name string) *prometheus.CounterVec {
	counter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "binlogo",
			Subsystem: "cluster1",
			Name:      name,
		},
		[]string{"pipeline", "node"},
	)
	prometheus.MustRegister(counter)
	return counter
}
