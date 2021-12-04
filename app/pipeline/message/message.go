package message

import (
	"encoding/json"
	"fmt"
	"github.com/jin06/binlogo/pkg/store/model/pipeline"
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
	Status         int16
	Filter         bool
	BinlogPosition *pipeline.Position
	Content        *Content
}

// New return a new message
func New() *Message {
	msg := &Message{
		Status: STATUS_NEW,
		Filter: false,
		Content: &Content{
			Head: &Head{},
			Data: nil,
		},
	}
	return msg
}

// Content of Message
type Content struct {
	Head *Head       `json:"head"`
	Data interface{} `json:"data"`
}

// Head of Message
type Head struct {
	Type     string             `json:"type"`
	Time     uint32             `json:"time"`
	Database string             `json:"database"`
	Table    string             `json:"table"`
	Position *pipeline.Position `json:"position"`
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
	if msg.Content == nil {
		return "", nil
	}
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
	if msg.BinlogPosition != nil {
		res += "BinlogPosition.File: " + fmt.Sprintf("%v\n", msg.BinlogPosition.BinlogFile)
		res += "BinlogPosition.Pos: " + fmt.Sprintf("%v\n", msg.BinlogPosition.BinlogPosition)
		res += "BinlogPosition.GTID: " + fmt.Sprintf("%v\n", msg.BinlogPosition.GTIDSet)
	}
	if msg.Content != nil {
		res += "Content.Head: " + fmt.Sprintf("%v\n", msg.Content.Head)
		res += "Content.Data: " + fmt.Sprintf("%v\n", msg.Content.Data)
	}
	return res
}

func (msg *Message) Table() string {
	return fmt.Sprintf("%s.%s", msg.Content.Head.Database, msg.Content.Head.Table)
}
