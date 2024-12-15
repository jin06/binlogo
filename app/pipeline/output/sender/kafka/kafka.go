package kafka

import (
	"encoding/json"
	"strings"

	"github.com/Shopify/sarama"
	message2 "github.com/jin06/binlogo/v2/app/pipeline/message"
	"github.com/jin06/binlogo/v2/pkg/store/model/pipeline"
	"github.com/sirupsen/logrus"
)

// Kafka send message to kafka
type Kafka struct {
	Kafka        *pipeline.Kafka
	SyncProducer sarama.SyncProducer
}

// New returns a new Kafka
func New(kafka pipeline.Kafka) (kaf *Kafka, err error) {
	kaf = &Kafka{Kafka: &kafka}
	err = kaf.init()
	return
}

func (s *Kafka) init() (err error) {
	addr := strings.Split(s.Kafka.Brokers, ",")
	cfg := sarama.NewConfig()
	cfg.Producer.Return.Successes = true
	if s.Kafka.RequiredAcks != nil {
		cfg.Producer.RequiredAcks = *s.Kafka.RequiredAcks
	}
	if s.Kafka.Compression != nil {
		cfg.Producer.Compression = *s.Kafka.Compression
	}
	if s.Kafka.Retries != nil {
		cfg.Producer.Retry.Max = *s.Kafka.Retries
	}

	producer, err := sarama.NewSyncProducer(addr, cfg)
	if err != nil {
		logrus.Error(err)
		return err
	}
	s.SyncProducer = producer
	return
}

// Send Logic and control
func (s *Kafka) Send(msg *message2.Message) (ok bool, err error) {
	return s.doSend(msg)
}

func (s *Kafka) doSend(msg *message2.Message) (ok bool, err error) {
	keyByte := []byte(msg.Content.Head.Database + "." + msg.Content.Head.Table)
	valByte, _ := json.Marshal(msg.Content)
	pMsg := sarama.ProducerMessage{
		Topic: s.Kafka.Topic,
		Key:   sarama.ByteEncoder(keyByte),
		Value: sarama.ByteEncoder(valByte),
	}
	par, off, err := s.SyncProducer.SendMessage(&pMsg)
	logrus.Debugf("Send to Kafka partition %d at offset %d\n", par, off)
	if err == nil {
		ok = true
	}
	return
}

type kafkaKey struct {
	database string `json:"database"`
	table    string `json:"table"`
}

func (s *Kafka) Close() error {
	if s.SyncProducer != nil {
		return s.SyncProducer.Close()
	}
	return nil
}
