package core

import "github.com/siddontang/go-mysql/replication"

type Event struct {
	BinlogEvent *replication.BinlogEvent
}
