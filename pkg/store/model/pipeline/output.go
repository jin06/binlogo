package pipeline

import "encoding/json"

// Output pipeline output
type Output struct {
	Sender *Sender `json:"sender" redis:"sender"`
}

func (m *Output) MarshalBinary() (data []byte, err error) {
	return json.Marshal(m)
}

func (m *Output) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, m)
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
