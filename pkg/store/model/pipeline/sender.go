package pipeline

const SENDER_TYPE_KAFKA = "kafka"
const SNEDER_TYPE_STDOUT = "stdout"

type Sender struct {
	Name   string  `json:"name"`
	Type   string  `json:"type"`
	Kafka  *Kafka  `json:"kafka"`
	Stdout *Stdout `json:"stdout"`
}

type Kafka struct {
	Brokers []string `json:"brokers"`
	Topic   string   `json:"topic"`
}

type Stdout struct {
}
