package rabbitmq

import (
	message2 "github.com/jin06/binlogo/app/pipeline/message"
	"github.com/jin06/binlogo/pkg/store/model/pipeline"
	"github.com/streadway/amqp"
)

// RabbitMQ send message to RabbitMQ
type RabbitMQ struct {
	RabbitMQ   *pipeline.RabbitMQ
	Connection *amqp.Connection
	Channel    *amqp.Channel
}

// New returns a new RabbitMQ object
func New(rq *pipeline.RabbitMQ) (r *RabbitMQ, err error) {
	r = &RabbitMQ{
		RabbitMQ:   rq,
		Connection: nil,
		Channel:    nil,
	}
	err = r.init()
	return
}

func (r *RabbitMQ) init() (err error) {
	err = r.resetConnection()
	if err != nil {
		return
	}
	err = r.exchangeDeclare()
	return
}

func (r *RabbitMQ) resetConnection() (err error) {
	if r.Connection != nil {
		_ = r.Connection.Close()
	}
	if r.Channel != nil {
		_ = r.Channel.Close()
	}
	r.Connection, err = amqp.Dial(r.RabbitMQ.Url)
	if err != nil {
		return
	}
	r.Channel, err = r.Connection.Channel()
	return
}

func (r *RabbitMQ) exchangeDeclare() (err error) {
	err = r.Channel.ExchangeDeclare(
		r.RabbitMQ.ExchangeName,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	return
}

// Send Logic and control
func (r *RabbitMQ) Send(msg *message2.Message) (ok bool, err error) {
	routingKey := msg.Content.Head.Database + "." + msg.Content.Head.Table
	body, _ := msg.JsonContent()
	rabbitMsg := amqp.Publishing{
		ContentType: "text/json",
		Body:        []byte(body),
	}
	err = r.Channel.Publish(r.RabbitMQ.ExchangeName, routingKey, false, false, rabbitMsg)
	if err == nil {
		ok = true
	}
	return
}
