package output

import (
	"context"
	message2 "github.com/jin06/binlogo/app/pipeline/message"
	sender2 "github.com/jin06/binlogo/app/pipeline/output/sender"
	"github.com/jin06/binlogo/app/pipeline/output/sender/http"
	kafka2 "github.com/jin06/binlogo/app/pipeline/output/sender/kafka"
	stdout2 "github.com/jin06/binlogo/app/pipeline/output/sender/stdout"
	"github.com/jin06/binlogo/pkg/store/dao/dao_pipe"
	"github.com/jin06/binlogo/pkg/store/model/pipeline"
	"github.com/sirupsen/logrus"
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
	case pipeline.SENDER_TYPE_KAFKA:
		fallthrough
	default:
		o.Sender, err = kafka2.New(
			o.Options.Output.Sender.Kafka,
		)
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

func (o *Output) handle(msg *message2.Message) {
	defer func() {
		err2 := recordPosition(msg)
		if err2 != nil {
			logrus.Error("Record position error, ", err2)
		}
	}()

	if !msg.Filter {
		ok, err := o.Sender.Send(msg)
		if err != nil {
			logrus.Error("Send message error: ", err)
			return
		}
		if !ok {
			logrus.Error("Send message failed")
			return
		}
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
			case msg := <- o.InChan:
				{
					o.handle(msg)
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
