package output

import (
	message2 "github.com/jin06/binlogo/app/pipeline/message"
	sender2 "github.com/jin06/binlogo/app/pipeline/output/sender"
	kafka2 "github.com/jin06/binlogo/app/pipeline/output/sender/kafka"
	stdout2 "github.com/jin06/binlogo/app/pipeline/output/sender/stdout"
	store2 "github.com/jin06/binlogo/pkg/store"
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

func recordPosition(msg *message2.Message) (bool, error) {
	logrus.Debugf(
		"Record new replication position, file %s, pos %v",
		msg.BinlogPosition.BinlogFile,
		msg.BinlogPosition.BinlogPosition,
	)
	return store2.Update(msg.BinlogPosition)
}

func (o *Output) doHandle() {
	for {
		logrus.Debug("Wait send message")
		var msg *message2.Message
		msg = <-o.InChan
		logrus.Debugf("Output read message: %v \n", *msg)
		if !msg.Filter {
			ok, err := o.Sender.Send(msg)
			if err != nil {
				logrus.Error("Send message error: ", err)
				continue
			}
			if !ok {
				logrus.Error("Send message failed")
				continue
			}
		}
		ok, err := recordPosition(msg)
		if err != nil {
			logrus.Error(err)
		}
		if !ok {
			logrus.Error("Record position error")
		}
	}
}

func (o *Output) handle() {
	go o.doHandle()
	return
}

func (o *Output) Run() (err error) {
	o.handle()
	return
}
