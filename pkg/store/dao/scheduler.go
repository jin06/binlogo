package dao

import (
	"context"
	"fmt"
	"github.com/jin06/binlogo/pkg/blog"
	"github.com/jin06/binlogo/pkg/store/etcd"
	"github.com/jin06/binlogo/pkg/store/model/scheduler"
	"go.etcd.io/etcd/clientv3"
)

//func UpdatePipelineBind(pb *scheduler.PipelineBind) (err error) {
//	ok, err := etcd.Update(pb)
//	if err != nil {
//		return
//	}
//	if !ok {
//		err = errors.New("update pipeline bind fail")
//	}
//	return
//}

func UpdatePipelineBind(pName string, nName string) (err error) {
	pb := scheduler.NewPipelineBindH()
	c := etcd.E.Client
	_ , err = etcd.E.GetH(pb)
	if err != nil {
		blog.Error(err)
	}
	pb.PipelineBind.Bindings[pName] = nName
	txn := c.Txn(context.Background()).If(clientv3.Compare(clientv3.CreateRevision(pb.Key()), "=", pb.Revision))
	txn = txn.Then(clientv3.OpPut(pb.Key(), pb.Val()))
	resp , err := txn.Commit()
	if err != nil {
		return
	}
	fmt.Println(resp.Succeeded)

	return
}
