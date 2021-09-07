package pipeline

type Pipeline struct {
	Input   Input
	Output  Output
	Filters []Filter
	Options
}

type Options struct {
	ID string
}

type Option func(*Options)

func OptionID(id string) Option {
	return func(ops *Options) {
		ops.ID = id
	}
}

func NewPipeline(opt ...Option) (p *Pipeline) {
	options := Options{}
	for _, v := range opt {
		v(&options)
	}
	p = &Pipeline{
		Options: options,
	}
	return
}
