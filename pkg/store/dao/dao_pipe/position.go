package dao_pipe

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/jin06/binlogo/pkg/etcd_client"
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
	_, err = etcd_client.Default().Put(context.TODO(), key, string(b))
	return
}

func GetPosition(pipeName string) (p *pipeline.Position, err error) {
	key := PositionPrefix() + "/" + pipeName
	res, err := etcd_client.Default().Get(context.TODO(), key)
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

func DeletePosition( name string) (err error) {
	if name == "" {
		return errors.New("empty name")
	}
	key := PositionPrefix() + "/" + name
	_, err = etcd_client.Default().Delete(context.TODO(), key)
	return
}
