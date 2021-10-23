package dao

import (
	"context"
	"github.com/jin06/binlogo/pkg/blog"
	"github.com/jin06/binlogo/pkg/store/etcd"
	"github.com/jin06/binlogo/pkg/store/model/pipeline"
	"go.etcd.io/etcd/clientv3"
	"sort"
)

func CreatePipeline(d *pipeline.Pipeline, opts ...clientv3.OpOption) (ok bool, err error) {
	return etcd.Create(d, opts...)
}

func UpdatePipeline(d *pipeline.Pipeline) (bool, error) {
	return etcd.Update(d)
}

func AllPipelines() (list []*pipeline.Pipeline, err error) {
	list = []*pipeline.Pipeline{}
	key := etcd.Prefix() + "/pipeline/"
	res, err := etcd.E.Client.Get(context.Background(), key, clientv3.WithPrefix())
	if err != nil {
		return
	}
	if len(res.Kvs) == 0 {
		return
	}
	for _, v := range res.Kvs {
		ele := &pipeline.Pipeline{
		}
		er := ele.Unmarshal(v.Value)
		if er != nil {
			blog.Error(er)
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

type PagePipeline struct {
	Page int
	Total int
	Data []*pipeline.Pipeline
}

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
	start := ( page -1 ) * size
	end := start + size
	if end > total {
		end = total
	}
	resList := all[start:end]

	res = &PagePipeline{
		Page : page,
		Total: len(all),
		Data: resList,
	}
	return
}
