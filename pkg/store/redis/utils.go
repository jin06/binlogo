package storeredis

import (
	"fmt"

	"github.com/jin06/binlogo/v2/configs"
)

const (
	NodeKey             = "node"
	CapacityKey         = "capacity"
	StatusKey           = "status"
	ClusterKey          = "cluster"
	ElectionKey         = "election"
	AllocatableKey      = "allocatable"
	RegisterKey         = "register"
	MasterKey           = "master"
	PipelineBindKey     = "pipeline_bind"
	ChannelKey          = "channel"
	PipelinesKey        = "pipelines"
	PipelineInstanceKey = "pipeline_instance"
	EventKey            = "event"
	PositionKey         = "position"
	RecordPostion       = "record_position"
)

func Prefix() string {
	return fmt.Sprintf("/%s/%s", configs.APP, configs.Default.ClusterName)
}

func GetRegisterKey(s string) string {
	return fmt.Sprintf("%s/%s", RegisterPrefix(), s)
}

func GetStatusKey(s string) string {
	return fmt.Sprintf("%s/%s", StatusPrefix(), s)
}

func RegisterPrefix() string {
	return fmt.Sprintf("%s/%s", NodePrefix(), RegisterKey)
}

func NodePrefix() string {
	return fmt.Sprintf("%s/%s", Prefix(), NodeKey)
}

func CapacityPrefix() string {
	return fmt.Sprintf("%s/%s", NodePrefix(), CapacityKey)
}

func StatusPrefix() string {
	return fmt.Sprintf("%s/%s", NodePrefix(), StatusKey)
}

// AllocatablePrefix returns prefix key of redis for allocatable
func AllocatablePrefix() string {
	return fmt.Sprintf("%s/%s", NodePrefix(), AllocatableKey)
}

func MasterPreifx() string {
	return fmt.Sprintf("%s/%s", NodePrefix(), MasterKey)
}

func PipelineBindPrefix() string {
	return fmt.Sprintf("%s/%s", Prefix(), PipelineBindKey)
}

func PipelineBindChan() string {
	return fmt.Sprintf("%s/%s", Prefix(), ChannelKey)
}

func NodeChan() string {
	return fmt.Sprintf("%s/%s", NodePrefix(), ChannelKey)
}

func PipelinesPrefix() string {
	return fmt.Sprintf("%s/%s", Prefix(), PipelinesKey)
}

func PipelineInstancePrefix() string {
	return fmt.Sprintf("%s/%s", Prefix(), PipelineInstanceKey)
}

func GetPipeInstanceKey(s string) string {
	return fmt.Sprintf("%s/%s", PipelineInstancePrefix(), s)
}

func EventPrefix() string {
	return fmt.Sprintf("%s/%s", Prefix(), EventKey)
}

func GetEventKey(s string) string {
	return fmt.Sprintf("%s/%s", EventPrefix(), s)
}

func PositionPrefix() string {
	return fmt.Sprintf("%s/%s", Prefix(), PositionKey)
}

func GetPositionKey(s string) string {
	return fmt.Sprintf("%s/%s", PositionPrefix(), s)
}

func RecordPositionPrefix() string {
	return fmt.Sprintf("%s/%s", Prefix(), RecordPostion)
}

func GetRecordPositionKey(s string) string {
	return fmt.Sprintf("%s/%s", RecordPositionPrefix(), s)
}
