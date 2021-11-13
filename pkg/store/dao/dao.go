package dao

import (
	"context"
	"errors"
	"github.com/coreos/etcd/clientv3"
	"github.com/jin06/binlogo/pkg/etcd_client"
	"github.com/jin06/binlogo/pkg/store/dao/dao_pipe"
	"github.com/jin06/binlogo/pkg/store/dao/dao_sche"
	"github.com/jin06/binlogo/pkg/store/model/scheduler"
)

func ClearOrDeleteBind(name string) (err error) {
	if name == "" {
		err = errors.New("empty name")
		return
	}
	res, err := etcd_client.Default().Get(context.TODO(), dao_sche.PipeBindPrefix())
	if err != nil {
		return
	}
	pb := scheduler.EmptyPipelineBind()
	var revision int64
	if len(res.Kvs) > 0 {
		revision = res.Kvs[0].CreateRevision
		err = pb.Unmarshal(res.Kvs[0].Value)
	}
	pipe, err := dao_pipe.GetPipeline(name)
	if err != nil {
		return
	}
	if pipe.ExpectRun() {
		pb.Bindings[name] = ""
	} else {
		delete(pb.Bindings, name)
	}
	txn := etcd_client.Default().Txn(context.TODO()).If(clientv3.Compare(clientv3.CreateRevision(dao_sche.PipeBindPrefix()), "=", revision))
	txn = txn.Then(clientv3.OpPut(dao_sche.PipeBindPrefix(), pb.Val()))
	_, err = txn.Commit()
	if err != nil {
		return
	}
	return
}
