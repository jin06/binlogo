package app

import (
	"context"

	"github.com/jin06/binlogo/app/server/console"
)

// RunConsole run gin
func RunConsole(c context.Context) (err error) {
	err = console.Run(c)
	return
}
