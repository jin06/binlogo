package dao_sche

import (
	"context"
	"fmt"
	"github.com/jin06/binlogo/configs"
	"github.com/jin06/binlogo/pkg/blog"
	"github.com/jin06/binlogo/pkg/store/etcd"
	"github.com/jin06/binlogo/pkg/store/model/scheduler"
	"github.com/spf13/viper"
	"go.etcd.io/etcd/clientv3"
)

func pipeBindPrefix() string {
	return fmt.Sprintf("/%s/%s/scheduler/pipeline_bind", configs.APP, viper.GetString("cluster.name"))
}

func GetPipelineBind() (pb *scheduler.PipelineBind, err error) {
	res, err := etcd.E.Client.Get(context.TODO(), pipeBindPrefix() )
	if err != nil {
		return
	}
	pb = scheduler.EmptyPipelineBind()
	for _, v := range res.Kvs {
		if err = pb.Unmarshal(v.Value); err != nil {
			return
		}
		break
	}
	return
}

func UpdatePipelineBindIfNotExist(pName string, nName string) (err error) {
	pb := scheduler.NewPipelineBindH()
	c := etcd.E.Client
	_, err = etcd.E.GetH(pb)
	if err != nil {
		blog.Error(err)
		return
	}
	if _, ok := pb.PipelineBind.Bindings[pName]; ok {
		return
	}
	pb.PipelineBind.Bindings[pName] = nName
	key := pipeBindPrefix()
	txn := c.Txn(context.Background()).If(clientv3.Compare(clientv3.CreateRevision(key), "=", pb.Revision))
	txn = txn.Then(clientv3.OpPut(key, pb.Val()))
	_, err = txn.Commit()

	return
}

func UpdatePipelineBind(pName string, nName string) (ok bool, err error) {
	res, err := etcd.E.Client.Get(context.Background(), pipeBindPrefix())
	if err != nil {
		return
	}
	pb := scheduler.EmptyPipelineBind()
	var revision int64
	if len(res.Kvs) > 0 {
		revision = res.Kvs[0].CreateRevision
		err = pb.Unmarshal(res.Kvs[0].Value)
	}
	pb.Bindings[pName] = nName
	txn := etcd.E.Client.Txn(context.Background()).If(clientv3.Compare(clientv3.CreateRevision(pipeBindPrefix()), "=", revision))
	txn = txn.Then(clientv3.OpPut(pipeBindPrefix(), pb.Val()))
	resp, err := txn.Commit()
	if err != nil {
		return
	}
	ok = resp.Succeeded
	return
}

func DeletePipelineBind(pName string) (ok bool, err error) {
	res, err := etcd.E.Client.Get(context.Background(), pipeBindPrefix())
	if err != nil {
		return
	}
	pb := scheduler.EmptyPipelineBind()
	var revision int64
	for _, v := range res.Kvs {
		err = pb.Unmarshal(v.Value)
		if err != nil {
			return
		}
		revision = v.CreateRevision
		break
	}
	delete(pb.Bindings, pName)
	txn := etcd.E.Client.Txn(context.Background()).If(clientv3.Compare(clientv3.CreateRevision(pipeBindPrefix()), "=", revision))
	txn = txn.Then(clientv3.OpPut(pipeBindPrefix(), pb.Val()))
	resp, err := txn.Commit()
	if err != nil {
		return
	}
	ok = resp.Succeeded
	return
}
