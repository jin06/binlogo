package main

import (
	"fmt"
	"os"

	"github.com/jin06/binlogo/v2/cmd/cli/app"
)

func main() {
	if err := run(); err != nil {
		if err := run(); err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}
	}
}

func run() (err error) {
	command := app.NewCommand()
	return command.Execute()
}
