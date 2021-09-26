package output

import (
	"fmt"
	"github.com/jin06/binlogo/pipeline/message"
	"github.com/jin06/binlogo/pipeline/output/sender"
	"github.com/sirupsen/logrus"
)

type Output struct {
	InChan  chan *message.Message
	Sender  sender.Sender
	Options *Options
}

func New(opts ...Option) (out *Output,err error) {
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
	return
}

func (o *Output) doHandle() {
	for {
		logrus.Debug("Output wait message ")
		msg := <-o.InChan
		fmt.Println("Output handle message : ")
		fmt.Println(msg.Json())
	}
}

func (o *Output) handle() {
	go o.doHandle()
	return
}

func (o *Output) Start() (err error) {
	o.handle()
	return
}
