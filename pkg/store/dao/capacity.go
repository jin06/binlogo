package dao

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/jin06/binlogo/pkg/store/etcd"
	"github.com/jin06/binlogo/pkg/store/model/node"
)

func UpdateCapacity(cap *node.Capacity, args ...Option) (err error) {
	opts := GetOptions(args...)
	key := opts.Key
	if key == "" {
		err = errors.New("empty key")
		return
	}
	ctx, cancel := context.WithTimeout(context.TODO(), etcd.E.Timeout)
	defer cancel()
	b, err := json.Marshal(cap)
	if err != nil {
		return
	}
	_, err = etcd.E.Client.Put(ctx, key, string(b))
	return
}
