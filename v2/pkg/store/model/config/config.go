package config

// Config global configs
type Config struct {
	Prometheus *Prometheus
}

// Prometheus configs
type Prometheus struct {
	Addr string
}
