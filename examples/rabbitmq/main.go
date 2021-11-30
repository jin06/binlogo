package main

import (
	"fmt"
	"github.com/streadway/amqp"
)

func main() {
	pipelineName := "test-rabbitmq"
	queueName := pipelineName + "-q"
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	defer ch.Close()
	_, err = ch.QueueDeclare(
		queueName, // queue name
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		panic(err)
	}
	routingKey := "*.*"
	err = ch.QueueBind(
		queueName,    // queue name
		routingKey,   // routing key
		pipelineName, // exchange name
		false,
		nil,
	)
	if err != nil {
		panic(err)
	}
	defer func() {
		ch.QueueDelete(
			queueName,
			false,
			false,
			true,
		)
	}()

	fmt.Println("Wait messages: ")

	msgs, errCon := ch.Consume(
		queueName, // queue
		"",        // consumer
		true,      // auto ack
		false,     // exclusive
		false,     // no local
		false,     // no wait
		nil,       // args
	)
	if errCon != nil {
		fmt.Println(errCon)
	}
	go func() {
		for msg := range msgs {
			fmt.Printf(" [x] %s", msg.Body)
		}
	}()
	select {}
}
