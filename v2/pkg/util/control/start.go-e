package control

import "time"

// Starter is control goroutine start with wait Duration, use in loop function
type Starter struct {
	key      string
	time     *time.Time
	duration time.Duration
	maxWait  time.Duration
	minWait  time.Duration
}

// New returns a Starter
func New(key string) *Starter {
	s := &Starter{
		key:      key,
		time:     nil,
		duration: 0,
		maxWait:  time.Second * 10,
		minWait:  0,
	}
	return s
}

// SetMaxWait set Starter max wait time
func (s *Starter) SetMaxWait(t time.Duration) {
	s.maxWait = t
}

// SetMinWait set Starter min wait time
func (s *Starter) SetMinWait(t time.Duration) {
	s.minWait = t
}

// SetTime set Starter time
func (s *Starter) SetTime(t *time.Time) {
	if t == nil {
		now := time.Now()
		s.time = &now
	} else {
		s.time = t
	}
}

// Wait returns wait duration
func (s *Starter) Wait() time.Duration {
	return s.duration
}
