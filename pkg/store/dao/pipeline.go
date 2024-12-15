package dao

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/jin06/binlogo/v2/pkg/etcdclient"
	"github.com/jin06/binlogo/v2/pkg/store/model/pipeline"
	store_redis "github.com/jin06/binlogo/v2/pkg/store/redis"
	clientv3 "go.etcd.io/etcd/client/v3"
)

// PipelinePrefix returns etcd prefix of pipeline info
func PipelinePrefix() string {
	return store_redis.Prefix() + "/pipeline"
}

func PipelinesKey() string {
	return store_redis.Prefix() + "/pipelines"
}

// GetPipeline get pipeline info from etcd
func GetPipeline(name string) (p *pipeline.Pipeline, err error) {
	key := PipelinePrefix() + "/" + name
	res, err := etcdclient.Default().Get(context.TODO(), key)
	if err != nil {
		return
	}
	if len(res.Kvs) == 0 {
		return
	}
	p = &pipeline.Pipeline{}
	if err = p.Unmarshal(res.Kvs[0].Value); err != nil {
		return
	}
	return
}

// CreatePipeline write pipeline info to etcd
func CreatePipeline(ctx context.Context, d *pipeline.Pipeline) (ok bool, err error) {
	cmd := store_redis.GetClient().HSet(ctx, PipelinesKey(), d.Name, d.Val())
	i, err := cmd.Result()
	if err != nil {
		return false, err
	}
	return i > 0, nil
}

// UpdatePipeline update pipeline info in etcd
func UpdatePipeline(pipeName string, opts ...pipeline.OptionPipeline) (ok bool, err error) {
	if pipeName == "" {
		return false, errors.New("empty pipeline name")
	}
	key := PipelinePrefix() + "/" + pipeName
	res, err := etcdclient.Default().Get(context.Background(), key)
	if err != nil {
		return
	}
	if len(res.Kvs) == 0 {
		return false, errors.New("empty pipeline")
	}
	revision := res.Kvs[0].CreateRevision
	pipe := &pipeline.Pipeline{}
	err = json.Unmarshal(res.Kvs[0].Value, pipe)
	if err != nil {
		return
	}
	for _, v := range opts {
		v(pipe)
	}
	txn := etcdclient.Default().Txn(context.Background()).
		If(clientv3.Compare(clientv3.CreateRevision(key), "=", revision)).
		Then(clientv3.OpPut(key, pipe.Val()))
	resp, err := txn.Commit()
	if err != nil {
		return
	}
	ok = resp.Succeeded
	return
}

// AllPipelines get all info of pipelines from etcd in array form
func AllPipelines(ctx context.Context) (list []*pipeline.Pipeline, err error) {
	var result map[string]string
	result, err = store_redis.GetClient().HGetAll(ctx, PipelinesKey()).Result()
	if err != nil {
		return
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

// AllPipelinesMap returns all pipelines from etcd in map form
func AllPipelinesMap(ctx context.Context) (mapping map[string]*pipeline.Pipeline, err error) {
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

// DeletePipeline delete pipeline info by name
func DeletePipeline(ctx context.Context, name string) (ok bool, err error) {
	return store_redis.Default.Delete(ctx, &pipeline.Pipeline{Name: name})
}

// DeleteCompletePipeline delete pipeline, contains pipeline info, pipeline position
func DeleteCompletePipeline(ctx context.Context, name string) (ok bool, err error) {
	if name == "" {
		err = errors.New("empty name")
		return
	}
	ok, err = DeletePipeline(ctx, name)
	if err != nil {
		return
	}
	_, err = DeletePosition(name)
	if err != nil {
		return
	}
	return
}
