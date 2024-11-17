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
const SENDER_TYPE_Elastic = "elastic"

// Sender output configuration
type Sender struct {
	Name     string    `json:"name" redis:"name"`
	Type     string    `json:"type" redis:"type"`
	Kafka    *Kafka    `json:"kafka" redis:"kafka"`
	Stdout   *Stdout   `json:"stdout" redis:"stdout"`
	Http     *Http     `json:"http" redis:"http"`
	RabbitMQ *RabbitMQ `json:"rabbitMQ" redis:"rabbitMQ"`
	Redis    *Redis    `json:"redis" redis:"redis"`
	RocketMQ *RocketMQ `json:"rocketMQ" redis:"rocketMQ"`
	// todo use interface for sender params
	Elastic Elastic `json:"elastic" redis:"elastic"`
}

// Kafka output configuration
type Kafka struct {
	Brokers      string                   `json:"brokers" redis:"brokers"`
	Topic        string                   `json:"topic" redis:"topic"`
	RequiredAcks *sarama.RequiredAcks     `json:"require_acks" redis:"require_acks"`
	Compression  *sarama.CompressionCodec `json:"compression" redis:"compression"`
	Retries      *int                     `json:"retries" redis:"retries"`
	Idepotent    *bool                    `json:"idepotent" redis:"idepotent"`
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

type Elastic struct {
	// Endpoints is elastic addresses
	Endpoints string `json:"endpoints"`
	Index     string `json:"index"`
	Type      string `json:"type"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	//TLSCA                    string `json:"tls_ca"`
	//TLSClientCertificateFile string `json:"tls_client_certificate_file"`
	//TLSClientKeyFile         string `json:"tls_client_key_file"`
}

// EmptyHttp return a new empty Http object
func EmptyHttp() *Http {
	return &Http{}
}
