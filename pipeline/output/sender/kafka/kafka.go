package kafka

import (
	"encoding/json"
	"github.com/Shopify/sarama"
	"github.com/jin06/binlogo/pipeline/message"
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
	err = kaf.Init()
	return
}

func (s *Kafka) Init() (err error) {
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

func (s *Kafka) Send(ch chan *message.Message) (ok bool, err error) {
	go s.doSend(ch)
	return
}

func (s *Kafka) doSend(ch chan *message.Message) (ok bool, err error) {
	defer func() {
		err2 := s.SyncProducer.Close()
		if err2 != nil {
			logrus.Error(err2.Error())
		}
	}()
	for {
		logrus.Debug("Wait send to kafka")
		msg := <-ch
		logrus.Debug("Read msg: ", msg)
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
		par, off, err1 := s.SyncProducer.SendMessage(&pMsg)
		if err1 != nil {
			logrus.Error("Failed to send message to kafka: ", err1)
		}
		logrus.Debugf("Send to Kafka partition %d at offset %d\n", par, off)
	}
	return
}

type kafkaKey struct {
	database string `json:"database"`
	table    string `json:"table"`
}
