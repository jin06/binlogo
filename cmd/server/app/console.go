package app

import (
	"context"
	"github.com/jin06/binlogo/app/server/console"
)

// RunConsole run gin
func RunConsole() (err error) {
	err = console.Run(context.Background())
	return
}
