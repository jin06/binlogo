package core

import (
	"fmt"
	"os"
)

type Exporter interface {
	Start(chan *Event) error
	Config()
	//Confirm() error
}

type Stand struct{}

func (s *Stand) Start(ch chan *Event) error {
	go func() {
		for {
			select {
			case event := <-ch:
				event.BinlogEvent.Dump(os.Stdout)
				event.BinlogEvent.Event.Dump(os.Stdout)
				fmt.Println(event.BinlogEvent.Header.EventType.String())
			}
		}
	}()
	return nil
}

func (s *Stand) Config() {
	return
}

func (s *Stand) Confirm() {

}
