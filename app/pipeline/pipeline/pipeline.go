package pipeline

import (
	"context"
	"errors"
	filter2 "github.com/jin06/binlogo/app/pipeline/filter"
	input2 "github.com/jin06/binlogo/app/pipeline/input"
	message2 "github.com/jin06/binlogo/app/pipeline/message"
	output2 "github.com/jin06/binlogo/app/pipeline/output"
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
	status   status
	ctx      context.Context
}

type status byte

const (
	STATUS_RUN  status = 1
	STATUS_STOP status = 2
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

	if !p.Options.Pipeline.ExpectRun() {
		err = errors.New("pipeline not expect run")
		return
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
		input2.OptionsPipeName(p.Options.Pipeline.Name),
		input2.OptionsNodeName(p.Options.Pipeline.Name),
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

func (p *Pipeline) Run(ctx context.Context) {
	p.runMutex.Lock()
	defer p.runMutex.Unlock()
	if p.status == STATUS_RUN {
		return
	}
	//logrus.Debug("mysql position", p.Input.Options.Position)
	myCtx, cancel := context.WithCancel(ctx)
	p.ctx = myCtx
	go func() {
		var err error
		defer func() {
			cancel()
		}()
		if err = p.Input.Run(myCtx); err != nil {
			return
		}
		if err = p.Filter.Run(myCtx); err != nil {
			return
		}
		if err = p.Output.Run(myCtx); err != nil {
			return
		}
		for {
			select {
			case <-ctx.Done():
				{
					return
				}
			case <-p.Input.Context().Done():
				{
					return
				}
			case <-p.Filter.Context().Done():
				{
					return
				}
			case <-p.Output.Context().Done():
				{
					return
				}
			}
		}
	}()
}

func (p *Pipeline) Context() context.Context {
	return p.ctx
}
