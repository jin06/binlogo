package main

import (
	"fmt"
	"github.com/Shopify/sarama"
)

func main() {
	addr := []string{
		"kafka-banana.shan.svc.cluster.local:9092",
	}
	cfg := sarama.NewConfig()
	cfg.ClientID = "test123"
	consume, err := sarama.NewConsumer(addr,cfg)
	if err != nil {
		panic(err)
	}
	par, err := consume.ConsumePartition("binlogo-test1", 0, sarama.OffsetNewest)
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
