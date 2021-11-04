package output

import (
	"context"
	message2 "github.com/jin06/binlogo/app/pipeline/message"
	sender2 "github.com/jin06/binlogo/app/pipeline/output/sender"
	kafka2 "github.com/jin06/binlogo/app/pipeline/output/sender/kafka"
	stdout2 "github.com/jin06/binlogo/app/pipeline/output/sender/stdout"
	"github.com/jin06/binlogo/pkg/store/dao/dao_pipe"
	"github.com/jin06/binlogo/pkg/store/model/pipeline"
	"github.com/sirupsen/logrus"
)

type Output struct {
	InChan  chan *message2.Message
	Sender  sender2.Sender
	Options *Options
}

func New(opts ...Option) (out *Output, err error) {
	options := &Options{}
	for _, v := range opts {
		v(options)
	}
	out = &Output{
		Options: options,
	}
	err = out.Init()
	return
}

func (o *Output) Init() (err error) {
	switch o.Options.Output.Sender.Type {
	case pipeline.SNEDER_TYPE_STDOUT:
		o.Sender, err = stdout2.New()
	case pipeline.SENDER_TYPE_KAFKA:
		fallthrough
	default:
		o.Sender, err = kafka2.New(
			&kafka2.Options{
				Kafka: o.Options.Output.Sender.Kafka,
			},
		)
	}

	return
}

func recordPosition(msg *message2.Message) (err error) {
	logrus.Debugf(
		"Record new replication position, file %s, pos %v",
		msg.BinlogPosition.BinlogFile,
		msg.BinlogPosition.BinlogPosition,
	)
	err = dao_pipe.UpdatePosition(msg.BinlogPosition)
	return
}

func (o *Output) doHandle() {
	logrus.Debug("Wait send message")
	var msg *message2.Message
	msg = <-o.InChan
	logrus.Debugf("Output read message: %v \n", *msg)
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
	err := recordPosition(msg)
	if err != nil {
		logrus.Error(err)
	}
}

func (o *Output) handle(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				{
					return
				}
			default:
				o.doHandle()
			}
		}
	}()
	return
}

func (o *Output) Run(ctx context.Context) (err error) {
	o.handle(ctx)
	return
}
