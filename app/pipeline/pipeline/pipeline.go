package pipeline

import (
	"context"
	"errors"
	filter2 "github.com/jin06/binlogo/app/pipeline/filter"
	input2 "github.com/jin06/binlogo/app/pipeline/input"
	message2 "github.com/jin06/binlogo/app/pipeline/message"
	output2 "github.com/jin06/binlogo/app/pipeline/output"
	"github.com/jin06/binlogo/pkg/blog"
	"github.com/sirupsen/logrus"
	"sync"
)

type Pipeline struct {
	Input    *input2.Input
	Output   *output2.Output
	Filter   *filter2.Filter
	OutChan  *OutChan
	Options  Options
	Role     Role
	cancel   context.CancelFunc
	runMutex sync.Mutex
	status   string
}

const (
	STATUS_RUN  = "running"
	STATUS_STOP = "stopped"
	STATUS_PENDING = "pending"
)

func New(opt ...Option) (p *Pipeline, err error) {
	options := Options{}
	for _, v := range opt {
		v(&options)
	}
	p = &Pipeline{
		Options:  options,
		runMutex: sync.Mutex{},
		status:   STATUS_STOP,
	}
	err = p.init()
	return
}

type OutChan struct {
	Input  chan *message2.Message
	Filter chan *message2.Message
	Out    chan *message2.Message
}

func (p *Pipeline) init() (err error) {
	if p.Options.Pipeline == nil {
		return errors.New("lack pipeline info")
	}

	p.initDataLine()
	if err = p.initInput(); err != nil {
		return
	}
	if err = p.initFilter(); err != nil {
		return
	}
	if err = p.initOutput(); err != nil {
		return
	}
	return
}

func (p *Pipeline) initDataLine() {
	capacity := 100000
	p.OutChan = &OutChan{
		Input:  make(chan *message2.Message, capacity),
		Filter: make(chan *message2.Message, capacity),
	}
}

func (p *Pipeline) initInput() (err error) {
	p.Input, err = input2.New(
		input2.OptionMysql(p.Options.Pipeline.Mysql),
		input2.OptionPosition(p.Options.Position),
		input2.OptionsPipeline(p.Options.Pipeline),
	)
	p.Input.OutChan = p.OutChan.Input
	return
}

func (p *Pipeline) initFilter() (err error) {
	p.Filter, err = filter2.New()
	p.Filter.InChan = p.OutChan.Input
	p.Filter.OutChan = p.OutChan.Filter
	return
}

func (p *Pipeline) initOutput() (err error) {
	p.Output, err = output2.New(output2.OptionOutput(p.Options.Pipeline.Output))
	p.Output.InChan = p.OutChan.Filter
	return
}

func (p *Pipeline) Run(ctx context.Context) (err error) {
	p.runMutex.Lock()
	defer p.runMutex.Unlock()
	if p.status == STATUS_RUN {
		return
	}
	logrus.Debug("mysql position", p.Input.Options.Position)
	sctx, cancel := context.WithCancel(ctx)
	p.cancel = cancel
	defer func() {
		if err != nil {
			blog.Info(err)
			p.cancel()
		} else {
			p.status = STATUS_RUN
		}
	}()
	if err = p.Input.Run(sctx); err != nil {
		return
	}
	if err = p.Filter.Run(sctx); err != nil {
		return
	}
	if err = p.Output.Run(sctx); err != nil {
		return
	}
	return
}

func (p *Pipeline) Stop() {
	p.runMutex.Lock()
	defer p.runMutex.Unlock()
	if p.status == STATUS_STOP {
		return
	}
	p.cancel()
	p.status = STATUS_STOP
}
