package pipeline

// RabbitMQ basic model for pipeline config
type RabbitMQ struct {
	// rabbitMQ url
	Url              string `json:"url"`
	ExchangeName     string `json:"exchange_name"`
}
