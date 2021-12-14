package dao

import (
	"context"
	"errors"

	clientv3 "go.etcd.io/etcd/client/v3"
	"github.com/jin06/binlogo/pkg/etcdclient"
	"github.com/jin06/binlogo/pkg/store/dao/dao_pipe"
	"github.com/jin06/binlogo/pkg/store/dao/dao_sche"
	"github.com/jin06/binlogo/pkg/store/model/scheduler"
)

// ClearOrDeleteBind clear or delete pipeline bind
// sets pipeline bind to blank if pipeline is expected to run, so pipeline will be scheduled to a node later
// delete pipeline bind if pipeline is expected to stop, pipeline will not be scheduled a node.
func ClearOrDeleteBind(name string) (err error) {
	if name == "" {
		err = errors.New("empty name")
		return
	}
	res, err := etcdclient.Default().Get(context.TODO(), dao_sche.PipeBindPrefix())
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
	txn := etcdclient.Default().Txn(context.TODO()).If(clientv3.Compare(clientv3.CreateRevision(dao_sche.PipeBindPrefix()), "=", revision))
	txn = txn.Then(clientv3.OpPut(dao_sche.PipeBindPrefix(), pb.Val()))
	_, err = txn.Commit()
	if err != nil {
		return
	}
	return
}

// DeleteCluster delete whole cluster
func DeleteCluster(clusterName string) (deleted int64, err error) {
	if clusterName == "" {
		err = errors.New("empty cluster name")
	}
	key := etcdclient.Prefix()
	res, err := etcdclient.Default().Delete(context.Background(), key)
	if err != nil {
		return
	}
	deleted = res.Deleted
	return
}
