package filter

import (
	"context"
	message2 "github.com/jin06/binlogo/app/pipeline/message"
	"github.com/jin06/binlogo/pkg/pipeline/tool"
	pipeline2 "github.com/jin06/binlogo/pkg/store/model/pipeline"
	"github.com/sirupsen/logrus"
	"strings"
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
	f.rulesTree = tree{
		DBBlack:    map[string]bool{},
		TableBlack: map[string]bool{},
		DBWhite:    map[string]bool{},
		TableWhite: map[string]bool{},
	}
	if f.Options.Pipe != nil {
		for _, v := range f.Options.Pipe.Filters {
			if !tool.FilterVerify(v) {
				continue
			}
			arr := strings.Split(v.Rule, ".")
			if len(arr) == 1 {
				if v.Type == pipeline2.FILTER_BLACK {
					f.rulesTree.DBBlack[v.Rule] = true
				} else {
					f.rulesTree.DBWhite[v.Rule] = true
				}
			} else {
				if v.Type == pipeline2.FILTER_BLACK {
					f.rulesTree.TableBlack[v.Rule] = true
				} else {
					f.rulesTree.TableWhite[v.Rule] = true
				}
			}

		}
	}

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
