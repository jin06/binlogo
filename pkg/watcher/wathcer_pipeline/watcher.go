package wathcer_pipeline

import (
	"context"
	"fmt"
	"github.com/jin06/binlogo/pkg/store/etcd"
	"github.com/jin06/binlogo/pkg/store/model/pipeline"
	"github.com/sirupsen/logrus"
)

type PipelineWatcher struct {
	Etcd    *etcd.ETCD
	queue   chan *pipeline.Pipeline
	key     string
}

func (w *PipelineWatcher) Run(ctx context.Context) {
	return
}

func New(key string) (p *PipelineWatcher, err  error){
	p = &PipelineWatcher{
		key : key,
	}
	err = p.Init()
	return
}

func (w *PipelineWatcher) Init() (err error) {
	w.Etcd = etcd.E
	return
}

func (w *PipelineWatcher) Watch(ctx context.Context) (p *pipeline.Pipeline) {
	p = &pipeline.Pipeline{}
	watchCh := w.Etcd.Client.Watch(ctx, w.key)

	go func() {
		for resp := range watchCh {
			err := resp.Err()
			if err != nil {
				logrus.Error(err)
			}
			for _, val := range resp.Events {
				fmt.Printf("%v \n", val)
			}
			select {
				case <- ctx.Done():{
					return
				}
			}
		}
	}()
	return
}
