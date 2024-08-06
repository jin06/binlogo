package input

// Options for Input Options
type Options struct {
	PipeName string
	NodeName string
}

// Option is function of configure Options
type Option func(options *Options)

// OptionsPipeName sets Options pipeline name
func OptionsPipeName(name string) Option {
	return func(options *Options) {
		options.PipeName = name
	}
}

// OptionsNodeName sets Options pipeline node name
func OptionsNodeName(name string) Option {
	return func(options *Options) {
		options.NodeName = name
	}
}
