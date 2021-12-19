package app

import (
	"context"
	"github.com/sirupsen/logrus"

	"github.com/jin06/binlogo/app/server/console"
)

// RunConsole run gin
func RunConsole(c context.Context) (err error) {
	logrus.Info("init console")
	err = console.Run(c)
	return
}
