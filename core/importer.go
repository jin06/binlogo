package core

import (
	"context"
	"github.com/siddontang/go-mysql/mysql"
	"github.com/siddontang/go-mysql/replication"
)

type Importer struct {
	Id         uint
	Syncer     *replication.BinlogSyncer
	SyncerCfg  *replication.BinlogSyncerConfig
	BinlogFile string
	BinlogPos  uint32
}

func (im *Importer) InitSyncer() {
	if im.SyncerCfg != nil {
		im.Syncer = replication.NewBinlogSyncer(*im.SyncerCfg)
	} else {
		panic("null Importer.SyncerCfg")
	}
}

func (im *Importer) Start() (ch chan *Event, err error) {
	im.InitSyncer()
	position := mysql.Position{im.BinlogFile, im.BinlogPos}
	streamer, err := im.Syncer.StartSync(position)
	if err != nil {
		return
	}
	ch = make(chan *Event, 10000)
	go func() {
		for {
			event, err := streamer.GetEvent(context.Background())
			if err != nil {
				panic(err.Error())
			}
			ch <- &Event{event}
		}
	}()
	return
}

func (im *Importer) updateBinlogFile(file string) {
	im.BinlogFile = file
}

func (im *Importer) updateBinlogPos(pos uint32) {
	im.BinlogPos = pos
}
