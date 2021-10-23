package main

import (
	"fmt"
	"github.com/jin06/binlogo/cmd/server/app"
	"os"
)

func main() {
	if err := run(); err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
}

func run() (err error) {
	command := app.NewCommand()
	return command.Execute()
}

