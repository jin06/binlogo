package output

import (
	"github.com/jin06/binlogo/pipeline/message"
	"github.com/jin06/binlogo/pipeline/output/sender"
	"github.com/jin06/binlogo/pipeline/output/sender/kafka"
	"github.com/jin06/binlogo/store/model"
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

func (o *Output) doHandle() {
	o.Sender.Send(o.InChan)
}

func (o *Output) handle() {
	go o.doHandle()
	return
}

func (o *Output) Run() (err error) {
	o.handle()
	return
}
