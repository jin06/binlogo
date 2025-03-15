package dao

import (
	"context"
	"time"

	"github.com/jin06/binlogo/v2/pkg/store/model"
	"github.com/jin06/binlogo/v2/pkg/store/model/node"
	"github.com/jin06/binlogo/v2/pkg/store/model/pipeline"
)

var myDao Dao

func Init() {
	myDao = NewDaoRedis()
}

type Dao interface {
	GetInstance(ctx context.Context, pname string) (ins *pipeline.Instance, err error)
	AllInstance(ctx context.Context) (all []*pipeline.Instance, err error)
	AllInstanceMap(ctx context.Context) (maps map[string]*pipeline.Instance, err error)
	// Compete to become the master node of the cluster.
	AcquireMasterLock(ctx context.Context, node *node.Node) error
	// Find out who the current master node is.
	GetMasterLock(ctx context.Context) (string, error)
	// Renew the lease of the master node.
	LeaseMasterLock(ctx context.Context) error
	RegisterNode(ctx context.Context, n *node.Node) (bool, error)
	LeaseNode(ctx context.Context, n *node.Node) error
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
	DeleteStatus(ctx context.Context, name string) error
	CreateOrUpdateStatus(ctx context.Context, nodeName string, conditions node.StatusConditions) (ok bool, err error)
	GetStatus(ctx context.Context, nodeName string) (s *node.Status, err error)
	LeaderNode(ctx context.Context) (node string, err error)
	UpdateAllocatable(ctx context.Context, al *node.Allocatable) (ok bool, err error)
	AllElections() (res []map[string]any, err error)
	GetPipelineBind(ctx context.Context) (*model.PipelineBind, error)
	UpdatePipelineBindIfNotExist(ctx context.Context, pName string, nName string) error
	UpdatePipelineBind(ctx context.Context, pName string, nName string) (bool, error)
	DeletePipelineBind(ctx context.Context, pName string) (ok bool, err error)
	GetPipeline(ctx context.Context, name string) (p *pipeline.Pipeline, err error)
	UpdatePipeline(ctx context.Context, name string, opts ...pipeline.OptionPipeline) (err error)
	AllPipelines(ctx context.Context) (list []*pipeline.Pipeline, err error)
	AllPipelinesMap(ctx context.Context) (mapping map[string]*pipeline.Pipeline, err error)
	ClearOrDeleteBind(ctx context.Context, name string) (err error)
	UpdateEvent(ctx context.Context, e *model.Event) error
	GetPosition(ctx context.Context, pipe string) (p *pipeline.Position, err error)
	UpdatePosition(ctx context.Context, p *pipeline.Position) error
	DeletePosition(ctx context.Context, name string) (err error)
	UpdateRecord(ctx context.Context, p *pipeline.RecordPosition) (err error)
	UpdateRecordSafe(ctx context.Context, pipe string, opts ...pipeline.OptionRecord) (err error)
	GetRecord(ctx context.Context, pipe string) (r *pipeline.RecordPosition, err error)
	RegisterInstance(ctx context.Context, ins *pipeline.Instance, exp time.Duration) error
	UnRegisterInstance(ctx context.Context, pipe string, n string) error
	LeaseInstance(ctx context.Context, pipe string, exp time.Duration) error
}

// ClearOrDeleteBind clear or delete pipeline bind
// sets pipeline bind to blank if pipeline is expected to run, so pipeline will be scheduled to a node later
// delete pipeline bind if pipeline is expected to stop, pipeline will not be scheduled a node.
func ClearOrDeleteBind(ctx context.Context, name string) (err error) {
	return myDao.ClearOrDeleteBind(ctx, name)
}
