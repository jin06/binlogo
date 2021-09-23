package output

import (
	"fmt"
	"github.com/jin06/binlogo/pipeline/message"
	"github.com/sirupsen/logrus"
)

type Output struct {
	InChan chan *message.Message
}

func New() (*Output, error) {
	out := &Output{}
	err := out.Init()
	return out, err
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
