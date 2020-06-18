package main

import (
	"context"
	_ "github.com/siddontang/go-mysql/driver"
	"github.com/siddontang/go-mysql/mysql"
	"github.com/siddontang/go-mysql/replication"
	"os"
	"time"
)

func main() {
	cfg := replication.BinlogSyncerConfig{
		ServerID: 100,
		Flavor:   "mysql",
		Host:     "127.0.0.1",
		Port:     3306,
		User:     "root",
		Password: "123456",
	}
	binlogFile := "mysql-bin.000040"
	binlogPos := uint32(4)
	position := mysql.Position{
		binlogFile,
		binlogPos,
	}
	//_ = position
	syncer := replication.NewBinlogSyncer(cfg)
	//nposition := syncer.GetNextPosition();
	//fmt.Print(nposition)
	//fmt.Print(nposition.Name)
	// Start sync with specified binlog file and position
	streamer, _ := syncer.StartSync(position)

	// or you can start a gtid replication like
	// streamer, _ := syncer.StartSyncGTID(gtidSet)
	// the mysql GTID set likes this "de278ad0-2106-11e4-9f8e-6edd0ca20947:1-2"
	// the mariadb GTID set likes this "0-1-100"

	for {
		ev, _ := streamer.GetEvent(context.Background())
		// Dump event
		ev.Dump(os.Stdout)
		//ev.Dump(os.Stdout)
	}

	// or we can use a timeout context
	for {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		ev, err := streamer.GetEvent(ctx)
		cancel()

		if err == context.DeadlineExceeded {
			// meet timeout
			continue
		}

		ev.Dump(os.Stdout)
	}
}
