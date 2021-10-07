package filter

type Options struct {
	Rules []Rule
}

type Option func(options *Options)

func OptionsRules(rules []Rule) Option {
	return func(options *Options) {
		options.Rules = rules
	}
}

