package dao_pipe

import (
	"context"
	"github.com/coreos/etcd/clientv3"
	"github.com/jin06/binlogo/pkg/store/etcd"
	"github.com/jin06/binlogo/pkg/store/model/pipeline"
	"github.com/sirupsen/logrus"
	"sort"
)

func PipelinePrefix() string {
	return etcd.Prefix() + "/pipeline/info"
}

func GetPipeline(name string) (p *pipeline.Pipeline , err error){
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

func UpdatePipeline(d *pipeline.Pipeline) (bool, error) {
	return etcd.Update(d)
}

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
		ele := &pipeline.Pipeline{
		}
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
