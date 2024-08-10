package app

import (
	"github.com/jin06/binlogo/v2/pkg/event"
	"github.com/sirupsen/logrus"
)

// RunEvent start event recorder goroutine
func RunEvent() {
	logrus.Info("init event")
	event.Init()
}
