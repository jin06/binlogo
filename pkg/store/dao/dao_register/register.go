package dao_register

import (
	"github.com/jin06/binlogo/pkg/etcd_client"
)

func PipePrefix() string {
	return etcd_client.Prefix() + "/pipeline/register"
}

func PipeInstancePrefix() string {
	return etcd_client.Prefix() + "/pipeline/instance"
}

