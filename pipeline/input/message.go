package input

import (
	"errors"
	"github.com/jin06/binlogo/pipeline/message"
	"github.com/siddontang/go-mysql/replication"
	"strconv"
)

func inputMessage(e *replication.BinlogEvent) (msg *message.Message, err error) {
	eventType := e.Header.EventType
	switch eventType {
	case replication.UPDATE_ROWS_EVENTv2:
		{
			return updateMessage(e)
		}
	case replication.WRITE_ROWS_EVENTv2:
		{
			return insertMessage(e)
		}
	case replication.DELETE_ROWS_EVENTv2:
		{
			return deleteMessage(e)
		}
	default:
	}
	return
}

func updateMessage(e *replication.BinlogEvent) (msg *message.Message, err error) {
	msg = &message.Message{}
	if val, ok := e.Event.(*replication.RowsEvent); ok {
		msg.Content = &message.Content{}
		msg.Content.Head = &message.Head{
			Type:     message.TYPE_UPDATE,
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

		msg.Content.Data = message.Update{
			Old: old,
			New: newer,
		}
	} else {
		err = errors.New("event type error: " + e.Header.EventType.String())
	}
	return
}

func insertMessage(e *replication.BinlogEvent) (msg *message.Message, err error) {
	msg = &message.Message{}
	if val, ok := e.Event.(*replication.RowsEvent); ok {
		msg.Content = &message.Content{}
		msg.Content.Head = &message.Head{
			Type:     message.TYPE_INSERT,
			Database: string(val.Table.Schema),
			Table:    string(val.Table.Table),
			Time:     e.Header.Timestamp,
		}
		newer := map[string]interface{}{}
		for col, cVal := range val.Rows[0] {
			newer["todo"+strconv.Itoa(col)] = cVal
		}

		msg.Content.Data = message.Insert{
			New: newer,
		}
	} else {
		err = errors.New("event type error: " + e.Header.EventType.String())
	}
	return
}

func deleteMessage(e *replication.BinlogEvent) (msg *message.Message, err error) {
	msg = &message.Message{}
	if val, ok := e.Event.(*replication.RowsEvent); ok {
		msg.Content = &message.Content{}
		msg.Content.Head = &message.Head{
			Type:     message.TYPE_DELETE,
			Database: string(val.Table.Schema),
			Table:    string(val.Table.Table),
			Time:     e.Header.Timestamp,
		}

		old := map[string]interface{}{}
		for col, cVal := range val.Rows[0] {
			old["todo"+strconv.Itoa(col)] = cVal
		}

		msg.Content.Data = message.Update{
			Old: old,
		}
	} else {
		err = errors.New("event type error: " + e.Header.EventType.String())
	}
	return
}
