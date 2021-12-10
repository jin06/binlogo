package filter

import (
	"context"

	message2 "github.com/jin06/binlogo/app/pipeline/message"
	"github.com/sirupsen/logrus"
)

// Filter filter message by rules
type Filter struct {
	InChan    chan *message2.Message
	OutChan   chan *message2.Message
	Options   *Options
	ctx       context.Context
	rulesTree tree
}

// New returns a new Filter
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
	if f.Options.Pipe == nil {
		f.rulesTree = newTree(nil)
	}
	f.rulesTree = newTree(f.Options.Pipe.Filters)
	return
}

func (f *Filter) filer(msg *message2.Message) (err error) {
	msg.Filter = f.rulesTree.isFilter(msg)
	return
}

func (f *Filter) handle(msg *message2.Message) {
	err := f.filer(msg)
	if err != nil {
		logrus.Error(err)
		return
	}
}

// Run Filter start working
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
			case msg := <-f.InChan:
				{
					f.handle(msg)
					f.OutChan <- msg
				}
			}
		}
	}()
	return
}

// Context return Filter's context for cancel goroutine
func (f *Filter) Context() context.Context {
	return f.ctx
}
