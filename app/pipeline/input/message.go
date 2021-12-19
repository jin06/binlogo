package input

import (
	"strconv"

	"github.com/go-mysql-org/go-mysql/canal"
	message2 "github.com/jin06/binlogo/app/pipeline/message"
	"github.com/jin06/binlogo/pkg/store/model/pipeline"
)

func emptyMessage() (msgs []*message2.Message) {
	msgs = []*message2.Message{}
	// msg.Filter = true
	// msg.Content.Head = &message2.Head{
	// Type: message2.TYPE_EMPTY.String(),
	// }
	// msg.Content.Data = ""
	return
}

func rowsMessage(e *canal.RowsEvent) (msgs []*message2.Message) {
	// msg = message2.New()
	// msgs = []*message2.Message{}
	switch e.Action {
	case canal.InsertAction:
		{
			return insert(e)
		}
	case canal.UpdateAction:
		{
			return update(e)
		}
	case canal.DeleteAction:
		{
			return delete(e)
		}
	default:
		return emptyMessage()
	}
}

func insert(e *canal.RowsEvent) (msgs []*message2.Message) {
	lengthRows := len(e.Rows)
	totalRows := lengthRows
	msgs = make([]*message2.Message, lengthRows)
	for i := 0; i < totalRows; i++ {
		msg := toMessage(e)
		//msg.Content.Head.Position.TotalRows = totalRows
		//msg.Content.Head.Position.ConsumeRows = i + 1
		newer := map[string]interface{}{}
		for key, val := range e.Rows[i] {
			var columnName string
			if len(e.Table.Columns) > key {
				columnName = e.Table.Columns[key].Name
			} else {
				columnName = strconv.Itoa(key)
			}
			newer[columnName] = val
		}
		msg.Content.Data = message2.Insert{New: newer}
		msgs[i] = msg
	}
	return
}

func update(e *canal.RowsEvent) (msgs []*message2.Message) {
	lengthRows := len(e.Rows)
	totalRows := lengthRows / 2
	msgs = make([]*message2.Message, totalRows)
	for i := 0; i < totalRows; i++ {
		msg := toMessage(e)
		//msg.Content.Head.Position.TotalRows = totalRows
		//msg.Content.Head.Position.ConsumeRows = i + 1
		old := map[string]interface{}{}
		newer := map[string]interface{}{}
		for key, val := range e.Rows[2*i] {
			var columnName string
			if len(e.Table.Columns) > key {
				columnName = e.Table.Columns[key].Name
			} else {
				columnName = strconv.Itoa(key)
			}
			old[columnName] = val
		}
		for key, val := range e.Rows[2*i+1] {
			var columnName string
			if len(e.Table.Columns) > key {
				columnName = e.Table.Columns[key].Name
			} else {
				columnName = strconv.Itoa(key)
			}
			newer[columnName] = val
		}
		msg.Content.Data = message2.Update{Old: old, New: newer}
		msgs[i] = msg
	}
	return
}

func delete(e *canal.RowsEvent) (msgs []*message2.Message) {
	lengthRows := len(e.Rows)
	totalRows := lengthRows
	msgs = make([]*message2.Message, totalRows)
	for i := 0; i < totalRows; i++ {
		msg := toMessage(e)
		//msg.Content.Head.Position.TotalRows = totalRows
		//msg.Content.Head.Position.ConsumeRows = i + 1
		old := map[string]interface{}{}
		for key, val := range e.Rows[i] {
			var columnName string
			if len(e.Table.Columns) > key {
				columnName = e.Table.Columns[key].Name
			} else {
				columnName = strconv.Itoa(key)
			}
			old[columnName] = val
		}
		msg.Content.Data = message2.Delete{Old: old}
		msgs[i] = msg
	}
	return
}

func toMessage(e *canal.RowsEvent) (msg *message2.Message) {
	msg = message2.New()
	msg.Content.Head = &message2.Head{
		Type:     mapType(e.Action),
		Database: e.Table.Schema,
		Table:    e.Table.Name,
		Time:     e.Header.Timestamp,
		Position: &pipeline.Position{},
	}
	return
}

func mapType(s string) (t string) {
	switch s {
	case canal.InsertAction:
		{
			t = message2.TYPE_INSERT.String()
		}
	case canal.UpdateAction:
		{
			t = message2.TYPE_INSERT.String()
		}
	case canal.DeleteAction:
		{
			t = message2.TYPE_DELETE.String()
		}
	default:
		{
			t = ""
		}
	}
	return
}
