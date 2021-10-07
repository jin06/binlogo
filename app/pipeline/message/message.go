package message

import (
	"encoding/json"
	"fmt"
	model2 "github.com/jin06/binlogo/pkg/store/model"
)

const (
	STATUS_NEW = iota
)

type STATUS int16

type Message struct {
	Status         int16
	Filter         bool
	BinlogPosition *model2.Position
	Content        *Content
}

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

type Content struct {
	Head *Head       `json:"head"`
	Data interface{} `json:"data"`
}

type Head struct {
	Type     string `json:"type"`
	Time     uint32 `json:"time"`
	Database string `json:"database"`
	Table    string `json:"table"`
}

func (msg *Message) Json() (string, error) {
	b, err := json.Marshal(msg)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

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

func (msg *Message) ToString() string {
	return fmt.Sprintf(
		`
	Status: %v
	Filter: %v
	BingloPosition.File: %v
	BingloPosition.Pos: %v
	Content.Head: %+v
	Content.Data: %+v
			`,
		msg.Status,
		msg.Filter,
		msg.BinlogPosition.BinlogFile,
		msg.BinlogPosition.BinlogPosition,
		*msg.Content.Head,
		msg.Content.Data,
	)
}
