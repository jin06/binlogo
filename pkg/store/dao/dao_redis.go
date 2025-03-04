package dao

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/jin06/binlogo/v2/configs"
	"github.com/jin06/binlogo/v2/internal/constant"
	"github.com/jin06/binlogo/v2/pkg/store/model"
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

func (d *DaoRedis) publish(ctx context.Context, channel string, message any) error {
	cmd := d.client().Publish(ctx, channel, message)
	return cmd.Err()
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

func getAllHashDatas[T any](ctx context.Context, client *redis.Client, prefix string) (list []*T, err error) {
	var keys []string
	if keys, err = scanKeysWithPrefix(ctx, client, prefix); err != nil {
		return
	}
	for _, key := range keys {
		cmd := client.HGetAll(ctx, key)
		if cmd.Err() != nil {
			logrus.Error("get hashdata error ", cmd.Err())
			continue
		}
		// var item T
		item := new(T)
		if err := cmd.Scan(item); err != nil {
			logrus.Error("scan hashdata error ", err)
			continue
		}
		list = append(list, item)
	}
	return
}

func (d *DaoRedis) RegisterInstance(ctx context.Context, ins *pipeline.Instance, exp time.Duration) error {
	key := storeredis.GetPipeInstanceKey(ins.PipelineName)
	err := d.client().Watch(ctx, func(tx *redis.Tx) error {
		i, err := tx.Exists(ctx, key).Result()
		if err != nil {
			return err
		}
		if i > 0 {
			return constant.ErrPipelineInstanceExists
		}
		_, err = tx.TxPipelined(ctx, func(p redis.Pipeliner) error {
			p.HMSet(ctx, key, ins)
			p.Expire(ctx, key, exp)
			return nil
		})
		return err
	}, key)
	return err
}

func (d *DaoRedis) UnRegisterInstance(ctx context.Context, pipe string, n string) error {
	key := storeredis.GetPipeInstanceKey(pipe)
	err := d.client().Watch(ctx, func(tx *redis.Tx) error {
		ins := &pipeline.Instance{}
		if err := tx.HGetAll(ctx, key).Scan(ins); err != nil {
			return err
		}
		return tx.Del(ctx, key).Err()
	}, key)
	return err
}

func (d *DaoRedis) LeaseInstance(ctx context.Context, pipe string, exp time.Duration) error {
	return d.client().Expire(ctx, storeredis.GetPipeInstanceKey(pipe), exp).Err()
}

func (d *DaoRedis) GetInstance(ctx context.Context, pipe string) (ins *pipeline.Instance, err error) {
	ins = &pipeline.Instance{}
	cmd := d.client().HGetAll(ctx, storeredis.GetPipeInstanceKey(pipe))
	if err = cmd.Err(); err != nil {
		return
	}
	err = cmd.Scan(ins)
	return
}

func (d *DaoRedis) AllInstance(ctx context.Context) (list []*pipeline.Instance, err error) {
	return getAllHashDatas[pipeline.Instance](ctx, d.client(), storeredis.PipelineInstancePrefix())
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
	var cursor uint64
	for {
		keys, nextCursor, err := d.client().Scan(ctx, cursor, storeredis.GetStatusKey("*"), 10).Result()
		if err != nil {
			return nil, err
		}
		for _, key := range keys {
			keyType, err := d.client().Type(ctx, key).Result()
			if err != nil {
				logrus.Error(err)
				continue
			}
			if keyType == "hash" {
				cmd := d.client().HGetAll(ctx, key)
				if cmd.Err() != nil {
					return nil, cmd.Err()
				}
				item := &node.Status{}
				err := cmd.Scan(item)
				if err != nil {
					logrus.Error(err)
					continue
				}
				list = append(list, item)
			}
		}
		if nextCursor == 0 {
			break
		}
		cursor = nextCursor

	}

	var result map[string]string
	// result, err = d.client().HGetAll(ctx, storeredis.StatusPrefix()).Result()
	// if err != nil {
	// return
	// }
	for _, v := range result {
		item := &node.Status{}
		if err := item.Unmarshal([]byte(v)); err != nil {
			continue
		}
		list = append(list, item)
	}
	return
}

func (d *DaoRedis) DeleteStatus(ctx context.Context, name string) error {
	key := storeredis.GetStatusKey(name)
	return d.client().Del(ctx, key).Err()
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

func (d *DaoRedis) CreateOrUpdateStatus(ctx context.Context, nodeName string, conditions node.StatusConditions) (ok bool, err error) {
	values := []any{}
	for name, val := range conditions {
		values = append(values, name, val)
	}
	if len(values) == 0 {
		return false, nil
	}
	i, err := d.client().HSet(ctx, storeredis.GetStatusKey(nodeName), values...).Result()
	if err != nil {
		return false, err
	}
	return (i > 0), nil
}

func (d *DaoRedis) GetStatus(ctx context.Context, nodeName string) (s *node.Status, err error) {
	cmd := d.client().HGetAll(ctx, storeredis.GetStatusKey(nodeName))
	if err = cmd.Err(); err != nil {
		if err != redis.Nil {
			return nil, err
		}
		return nil, nil
	}
	s = &node.Status{}
	err = cmd.Scan(s)
	return
}

func (d *DaoRedis) LeaderNode(ctx context.Context) (node string, err error) {
	str, err := d.client().Get(ctx, storeredis.MasterPreifx()).Result()
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

func (d *DaoRedis) GetPipelineBind(ctx context.Context) (*model.PipelineBind, error) {
	bindings, err := d.client().HGetAll(ctx, storeredis.PipelineBindPrefix()).Result()
	if err != nil {
		return nil, err
	}
	return &model.PipelineBind{
		Bindings: bindings,
	}, nil
}

func (d *DaoRedis) getPipelineBindNode(ctx context.Context, name string) (nodeName string, err error) {
	return d.client().HGet(ctx, storeredis.PipelineBindPrefix(), name).Result()
}

func (d *DaoRedis) setPipeLineBindNode(ctx context.Context, pName string, nodeName string) error {
	return d.client().HSet(ctx, storeredis.PipelineBindPrefix(), pName, nodeName).Err()
}

func (d *DaoRedis) delPipelineBindNode(ctx context.Context, pName string) error {
	return d.client().HDel(ctx, storeredis.PipelineBindPrefix(), pName).Err()
}

func (d *DaoRedis) publishPelineBindChange(ctx context.Context) error {
	data, err := d.GetPipelineBind(ctx)
	if err != nil {
		return err
	}
	return d.client().Publish(ctx, storeredis.PipelineBindChan(), data).Err()
}

func (d *DaoRedis) UpdatePipelineBindIfNotExist(ctx context.Context, pName string, nName string) error {
	if err := d.client().HSetNX(ctx, storeredis.PipelineBindPrefix(), pName, nName).Err(); err != nil {
		return err
	}
	d.publishPelineBindChange(ctx)
	return nil
}

func (d *DaoRedis) UpdatePipelineBind(ctx context.Context, pName string, nName string) (bool, error) {
	i, err := d.client().HSet(ctx, storeredis.PipelineBindPrefix(), pName, nName).Result()
	d.publishPelineBindChange(ctx)
	return i > 0, err
}

func (d *DaoRedis) DeletePipelineBind(ctx context.Context, pName string) (bool, error) {
	i, err := d.client().HDel(ctx, storeredis.PipelineBindPrefix(), pName).Result()
	d.publishPelineBindChange(ctx)
	return i > 0, err
}

func (d *DaoRedis) WatchPipelinBind(ctx context.Context) chan model.PipelineBind {
	pubsub := d.client().Subscribe(ctx, storeredis.PipelineBindChan())
	ch := make(chan model.PipelineBind)
	go func() {
		msgch := pubsub.Channel()
		for {
			select {
			case msg := <-msgch:
				fmt.Println(msg)
			}
		}
	}()
	return ch
}

func (d *DaoRedis) GetPipeline(ctx context.Context, name string) (p *pipeline.Pipeline, err error) {
	raw, err := d.client().HGet(ctx, storeredis.PipelinesPrefix(), name).Bytes()
	fmt.Println(storeredis.PipelinesPrefix())
	if err != nil {
		return nil, err
	}

	p = &pipeline.Pipeline{}
	err = p.Unmarshal(raw)
	return
}

func (d *DaoRedis) UpdatePipeline(ctx context.Context, name string, opts ...pipeline.OptionPipeline) (err error) {
	defer func() {
		if err == redis.Nil {
			err = nil
		}
	}()
	p, err := d.GetPipeline(ctx, name)
	if err != nil {
		return err
	}
	for _, v := range opts {
		v(p)
	}
	cmd := d.client().HSet(ctx, PipelinesKey(), name, p.Val())
	return cmd.Err()
}

func (d *DaoRedis) AllPipelines(ctx context.Context) (list []*pipeline.Pipeline, err error) {
	result, err := storeredis.GetClient().HGetAll(ctx, storeredis.PipelinesPrefix()).Result()
	if err != nil {
		return nil, err
	}
	for _, v := range result {
		pipe := &pipeline.Pipeline{}
		if err := pipe.Unmarshal([]byte(v)); err != nil {
			continue
		}
		list = append(list, pipe)
	}
	return
}

func (d *DaoRedis) AllPipelinesMap(ctx context.Context) (mapping map[string]*pipeline.Pipeline, err error) {
	list, err := AllPipelines(ctx)
	if err != nil {
		return
	}
	mapping = map[string]*pipeline.Pipeline{}
	for i := 0; i < len(list); i++ {
		mapping[list[i].Name] = list[i]
	}
	return
}

func (d *DaoRedis) ClearOrDeleteBind(ctx context.Context, name string) (err error) {
	pipe, err := d.GetPipeline(ctx, name)
	if err != nil {
		return err
	}
	if pipe.ExpectRun() {
		err = d.setPipeLineBindNode(ctx, name, "")
	} else {
		err = d.delPipelineBindNode(ctx, name)
	}
	return
}

func (d *DaoRedis) UpdateEvent(ctx context.Context, e *model.Event) error {
	key := storeredis.GetEventKey(e.K)
	expire := time.Duration(configs.Default.EventExpire) * time.Second
	if expire <= 0 {
		expire = time.Second * 86400
	}
	pipe := d.client().Pipeline()
	pipe.HSet(ctx, key, e)
	pipe.Expire(ctx, key, expire)
	_, err := pipe.Exec(ctx)
	return err
}

func (d *DaoRedis) GetPosition(ctx context.Context, pipe string) (p *pipeline.Position, err error) {
	cmd := d.client().HGetAll(ctx, storeredis.GetPositionKey(pipe))
	if err = cmd.Err(); err != nil {
		return
	}
	p = &pipeline.Position{}
	err = cmd.Scan(p)
	return
}

func (d *DaoRedis) UpdatePosition(ctx context.Context, p *pipeline.Position) error {
	return d.client().HMSet(ctx, storeredis.GetPositionKey(p.PipelineName), p).Err()
}

func (d *DaoRedis) DeletePosition(ctx context.Context, name string) (err error) {
	return d.client().Del(ctx, storeredis.GetPositionKey(name)).Err()
}

func (d *DaoRedis) UpdateRecord(ctx context.Context, p *pipeline.RecordPosition) (err error) {
	return d.client().HMSet(ctx, storeredis.GetRecordPositionKey(p.PipelineName), p).Err()
}

func (d *DaoRedis) UpdateRecordSafe(ctx context.Context, pipe string, opts ...pipeline.OptionRecord) (err error) {
	r, err := d.GetRecord(ctx, pipe)
	if err != nil {
		return err
	}
	r.PipelineName = pipe
	for _, v := range opts {
		v(r)
	}
	return d.UpdateRecord(ctx, r)
}

func (d *DaoRedis) GetRecord(ctx context.Context, pipe string) (r *pipeline.RecordPosition, err error) {
	cmd := d.client().HGetAll(ctx, storeredis.GetRecordPositionKey(pipe))
	if err = cmd.Err(); err != nil {
		if err != redis.Nil {
			return
		}
		return nil, nil
	}
	r = pipeline.NewRecordPosition()
	err = cmd.Scan(r)
	return
}
