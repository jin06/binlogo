package dao_node

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/jin06/binlogo/pkg/store/etcd"
	"github.com/jin06/binlogo/pkg/store/model/node"
)

func AllocatablePrefix() string {
	return etcd.Prefix() + "/node/allocatable"
}
func UpdateAllocatable(al *node.Allocatable, args ...Option) (err error) {
	opts := GetOptions(args...)
	key := opts.Key
	if key == "" {
		if opts.NodeName != "" {
			key = AllocatablePrefix() + "/" + opts.NodeName
		}
	}
	if key == "" {
		err = errors.New("empty key")
		return
	}
	ctx, cancel := context.WithTimeout(context.TODO(), etcd.E.Timeout)
	defer cancel()
	b, err := json.Marshal(al)
	if err != nil {
		return
	}
	_, err = etcd.E.Client.Put(ctx, key, string(b))
	return

}
