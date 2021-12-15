package main

import (
	"fmt"

	"github.com/Shopify/sarama"
)

func main() {
	addr := []string{
		//"127.0.0.1:9092",
		"kafka-banana.shan.svc.cluster.local:9092",
	}
	topic := "product_incr_release2"
	cfg := sarama.NewConfig()
	cfg.ClientID = "test123"
	consume, err := sarama.NewConsumer(addr, cfg)
	if err != nil {
		panic(err)
	}
	par, err := consume.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		panic(err)
	}
	for {
		fmt.Println("wait message")
		select {
		case msg := <-par.Messages():
			fmt.Println(string(msg.Value))
		}
	}
}
