package input

import (
	"errors"
	"github.com/go-mysql-org/go-mysql/canal"
	"github.com/go-mysql-org/go-mysql/replication"
	message2 "github.com/jin06/binlogo/app/pipeline/message"
	"strconv"
)

func inputMessage(e *replication.BinlogEvent) (msg *message2.Message, err error) {
	eventType := e.Header.EventType
	msg = message2.New()
	//pos := &model.Position{}
	switch eventType {
	case replication.UPDATE_ROWS_EVENTv2:
		{
			err = updateMessage(e, msg)
		}
	case replication.WRITE_ROWS_EVENTv2:
		{
			err = insertMessage(e, msg)
		}
	case replication.DELETE_ROWS_EVENTv2:
		{
			err = deleteMessage(e, msg)
		}
	default:
		emptyMessage(msg)
	}
	return
}

func updateMessage(e *replication.BinlogEvent, msg *message2.Message) (err error) {
	if val, ok := e.Event.(*replication.RowsEvent); ok {
		msg.Content.Head = &message2.Head{
			Type:     message2.TYPE_UPDATE.String(),
			Database: string(val.Table.Schema),
			Table:    string(val.Table.Table),
			Time:     e.Header.Timestamp,
		}

		old := map[string]interface{}{}
		newer := map[string]interface{}{}

		for col, cVal := range val.Rows[0] {
			old["todo"+strconv.Itoa(col)] = cVal
		}
		for col, cVal := range val.Rows[1] {
			newer["todo"+strconv.Itoa(col)] = cVal
		}

		msg.Content.Data = message2.Update{
			Old: old,
			New: newer,
		}
	} else {
		err = errors.New("event type error: " + e.Header.EventType.String())
	}
	return
}

func insertMessage(e *replication.BinlogEvent, msg *message2.Message) (err error) {
	if val, ok := e.Event.(*replication.RowsEvent); ok {
		msg.Content.Head = &message2.Head{
			Type:     message2.TYPE_INSERT.String(),
			Database: string(val.Table.Schema),
			Table:    string(val.Table.Table),
			Time:     e.Header.Timestamp,
		}
		newer := map[string]interface{}{}
		for col, cVal := range val.Rows[0] {
			newer["todo"+strconv.Itoa(col)] = cVal
		}

		msg.Content.Data = message2.Insert{
			New: newer,
		}
	} else {
		err = errors.New("event type error: " + e.Header.EventType.String())
	}
	return
}

func deleteMessage(e *replication.BinlogEvent, msg *message2.Message) (err error) {
	if val, ok := e.Event.(*replication.RowsEvent); ok {
		msg.Content.Head = &message2.Head{
			Type:     message2.TYPE_DELETE.String(),
			Database: string(val.Table.Schema),
			Table:    string(val.Table.Table),
			Time:     e.Header.Timestamp,
		}

		old := map[string]interface{}{}
		for col, cVal := range val.Rows[0] {
			old["todo"+strconv.Itoa(col)] = cVal
		}

		msg.Content.Data = message2.Delete{
			Old: old,
		}
	} else {
		err = errors.New("event type error: " + e.Header.EventType.String())
	}
	return
}

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
