package input

type Options struct {
	PipeName string
	NodeName string
}

type Option func(options *Options)

func OptionsPipeName(name string) Option {
	return func(options *Options) {
		options.PipeName = name
	}
}

func OptionsNodeName(name string) Option {
	return func(options *Options) {
		options.NodeName = name
	}
}
