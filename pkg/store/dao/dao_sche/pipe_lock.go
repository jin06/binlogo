package dao_sche

import "github.com/jin06/binlogo/pkg/store/etcd"

func PipelineLockPrefix() string {
	return etcd.E.Prefix + "/scheduler/pipeline_lock"
}
