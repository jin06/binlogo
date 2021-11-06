package main

import (
	"fmt"
	"github.com/go-mysql-org/go-mysql/canal"
	"github.com/go-mysql-org/go-mysql/mysql"
	"github.com/go-mysql-org/go-mysql/replication"
	"github.com/siddontang/go-log/log"
	"os"
)

type MyEventHandler struct {
	canal.DummyEventHandler
	canal *canal.Canal
}

func (h *MyEventHandler) OnRow(e *canal.RowsEvent) error {
	log.Infof("%s %v\n", e.Action, e.Rows)
	log.Infof("$f", e.Table.Schema)
	log.Infof("$f", e.Table.Name)
	log.Infof("$f", e.Table.Columns)
	log.Infoln(e.Header.LogPos)
	log.Infoln(e.Header)
	log.Infoln("synced position", h.canal.SyncedPosition())
	e.Header.Dump(os.Stdout)
	return nil
	//return nil
}
func (h *MyEventHandler) OnPosSynced(pos mysql.Position, set mysql.GTIDSet, force bool) error {
	log.Infoln(pos)
	log.Infoln(set)
	return nil
}

func (h *MyEventHandler) OnRotate(e *replication.RotateEvent) error {
	log.Infoln("next log name", e.NextLogName)
	return nil
}

func (h *MyEventHandler) OnXID(p mysql.Position) error {
	log.Infoln("xid event", p)
	return nil
}

func (h *MyEventHandler) String() string {
	return "MyEventHandler"
}

func main() {
	//cfg := canal.NewDefaultConfig()
	cfg := &canal.Config{}
	cfg.Addr = "127.0.0.1:13306"
	cfg.User = "root"
	cfg.Password = "123456"
	cfg.ServerID = 1111

	// We only care table canal_test in test db
	//cfg.Dump.ExecutionPath = "./abc.dump"
	//cfg.Dump.TableDB = "test"
	//cfg.Dump.Tables = []string{"canal_test"}

	c, err := canal.NewCanal(cfg)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(123)

	// Register a handler to handle RowsEvent
	// Start canal
	//c.Run()
	c.SetEventHandler(&MyEventHandler{canal: c})

	c.RunFrom(mysql.Position{Name: "mysql-bin.000016", Pos: uint32(0)})
}
