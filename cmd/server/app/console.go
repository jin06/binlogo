package app

import (
	"context"
	"github.com/jin06/binlogo/app/server/console"
)

func RunConsole() (err error) {
	err = console.Run(context.Background())
	return
}
