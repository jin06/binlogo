package dao

import (
	"context"
	"errors"
	"fmt"

	"github.com/jin06/binlogo/v2/configs"
	"github.com/jin06/binlogo/v2/pkg/etcdclient"
	"github.com/jin06/binlogo/v2/pkg/store/model"
	"github.com/jin06/binlogo/v2/pkg/store/model/node"
	"github.com/jin06/binlogo/v2/pkg/store/model/pipeline"
	clientv3 "go.etcd.io/etcd/client/v3"
)

var myDao Dao

func Init() {
	myDao = NewDaoRedis()
}

type Dao interface {
	GetInstance(ctx context.Context, pname string) (ins *pipeline.Instance, err error)
	AllInstance(ctx context.Context) (all []*pipeline.Instance, err error)
	AllInstanceMap(ctx context.Context) (maps map[string]*pipeline.Instance, err error)
	GetNode(ctx context.Context, name string) (*node.Node, error)
	AllNodes(ctx context.Context) (list []*node.Node, err error)
	UpdateNode(ctx context.Context, name string, opts ...node.NodeOption) (bool, error)
	UpdateNodeIP(ctx context.Context, name string, ip string) (ok bool, err error)
	UpdateCapacity(ctx context.Context, cap *node.Capacity) (bool, error)
	RefreshNode(ctx context.Context, n *node.Node) (success bool, err error)
	AllCapacity(ctx context.Context) (list []*node.Capacity, err error)
	CapacityMap(ctx context.Context) (mapping map[string]*node.Capacity, err error)
	AllStatus(ctx context.Context) (list []*node.Status, err error)
	StatusMap(ctx context.Context) (mapping map[string]*node.Status, err error)
	LeaderNode(ctx context.Context) (node string, err error)
}

func Prefix() string {
	return fmt.Sprintf("/%s/%s", configs.APP, configs.Default.ClusterName)
}

// ClearOrDeleteBind clear or delete pipeline bind
// sets pipeline bind to blank if pipeline is expected to run, so pipeline will be scheduled to a node later
// delete pipeline bind if pipeline is expected to stop, pipeline will not be scheduled a node.
func ClearOrDeleteBind(name string) (err error) {
	if name == "" {
		err = errors.New("empty name")
		return
	}
	res, err := etcdclient.Default().Get(context.TODO(), PipeBindPrefix())
	if err != nil {
		return
	}
	pb := model.EmptyPipelineBind()
	var revision int64
	if len(res.Kvs) > 0 {
		revision = res.Kvs[0].CreateRevision
		err = pb.Unmarshal(res.Kvs[0].Value)
	}
	pipe, err := GetPipeline(name)
	if err != nil {
		return
	}
	if pipe.ExpectRun() {
		pb.Bindings[name] = ""
	} else {
		delete(pb.Bindings, name)
	}
	txn := etcdclient.Default().Txn(context.TODO()).If(clientv3.Compare(clientv3.CreateRevision(PipeBindPrefix()), "=", revision))
	txn = txn.Then(clientv3.OpPut(PipeBindPrefix(), pb.Val()))
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
