package dao_cluster

import (
	"context"
	"errors"

	"github.com/jin06/binlogo/v2/pkg/etcdclient"
)

// ElectionPrefix returns etcd prefix
func ElectionPrefix() string {
	return etcdclient.Prefix() + "/cluster/election"
}

// GetElection get election value
func GetElection(leaseId string) (res *string, err error) {
	if leaseId == "" {
		err = errors.New("empty lease id")
		return
	}
	key := ElectionPrefix() + "/" + leaseId
	resp, err := etcdclient.Default().Get(context.Background(), key)
	if err != nil {
		return
	}
	if len(resp.Kvs) == 0 {
		return
	}
	val := string(resp.Kvs[0].Value)
	res = &val
	return
}
