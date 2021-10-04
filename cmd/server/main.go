package main

import (
	"fmt"
	"github.com/jin06/binlogo/cmd/server/app"
	"github.com/jin06/binlogo/config"
	"github.com/jin06/binlogo/store/etcd"
	"os"
)

func main() {
	config.InitViperFromFile()
	etcd.DefaultETCD()

	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func run() (err error) {
	command := app.NewCommand()
	return command.Execute()
}
