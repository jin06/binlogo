package input

import (
	"github.com/go-mysql-org/go-mysql/canal"
	"github.com/go-mysql-org/go-mysql/mysql"
	"github.com/go-mysql-org/go-mysql/replication"
	"github.com/jin06/binlogo/app/pipeline/message"
	"github.com/jin06/binlogo/pkg/store/model/pipeline"
)

type canalHandler struct {
	canal.DummyEventHandler
	ch           chan *message.Message
	pipe         *pipeline.Pipeline
	msg          *message.Message
}

func (h *canalHandler) OnRow(e *canal.RowsEvent) error {
	msg := rowsMessage(e)
	h.msg = msg
	return nil
}
func (h *canalHandler) OnPosSynced(pos mysql.Position, set mysql.GTIDSet, force bool) error {
	if h.msg != nil {
		h.msg.BinlogPosition = &pipeline.Position{
			BinlogFile:     pos.Name,
			BinlogPosition: pos.Pos,
			PipelineName:   h.pipe.Name,
		}
		//fmt.Println("on pos synced" ,set)
		if set != nil {
			h.msg.BinlogPosition.GTIDSet = set.String()
		}

		h.ch <- h.msg
		h.msg = nil
	}
	return nil
}

func (h *canalHandler) OnRotate(e *replication.RotateEvent) error {
	return nil
}

func (h *canalHandler) OnXID(p mysql.Position) error {
	//if h.msg != nil {
	//	fmt.Println("on xid ", p)
	//}
	return nil
}

func (h *canalHandler) String() string {
	return "MyEventHandler"
}

func (h *canalHandler) OnGTID(set mysql.GTIDSet) (err error) {
	return
}
