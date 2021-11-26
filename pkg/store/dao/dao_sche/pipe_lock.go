package dao_sche

import (
	"github.com/jin06/binlogo/pkg/etcd_client"
)

// PipelineLockPrefix returns etcd prefix of lock
func PipelineLockPrefix() string {
	return etcd_client.Prefix() + "/scheduler/pipeline_lock"
}
