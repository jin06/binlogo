package input

import (
	"testing"
	"time"

	"github.com/go-mysql-org/go-mysql/canal"
	"github.com/go-mysql-org/go-mysql/replication"
	"github.com/go-mysql-org/go-mysql/schema"
	"github.com/jin06/binlogo/v2/app/pipeline/message"
)

func TestRowsMessage(t *testing.T) {
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
	msg := rowsMessage(rowsEvent)
	if msg.Content.Head.Type != "insert" {
		t.Fail()
	}
	if val, ok := msg.Content.Data.(message.Insert); !ok {
		t.Fail()
	} else {
		if val2, ok2 := val.New["id"].(int); !ok2 {
			t.Fail()
		} else {
			if val2 != 10001 {
				t.Fail()
			}
		}
	}

	rowsEvent.Action = canal.UpdateAction
	msg = rowsMessage(rowsEvent)
	if _, ok := msg.Content.Data.(message.Update); !ok {
		t.Fail()
	}
	rowsEvent.Action = canal.DeleteAction
	msg = rowsMessage(rowsEvent)
	if _, ok := msg.Content.Data.(message.Delete); !ok {
		t.Fail()
	}
}
