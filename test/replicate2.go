package main

import (
	"context"
	repl2 "github.com/jin06/binlogo/pkg/mysql/repl"
	_ "github.com/siddontang/go-mysql/driver"
	"github.com/siddontang/go-mysql/replication"
	"os"
	"time"
)

func main() {
	syncer := repl2.Dumper{

	}
	syncer.Position = repl2.Position{
		"mysql-bin.000001",
		uint32(120),
	}

	cfg := replication.BinlogSyncerConfig{
		ServerID: 1001,
		Flavor:   "mysql",
		//Host:     "127.0.0.1",
		//Port:     3306,
		//User:     "root",
		//Password: "123456",
		Host:     "127.0.0.1",
		Port:     13306,
		User:     "root",
		Password: "123456",
	}

	syncer.BinlogSyncer = replication.NewBinlogSyncer(cfg)
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
