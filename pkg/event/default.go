package event

import (
	"context"
	"github.com/jin06/binlogo/pkg/store/model/event"
)

var DefaultRecorder *Recorder

func Init() {
	DefaultRecorder , _ = New()
	DefaultRecorder.Loop(context.Background())
}

func Event(e *event.Event) {
	//fmt.Println(e)
	go DefaultRecorder.Event(e)
}

