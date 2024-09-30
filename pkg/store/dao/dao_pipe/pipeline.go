package dao_pipe

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
	return etcdclient.Prefix() + "/pipeline/info"
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
	key := d.Key()
	if err := store_redis.GetClient().HSet(ctx, key, d).Err(); err != nil {
		return false, err
	}
	return
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
func AllPipelines(ctx context.Context) ([]*pipeline.Pipeline, error) {
	list := []*pipeline.Pipeline{}
	// if err := store_redis.Default.List(ctx, list); err != nil {
	// 	return nil, err
	// }
	return list, nil
}

//type pipelineSlice []*pipeline.Pipeline
//
//func (ps pipelineSlice) Len() int { return len(ps) }
//
//func (ps pipelineSlice) Swap(i, j int) { ps[i], ps[j] = ps[j], ps[i] }
//
//func (ps pipelineSlice) Less(i, j int) bool {
//	return ps[i].CreateTime.Before(ps[i].CreateTime)
//}

// PagePipeline for pipeline pages
//type PagePipeline struct {
//	Page  int
//	Total int
//	Data  []*pipeline.Pipeline
//}

// PagePipelines get a PagePipeline by page and size
//func PagePipelines(page int, size int) (res *PagePipeline, err error) {
//	if page < 1 {
//		page = 1
//	}
//	if size < 1 {
//		size = 10
//	}
//	all, err := AllPipelines()
//	if err != nil {
//		return
//	}
//	sort.Sort(pipelineSlice(all))
//
//	total := len(all)
//	totalPage := total / size
//	if page > totalPage {
//		page = totalPage
//	}
//	start := (page - 1) * size
//	end := start + size
//	if end > total {
//		end = total
//	}
//	resList := all[start:end]
//
//	res = &PagePipeline{
//		Page:  page,
//		Total: len(all),
//		Data:  resList,
//	}
//	return
//}

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
func DeletePipeline(name string) (ok bool, err error) {
	if name == "" {
		err = errors.New("empty name")
		return
	}
	key := PipelinePrefix() + "/" + name
	res, err := etcdclient.Default().Delete(context.Background(), key)
	if err != nil {
		return
	}
	if res.Deleted > 0 {
		ok = true
	}
	return
}

// DeleteCompletePipeline delete pipeline, contains pipeline info, pipeline position
func DeleteCompletePipeline(name string) (ok bool, err error) {
	if name == "" {
		err = errors.New("empty name")
		return
	}
	ok, err = DeletePipeline(name)
	if err != nil {
		return
	}
	_, err = DeletePosition(name)
	if err != nil {
		return
	}
	return
}
