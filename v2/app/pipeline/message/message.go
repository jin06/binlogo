package message

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/jin06/binlogo/v2/pkg/store/model/pipeline"
)

const (
	// STATUS_NEW a new message
	STATUS_NEW = iota
	// STATUS_SEND message has be sended
	STATUS_SEND
	// STATUS_RECORD message has be recorded
	STATUS_RECORD
)

// STATUS for message control
type STATUS int16

// Message basic message for pipeline handle
type Message struct {
	Status  int16
	Filter  bool
	Content Content
}

// New return a new message
func New() *Message {
	msg := &Message{
		Status: STATUS_NEW,
		Filter: false,
		Content: Content{
			Head: Head{},
			Data: nil,
		},
	}
	return msg
}

// Content of Message
type Content struct {
	Head Head        `json:"head"`
	Data interface{} `json:"data"`
}

func (c *Content) reset() {
	c.Data = nil
	c.Head.reset()
}

// Head of Message
type Head struct {
	Type     string            `json:"type"`
	Time     uint32            `json:"time"`
	Database string            `json:"database"`
	Table    string            `json:"table"`
	Position pipeline.Position `json:"position"`
}

func (h *Head) reset() {
	h.Type = ""
	h.Time = 0
	h.Database = ""
	h.Table = ""
	h.Position.Reset()
}

// Json marshal message to json data
func (msg *Message) Json() (string, error) {
	b, err := json.Marshal(msg)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// JsonContent only marshal message's content to josn data
func (msg *Message) JsonContent() (string, error) {
	//if msg.Content == nil {
	//	return "", nil
	//}
	b, err := json.Marshal(msg.Content)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// ToString returns message's string data
func (msg *Message) ToString() string {
	res := "Status: " + fmt.Sprintf("%v\n", msg.Status)
	res += "Filter: " + fmt.Sprintf("%v\n", msg.Filter)
	//if msg.Content != nil {
	res += "BinlogPosition.File: " + fmt.Sprintf("%v\n", msg.Content.Head.Position.BinlogFile)
	res += "BinlogPosition.Pos: " + fmt.Sprintf("%v\n", msg.Content.Head.Position.BinlogPosition)
	res += "BinlogPosition.GTID: " + fmt.Sprintf("%v\n", msg.Content.Head.Position.GTIDSet)
	res += "Content.Head: " + fmt.Sprintf("%v\n", msg.Content.Head)
	res += "Content.Data: " + fmt.Sprintf("%v\n", msg.Content.Data)
	//}
	return res
}

// Table returns table with database
func (msg *Message) Table() string {
	return fmt.Sprintf("%s.%s", msg.Content.Head.Database, msg.Content.Head.Table)
}

// reset
func (msg *Message) reset() {
	msg.Status = STATUS_NEW
	msg.Filter = false
	msg.Content.reset()
}

// Pool reuse message object
var Pool = sync.Pool{
	New: func() interface{} {
		return &Message{}
	},
}

// Get get a message object from Pool
// func Get() *Message {
// 	msg := Pool.Get().(*Message)
// 	msg.reset()
// 	return msg
// }

// Put put a message to Pool
// func Put(msg *Message) {
// 	Pool.Put(msg)
// }

func Get() *Message {
	return &Message{}
}

func Put(msg *Message) {}
