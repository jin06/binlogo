package input

import (
	"github.com/go-mysql-org/go-mysql/canal"
	message2 "github.com/jin06/binlogo/app/pipeline/message"
	"strconv"
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
	newer := map[string]interface{}{}
	for key, val := range e.Rows[0] {
		var columnName string
		if len(e.Table.Columns) > key {
			columnName = e.Table.Columns[key].Name
		} else {
			columnName = strconv.Itoa(key)
		}
		newer[columnName] = val
	}

	msg.Content.Data = message2.Insert{
		New: newer,
	}
}

func update(e *canal.RowsEvent, msg *message2.Message) {
	msg.Content.Head = &message2.Head{
		Type:     message2.TYPE_UPDATE.String(),
		Database: e.Table.Schema,
		Table:    e.Table.Name,
		Time:     e.Header.Timestamp,
	}
	old := map[string]interface{}{}
	newer := map[string]interface{}{}
	if len(e.Rows) == 2 {
		for key, val := range e.Rows[0] {
			var columnName string
			if len(e.Table.Columns) > key {
				columnName = e.Table.Columns[key].Name
			} else {
				columnName = strconv.Itoa(key)
			}
			old[columnName] = val
		}
		for key, val := range e.Rows[1] {
			var columnName string
			if len(e.Table.Columns) > key {
				columnName = e.Table.Columns[key].Name
			} else {
				columnName = strconv.Itoa(key)
			}
			newer[columnName] = val
		}
	}

	msg.Content.Data = message2.Update{
		Old: old,
		New: newer,
	}
}

func delete(e *canal.RowsEvent, msg *message2.Message) {
	msg.Content.Head = &message2.Head{
		Type:     message2.TYPE_DELETE.String(),
		Database: e.Table.Schema,
		Table:    e.Table.Name,
		Time:     e.Header.Timestamp,
	}
	old := map[string]interface{}{}
	for key, val := range e.Rows[0] {
		var columnName string
		if len(e.Table.Columns) > key {
			columnName = e.Table.Columns[key].Name
		} else {
			columnName = strconv.Itoa(key)
		}
		old[columnName] = val
	}

	msg.Content.Data = message2.Delete{
		Old: old,
	}
}
