package output

import (
	"context"
	"github.com/jin06/binlogo/configs"
	"github.com/jin06/binlogo/pkg/store/model/pipeline"
	"testing"
	"time"
)

func TestOutput(t *testing.T) {
	configs.InitGoTest()
	outModel := &pipeline.Output{
		Sender: &pipeline.Sender{
			Type:   pipeline.SNEDER_TYPE_STDOUT,
			Stdout: &pipeline.Stdout{},
		},
	}
	out1, err := New(OptionOutput(outModel), OptionPipeName("go_test_pipe"))
	if err != nil {
		t.Error(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	err = out1.Run(ctx)
	if err != nil {
		t.Error(err)
	}
	cancel()
	outModel.Sender.Type = pipeline.SENDER_TYPE_KAFKA
	outModel.Sender.Kafka = &pipeline.Kafka{
		Brokers:      "127.0.0.1:2000",
		Topic:        "test",
		RequiredAcks: nil,
		Compression:  nil,
		Retries:      nil,
		Idepotent:    nil,
	}
	outModel.Sender.Type = pipeline.SNEDER_TYPE_HTTP
	outModel.Sender.Http = &pipeline.Http{
		API:     "http://127.0.0.1:1999/event",
		Retries: 3,
	}
	out3, _ := New(OptionOutput(outModel))
	err = out3.Run(context.Background())
	outModel.Sender.Type = pipeline.SENDER_TYPE_REDIS
	outModel.Sender.Redis = &pipeline.Redis{
		Addr:     "127.0.0.1:16379",
		UserName: "",
		Password: "",
		DB:       0,
		List:     "test123",
	}
	out4, _ := New(OptionOutput(outModel))
	err = out4.Run(context.Background())
	if err != nil {
		t.Fail()
	}

	outModel.Sender.Type = pipeline.SNEDER_TYPE_RABBITMQ
	outModel.Sender.RabbitMQ = &pipeline.RabbitMQ{
		Url:          "amqp://guest:guest@localhost:5672/",
		ExchangeName: "test123",
	}
	out5, _ := New(OptionOutput(outModel))
	err = out5.Run(context.Background())
	if err != nil {
		t.Error(err)
	}

	time.Sleep(time.Millisecond * 100)
	if err != nil {
		t.Error(err)
	}
}
