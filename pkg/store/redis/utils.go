package store_redis

import (
	"fmt"

	"github.com/jin06/binlogo/v2/configs"
)

const (
	NodeKey     = "node"
	CapacityKey = "node_capacity"
	StatusKey   = "status"
	ClusterKey  = "cluster"
	ElectionKey = "election"
)

func Prefix() string {
	return fmt.Sprintf("/%s/%s", configs.APP, configs.Default.ClusterName)
}

func NodePrefix() string {
	return fmt.Sprintf("%s/%s", Prefix(), NodeKey)
}

func CapacityPrefix() string {
	return fmt.Sprintf("%s/%s", Prefix(), CapacityKey)
}

func StatusPrefix() string {
	return fmt.Sprintf("%s/%s", Prefix(), StatusKey)
}

func ClusterPrefix() string {
	return fmt.Sprintf("%s/%s", Prefix(), ClusterKey)
}

func ElectionPrefix() string {
	return fmt.Sprintf("%s/%s", ClusterPrefix(), ElectionKey)
}
