package dao

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/jin06/binlogo/v2/configs"
	"github.com/jin06/binlogo/v2/pkg/store/model/node"
	"github.com/jin06/binlogo/v2/pkg/store/model/pipeline"
	storeredis "github.com/jin06/binlogo/v2/pkg/store/redis"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

func NewDaoRedis() *DaoRedis {
	return &DaoRedis{
		store:       storeredis.Default,
		registerTTL: time.Second * 5,
		masterTTL:   time.Second * 5,
	}
}

type DaoRedis struct {
	store       *storeredis.Redis
	registerTTL time.Duration
	masterTTL   time.Duration
}

func (d *DaoRedis) client() *redis.Client {
	return d.store.GetClient()
}

func (d *DaoRedis) prefix() string {
	return "/binlogo"
}

func (d *DaoRedis) clusterPrefix() string {
	return d.prefix() + "/" + configs.Default.ClusterName
}

func (d *DaoRedis) instancePrefix() string {
	return d.clusterPrefix() + "/pipeline/instance"
}

func (d *DaoRedis) instanceKey(name string) string {
	return fmt.Sprintf("%s/%s", d.instancePrefix(), name)
}

func scanKeysWithPrefix(ctx context.Context, client *redis.Client, prefix string) ([]string, error) {
	var cursor uint64 = 0
	var keys []string
	for {
		var scanResult []string
		var err error
		scanResult, cursor, err = client.Scan(ctx, cursor, prefix+"*", 0).Result()
		if err != nil {
			return nil, err
		}
		keys = append(keys, scanResult...)
		if cursor == 0 {
			break
		}
	}
	return keys, nil
}

func getAllHashDatas[T any](ctx context.Context, client *redis.Client, prefix string) (list []T, err error) {
	var keys []string
	list = []T{}
	if keys, err = scanKeysWithPrefix(ctx, client, prefix); err != nil {
		return
	}
	for _, key := range keys {
		cmd := client.HGetAll(ctx, key)
		if cmd.Err() != nil {
			logrus.Error("get hashdata error ", cmd.Err())
			continue
		}
		var item T
		if err := cmd.Scan(item); err != nil {
			logrus.Error("scan hashdata error ", err)
			continue
		}
		list = append(list, item)
	}
	return
}

func (d *DaoRedis) GetInstance(ctx context.Context, pname string) (ins *pipeline.Instance, err error) {
	ins = &pipeline.Instance{}
	cmd := d.client().Get(ctx, d.instanceKey(pname))
	if err = cmd.Err(); err != nil {
		return
	}
	err = cmd.Scan(ins)
	return
}

func (d *DaoRedis) AllInstance(ctx context.Context) (list []*pipeline.Instance, err error) {
	return getAllHashDatas[*pipeline.Instance](ctx, d.client(), d.instancePrefix())
}

func (d *DaoRedis) AllInstanceMap(ctx context.Context) (maps map[string]*pipeline.Instance, err error) {
	maps = map[string]*pipeline.Instance{}
	list, err := d.AllInstance(ctx)
	if err != nil {
		return
	}
	for _, v := range list {
		maps[v.PipelineName] = v
	}
	return
}

// Compete to become the master node of the cluster by obtaining a key.
// The one who gets it first will become the master node.
func (d *DaoRedis) AcquireMasterLock(ctx context.Context, node *node.Node) error {
	return d.store.GetClient().SetNX(ctx, storeredis.MasterPreifx(), node.Name, d.masterTTL).Err()
}

func (d *DaoRedis) GetMasterLock(ctx context.Context) (string, error) {
	val, err := d.store.GetClient().Get(ctx, storeredis.MasterPreifx()).Result()
	if err == redis.Nil {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return val, nil
}

func (d *DaoRedis) LeaseMasterLock(ctx context.Context) error {
	return d.store.GetClient().Expire(ctx, storeredis.MasterPreifx(), d.masterTTL).Err()
}

func (d *DaoRedis) RegisterNode(ctx context.Context, n *node.Node) (bool, error) {
	return d.store.GetClient().SetNX(ctx, storeredis.GetRegisterKey(n.Name), n.Name, d.registerTTL).Result()
}

func (d *DaoRedis) LeaseNode(ctx context.Context, n *node.Node) error {
	return d.store.GetClient().SetEx(ctx, storeredis.GetRegisterKey(n.Name), n.Name, d.registerTTL).Err()
}

func (d *DaoRedis) CreateNodeIfNotExist(ctx context.Context, n *node.Node) (err error) {
	if n == nil {
		return errors.New("empty node")
	}
	if len(n.Name) == 0 {
		return errors.New("empty node name")
	}
	_, err = d.store.Create(ctx, n)
	return
}

func (d *DaoRedis) GetNode(ctx context.Context, name string) (n *node.Node, err error) {
	var str string
	str, err = d.client().HGet(ctx, storeredis.NodePrefix(), name).Result()
	if err == redis.Nil {
		err = nil
	}
	if err != nil {
		return
	}
	if len(str) > 0 {
		n = &node.Node{}
		if err = json.Unmarshal([]byte(str), n); err != nil {
			return
		}
	}
	return
}

func (d *DaoRedis) GetNodeNoEmpty(ctx context.Context, name string) (n *node.Node, err error) {
	n, err = d.GetNode(ctx, name)
	if err != nil {
		return
	}
	if n == nil {
		n = &node.Node{}
	}
	return
}

func (d *DaoRedis) AllNodes(ctx context.Context) (list []*node.Node, err error) {
	var result map[string]string
	result, err = d.client().HGetAll(ctx, storeredis.NodePrefix()).Result()
	if err != nil {
		return
	}
	for _, v := range result {
		pipe := &node.Node{}
		if err := pipe.Unmarshal([]byte(v)); err != nil {
			continue
		}
		list = append(list, pipe)
	}
	return
}

func (d *DaoRedis) RefreshNode(ctx context.Context, n *node.Node) (ok bool, err error) {
	var old *node.Node
	old, err = d.GetNode(ctx, n.Name)
	if err != nil {
		return
	}
	if old != nil {
		n.CreateTime = old.CreateTime
	}
	var i int64
	if i, err = d.client().HSet(ctx, storeredis.NodePrefix(), n.Name, n.Val()).Result(); err != nil {
		return
	} else {
		ok = (i > 0)
	}
	return
}

func (d *DaoRedis) UpdateNode(ctx context.Context, name string, opts ...node.NodeOption) (bool, error) {
	return false, nil
}

func (d *DaoRedis) UpdateNodeIP(ctx context.Context, name string, ip string) (ok bool, err error) {
	var old *node.Node
	if old, err = d.GetNodeNoEmpty(ctx, name); err != nil {
		return
	}
	old.IP = ip
	old.UpdateTime = time.Now()
	var i int64
	i, err = d.client().HSet(ctx, storeredis.NodePrefix(), name, old.Val()).Result()
	if err != nil {
		return false, err
	}
	return (i > 0), nil
}

func (d *DaoRedis) UpdateCapacity(ctx context.Context, cap *node.Capacity) (bool, error) {
	i, err := d.client().HSet(ctx, storeredis.CapacityPrefix(), cap.NodeName, cap.Val()).Result()
	return (i > 0), err
}

func (d *DaoRedis) AllCapacity(ctx context.Context) (list []*node.Capacity, err error) {
	var result map[string]string
	result, err = d.client().HGetAll(ctx, storeredis.CapacityPrefix()).Result()
	if err != nil {
		return
	}
	for _, v := range result {
		item := &node.Capacity{}
		if err := item.Unmarshal([]byte(v)); err != nil {
			continue
		}
		list = append(list, item)
	}
	return
}

func (d *DaoRedis) CapacityMap(ctx context.Context) (mapping map[string]*node.Capacity, err error) {
	list, err := d.AllCapacity(ctx)
	if err != nil {
		return nil, err
	}
	mapping = map[string]*node.Capacity{}
	for _, item := range list {
		mapping[item.NodeName] = item
	}
	return
}

func (d *DaoRedis) AllStatus(ctx context.Context) (list []*node.Status, err error) {
	var result map[string]string
	result, err = d.client().HGetAll(ctx, storeredis.StatusPrefix()).Result()
	if err != nil {
		return
	}
	for _, v := range result {
		item := &node.Status{}
		if err := item.Unmarshal([]byte(v)); err != nil {
			continue
		}
		list = append(list, item)
	}
	return
}

func (d *DaoRedis) StatusMap(ctx context.Context) (mapping map[string]*node.Status, err error) {
	list, err := d.AllStatus(ctx)
	if err != nil {
		return nil, err
	}
	mapping = map[string]*node.Status{}
	for _, item := range list {
		mapping[item.NodeName] = item
	}
	return
}

func (d *DaoRedis) CreateOrUpdateStatus(ctx context.Context, nodeName string, opts ...node.StatusOption) (ok bool, err error) {
	status, err := d.GetStatus(ctx, nodeName)
	if err != nil {
		return false, err
	}
	if status == nil {
		status = &node.Status{}
	}
	for _, v := range opts {
		v(status)
	}
	i, err := d.client().HSet(ctx, storeredis.StatusPrefix(), nodeName, status.Val()).Result()
	if err != nil {
		return false, err
	}
	return (i > 0), nil
}

func (d *DaoRedis) GetStatus(ctx context.Context, nodeName string) (s *node.Status, err error) {
	var str string
	if str, err = d.client().HGet(ctx, storeredis.StatusPrefix(), nodeName).Result(); err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return
	}
	s = &node.Status{}
	if err = s.Unmarshal([]byte(str)); err != nil {
		return
	}
	return
}

func (d *DaoRedis) LeaderNode(ctx context.Context) (node string, err error) {
	str, err := d.client().Get(ctx, storeredis.ElectionPrefix()).Result()
	if err == redis.Nil {
		err = nil
	}
	return str, err
}

func (d *DaoRedis) AllElections() (res []map[string]any, err error) {
	return []map[string]any{}, nil
}

func (d *DaoRedis) UpdateAllocatable(ctx context.Context, al *node.Allocatable) (ok bool, err error) {
	i, err := d.client().HSet(ctx, storeredis.AllocatablePrefix(), al.NodeName, al.Val()).Result()
	return (i > 0), err
}
