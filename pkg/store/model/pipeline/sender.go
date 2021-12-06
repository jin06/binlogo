package pipeline

import (
	"github.com/Shopify/sarama"
)

const SENDER_TYPE_KAFKA = "kafka"
const SNEDER_TYPE_STDOUT = "stdout"
const SNEDER_TYPE_HTTP = "http"
const SNEDER_TYPE_RABBITMQ = "rabbitMQ"
const SENDER_TYPE_REDIS = "redis"
const SENDER_TYPE_ROCKETMQ = "rocketMQ"

// Sender output configuration
type Sender struct {
	Name     string    `json:"name"`
	Type     string    `json:"type"`
	Kafka    *Kafka    `json:"kafka"`
	Stdout   *Stdout   `json:"stdout"`
	Http     *Http     `json:"http"`
	RabbitMQ *RabbitMQ `json:"rabbitMQ"`
	Redis    *Redis    `json:"redis"`
	RocketMQ *RocketMQ `json:"rocketMQ"`
}

// Kafka output configuration
type Kafka struct {
	Brokers      string                   `json:"brokers"`
	Topic        string                   `json:"topic"`
	RequiredAcks *sarama.RequiredAcks     `json:"require_acks"`
	Compression  *sarama.CompressionCodec `json:"compression"`
	Retries      *int                     `json:"retries"`
	Idepotent    *bool                    `json:"idepotent"`
}

// Stdout output configuration
type Stdout struct {
}

// Http output configuration
type Http struct {
	API     string `json:"api"`
	Retries int    `json:"retries"`
}

// RabbitMQ basic model for pipeline config
type RabbitMQ struct {
	// rabbitMQ url
	Url          string `json:"url"`
	ExchangeName string `json:"exchange_name"`
}

// Redis output configuration
type Redis struct {
	Addr     string `json:"address"`
	UserName string `json:"username"`
	Password string `json:"password"`
	DB       int    `json:"db"`
	List     string `json:"list"`
}

// RocketMQ aliyun rocketmq configuration
type RocketMQ struct {
	// Endpoint http endpoint
	Endpoint string `json:"endpoint"`
	// TopicName topic name
	TopicName string `json:"topic_name"`
	// Instance name
	InstanceId string `json:"instance_id"`
	AccessKey  string `json:"access_key"`
	SecretKey  string `json:"secret_key"`
}
