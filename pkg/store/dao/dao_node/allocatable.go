package dao_node

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/jin06/binlogo/pkg/etcd_client"
	"github.com/jin06/binlogo/pkg/store/model/node"
)

func AllocatablePrefix() string {
	return etcd_client.Prefix() + "/node/allocatable"
}
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
