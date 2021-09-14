package input

import (
	"context"
	"errors"
	"fmt"
	"github.com/jin06/binlogo/pipeline/message"
	"github.com/siddontang/go-mysql/mysql"
	"github.com/siddontang/go-mysql/replication"
	"strconv"
)

type Input struct {
	syncer   *replication.BinlogSyncer
	streamer *replication.BinlogStreamer
	Ch       chan message.Message
	Options  *Options
}

func (r *Input) Start() (err error) {
	err = r.connect()
	if err != nil {
		return
	}
	err = r.sync()
	return
}

func (r *Input) connect() (err error) {
	binlogFile := r.Options.Position.BinlogFile
	if binlogFile == "" {
		return errors.New("empty binlog file")
	}
	i, err := strconv.Atoi(r.Options.Position.BinlogPosition)
	if err != nil {
		return
	}
	binlogPos := uint32(i)

	pos := mysql.Position{
		binlogFile,
		binlogPos,
	}
	streamer, err := r.syncer.StartSync(pos)
	if err != nil {
		return
	}
	r.streamer = streamer
	return
}

func (r *Input) sync() (err error) {
	go func() {
		//ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		ctx := context.Background()
		ev, _ := r.streamer.GetEvent(ctx)
		fmt.Println(ev.Event)
	}()
	return
}


func (r *Input) DataLine() chan message.Message {
	return r.Ch
}
