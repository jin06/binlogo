package dao_sche

import (
	"context"
	"github.com/coreos/etcd/clientv3"
	"github.com/jin06/binlogo/pkg/etcd_client"
	"github.com/jin06/binlogo/pkg/store/etcd"
	"github.com/jin06/binlogo/pkg/store/model/scheduler"
)

func PipeBindPrefix() string {
	return etcd.Prefix() + "/scheduler/pipeline_bind"
}

func GetPipelineBind() (pb *scheduler.PipelineBind, err error) {
	res, err := etcd.E.Client.Get(context.TODO(), PipeBindPrefix() )
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

func UpdatePipelineBindIfNotExist(pName string, nName string) (err error) {
	res, err := etcd_client.Default().Get(context.TODO(), PipeBindPrefix())
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
	txn := etcd_client.Default().Txn(context.TODO()).
		If(clientv3.Compare(clientv3.CreateRevision(PipeBindPrefix()), "=", revision)).
		Then(clientv3.OpPut(PipeBindPrefix(), pb.Val()))
	_, err = txn.Commit()
	if err != nil {
		return
	}
	return
}

func UpdatePipelineBind(pName string, nName string) (ok bool, err error) {
	res, err := etcd_client.Default().Get(context.TODO(), PipeBindPrefix())
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
	txn := etcd_client.Default().Txn(context.TODO()).
		If(clientv3.Compare(clientv3.CreateRevision(PipeBindPrefix()), "=", revision)).
		Then(clientv3.OpPut(PipeBindPrefix(), pb.Val()))
	resp, err := txn.Commit()
	if err != nil {
		return
	}
	ok = resp.Succeeded
	return
}

func DeletePipelineBind(pName string) (ok bool, err error) {
	res, err := etcd_client.Default().Get(context.TODO(), PipeBindPrefix())
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
	txn := etcd_client.Default().Txn(context.TODO()).
		If(clientv3.Compare(clientv3.CreateRevision(PipeBindPrefix()), "=", revision)).
		Then(clientv3.OpPut(PipeBindPrefix(), pb.Val()))
	resp, err := txn.Commit()
	if err != nil {
		return
	}
	ok = resp.Succeeded
	return
}


