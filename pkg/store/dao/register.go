package dao

import (
	"github.com/jin06/binlogo/v2/pkg/etcdclient"
)

// PipeInstancePrefix returns etcd prefix of pipeline instance
func PipeInstancePrefix() string {
	return etcdclient.Prefix() + "/pipeline/instance"
}
