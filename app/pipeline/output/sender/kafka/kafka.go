package kafka

import (
	"encoding/json"
	"github.com/Shopify/sarama"
	message2 "github.com/jin06/binlogo/app/pipeline/message"
	"github.com/sirupsen/logrus"
)

type Kafka struct {
	Options      *Options
	SyncProducer sarama.SyncProducer
}

func New(options *Options) (kaf *Kafka, err error) {
	kaf = &Kafka{
		Options: options,
	}
	err = kaf.init()
	return
}

func (s *Kafka) init() (err error) {
	addr := s.Options.Kafka.Brokers
	cfg := sarama.NewConfig()
	cfg.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer(addr, cfg)
	if err != nil {
		logrus.Error(err)
		panic(err)
	}
	s.SyncProducer = producer
	return
}

func (s *Kafka) Send(msg *message2.Message) (ok bool, err error) {
	return s.doSend(msg)
}

func (s *Kafka) doSend(msg *message2.Message) (ok bool, err error) {
	key := kafkaKey{
		database: msg.Content.Head.Database,
		table:    msg.Content.Head.Table,
	}
	keyByte, _ := json.Marshal(key)
	valByte, _ := json.Marshal(msg.Content)
	pMsg := sarama.ProducerMessage{
		Topic: s.Options.Kafka.Topic,
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
