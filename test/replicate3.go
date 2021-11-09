package main

import (
	"context"
	"fmt"
	repl2 "github.com/jin06/binlogo/pkg/mysql/repl"
	_ "github.com/go-mysql-org/go-mysql/driver"
	"os"
	"time"
)

func main() {
	var cfg = repl2.Config{
		ServerID: 1001,
		Master: repl2.Master{
			Flavor:   "mysql",
			Host:     "127.0.0.1",
			Port:     13306,
			User:     "root",
			Password: "123456",
		},
		Position: repl2.Position{
			File:     "mysql-bin.000001",
			Position: uint32(120),
		},
	}
	syncer := repl2.NewSyncer(cfg)
	syncer.Start()

	for {
		ev, _ := syncer.BinlogStreamer.GetEvent(context.TODO())
		// Dump event
		ev.Dump(os.Stdout)
		fmt.Println(syncer.BinlogSyncer.GetNextPosition())
		//ev.Dump(os.Stdout)
	}

	for {
		ctx, cancel := context.WithTimeout(context.TODO(), 2*time.Second)
		ev, err := syncer.BinlogStreamer.GetEvent(ctx)
		cancel()

		if err == context.DeadlineExceeded {
			// meet timeout
			continue
		}

		ev.Dump(os.Stdout)
	}
}
