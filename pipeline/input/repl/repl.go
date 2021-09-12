package repl

import (
	"github.com/jin06/binlogo/pipeline/message"
	"github.com/siddontang/go-mysql/replication"
)

func New() *Repl {
	r := &Repl{}
	return r
}

type Repl struct {
	Syncer *replication.BinlogSyncer
	Streamer *replication.BinlogStreamer
	DataLine chan message.Message
}

func (repl *Repl) Start() {

}
