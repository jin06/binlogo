package output

import (
	"context"
	"errors"
	"time"

	message2 "github.com/jin06/binlogo/app/pipeline/message"
	sender2 "github.com/jin06/binlogo/app/pipeline/output/sender"
	"github.com/jin06/binlogo/app/pipeline/output/sender/http"
	kafka2 "github.com/jin06/binlogo/app/pipeline/output/sender/kafka"
	"github.com/jin06/binlogo/app/pipeline/output/sender/rabbitmq"
	"github.com/jin06/binlogo/app/pipeline/output/sender/redis"
	"github.com/jin06/binlogo/app/pipeline/output/sender/rocketmq"
	stdout2 "github.com/jin06/binlogo/app/pipeline/output/sender/stdout"
	"github.com/jin06/binlogo/configs"
	"github.com/jin06/binlogo/pkg/event"
	prometheus2 "github.com/jin06/binlogo/pkg/promethues"
	"github.com/jin06/binlogo/pkg/store/dao/dao_pipe"
	event2 "github.com/jin06/binlogo/pkg/store/model/event"
	"github.com/jin06/binlogo/pkg/store/model/pipeline"
	"github.com/prometheus/client_golang/prometheus"
)

// Output handle message output
// depends pipeline config, send message to stdout、kafka, etc.
type Output struct {
	InChan  chan *message2.Message
	Sender  sender2.Sender
	Options *Options
	ctx     context.Context
}

// New return a Output object
func New(opts ...Option) (out *Output, err error) {
	options := &Options{}
	for _, v := range opts {
		v(options)
	}
	out = &Output{
		Options: options,
	}
	return
}

func (o *Output) init() (err error) {
	switch o.Options.Output.Sender.Type {
	case pipeline.SNEDER_TYPE_STDOUT:
		o.Sender, err = stdout2.New()
	case pipeline.SNEDER_TYPE_HTTP:
		o.Sender, err = http.New(o.Options.Output.Sender.Http)
	case pipeline.SNEDER_TYPE_RABBITMQ:
		o.Sender, err = rabbitmq.New(o.Options.Output.Sender.RabbitMQ)
	case pipeline.SENDER_TYPE_REDIS:
		o.Sender, err = redis.New(o.Options.Output.Sender.Redis)
	case pipeline.SENDER_TYPE_KAFKA:
		o.Sender, err = kafka2.New(o.Options.Output.Sender.Kafka)
	case pipeline.SENDER_TYPE_ROCKETMQ:
		o.Sender, err = rocketmq.New(o.Options.Output.Sender.RocketMQ)
	default:
		o.Sender, err = stdout2.New()
	}

	return
}

func recordPosition(msg *message2.Message) (err error) {
	//logrus.Debugf(
	//	"Record new replication position, file %s, pos %v",
	//	msg.BinlogPosition.BinlogFile,
	//	msg.BinlogPosition.BinlogPosition,
	//)
	err = dao_pipe.UpdatePosition(msg.BinlogPosition)
	return
}

func (o *Output) loopHandle(ctx context.Context, msg *message2.Message) (err error) {
	for i := 0; i < 3; i++ {
		err = o.handle(msg)
		if err == nil {
			return
		}
		select {
		case <-ctx.Done():
			{
				return
			}
		case <-time.Tick(time.Second):
			{
				continue
			}
		}
	}
	return
}

func (o *Output) handle(msg *message2.Message) (err error) {
	defer func() {
		if err != nil {
			event.Event(event2.NewErrorPipeline(o.Options.PipelineName, "send message error: "+err.Error()))
		}
	}()
	switch msg.Status {
	case message2.STATUS_RECORD:
		{
			return
		}
	case message2.STATUS_SEND:
		{
			err = recordPosition(msg)
			if err == nil {
				msg.Status = message2.STATUS_RECORD
			}
			return
		}
	default:
		if !msg.Filter {
			var ok bool
			ok, err = o.Sender.Send(msg)
			if err == nil {
				if !ok {
					err = errors.New("send message failed")
				}
			}
			if err == nil {
				msg.Status = message2.STATUS_SEND
			}
		}
		if err == nil {
			err = recordPosition(msg)
			if err == nil {
				msg.Status = message2.STATUS_RECORD
			}
		}
		return
	}
}

// Run start Output to send message
func (o *Output) Run(ctx context.Context) (err error) {
	err = o.init()
	if err != nil {
		return
	}
	myCtx, cancel := context.WithCancel(ctx)
	o.ctx = myCtx
	go func() {
		counterSend := prometheus2.RegisterPipelineCounter(o.Options.PipelineName, "message_send")
		defer func() {
			prometheus.Unregister(counterSend)
			cancel()
		}()
		for {
			select {
			case <-ctx.Done():
				{
					return
				}
			case msg := <-o.InChan:
				{
					counterSend.With(prometheus.Labels{"pipeline": o.Options.PipelineName, "node": configs.NodeName}).Inc()
					if err1 := o.loopHandle(ctx, msg); err1 != nil {
						return
					}
				}
			}
		}
	}()
	return
}

// Context return Output's context
func (o *Output) Context() context.Context {
	return o.ctx
}
