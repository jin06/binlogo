package core

import (
	"os"
)

type Exporter interface {
	Start(chan *Event) error
	Config()
}

type Stand struct{}

func (s *Stand) Start(ch chan *Event) error {
	go func() {
		for {
			select {
			case event := <-ch:
				event.BinlogEvent.Dump(os.Stdout)
			}
		}
	}()
	return nil
}

func (s *Stand) Config() {
	return
}
