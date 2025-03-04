package output

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/jin06/binlogo/v2/app/pipeline/output/sender/elastic"
	"github.com/sirupsen/logrus"

	"github.com/jin06/binlogo/v2/app/pipeline/message"
	"github.com/jin06/binlogo/v2/app/pipeline/output/sender"
	"github.com/jin06/binlogo/v2/app/pipeline/output/sender/http"
	"github.com/jin06/binlogo/v2/app/pipeline/output/sender/kafka"
	"github.com/jin06/binlogo/v2/app/pipeline/output/sender/rabbitmq"
	"github.com/jin06/binlogo/v2/app/pipeline/output/sender/redis"
	"github.com/jin06/binlogo/v2/app/pipeline/output/sender/rocketmq"
	"github.com/jin06/binlogo/v2/app/pipeline/output/sender/stdout"
	"github.com/jin06/binlogo/v2/configs"
	"github.com/jin06/binlogo/v2/pkg/event"
	"github.com/jin06/binlogo/v2/pkg/promeths"
	"github.com/jin06/binlogo/v2/pkg/store/dao"
	event_store "github.com/jin06/binlogo/v2/pkg/store/model"
	"github.com/jin06/binlogo/v2/pkg/store/model/pipeline"
	"github.com/prometheus/client_golang/prometheus"
)

// Output handle message output
// depends pipeline config, send message to stdout、kafka, etc.
type Output struct {
	InChan  chan *message.Message
	Sender  sender.Sender
	Options *Options
	record  *pipeline.RecordPosition

	closed       chan struct{}
	completeOnce sync.Once
	closing      chan struct{}
	closeOnce    sync.Once

	recordMutex  sync.Mutex
	recordSynced bool
}

// New return a Output object
func New(opts ...Option) (out *Output, err error) {
	options := &Options{}
	for _, v := range opts {
		v(options)
	}
	out = &Output{
		Options:      options,
		closed:       make(chan struct{}),
		closing:      make(chan struct{}),
		recordSynced: false,
	}
	return
}

func (o *Output) init() (err error) {
	switch o.Options.Output.Sender.Type {
	case pipeline.SNEDER_TYPE_STDOUT:
		o.Sender, err = stdout.New()
	case pipeline.SNEDER_TYPE_HTTP:
		o.Sender, err = http.New(o.Options.Output.Sender.Http)
	case pipeline.SNEDER_TYPE_RABBITMQ:
		o.Sender, err = rabbitmq.New(o.Options.Output.Sender.RabbitMQ)
	case pipeline.SENDER_TYPE_REDIS:
		o.Sender, err = redis.New(o.Options.Output.Sender.Redis)
	case pipeline.SENDER_TYPE_KAFKA:
		o.Sender, err = kafka.New(o.Options.Output.Sender.Kafka)
	case pipeline.SENDER_TYPE_ROCKETMQ:
		o.Sender, err = rocketmq.New(o.Options.Output.Sender.RocketMQ)
	case pipeline.SENDER_TYPE_Elastic:
		o.Sender, err = elastic.New(o.Options.Output.Sender.Elastic)
	default:
		o.Sender, err = stdout.New()
	}

	return
}

func (o *Output) loopHandle(ctx context.Context, msg *message.Message) (err error) {
	for i := 0; i < 3; i++ {
		if err = o.handle(msg); err == nil {
			return
		}
		promeths.MessageSendErrCounter.With(prometheus.Labels{"pipeline": o.Options.PipelineName, "node": configs.GetNodeName()}).Inc()
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		select {
		case <-ctx.Done():
			{
				return
			}
		case <-ticker.C:
			{
				continue
			}
		}
	}
	return
}

func (o *Output) handle(msg *message.Message) (err error) {
	defer func() {
		if err != nil {
			event.Event(event_store.NewErrorPipeline(o.Options.PipelineName, "send message error: "+err.Error()))
		}
	}()
	var ok bool
	if !msg.Filter {
		ok, err = o.Sender.Send(msg)
		if err != nil {
			return
		}
		if !ok {
			err = errors.New("send message failed")
		}
	}
	if err != nil {
		return
	}
	o.updateRecord()
	return
	// switch msg.Status {
	// case message.STATUS_RECORD:
	// 	{
	// 		return
	// 	}
	// case message.STATUS_SEND:
	// 	{
	// 		err = o.syncRecord()
	// 		if err == nil {
	// 			msg.Status = message.STATUS_RECORD
	// 		}
	// 		return
	// 	}
	// default:
	// 	if !msg.Filter {
	// 		var ok bool
	// 		ok, err = o.Sender.Send(msg)
	// 		if err == nil {
	// 			if !ok {
	// 				err = errors.New("send message failed")
	// 			}
	// 		}
	// 		if err == nil {
	// 			msg.Status = message.STATUS_SEND
	// 		}
	// 	}
	// 	if err == nil {
	// 		err = o.syncRecord()
	// 		if err == nil {
	// 			msg.Status = message.STATUS_RECORD
	// 		}
	// 	}
	// 	return
	// }
}

// Run start Output to send message
func (o *Output) Run(ctx context.Context) (err error) {
	defer o.CompleteClose()
	defer o.Close()
	if err = o.init(); err != nil {
		return
	}
	go o.loopSync(ctx)
	for {
		select {
		case <-o.closing:
			return
		case <-ctx.Done():
			{
				return
			}
		case msg := <-o.InChan:
			{
				check, errPrepare := o.prepareRecord(ctx, msg)
				if errPrepare != nil {
					message.Put(msg)
					return
				}
				if !check {
					message.Put(msg)
					continue
				}
				if loopErr := o.loopHandle(ctx, msg); loopErr != nil {
					message.Put(msg)
					return
				}
				promeths.MessageSendCounter.With(prometheus.Labels{"pipeline": o.Options.PipelineName, "node": configs.GetNodeName()}).Inc()
				pass := uint32(time.Now().Unix()) - msg.Content.Head.Time
				promeths.MessageSendHistogram.With(prometheus.Labels{"pipeline": o.Options.PipelineName, "node": configs.GetNodeName()}).Observe(float64(pass))
				message.Put(msg)
			}
		}
	}
}

func (o *Output) prepareRecord(ctx context.Context, msg *message.Message) (pass bool, err error) {
	o.recordMutex.Lock()
	defer o.recordMutex.Unlock()
	pass = true
	if o.record == nil {
		o.record, err = dao.GetRecord(ctx, o.Options.PipelineName)
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

func (o *Output) loopSync(c context.Context) {
	defer o.Close()
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-o.closing:
			return
		case <-c.Done():
			{
				return
			}
		case <-ticker.C:
			{
				o.syncRecord(c)
			}
		}
	}
}

func (o *Output) updateRecord() {
	o.recordMutex.Lock()
	defer o.recordMutex.Unlock()
	o.recordSynced = true
}

func (o *Output) syncRecord(ctx context.Context) (err error) {
	o.recordMutex.Lock()
	defer o.recordMutex.Unlock()
	if o.recordSynced {
		err = dao.UpdateRecord(ctx, o.record)
	}
	return
}

func (o *Output) Closed() chan struct{} {
	return o.closed
}

func (o *Output) Close() {
	o.closeOnce.Do(func() {
		if o.Sender != nil {
			o.Sender.Close()
		}
		close(o.closing)
	})
}

func (o *Output) CompleteClose() {
	o.completeOnce.Do(func() {
		logrus.WithField("Pipeline Name", o.Options.PipelineName).Debug("Output closed")
		close(o.closed)
	})
}
