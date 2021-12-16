package input

import (
	"strconv"

	"github.com/go-mysql-org/go-mysql/canal"
	"github.com/jin06/binlogo/app/pipeline/message"
	message2 "github.com/jin06/binlogo/app/pipeline/message"
)

func emptyMessage(msg *message2.Message) {
	msg.Filter = true
	msg.Content.Head = &message2.Head{
		Type: message2.TYPE_EMPTY.String(),
	}
	msg.Content.Data = ""
	return
}

func rowsMessage(e *canal.RowsEvent) (msg *message2.Message) {
	msg = message2.New()
	switch e.Action {
	case canal.InsertAction:
		{
			insert(e, msg)
		}
	case canal.UpdateAction:
		{
			update(e, msg)
		}
	case canal.DeleteAction:
		{
			delete(e, msg)
		}
	default:
		emptyMessage(msg)
	}
	return
}

func insert(e *canal.RowsEvent, msg *message2.Message) {
	msg.Content.Head = &message2.Head{
		Type:     message2.TYPE_INSERT.String(),
		Database: e.Table.Schema,
		Table:    e.Table.Name,
		Time:     e.Header.Timestamp,
	}
	arr := make([]message.Insert, len(e.Rows))
	for k, v := range e.Rows {
		newer := map[string]interface{}{}
		for key, val := range v {
			var columnName string
			if len(e.Table.Columns) > key {
				columnName = e.Table.Columns[key].Name
			} else {
				columnName = strconv.Itoa(key)
			}
			newer[columnName] = val
		}
		arr[k] = message2.Insert{New: newer}
	}
	msg.Content.Data = arr
}

func update(e *canal.RowsEvent, msg *message2.Message) {
	msg.Content.Head = &message2.Head{
		Type:     message2.TYPE_UPDATE.String(),
		Database: e.Table.Schema,
		Table:    e.Table.Name,
		Time:     e.Header.Timestamp,
	}

	length := len(e.Rows)
	arr := make([]message2.Update, length/2)

	for i := 0; i < length; i = i + 2 {
		old := map[string]interface{}{}
		newer := map[string]interface{}{}
		for key, val := range e.Rows[i] {
			var columnName string
			if len(e.Table.Columns) > key {
				columnName = e.Table.Columns[key].Name
			} else {
				columnName = strconv.Itoa(key)
			}
			old[columnName] = val
		}
		for key, val := range e.Rows[i+1] {
			var columnName string
			if len(e.Table.Columns) > key {
				columnName = e.Table.Columns[key].Name
			} else {
				columnName = strconv.Itoa(key)
			}
			newer[columnName] = val
		}
		arr[i/2] = message2.Update{Old: old, New: newer}
	}
	msg.Content.Data = arr
}

func delete(e *canal.RowsEvent, msg *message2.Message) {
	msg.Content.Head = &message2.Head{
		Type:     message2.TYPE_DELETE.String(),
		Database: e.Table.Schema,
		Table:    e.Table.Name,
		Time:     e.Header.Timestamp,
	}
	length := len(e.Rows)
	arr := make([]message2.Delete, length)
	for k, v := range e.Rows {
		old := map[string]interface{}{}
		for key, val := range v {
			var columnName string
			if len(e.Table.Columns) > key {
				columnName = e.Table.Columns[key].Name
			} else {
				columnName = strconv.Itoa(key)
			}
			old[columnName] = val
		}
		arr[k] = message2.Delete{Old: old}
	}

	msg.Content.Data = arr
}
