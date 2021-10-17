package dao

import (
	"context"
	"github.com/jin06/binlogo/pkg/blog"
	"github.com/jin06/binlogo/pkg/store/etcd"
	"github.com/jin06/binlogo/pkg/store/model"
	"go.etcd.io/etcd/clientv3"
)

func AllNodes() (list []*model.Node, err error){

	list = []*model.Node{}
	key := etcd.E.Prefix + "/nodes"
	res, err := etcd.E.Client.Get(context.Background(), key, clientv3.WithPrefix(), clientv3.WithFromKey())
	if err != nil {
		return
	}
	if len(res.Kvs ) == 0 {
		return
	}
	for _, v := range res.Kvs {
		ele := &model.Node{}
		er := ele.Unmarshal(v.Value)
		if er != nil {
			blog.Error(er)
			continue
		}
		list = append(list, ele)
	}
	return
}
