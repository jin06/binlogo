package core

import "fmt"

type Filter struct {
	Line chan *Event
}

func (f *Filter) Start(in chan *Event) (out chan *Event, err error) {
	f.Line = out
	out = make(chan *Event, 10000)
	go func() {
		for {
			select {
			case event := <-in:
				fmt.Println("filter read a event", event.BinlogEvent.Header.EventType.String())
				out <- event
			}
		}
	}()
	return
}
