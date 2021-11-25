package message

import (
	"encoding/json"
	"fmt"
	"github.com/jin06/binlogo/pkg/store/model/pipeline"
)

const (
	// STATUS_NEW a new message
	STATUS_NEW = iota
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
	return fmt.Sprintf(
		`
	Status: %v
	Filter: %v
	BingloPosition.File: %v
	BingloPosition.Pos: %v
	BingloPosition.GTID: %v
	Content.Head: %+v
	Content.Data: %+v
			`,
		msg.Status,
		msg.Filter,
		msg.BinlogPosition.BinlogFile,
		msg.BinlogPosition.BinlogPosition,
		msg.BinlogPosition.GTIDSet,
		*msg.Content.Head,
		msg.Content.Data,
	)
}

func (msg *Message) Table() string {
	return fmt.Sprintf("%s.%s", msg.Content.Head.Database, msg.Content.Head.Table)
}
