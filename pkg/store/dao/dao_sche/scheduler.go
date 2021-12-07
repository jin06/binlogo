package dao_sche

import (
	"context"
	"github.com/coreos/etcd/clientv3"
	"github.com/jin06/binlogo/pkg/etcdclient"
	"github.com/jin06/binlogo/pkg/store/model/scheduler"
)

// PipeBindPrefix returns etcd prefix of pipeline bind
func PipeBindPrefix() string {
	return etcdclient.Prefix() + "/scheduler/pipeline_bind"
}

// SchedulerPrefix etcd prefix of scheduler
func SchedulerPrefix() string {
	return etcdclient.Prefix() + "/scheduler"
}

// GetPipelineBind get pipeline bind from etcd
func GetPipelineBind() (pb *scheduler.PipelineBind, err error) {
	res, err := etcdclient.Default().Get(context.TODO(), PipeBindPrefix())
	if err != nil {
		return
	}
	pb = scheduler.EmptyPipelineBind()
	if len(res.Kvs) == 0 {
		return
	}
	err = pb.Unmarshal(res.Kvs[0].Value)
	return
}

// UpdatePipelineBindIfNotExist update pipeline bind if the pipeline not exist in pipeline bind
func UpdatePipelineBindIfNotExist(pName string, nName string) (err error) {
	res, err := etcdclient.Default().Get(context.TODO(), PipeBindPrefix())
	if err != nil {
		return
	}
	pb := scheduler.EmptyPipelineBind()
	var revision int64
	if len(res.Kvs) > 0 {
		revision = res.Kvs[0].CreateRevision
		err = pb.Unmarshal(res.Kvs[0].Value)
	}
	if _, ok := pb.Bindings[pName]; ok {
		return
	}
	pb.Bindings[pName] = nName
	txn := etcdclient.Default().Txn(context.TODO()).
		If(clientv3.Compare(clientv3.CreateRevision(PipeBindPrefix()), "=", revision)).
		Then(clientv3.OpPut(PipeBindPrefix(), pb.Val()))
	_, err = txn.Commit()
	if err != nil {
		return
	}
	return
}

// UpdatePipelineBind update pipeline bind ignore exists
func UpdatePipelineBind(pName string, nName string) (ok bool, err error) {
	res, err := etcdclient.Default().Get(context.TODO(), PipeBindPrefix())
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
	txn := etcdclient.Default().Txn(context.TODO()).
		If(clientv3.Compare(clientv3.CreateRevision(PipeBindPrefix()), "=", revision)).
		Then(clientv3.OpPut(PipeBindPrefix(), pb.Val()))
	resp, err := txn.Commit()
	if err != nil {
		return
	}
	ok = resp.Succeeded
	return
}

// DeletePipelineBind delete pipeline bind in etcd
func DeletePipelineBind(pName string) (ok bool, err error) {
	res, err := etcdclient.Default().Get(context.TODO(), PipeBindPrefix())
	if err != nil {
		return
	}
	pb := scheduler.EmptyPipelineBind()
	var revision int64
	if len(res.Kvs) > 0 {
		revision = res.Kvs[0].CreateRevision
		err = pb.Unmarshal(res.Kvs[0].Value)
		if err != nil {
			return
		}
	}
	delete(pb.Bindings, pName)
	txn := etcdclient.Default().Txn(context.TODO()).
		If(clientv3.Compare(clientv3.CreateRevision(PipeBindPrefix()), "=", revision)).
		Then(clientv3.OpPut(PipeBindPrefix(), pb.Val()))
	resp, err := txn.Commit()
	if err != nil {
		return
	}
	ok = resp.Succeeded
	return
}
