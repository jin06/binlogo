package input

import (
	"github.com/go-mysql-org/go-mysql/canal"
	"github.com/go-mysql-org/go-mysql/mysql"
	"github.com/go-mysql-org/go-mysql/replication"
	"github.com/jin06/binlogo/app/pipeline/message"
	"github.com/jin06/binlogo/pkg/store/model/pipeline"
	"github.com/siddontang/go-log/log"
	"github.com/sirupsen/logrus"
)

type canalHandler struct {
	canal.DummyEventHandler
	ch           chan *message.Message
	positionFile string
	_GTIDSet     string
	pipe         *pipeline.Pipeline
}

func (h *canalHandler) OnRow(e *canal.RowsEvent) error {
	//logrus.Warnln(e)
	msg := rowsMessage(e)
	msg.BinlogPosition = &pipeline.Position{
		BinlogFile:     h.positionFile,
		BinlogPosition: e.Header.LogPos,
		PipelineName:   h.pipe.Name,
		GTIDSet:        h._GTIDSet,
	}
	logrus.Errorf("%v", e.Header)
	h.ch <- msg
	return nil
}
func (h *canalHandler) OnPosSynced(pos mysql.Position, set mysql.GTIDSet, force bool) error {
	return nil
}

func (h *canalHandler) OnRotate(e *replication.RotateEvent) error {
	h.positionFile = string(e.NextLogName)
	return nil
}

func (h *canalHandler) OnXID(p mysql.Position) error {
	log.Infoln("xid event", p)
	return nil
}

func (h *canalHandler) String() string {
	return "MyEventHandler"
}

func (h *canalHandler) OnGTID(set mysql.GTIDSet) (err error) {
	set.String()
	return
}

