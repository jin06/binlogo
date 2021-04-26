package main

import (
	"context"
	"github.com/jin06/binlogo/mysql/repl"
	_ "github.com/siddontang/go-mysql/driver"
	"os"
	"time"
)

func main() {
	var cfg = repl.Config{
		ServerID: 1001,
		Master: repl.Master{
			Flavor:   "mysql",
			Host:     "127.0.0.1",
			Port:     13306,
			User:     "root",
			Password: "123456",
		},
		Position: repl.Position{
			File:     "mysql-bin.000001",
			Position: uint32(120),
		},
	}
	syncer := repl.NewSyncer(cfg)
	syncer.Start()

	for {
		ev, _ := syncer.BinlogStreamer.GetEvent(context.Background())
		// Dump event
		ev.Dump(os.Stdout)
		//ev.Dump(os.Stdout)
	}

	for {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		ev, err := syncer.BinlogStreamer.GetEvent(ctx)
		cancel()

		if err == context.DeadlineExceeded {
			// meet timeout
			continue
		}

		ev.Dump(os.Stdout)
	}
}
