package filter

import (
	"context"
	message2 "github.com/jin06/binlogo/app/pipeline/message"
	"github.com/sirupsen/logrus"
)

type Filter struct {
	InChan  chan *message2.Message
	OutChan chan *message2.Message
	Options *Options
	Ctx     context.Context
}

func New(opts ...Option) (filter *Filter, err error) {
	options := &Options{}
	for _, v := range opts {
		v(options)
	}
	filter = &Filter{
		Options: options,
	}
	err = filter.Init()
	return
}

func (f *Filter) Init() (err error) {
	return
}

type RuleType byte

var (
	RULE_TYPE_BLACK RuleType = 1
	RULE_TYPE_WHITE RuleType = 2
)

type Rule struct {
	Type     RuleType
	Database string
	Table    string
	MsgType  message2.MessageType
}

func (f *Filter) filer(msg *message2.Message) (err error) {
	if len(f.Options.Rules) > 0 {
		for _, v := range f.Options.Rules {
			switch v.Type {
			case RULE_TYPE_BLACK:
				{

				}
			case RULE_TYPE_WHITE:
				{

				}
			}
		}
	}
	return
}

func (f *Filter) doHandle(ctx context.Context) {
	for {
		logrus.Debug("Filter wait message")
		msg := <-f.InChan
		err := f.filer(msg)
		if err != nil {
			logrus.Error(err)
			continue
		}
		logrus.Debug("Filter message: ", msg)
		f.OutChan <- msg
	}
}

func (f *Filter) Handle() {
	go f.doHandle(f.Ctx)
	return
}

func (f *Filter) Run() (err error) {
	f.Handle()
	return
}