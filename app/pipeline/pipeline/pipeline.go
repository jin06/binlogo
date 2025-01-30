package pipeline

import (
	"context"
	"errors"
	"sync"

	"github.com/jin06/binlogo/v2/app/pipeline/filter"
	"github.com/jin06/binlogo/v2/app/pipeline/input"
	"github.com/jin06/binlogo/v2/app/pipeline/message"
	"github.com/jin06/binlogo/v2/app/pipeline/output"
	"github.com/jin06/binlogo/v2/pkg/event"
	event_store "github.com/jin06/binlogo/v2/pkg/store/model"
)

type status byte

const (
	STATUS_RUN  status = 1
	STATUS_STOP status = 2
)

// OutChan  is used to deliver messages in different components
type outChan struct {
	input  chan *message.Message
	filter chan *message.Message
	out    chan *message.Message
}

func (o *outChan) Close() {
	if o.input != nil {
		close(o.input)
	}
	if o.filter != nil {
		close(o.filter)
	}
	if o.out != nil {
		close(o.out)
	}
}

// New returns a new pipeline
func New(opt ...Option) (p *Pipeline, err error) {
	options := Options{}
	for _, v := range opt {
		v(&options)
	}
	p = &Pipeline{
		Options:  options,
		runMutex: sync.Mutex{},
		closed:   make(chan struct{}),
		closing:  make(chan struct{}),
	}
	// p.log = logrus.WithField("mod", "pipeline").WithField("name", options.Pipeline.Name)
	err = p.init()
	return
}

// Pipeline for handle message
// contains input, filter, output
type Pipeline struct {
	Input  *input.Input
	Output *output.Output
	Filter *filter.Filter

	Options  Options
	outChan  *outChan
	runMutex sync.Mutex
	// log      *logrus.Entry

	closed       chan struct{}
	closeOnce    sync.Once
	closing      chan struct{}
	completeOnce sync.Once

	isClosed bool
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
	p.outChan = &outChan{
		input:  make(chan *message.Message, capacity),
		filter: make(chan *message.Message, capacity),
	}
}

func (p *Pipeline) initInput() (err error) {
	p.Input, err = input.New(
		input.OptionsPipeName(p.Options.Pipeline.Name),
		input.OptionsNodeName(p.Options.Pipeline.Name),
	)
	p.Input.OutChan = p.outChan.input
	return
}

func (p *Pipeline) initFilter() (err error) {
	p.Filter, err = filter.New(filter.WithPipe(p.Options.Pipeline))
	p.Filter.InChan = p.outChan.input
	p.Filter.OutChan = p.outChan.filter
	return
}

func (p *Pipeline) initOutput() (err error) {
	p.Output, err = output.New(
		output.OptionOutput(p.Options.Pipeline.Output),
		output.OptionPipeName(p.Options.Pipeline.Name),
		output.OptionMysqlMode(p.Options.Pipeline.Mysql.Mode),
	)
	p.Output.InChan = p.outChan.filter
	return
}

// Run pipeline start working
func (p *Pipeline) Run(ctx context.Context) {
	defer p.CompleteClose()
	defer p.Close()
	defer p.Options.Log.Infof("Pipeline stream stopped: %s", p.Options.Pipeline.Name)
	p.Options.Log.Infof("Pipeline stream run: %s", p.Options.Pipeline.Name)

	//logrus.Debug("mysql position", p.Input.Options.Position)
	go func() {
		if err := p.Input.Run(ctx); err != nil {
			event.Event(event_store.NewErrorPipeline(p.Options.Pipeline.Name, "Start error: "+err.Error()))
			p.Close()
			return
		}
	}()
	go func() {
		if err := p.Filter.Run(ctx); err != nil {
			event.Event(event_store.NewErrorPipeline(p.Options.Pipeline.Name, "Start error: "+err.Error()))
			p.Close()
			return
		}
	}()
	go func() {
		if err := p.Output.Run(ctx); err != nil {
			event.Event(event_store.NewErrorPipeline(p.Options.Pipeline.Name, "Start error: "+err.Error()))
			p.Close()
			return
		}
	}()

	event.Event(event_store.NewInfoPipeline(p.Options.Pipeline.Name, "Start succeeded"))
	for {
		select {
		case <-ctx.Done():
			return
		case <-p.Input.Closed():
			p.Options.Log.Info("input closed")
			return
		case <-p.Filter.Closed():
			p.Options.Log.Info("filter closed")
			return
		case <-p.Output.Closed():
			p.Options.Log.Info("output closed")
			return
		case <-p.closing:
			return
		}
	}
}

func (p *Pipeline) Close() (err error) {
	p.closeOnce.Do(func() {
		p.Input.Close()
		p.Filter.Close()
		p.Output.Close()
		if p.outChan != nil {
			p.outChan.Close()
		}
		close(p.closing)
	})
	return
}

func (p *Pipeline) CompleteClose() {
	p.completeOnce.Do(func() {
		<-p.Input.Closed()
		<-p.Filter.Closed()
		<-p.Output.Closed()
		p.isClosed = true
		close(p.closed)
	})
}

func (p *Pipeline) IsClosed() bool {
	return p.isClosed
}

func (p *Pipeline) Closed() chan struct{} {
	return p.closed
}
