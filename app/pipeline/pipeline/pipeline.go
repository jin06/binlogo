package pipeline

import (
	"errors"
	filter2 "github.com/jin06/binlogo/app/pipeline/filter"
	input2 "github.com/jin06/binlogo/app/pipeline/input"
	message2 "github.com/jin06/binlogo/app/pipeline/message"
	output2 "github.com/jin06/binlogo/app/pipeline/output"
	store2 "github.com/jin06/binlogo/pkg/store"
	model2 "github.com/jin06/binlogo/pkg/store/model"
	"github.com/siddontang/go-log/log"
	"github.com/sirupsen/logrus"
)

type Pipeline struct {
	Input       *input2.Input
	Output      *output2.Output
	Filter      *filter2.Filter
	OutChan     *OutChan
	ControlLine *ControlLine
	Options     Options
	Role        Role
}

func NewFromStore(name string) (p *Pipeline, err error) {
	pipeModel := &model2.Pipeline{
		Name: name,
	}
	ok, err := store2.Get(pipeModel)
	if err != nil {
		return
	}
	if !ok {
		err = errors.New("No pipeline data ")
		return
	}
	posModel := &model2.Position{
		PipelineName: name,
	}
	ok, err = store2.Get(posModel)
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
	Input  chan *message2.Message
	Filter chan *message2.Message
	Out    chan *message2.Message
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
		Input:  make(chan *message2.Message, capacity),
		Filter: make(chan *message2.Message, capacity),
	}
}

func (p *Pipeline) initInput() (err error) {
	p.Input, err = input2.New(
		input2.OptionMysql(p.Options.Pipeline.Mysql),
		input2.OptionPosition(p.Options.Position),
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

