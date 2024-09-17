package rabbitmq

import (
	"testing"

	message2 "github.com/jin06/binlogo/v2/app/pipeline/message"
	"github.com/jin06/binlogo/v2/pkg/store/model/pipeline"
)

func TestRabbitmq(t *testing.T) {
	cfg := pipeline.RabbitMQ{
		Url:          "amqp://guest:guest@localhost:5672/",
		ExchangeName: "go_test_message",
	}
	r, err := New(&cfg)
	if err != nil {
		t.Error(err)
	}
	msg := message2.New()
	ok, err := r.Send(msg)
	if err != nil {
		t.Error(err)
	}
	if !ok {
		t.Fail()
	}
}
