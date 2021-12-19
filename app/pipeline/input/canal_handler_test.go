package input

import (
	"testing"
	"time"

	"github.com/go-mysql-org/go-mysql/canal"
	"github.com/go-mysql-org/go-mysql/mysql"
	"github.com/go-mysql-org/go-mysql/replication"
	"github.com/go-mysql-org/go-mysql/schema"
	"github.com/jin06/binlogo/app/pipeline/message"
	"github.com/jin06/binlogo/pkg/store/model/pipeline"
)

func TestCanalHandler(t *testing.T) {
	handler := &canalHandler{
		DummyEventHandler: canal.DummyEventHandler{},
		ch:                make(chan *message.Message, 10),
		pipe:              &pipeline.Pipeline{Name: "go_test_pipeline"},
	}
	t.Log(handler.String())
	rowsEvent := &canal.RowsEvent{
		Header: &replication.EventHeader{
			Timestamp: uint32(time.Now().Unix()),
		},
		Action: canal.InsertAction,
		Table: &schema.Table{
			Schema: "database1",
			Name:   "table1",
			Columns: []schema.TableColumn{
				{
					Name: "id",
					Type: schema.TYPE_NUMBER,
				},
			},
			Indexes:         []*schema.Index{},
			PKColumns:       []int{},
			UnsignedColumns: []int{},
		},
		Rows: [][]interface{}{
			{
				10001,
			},
		},
	}
	err := handler.OnRow(rowsEvent)
	if err != nil {
		t.Error(err)
	}
	err = handler.OnTableChanged("", "")
	if err != nil {
		t.Error(err)
	}
	position := mysql.Position{
		Name: "",
		Pos:  0,
	}
	err = handler.OnPosSynced(position, nil, false)
	if err != nil {
		t.Error(err)
	}
	err = handler.OnRotate(nil)
	if err != nil {
		t.Error(err)
	}
	if err = handler.OnXID(mysql.Position{}); err != nil {
		t.Error(err)
	}
	if err = handler.OnGTID(nil); err != nil {
		t.Error(err)
	}
	if err = handler.OnDDL(mysql.Position{}, &replication.QueryEvent{}); err != nil {
		t.Error(err)
	}
}
