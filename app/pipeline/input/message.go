package input

import (
	"strconv"

	"github.com/go-mysql-org/go-mysql/canal"
	"github.com/jin06/binlogo/v2/app/pipeline/message"
)

func emptyMessage() (msgs []*message.Message) {
	msgs = []*message.Message{}
	// msg.Filter = true
	// msg.Content.Head = &message2.Head{
	// Type: message2.TYPE_EMPTY.String(),
	// }
	// msg.Content.Data = ""
	return
}

func rowsMessage(e *canal.RowsEvent) (msgs []*message.Message) {
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

func insert(e *canal.RowsEvent) (msgs []*message.Message) {
	lengthRows := len(e.Rows)
	totalRows := lengthRows
	msgs = make([]*message.Message, lengthRows)
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
		msg.Content.Data = message.Insert{New: newer}
		msgs[i] = msg
	}
	return
}

func update(e *canal.RowsEvent) (msgs []*message.Message) {
	lengthRows := len(e.Rows)
	totalRows := lengthRows / 2
	msgs = make([]*message.Message, totalRows)
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
		msg.Content.Data = message.Update{Old: old, New: newer}
		msgs[i] = msg
	}
	return
}

func delete(e *canal.RowsEvent) (msgs []*message.Message) {
	lengthRows := len(e.Rows)
	totalRows := lengthRows
	msgs = make([]*message.Message, totalRows)
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
		msg.Content.Data = message.Delete{Old: old}
		msgs[i] = msg
	}
	return
}

func toMessage(e *canal.RowsEvent) (msg *message.Message) {
	//msg = message2.New()
	msg = message.Get()
	//msg.Content.Head = &message2.Head{
	//	Type:     mapType(e.Action),
	//	Database: e.Table.Schema,
	//	Table:    e.Table.Name,
	//	Time:     e.Header.Timestamp,
	//	Position: &pipeline.Position{},
	//}
	msg.Content.Head.Type = mapType(e.Action)
	msg.Content.Head.Database = e.Table.Schema
	msg.Content.Head.Table = e.Table.Name
	msg.Content.Head.Time = e.Header.Timestamp
	return
}

func mapType(s string) (t string) {
	switch s {
	case canal.InsertAction:
		{
			t = message.TYPE_INSERT.String()
		}
	case canal.UpdateAction:
		{
			t = message.TYPE_INSERT.String()
		}
	case canal.DeleteAction:
		{
			t = message.TYPE_DELETE.String()
		}
	default:
		{
			t = ""
		}
	}
	return
}
