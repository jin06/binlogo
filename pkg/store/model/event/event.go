package event

import (
	"github.com/jin06/binlogo/configs"
	"github.com/jin06/binlogo/pkg/util/key"
	"net"
	"time"
)

// Event generated by pipelines and nodes in the cluster.
// Events include common information and error information,
// which can be easily located when an error occurs
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

// ResourceType pipeline or node or etc.
type ResourceType string

const (
	// PIPELINE resource type
	PIPELINE ResourceType = "pipeline"
	// ResourceType resource type
	NODE ResourceType = "node"
	// CLUSTER resource type
	CLUSTER ResourceType = "cluster"
)

// Type event level
type Type string

const (
	// WARN level
	WARN Type = "warn"
	// ERROR level
	ERROR Type = "error"
	// INFO level
	INFO Type = "info"
)

// MessageType todo
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

// NewInfoNode returns a new node event of info level
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

// NewInfoPipeline returns a new pipeline info event
func NewInfoPipeline(name string, msg string) (e *Event) {
	e = newPipeline(name, msg)
	e.Type = INFO
	return
}

// NewWarnPipeline returns a new pipeline warn event
func NewWarnPipeline(name string, msg string) (e *Event) {
	e = newPipeline(name, msg)
	e.Type = WARN
	return
}

// NewErrorPipeline returns a new pipeline error event
func NewErrorPipeline(name string, msg string) (e *Event) {
	e = newPipeline(name, msg)
	e.Type = ERROR
	return
}
