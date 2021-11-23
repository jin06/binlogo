package app

import "github.com/jin06/binlogo/pkg/event"

// RunEvent start event recorder goroutine
func RunEvent() {
	event.Init()
}
