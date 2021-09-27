package output

import (
	"github.com/jin06/binlogo/pipeline/message"
	"github.com/jin06/binlogo/pipeline/output/sender"
	"github.com/jin06/binlogo/pipeline/output/sender/kafka"
	"github.com/jin06/binlogo/store"
	"github.com/jin06/binlogo/store/model"
	"github.com/sirupsen/logrus"
)

type Output struct {
	InChan  chan *message.Message
	Sender  sender.Sender
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
	case model.SENDER_TYPE_KAFKA:
		fallthrough
	default:
		o.Sender, err = kafka.New(
			&kafka.Options{
				Kafka: o.Options.Output.Sender.Kafka,
			},
		)
	}

	return
}

func recordPosition(msg *message.Message) (bool, error) {
	return store.Update(msg.BinlogPosition)
}

func (o *Output) doHandle() {
	for {
		logrus.Debug("Wait send message")
		var msg *message.Message
		msg = <-o.InChan
		logrus.Debug("Output read message: ", msg)
		if msg.Filter {
			continue
		}
		ok, err := o.Sender.Send(msg)
		if err != nil {
			logrus.Error("Send message error: ", err)
			continue
		}
		if !ok {
			logrus.Error("Send message failed")
			continue
		}
		ok, err = recordPosition(msg)
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
