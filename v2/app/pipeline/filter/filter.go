package filter

import (
	"context"
	"sync"

	"github.com/jin06/binlogo/app/pipeline/message"
	"github.com/jin06/binlogo/configs"
	"github.com/jin06/binlogo/pkg/promeths"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

// Filter filter message by rules
type Filter struct {
	InChan    chan *message.Message
	OutChan   chan *message.Message
	Options   *Options
	rulesTree tree
	closed    chan struct{}
	closeOnce sync.Once
}

// New returns a new Filter
func New(opts ...Option) (filter *Filter, err error) {
	options := &Options{}
	for _, v := range opts {
		v(options)
	}
	filter = &Filter{
		Options: options,
		closed:  make(chan struct{}),
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

func (f *Filter) filer(msg *message.Message) (err error) {
	msg.Filter = f.rulesTree.isFilter(msg)
	return
}

func (f *Filter) handle(msg *message.Message) {
	err := f.filer(msg)
	if err != nil {
		logrus.Error(err)
		return
	}
	if msg.Filter {
		promeths.MessageFilterCounter.With(prometheus.Labels{"pipeline": f.Options.Pipe.Name, "node": configs.NodeName}).Inc()
	}
}

// Run Filter start working
func (f *Filter) Run(ctx context.Context) (err error) {
	defer f.close()
	if err = f.init(); err != nil {
		return
	}
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-f.InChan:
			{
				f.handle(msg)
				f.OutChan <- msg
			}
		}
	}
}

func (f *Filter) Closed() chan struct{} {
	return f.closed
}

func (f *Filter) close() {
	f.closeOnce.Do(func() {
		close(f.closed)
	})
}
