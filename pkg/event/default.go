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

// EventErrorPipeline record a pipeline error event
func EventErrorPipeline(name string, msg string) {
	e := event.NewErrorPipeline(name, msg)
	Event(e)
}

// EventInfoPipeline record a pipeline info event
func EventInfoPipeline(name string, msg string) {
	e := event.NewInfoPipeline(name, msg)
	Event(e)
}
