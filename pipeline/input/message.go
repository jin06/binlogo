package input

import (
	"errors"
	"fmt"
	"github.com/jin06/binlogo/pipeline/message"
	"github.com/siddontang/go-mysql/replication"
	"strconv"
)

func inputMessage(e *replication.BinlogEvent) (msg *message.Message, err error) {
	eventType := e.Header.EventType
	msg = &message.Message{}
	switch eventType {
	case replication.UPDATE_ROWS_EVENTv2:
		{
			if val, ok := e.Event.(*replication.RowsEvent); ok {
				fmt.Println("update_rows_eventv2 :")
				//val.Table.Dump(os.Stdout)
				//fmt.Println(val.Table.ColumnNameString())
				//val.Table.Dump(os.Stdout)
				fmt.Println(val.Rows)
				//time.Sleep(10 * time.Second)
				msg.Content = &message.Content{}
				msg.Content.Head = &message.Head{
					Type:     message.TYPE_INSERT,
					Database: string(val.Table.Schema),
					Table:    string(val.Table.Table),
					Time:     e.Header.Timestamp,
				}

				old, newer := map[string]interface{}{}, map[string]interface{}{}
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

				return
			} else {
				err = errors.New("event type error: " + eventType.String())
				return
			}
		}
	case replication.WRITE_ROWS_EVENTv2:
		{

		}
	}
	return
}

func updateMessage(e *replication.BinlogEvent) (msg *message.Message, err error ) {

	return
}
