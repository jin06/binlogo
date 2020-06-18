package core

import (
	"context"
	"fmt"
	"github.com/siddontang/go-mysql/mysql"
	"github.com/siddontang/go-mysql/replication"
	"os"
)

type Importer struct {
	Id uint
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

func (im *Importer) Start() {
	position := mysql.Position{im.BinlogFile,im.BinlogPos}
	streamer, err := im.Syncer.StartSync(position)
	if err != nil {
		panic(err.Error())
	}
	for {
		event , err := streamer.GetEvent(context.Background())
		if err != nil {
			panic(err.Error())
		}
		//fmt.Println(string(event.RawData))
		fmt.Println("@@")
		event.Dump(os.Stdout)
		//fmt.Println(event.Header.EventType)
	}
}

func (im *Importer) updateBinlogFile(file string) {
	im.BinlogFile = file
}

func (im *Importer) updateBinlogPos(pos uint32) {
	im.BinlogPos = pos
}

