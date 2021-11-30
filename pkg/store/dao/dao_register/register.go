package dao_register

import (
	"github.com/jin06/binlogo/pkg/etcdclient"
)

//func PipePrefix() string {
//	return etcdclient.Prefix() + "/pipeline/register"
//}

// PipeInstancePrefix returns etcd prefix of pipeline instance
func PipeInstancePrefix() string {
	return etcdclient.Prefix() + "/pipeline/instance"
}
