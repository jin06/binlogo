package pipeline

// Output pipeline output
type Output struct {
	Sender *Sender `json:"sender"`
}

// EmptyOutput return a new empty Output object
func EmptyOutput() *Output {
	return &Output{
		Sender: &Sender{
			Name:     "",
			Type:     "",
			Kafka:    nil,
			Stdout:   nil,
			Http:     &Http{},
			RabbitMQ: nil,
			Redis:    nil,
			RocketMQ: nil,
		},
	}
}
