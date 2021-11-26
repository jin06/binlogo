package dao_node

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/jin06/binlogo/pkg/etcd_client"
	"github.com/jin06/binlogo/pkg/store/model/node"
)

// AllocatablePrefix returns prefix key of etcd for allocatable
func AllocatablePrefix() string {
	return etcd_client.Prefix() + "/node/allocatable"
}

// UpdateAllocatable update allocatable data to etcd
func UpdateAllocatable(al *node.Allocatable) (err error) {
	if al.NodeName == "" {
		return errors.New("empty name")
	}
	key := AllocatablePrefix() + "/" + al.NodeName
	if key == "" {
		err = errors.New("empty key")
		return
	}
	b, err := json.Marshal(al)
	if err != nil {
		return
	}
	_, err = etcd_client.Default().Put(context.Background(), key, string(b))
	return
}

// DeleteAllocatable delete allocatable data in etcd
func DeleteAllocatable(nodeName string) (ok bool, err error) {
	if nodeName == "" {
		err = errors.New("empty name")
		return
	}
	key := AllocatablePrefix() + "/" + nodeName
	res, err := etcd_client.Default().Delete(context.Background(), key)
	if err != nil {
		return
	}
	if res.Deleted > 0 {
		ok = true
	}
	return
}
