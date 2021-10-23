package app

import (
	"context"
	"github.com/jin06/binlogo/app/server/console"
)

func RunConsole() (err error) {
	ctx := context.TODO()
	err = console.Run(ctx)
	return
}
