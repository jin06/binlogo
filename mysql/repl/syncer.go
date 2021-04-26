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
	//Streamer     *Streamer
}

func (s *Syncer) Start() {
	pos := s.binlogPosition()
	fmt.Println(pos)
	s.BinlogStreamer, _ = s.BinlogSyncer.StartSync(pos)
	return
}

func (s *Syncer) Read() {

}

func (s *Syncer) initBinlogSyncer() {
	c := s.Config.BinlogSyncerConfig()
	s.BinlogSyncer = replication.NewBinlogSyncer(c)
}

func (s *Syncer) binlogPosition() mysql.Position {
	return s.Config.Position.BinlogPosition()
}

func NewSyncer(config Config) *Syncer{
	syncer := new(Syncer)
	syncer.Config = config
	syncer.initBinlogSyncer()
	return syncer
}