package model

const SENDER_TYPE_KAFKA = "kafka"

type Sender struct {
	Name  string `json:"sender"`
	Type  string `json:"type"`
	Kafka *Kafka `json:"kafka"`
}

type Kafka struct {
	Brokers []string `json:"brokers"`
	Topic   string   `json:"topic"`
}
