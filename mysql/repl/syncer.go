package repl

import (
	"fmt"
	"github.com/siddontang/go-mysql/mysql"
	"github.com/siddontang/go-mysql/replication"
)

type Syncer struct {
	BinlogSyncer *replication.BinlogSyncer
	BinlogStreamer *replication.BinlogStreamer
	Config       Config
	Position     Position
	//Streamer     *Streamer
}

func (s *Syncer) Start() {
	position := mysql.Position{
		s.Position.File,
		s.Position.Position,
	}
	fmt.Println(position)
	s.BinlogStreamer, _ = s.BinlogSyncer.StartSync(position)
	return
}

func (s *Syncer) Read() {

}
