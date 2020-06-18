package core

import (
	"github.com/siddontang/go-mysql/replication"
)

type Importer struct {
	Syncer     *replication.BinlogSyncer
	SyncerCfg  *replication.BinlogSyncerConfig
	BinlogFile string
	BinlogPos  uint32
}

func (importer *Importer) InitSyncer() {
	if importer.SyncerCfg != nil {
		importer.Syncer = replication.NewBinlogSyncer(*importer.SyncerCfg)
	} else {
		panic("null Importer.SyncerCfg")
	}
}

func (importer *Importer) updateBinlogFile(file string) {
	importer.BinlogFile = file
}

func (importer *Importer) updateBinlogPos(pos uint32) {
	importer.BinlogPos = pos
}
