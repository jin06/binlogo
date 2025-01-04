package storeredis

import (
	"fmt"

	"github.com/jin06/binlogo/v2/configs"
)

const (
	NodeKey         = "node"
	CapacityKey     = "capacity"
	StatusKey       = "status"
	ClusterKey      = "cluster"
	ElectionKey     = "election"
	AllocatableKey  = "allocatable"
	RegisterKey     = "register"
	MasterKey       = "master"
	PipelineBindKey = "pipeline_bind"
)

func Prefix() string {
	return fmt.Sprintf("/%s/%s", configs.APP, configs.Default.ClusterName)
}

func GetRegisterKey(s string) string {
	return fmt.Sprintf("%s/%s", RegisterPrefix(), s)
}

func RegisterPrefix() string {
	return fmt.Sprintf("%s/%s", NodePrefix(), RegisterKey)
}

func ClusterPrefix() string {
	return fmt.Sprintf("%s/%s", Prefix(), ClusterKey)
}

func ElectionPrefix() string {
	return fmt.Sprintf("%s/%s", ClusterPrefix(), ElectionKey)
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
	return fmt.Sprintf("%s/%s", ClusterPrefix(), PipelineBindKey)
}
