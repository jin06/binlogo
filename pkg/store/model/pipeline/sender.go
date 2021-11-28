package pipeline

import (
	"github.com/Shopify/sarama"
)

const SENDER_TYPE_KAFKA = "kafka"
const SNEDER_TYPE_STDOUT = "stdout"
const SNEDER_TYPE_HTTP = "http"
const SNEDER_TYPE_RABBITMQ = "rabbitMQ"

type Sender struct {
	Name     string    `json:"name"`
	Type     string    `json:"type"`
	Kafka    *Kafka    `json:"kafka"`
	Stdout   *Stdout   `json:"stdout"`
	Http     *Http     `json:"http"`
	RabbitMQ *RabbitMQ `json:"rabbitMQ"`
}

type Kafka struct {
	Brokers      string                   `json:"brokers"`
	Topic        string                   `json:"topic"`
	RequiredAcks *sarama.RequiredAcks     `json:"require_acks"`
	Compression  *sarama.CompressionCodec `json:"compression"`
	Retries      *int                     `json:"retries"`
	Idepotent    *bool                    `json:"idepotent"`
}

type Stdout struct {
}

type Http struct {
	API     string `json:"api"`
	Retries int    `json:"retries"`
}
