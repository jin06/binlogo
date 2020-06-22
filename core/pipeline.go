package core

import "context"

type Pipeline struct {
	Importer *Importer
	Exporter Exporter
	Filters  []Filter
	Context  context.Context
	SourceLine chan *Event
	FilterLine chan *Event
}

func NewPipeline(context context.Context) *Pipeline {
	return &Pipeline{Context: context}
}

func newPipeline() *Pipeline {
	p := new(Pipeline)
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
	ch, err := p.Importer.Start()
	if err != nil {
		panic(err)
	}
	p.SourceLine = ch
	err =  p.Exporter.Start(p.SourceLine)
	return err
}

