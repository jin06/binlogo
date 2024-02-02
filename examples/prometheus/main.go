package main

import (
	"net/http"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	counter := prometheus.NewCounterVec(
		prometheus.CounterOpts{Namespace: "binlogo", Subsystem: "cluster1", Name: "test2", Help: "just test"},
		[]string{"pipeline", "node"},
	)
	prometheus.MustRegister(counter)
	go func() {
		for range time.Tick(time.Millisecond * 100) {
			counter.With(prometheus.Labels{"pipeline": "p1", "node": "n1"}).Inc()
		}
	}()
	gauger := prometheus.NewGauge(prometheus.GaugeOpts{Namespace: "binlogo", Subsystem: "cluster1", Name: "gauger", Help: "gauger"})
	go func() {
		for range time.Tick(time.Millisecond * 100) {
			gauger.Add(1)
		}
	}()
	prometheus.MustRegister(gauger)
	h := prometheus.NewHistogram(prometheus.HistogramOpts{Namespace: "binlogo", Subsystem: "cluster1", Name: "histogram2", Help: "his", Buckets: []float64{1.00, 5.00, 10.00}})
	prometheus.MustRegister(h)
	go func() {
		for {
			for i := 0; i < 10; i++ {
				h.Observe(float64(i) + 0)
				time.Sleep(time.Second)
			}
		}
	}()
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":2112", nil)
	strings.FieldsFunc()
}
