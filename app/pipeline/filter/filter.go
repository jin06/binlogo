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
	ctx     context.Context
}

func New(opts ...Option) (filter *Filter, err error) {
	options := &Options{}
	for _, v := range opts {
		v(options)
	}
	filter = &Filter{
		Options: options,
	}
	return
}

func (f *Filter) init() (err error) {
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

func (f *Filter) handle(msg *message2.Message) {
	logrus.Debug("Filter wait message")
	err := f.filer(msg)
	if err != nil {
		logrus.Error(err)
		return
	}
	logrus.Debug("Filter message: ", msg)
	//f.OutChan <- msg
}

func (f *Filter) Run(ctx context.Context) (err error) {
	err = f.init()
	if err != nil {
		return
	}
	myCtx, c := context.WithCancel(ctx)
	f.ctx = myCtx
	go func() {
		defer func() {
			c()
		}()
		for {
			select {
			case <-ctx.Done():
				{
					return
				}
				case msg := <-f.InChan:{
					f.handle(msg)
					f.OutChan <- msg
				}
			}
		}
	}()
	return
}

func (f *Filter) Context() context.Context {
	return f.ctx
}
