package input

import (
	"strconv"

	"github.com/go-mysql-org/go-mysql/canal"
	message2 "github.com/jin06/binlogo/app/pipeline/message"
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
	msgs = make([]*message2.Message, len(e.Rows))
	for k, v := range e.Rows {
		msg := message2.New()
		msg.Content.Head = &message2.Head{
			Type:     message2.TYPE_INSERT.String(),
			Database: e.Table.Schema,
			Table:    e.Table.Name,
			Time:     e.Header.Timestamp,
		}
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
		msg.Content.Data = message2.Insert{New: newer}
		msgs[k] = msg
	}
	return
}

func update(e *canal.RowsEvent) (msgs []*message2.Message) {
	length := len(e.Rows)
	msgs = make([]*message2.Message, length/2)
	for i := 0; i < length; i = i + 2 {
		msg := message2.New()
		msg.Content.Head = &message2.Head{
			Type:     message2.TYPE_UPDATE.String(),
			Database: e.Table.Schema,
			Table:    e.Table.Name,
			Time:     e.Header.Timestamp,
		}
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
		msg.Content.Data = message2.Update{Old: old, New: newer}
		msgs[i/2] = msg
	}
	return
}

func delete(e *canal.RowsEvent) (msgs []*message2.Message) {
	length := len(e.Rows)
	msgs = make([]*message2.Message, length)
	for k, v := range e.Rows {
		msg := message2.New()
		msg.Content.Head = &message2.Head{
			Type:     message2.TYPE_DELETE.String(),
			Database: e.Table.Schema,
			Table:    e.Table.Name,
			Time:     e.Header.Timestamp,
		}
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
		msg.Content.Data = message2.Delete{Old: old}
		msgs[k] = msg
	}
	return
}
