package main

import (
	"fmt"
	"github.com/jin06/binlogo/core"
	"github.com/siddontang/go-mysql/replication"
)

func main()  {
	importer := core.Importer{}
	importer.SyncerCfg = &replication.BinlogSyncerConfig{
		ServerID: 100,
		Flavor:   "mysql",
		Host:     "127.0.0.1",
		Port:     3306,
		User:     "root",
		Password: "123456",
	}
	importer.BinlogFile = "mysql-bin.000040"
	importer.BinlogPos =  4
	pipe := new(core.Pipeline)

	pipe.Importer = &importer
	pipe.Exporter = &core.Stand{}
	err := pipe.Run()
	fmt.Println(err)
	select {

	}

}