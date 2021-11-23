package dao_sche

import "github.com/jin06/binlogo/pkg/store/etcd"

// PipelineLockPrefix returns etcd prefix of lock
func PipelineLockPrefix() string {
	return etcd.E.Prefix + "/scheduler/pipeline_lock"
}
