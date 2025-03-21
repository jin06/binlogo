package model

import (
	"encoding/json"
	"fmt"
	"net"
	"time"

	"github.com/jin06/binlogo/v2/configs"
	"github.com/jin06/binlogo/v2/pkg/util/key"
)

// Event generated by pipelines and nodes in the cluster.
// Events include common information and error information,
// which can be easily located when an error occurs
type Event struct {
	K            string       `json:"key" redis:"key"`
	Type         Type         `json:"type" redis:"type"`
	ResourceType ResourceType `json:"resource_type" redis:"resource_type"`
	ResourceName string       `json:"resource_name" redis:"resource_name"`
	FirstTime    time.Time    `json:"first_time" redis:"first_time"`
	LastTime     time.Time    `json:"last_time" redis:"last_time"`
	Count        int          `json:"count" redis:"count"`
	MessageType  MessageType  `json:"type_message" redis:"type_message"`
	Message      string       `json:"message" redis:"message"`
	NodeName     string       `json:"node_name" redis:"node_name"`
	NodeIP       net.IP       `json:"node_ip" redis:"node_ip"`
}

func (e *Event) Key() string {
	return fmt.Sprintf("event/%s", e.K)
}

func (e *Event) Val() string {
	b, _ := json.Marshal(e)
	return string(b)
}

func (e *Event) Unmarshal(data []byte) error {
	return json.Unmarshal(data, e)
}

// ResourceType pipeline or node or etc.
type ResourceType = string

const (
	// PIPELINE resource type
	PIPELINE ResourceType = "pipeline"
	// ResourceType resource type
	NODE ResourceType = "node"
	// CLUSTER resource type
	CLUSTER ResourceType = "cluster"
)

// Type event level
type Type = string

const (
	// WARN level
	WARN Type = "warn"
	// ERROR level
	ERROR Type = "error"
	// INFO level
	INFO Type = "info"
)

// MessageType todo
type MessageType = string

const (
	// MSG_EMPTY default type
	MSG_EMPTY MessageType = ""
)

func baseEvent() (e *Event) {
	e = &Event{
		K:        key.GetKey(),
		NodeName: configs.GetNodeName(),
		NodeIP:   configs.NodeIP,
	}
	return
}

func newNode(msg string) (e *Event) {
	e = baseEvent()
	e.ResourceType = NODE
	e.ResourceName = configs.GetNodeName()
	e.Message = msg
	return
}

func newPipeline(name string, msg string) (e *Event) {
	e = baseEvent()
	e.ResourceType = PIPELINE
	e.ResourceName = name
	e.Message = msg
	return
}

func newCluster(msg string) (e *Event) {
	e = baseEvent()
	e.ResourceType = CLUSTER
	e.ResourceName = configs.Default.ClusterName
	e.Message = msg
	return
}

// NewInfoNode returns a new node event of info level
func NewInfoNode(msg string) (e *Event) {
	e = newNode(msg)
	e.Type = INFO
	return
}

// NewInfoPipeline returns a new pipeline info event
func NewInfoPipeline(name string, msg string) (e *Event) {
	e = newPipeline(name, msg)
	e.Type = INFO
	return
}

// NewInfoCluster returns a new cluster info event
func NewInfoCluster(msg string) (e *Event) {
	e = newCluster(msg)
	e.Type = INFO
	return
}

// NewWarnPipeline returns a new pipeline warn event
func NewWarnPipeline(name string, msg string) (e *Event) {
	e = newPipeline(name, msg)
	e.Type = WARN
	return
}

// NewWarnNode returns a new node warn event
func NewWarnNode(msg string) (e *Event) {
	e = newNode(msg)
	e.Type = WARN
	return
}

// NewWarnCluster returns a new cluster warn event
func NewWarnCluster(msg string) (e *Event) {
	e = newCluster(msg)
	e.Type = WARN
	return
}

// NewErrorPipeline returns a new pipeline error event
func NewErrorPipeline(name string, msg string) (e *Event) {
	e = newPipeline(name, msg)
	e.Type = ERROR
	return
}

// NewErrorNode returns a new node error event
func NewErrorNode(msg string) (e *Event) {
	e = newNode(msg)
	e.Type = ERROR
	return
}

// NewErrorCluster returns a new cluster error event
func NewErrorCluster(msg string) (e *Event) {
	e = newCluster(msg)
	e.Type = ERROR
	return
}
