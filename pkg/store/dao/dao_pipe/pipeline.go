package dao_pipe

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/coreos/etcd/clientv3"
	"github.com/jin06/binlogo/pkg/etcd_client"
	"github.com/jin06/binlogo/pkg/store/etcd"
	"github.com/jin06/binlogo/pkg/store/model/pipeline"
	"github.com/sirupsen/logrus"
	"sort"
)

// PipelinePrefix returns etcd prefix of pipeline info
func PipelinePrefix() string {
	return etcd_client.Prefix() + "/pipeline/info"
}

// GetPipeline get pipeline info from etcd
func GetPipeline(name string) (p *pipeline.Pipeline, err error) {
	key := PipelinePrefix() + "/" + name
	res, err := etcd.E.Client.Get(context.TODO(), key)
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
func CreatePipeline(d *pipeline.Pipeline, opts ...clientv3.OpOption) (ok bool, err error) {
	key := PipelinePrefix() + "/" + d.Name
	txn := etcd.E.Client.Txn(context.TODO())
	txn = txn.If(clientv3.Compare(clientv3.CreateRevision(key), "=", int64(0)))
	txn = txn.Then(clientv3.OpPut(key, d.Val(), opts...))
	resp, err := txn.Commit()
	if err != nil {
		return
	}
	ok = resp.Succeeded
	return
}

// UpdatePipeline update pipeline info in etcd
func UpdatePipeline(pipeName string, opts ...pipeline.OptionPipeline) (ok bool, err error) {
	if pipeName == "" {
		return false, errors.New("empty pipeline name")
	}
	key := PipelinePrefix() + "/" + pipeName
	res, err := etcd_client.Default().Get(context.Background(), key)
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
	txn := etcd_client.Default().Txn(context.Background()).
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
func AllPipelines() (list []*pipeline.Pipeline, err error) {
	list = []*pipeline.Pipeline{}
	key := PipelinePrefix()
	res, err := etcd.E.Client.Get(context.TODO(), key, clientv3.WithPrefix())
	if err != nil {
		return
	}
	if len(res.Kvs) == 0 {
		return
	}
	for _, v := range res.Kvs {
		ele := &pipeline.Pipeline{}
		er := ele.Unmarshal(v.Value)
		if er != nil {
			logrus.Error(er)
			continue
		}
		list = append(list, ele)
	}
	return
}

type pipelineSlice []*pipeline.Pipeline

func (ps pipelineSlice) Len() int { return len(ps) }

func (ps pipelineSlice) Swap(i, j int) { ps[i], ps[j] = ps[j], ps[i] }

func (ps pipelineSlice) Less(i, j int) bool {
	return ps[i].CreateTime.Before(ps[i].CreateTime)
}

// PagePipeline for pipeline pages
type PagePipeline struct {
	Page  int
	Total int
	Data  []*pipeline.Pipeline
}

// PagePipelines get a PagePipeline by page and size
func PagePipelines(page int, size int) (res *PagePipeline, err error) {
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 10
	}
	all, err := AllPipelines()
	if err != nil {
		return
	}
	sort.Sort(pipelineSlice(all))

	total := len(all)
	totalPage := total / size
	if page > totalPage {
		page = totalPage
	}
	start := (page - 1) * size
	end := start + size
	if end > total {
		end = total
	}
	resList := all[start:end]

	res = &PagePipeline{
		Page:  page,
		Total: len(all),
		Data:  resList,
	}
	return
}

// AllPipelinesMap returns all pipelines from etcd in map form
func AllPipelinesMap() (mapping map[string]*pipeline.Pipeline, err error) {
	list, err := AllPipelines()
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
func DeletePipeline(name string) (err error) {
	if name == "" {
		return errors.New("empty name")
	}
	key := PipelinePrefix() + "/" + name
	_, err = etcd_client.Default().Delete(context.Background(), key)
	if err != nil {
		return
	}
	return
}

// DeleteCompletePipeline delete pipeline, contains pipeline info, pipeline position
func DeleteCompletePipeline(name string) (err error) {
	if name == "" {
		return errors.New("empty name")
	}
	err = DeletePipeline(name)
	if err != nil {
		return
	}
	err = DeletePosition(name)
	if err != nil {
		return
	}
	return
}
