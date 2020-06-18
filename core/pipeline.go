package core

type Pipeline struct {
	Importer Importer
	Exporter Exporter
	Filters  []Filter
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
