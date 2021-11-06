package manager_pipe

import (
	"context"
	"errors"
	"github.com/jin06/binlogo/app/pipeline/pipeline"
	"github.com/jin06/binlogo/pkg/blog"
	"github.com/jin06/binlogo/pkg/mutex"
	"github.com/jin06/binlogo/pkg/store/dao/dao_pipe"
	"github.com/jin06/binlogo/pkg/store/dao/dao_sche"
	pipeline2 "github.com/jin06/binlogo/pkg/store/model/pipeline"
)

type instance struct {
	pName string
	p     *pipeline.Pipeline
	m     *mutex.Mutex
}

func newInstance(name string) (ins *instance, err error) {
	pipeModel, err := dao_pipe.GetPipeline(name)
	if err != nil {
		return
	}
	if pipeModel == nil {
		err = errors.New("no pipeline: " + name)
		return
	}
	posModel, err := dao_pipe.GetPosition(name)
	if err != nil {
		return
	}
	if posModel == nil {
		posModel = &pipeline2.Position{
			//PipelineName: name,
		}
	}
	pipe, err := pipeline.New(pipeline.OptionPipeline(pipeModel), pipeline.OptionPosition(posModel))
	if err != nil {
		return
	}
	ins = &instance{p: pipe, pName: name}
	ins.m, err = mutex.New(dao_sche.PipelineLockPrefix() + "/" + name)
	if err != nil {
		return
	}
	return
}

func (i *instance) start() (err error) {
	err = i.m.Lock()
	if err != nil {
		return
	}

	defer func() {
		if err != nil {
			i.m.Unlock()
		}
	}()
	blog.Info("Start pipeline, ", i.pName)
	err = i.p.Run(context.Background())
	return
}

func (i *instance) stop() (err error) {
	i.p.Stop()
	err = i.m.Unlock()
	return
}
