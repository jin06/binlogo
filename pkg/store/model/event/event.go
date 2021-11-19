package event

import (
	"github.com/jin06/binlogo/configs"
	"github.com/jin06/binlogo/pkg/util/key"
	"net"
	"time"
)

type Event struct {
	Key          string       `json:"key"`
	Type         Type         `json:"type"`
	ResourceType ResourceType `json:"resource_type"`
	ResourceName string       `json:"resource_name"`
	FirstTime    time.Time    `json:"first_time"`
	LastTime     time.Time    `json:"last_time"`
	Count        int          `json:"count"`
	MessageType  MessageType  `json:"type_message"`
	Message      string       `json:"message"`
	NodeName     string       `json:"node_name"`
	NodeIP       net.IP       `json:"node_ip"`
}

type ResourceType string

const (
	PIPELINE ResourceType = "pipeline"
	NODE     ResourceType = "node"
	CLUSTER  ResourceType = "cluster"
)

type Type string

const (
	WARN  Type = "warn"
	ERROR Type = "error"
	INFO  Type = "info"
)

type MessageType string

func baseEvent() (e *Event) {
	e = &Event{
		Key:      key.GetKey(),
		NodeName: configs.NodeName,
		NodeIP:   configs.NodeIP,
	}
	return
}

func newNode(msg string) (e *Event) {
	e = baseEvent()
	e.ResourceType = NODE
	e.ResourceName = configs.NodeName
	e.Message = msg
	return
}

func NewInfoNode(msg string) (e *Event) {
	e = newNode(msg)
	e.Type = INFO
	return
}

func newPipeline(name string, msg string) (e *Event) {
	e = baseEvent()
	e.ResourceType = PIPELINE
	e.ResourceName = name
	e.Message = msg
	return
}

func NewInfoPipeline(name string, msg string) (e *Event) {
	e = newPipeline(name, msg)
	e.Type = INFO
	return
}

func NewWarnPipeline(name string, msg string) (e *Event) {
	e = newPipeline(name, msg)
	e.Type = WARN
	return
}

func NewErrorPipeline(name string, msg string) (e *Event) {
	e = newPipeline(name, msg)
	e.Type = ERROR
	return
}
