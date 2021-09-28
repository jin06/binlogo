package pipeline

import (
	"errors"
	"github.com/jin06/binlogo/pipeline/filter"
	"github.com/jin06/binlogo/pipeline/input"
	"github.com/jin06/binlogo/pipeline/message"
	"github.com/jin06/binlogo/pipeline/output"
	"github.com/jin06/binlogo/store"
	"github.com/jin06/binlogo/store/model"
	"github.com/siddontang/go-log/log"
	"github.com/sirupsen/logrus"
)

type Pipeline struct {
	Input       *input.Input
	Output      *output.Output
	Filter      *filter.Filter
	OutChan     *OutChan
	ControlLine *ControlLine
	Options     Options
	Role        Role
}

func NewFromStore(name string) (p *Pipeline, err error) {
	pipeModel := &model.Pipeline{
		Name: name,
	}
	ok, err := store.Get(pipeModel)
	if err != nil {
		return
	}
	if !ok {
		err = errors.New("No pipeline data ")
		return
	}
	posModel := &model.Position{
		PipelineName: name,
	}
	ok, err = store.Get(posModel)
	if err != nil {
		return
	}
	if !ok {
		//errors.New("No position data ")
		log.Warn("No position data")
		//return
	}
	logrus.Debug("Position data: ", posModel)

	return New(
		OptionPipeline(pipeModel),
		OptionPosition(posModel),
	)
}

func New(opt ...Option) (p *Pipeline, err error) {
	options := Options{}
	for _, v := range opt {
		v(&options)
	}
	p = &Pipeline{
		Options: options,
	}
	err = p.Init()
	return
}

type OutChan struct {
	Input  chan *message.Message
	Filter chan *message.Message
	Out    chan *message.Message
}

type ControlLine struct {
}

func (p *Pipeline) Init() (err error) {
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
		Input:  make(chan *message.Message, capacity),
		Filter: make(chan *message.Message, capacity),
	}
}

func (p *Pipeline) initInput() (err error) {
	p.Input, err = input.New(
		input.OptionMysql(p.Options.Pipeline.Mysql),
		input.OptionPosition(p.Options.Position),
	)
	p.Input.OutChan = p.OutChan.Input
	return
}

func (p *Pipeline) initFilter() (err error) {
	p.Filter, err = filter.New()
	p.Filter.InChan = p.OutChan.Input
	p.Filter.OutChan = p.OutChan.Filter
	return
}

func (p *Pipeline) initOutput() (err error) {
	p.Output, err = output.New(output.OptionOutput(p.Options.Pipeline.Output))
	p.Output.InChan = p.OutChan.Filter
	return
}

func (p *Pipeline) Run() (err error) {
	logrus.Debug("mysql position", p.Input.Options.Position)
	if err = p.Input.Run(); err != nil {
		return
	}
	if err = p.Filter.Run(); err != nil {
		return
	}
	if err = p.Output.Run(); err != nil {
		return
	}
	return
}

func (p *Pipeline) IsLeader() bool {
	return p.Role == RoleLeader
}

func (p *Pipeline) IsFollower() bool {
	return p.Role == RoleFollower
}

