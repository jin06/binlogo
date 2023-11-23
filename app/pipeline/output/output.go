package output

import (
	"context"
	"errors"
	"github.com/jin06/binlogo/app/pipeline/output/sender/elastic"
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
	"github.com/jin06/binlogo/pkg/promeths"
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
	record  *pipeline.RecordPosition
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
	case pipeline.SENDER_TYPE_Elastic:
		o.Sender, err = elastic.New(o.Options.Output.Sender.Elastic)
	default:
		o.Sender, err = stdout2.New()
	}

	return
}

func (o *Output) loopHandle(ctx context.Context, msg *message2.Message) (err error) {
	for i := 0; i < 3; i++ {
		err = o.handle(msg)
		if err != nil {
			promeths.MessageSendErrCounter.With(prometheus.Labels{"pipeline": o.Options.PipelineName, "node": configs.NodeName}).Inc()
		} else {
			return
		}
		// time.Sleep(time.Second)
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
			err = o.syncRecord()
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
			err = o.syncRecord()
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
		defer func() {
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
					check, errPrepare := o.prepareRecord(msg)
					if errPrepare != nil {
						message2.Put(msg)
						return
					}
					if !check {
						message2.Put(msg)
						continue
					}
					if err1 := o.loopHandle(ctx, msg); err1 != nil {
						message2.Put(msg)
						return
					}
					promeths.MessageSendCounter.With(prometheus.Labels{"pipeline": o.Options.PipelineName, "node": configs.NodeName}).Inc()
					pass := uint32(time.Now().Unix()) - msg.Content.Head.Time
					promeths.MessageSendHistogram.With(prometheus.Labels{"pipeline": o.Options.PipelineName, "node": configs.NodeName}).Observe(float64(pass))
					message2.Put(msg)
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

func (o *Output) prepareRecord(msg *message2.Message) (pass bool, err error) {
	pass = true
	if o.record == nil {
		o.record, err = dao_pipe.GetRecord(o.Options.PipelineName)
		if err != nil {
			return
		}
		if o.record == nil {
			o.record = pipeline.NewRecordPosition(pipeline.WithPipelineName(o.Options.PipelineName))
			//o.record = &pipeline.RecordPosition{
			//	PipelineName: o.Options.PipelineName,
			//	Pre:         &pipeline.Position{},
			//	Now:         &pipeline.Position{},
			//}
			*o.record.Pre = msg.Content.Head.Position
			*o.record.Now = msg.Content.Head.Position
			return
		}
	}
	if msg.Content.Head.Position.TotalRows == msg.Content.Head.Position.ConsumeRows {
		//o.record = &pipeline.RecordPosition{
		//	PipelineName: o.Options.PipelineName,
		//Pre:          &msg.Content.Head.Position,
		//Now:          &msg.Content.Head.Position,
		//}
		o.record = pipeline.NewRecordPosition(pipeline.WithPipelineName(o.Options.PipelineName))
		*o.record.Pre = msg.Content.Head.Position
		*o.record.Now = msg.Content.Head.Position
		return
	}
	if o.record.Pre == nil {
		//o.record.Pre = &msg.Content.Head.Position
		o.record.Pre = &pipeline.Position{}
		*o.record.Pre = msg.Content.Head.Position
	}
	if o.record.Now == nil {
		//o.record.Now = &msg.Content.Head.Position
		o.record.Now = &pipeline.Position{}
		*o.record.Now = msg.Content.Head.Position
		return
	}
	switch o.Options.MysqlMode {
	case pipeline.MODE_GTID:
		{
			if o.record.Now.GTIDSet != msg.Content.Head.Position.GTIDSet {
				*o.record.Now = msg.Content.Head.Position
				//o.record.Now = msg.Content.Head.Position
				return
			}
		}
	case pipeline.MODE_POSITION:
		{
			if o.record.Now.BinlogFile != msg.Content.Head.Position.BinlogFile ||
				o.record.Now.BinlogPosition != msg.Content.Head.Position.BinlogPosition {
				//o.record.Now = msg.Content.Head.Position
				*o.record.Now = msg.Content.Head.Position
				return
			}
		}
	}
	if o.record.Now.ConsumeRows < msg.Content.Head.Position.ConsumeRows {
		//o.record.Now = msg.Content.Head.Position
		*o.record.Now = msg.Content.Head.Position
	} else {
		pass = false
	}
	return
}

func (o *Output) syncRecord() (err error) {
	err = dao_pipe.UpdateRecord(o.record)
	return
}

func (o *Output) close() error {
	return o.Sender.Close()
}
