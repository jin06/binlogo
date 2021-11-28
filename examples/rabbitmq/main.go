package main

import (
	"fmt"
	"github.com/streadway/amqp"
)

func main() {
	pipelineName := "pipe-rabbitmq"
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
	msgs, err := ch.Consume(
		queueName, // queue
		"",        // consumer
		true,      // auto ack
		false,     // exclusive
		false,     // no local
		false,     // no wait
		nil,       // args
	)

	fmt.Println("Wait messages: ")

	for d := range msgs {
		fmt.Printf(" [x] %s", d.Body)
	}
}
