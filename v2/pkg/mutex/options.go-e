package mutex

import "time"

type Option func(*Mutex)

func WithTimeout(t time.Duration) Option {
	return func(m *Mutex) {
		m.timeout = t
	}
}
