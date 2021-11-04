package dao_pipe

import (
	"context"
	"encoding/json"
	"github.com/jin06/binlogo/pkg/store/etcd"
	"github.com/jin06/binlogo/pkg/store/model/pipeline"
)

func PositionPrefix() string {
	return etcd.Prefix() + "/pipeline/position"
}

func UpdatePosition(p *pipeline.Position) (err error) {
	key := PositionPrefix() + "/" + p.PipelineName
	b, err := json.Marshal(p)
	if err != nil {
		return
	}
	_, err = etcd.E.Client.Put(context.Background(), key, string(b))
	return
}

func GetPosition(pipeName string) (p *pipeline.Position, err error) {
	key := PositionPrefix() + "/" + pipeName
	res, err := etcd.E.Client.Get(context.Background(), key)
	if err != nil {
		return
	}
	if len(res.Kvs) == 0 {
		return
	}
	p = &pipeline.Position{}
	if err = p.Unmarshal(res.Kvs[0].Value); err != nil {
		return
	}
	return
}
