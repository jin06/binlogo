package repl

import (
	"fmt"
	"github.com/siddontang/go-mysql/mysql"
	"github.com/siddontang/go-mysql/replication"
)

type Dumper struct {
	BinlogSyncer *replication.BinlogSyncer
	BinlogStreamer *replication.BinlogStreamer
	Config       Config
	//Streamer     *Streamer
}

func (s *Dumper) Start() {
	pos := s.binlogPosition()
	fmt.Println(pos)
	s.BinlogStreamer, _ = s.BinlogSyncer.StartSync(pos)
	return
}

func (s *Dumper) Read() {

}

func (s *Dumper) initBinlogSyncer() {
	c := s.Config.BinlogSyncerConfig()
	s.BinlogSyncer = replication.NewBinlogSyncer(c)
}

func (s *Dumper) binlogPosition() mysql.Position {
	return s.Config.Position.BinlogPosition()
}

func NewSyncer(config Config) *Dumper{
	syncer := new(Dumper)
	syncer.Config = config
	syncer.initBinlogSyncer()
	return syncer
}