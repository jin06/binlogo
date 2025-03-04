package input

import (
	"github.com/go-mysql-org/go-mysql/canal"
	"github.com/go-mysql-org/go-mysql/mysql"
	"github.com/go-mysql-org/go-mysql/replication"
	"github.com/jin06/binlogo/v2/app/pipeline/message"
	"github.com/jin06/binlogo/v2/configs"
	"github.com/jin06/binlogo/v2/pkg/promeths"
	"github.com/jin06/binlogo/v2/pkg/store/model/pipeline"
	"github.com/prometheus/client_golang/prometheus"
)

type canalHandler struct {
	canal.DummyEventHandler
	ch       chan *message.Message
	pipe     *pipeline.Pipeline
	messages []*message.Message
}

func (h *canalHandler) OnRow(e *canal.RowsEvent) error {
	// fmt.Println(e.Rows)
	if h.messages == nil {
		h.messages = []*message.Message{}
	}
	// fmt.Println("---> ", len(e.Rows))
	// fmt.Println(e.Header.LogPos)
	msgs := rowsMessage(e)
	// h.msg = msg
	h.messages = append(h.messages, msgs...)

	return nil
}
func (h *canalHandler) OnTableChanged(schema string, table string) error {
	//fmt.Println(schema)
	//fmt.Println(table)
	return nil
}
func (h *canalHandler) OnPosSynced(pos mysql.Position, set mysql.GTIDSet, force bool) error {
	defer func() {
		h.messages = nil
	}()
	// fmt.Println("on pos synced", set)
	if h.messages == nil {
		return nil
	}
	total := len(h.messages)
	for i := 0; i < total; i++ {
		msg := h.messages[i]
		msg.Content.Head.Position.BinlogPosition = pos.Pos
		msg.Content.Head.Position.BinlogFile = pos.Name
		msg.Content.Head.Position.PipelineName = h.pipe.Name
		msg.Content.Head.Position.TotalRows = total
		msg.Content.Head.Position.ConsumeRows = i + 1
		//fmt.Println("on pos synced" ,set)
		if set != nil {
			msg.Content.Head.Position.GTIDSet = set.String()
		}
		h.ch <- msg
		promeths.MessageTotalCounter.With(prometheus.Labels{"pipeline": h.pipe.Name, "node": configs.GetNodeName()}).Inc()
	}
	return nil
}

func (h *canalHandler) OnRotate(e *replication.RotateEvent) error {
	return nil
}

func (h *canalHandler) OnXID(p mysql.Position) error {
	// fmt.Println("-------> on xid")
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

func (h *canalHandler) OnDDL(nextPos mysql.Position, queryEvent *replication.QueryEvent) error {
	//fmt.Println(nextPos, queryEvent)
	//fmt.Println(string(queryEvent.Query))
	//queryEvent.Dump(os.Stdout)
	return nil
}
