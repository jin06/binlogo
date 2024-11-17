package pipeline

import (
	"context"
	"errors"
	"sync"

	"github.com/sirupsen/logrus"

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

func (o *outChan) close() {
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
		status:   STATUS_STOP,
		closed:   make(chan struct{}),
	}
	p.log = logrus.WithField("mod", "pipeline").WithField("name", options.Pipeline.Name)
	err = p.init()
	return
}

// Pipeline for handle message
// contains input, filter, output
type Pipeline struct {
	Input     *input.Input
	Output    *output.Output
	Filter    *filter.Filter
	Options   Options
	outChan   *outChan
	runMutex  sync.Mutex
	status    status
	closed    chan struct{}
	closeOnce sync.Once
	log       *logrus.Entry
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
	defer p.close()
	defer logrus.Info("pipeline stream stopped: ", p.Options.Pipeline.Name)
	logrus.Info("pipeline stream run: ", p.Options.Pipeline.Name)

	//logrus.Debug("mysql position", p.Input.Options.Position)
	stx, cancel := context.WithCancel(ctx)
	defer cancel()
	go func() {
		if err := p.Input.Run(stx); err != nil {
			event.Event(event_store.NewErrorPipeline(p.Options.Pipeline.Name, "Start error: "+err.Error()))
			return
		}
	}()
	go func() {
		if err := p.Filter.Run(stx); err != nil {
			event.Event(event_store.NewErrorPipeline(p.Options.Pipeline.Name, "Start error: "+err.Error()))
			return
		}
	}()
	go func() {
		if err := p.Output.Run(stx); err != nil {
			event.Event(event_store.NewErrorPipeline(p.Options.Pipeline.Name, "Start error: "+err.Error()))
			return
		}
	}()

	event.Event(event_store.NewInfoPipeline(p.Options.Pipeline.Name, "Start succeeded"))
	for {
		select {
		case <-ctx.Done():
			return
		case <-p.Input.Closed():
			logrus.Info("input closed")
			return
		case <-p.Filter.Closed():
			logrus.Info("filter closed")
			return
		case <-p.Output.Closed():
			logrus.Info("output closed")
			return
		}
	}
}

func (p *Pipeline) close() (err error) {
	p.closeOnce.Do(func() {
		<-p.Input.Closed()
		<-p.Filter.Closed()
		<-p.Output.Closed()
		if p.outChan != nil {
			p.outChan.close()
		}
		close(p.closed)
	})
	return
}
