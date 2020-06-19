package core

import "context"

type Pipeline struct {
	Importer Importer
	Exporter Exporter
	Filters  []Filter
	line     chan *Event
	Context  context.Context
}

func NewPipeline(context context.Context) *Pipeline {
	return &Pipeline{Context: context}
}

func newPipeline() *Pipeline {
	p := new(Pipeline)
	p.line = make(chan *Event, 10000)
	return p
}

func (p *Pipeline) check() error {
	return nil
}

func (p *Pipeline) ChangeIm(importer *Importer) error {
	return nil
}

func (p *Pipeline) ChangeEx(exporter *Exporter) error {
	return nil
}

func (p *Pipeline) Run() error {
	return nil
}
