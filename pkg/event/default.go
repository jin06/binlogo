package event

import (
	"context"
	"github.com/jin06/binlogo/pkg/store/model/event"
)

// DefaultRecorder global default event recorder
var DefaultRecorder *Recorder

// Init generate a recorder for global use
func Init() {
	DefaultRecorder, _ = New()
	DefaultRecorder.Loop(context.Background())
}

// Event client call this function to record a event
func Event(e *event.Event) {
	//fmt.Println(e)
	if DefaultRecorder != nil {
		go DefaultRecorder.Event(e)
	}
}
