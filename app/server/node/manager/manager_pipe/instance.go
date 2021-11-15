package manager_pipe

import (
	"context"
	"errors"
	"github.com/jin06/binlogo/app/pipeline/pipeline"
	"github.com/jin06/binlogo/pkg/register"
	"github.com/jin06/binlogo/pkg/store/dao/dao_pipe"
	"github.com/jin06/binlogo/pkg/store/dao/dao_register"
	pipeline2 "github.com/jin06/binlogo/pkg/store/model/pipeline"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

type instance struct {
	pipeName string
	nodeName string
	pipeIns  *pipeline.Pipeline
	pipeInfo *pipeline2.Pipeline
	pipeReg  *register.Register
	cancel   context.CancelFunc
	status   status
	mutex    sync.Mutex
}

type status byte

const (
	STATUS_RUN  status = 2
	STATUS_STOP status = 4
)

func newInstance(pipeName string, nodeName string) *instance {
	ins := &instance{
		pipeName: pipeName,
		nodeName: nodeName,
		mutex:    sync.Mutex{},
		status:   STATUS_STOP,
	}
	return ins
}

func (i *instance) init() (err error) {
	pipeInfo, err := dao_pipe.GetPipeline(i.pipeName)
	if err != nil {
		return
	}
	if pipeInfo == nil {
		err = errors.New("no pipeline: " + i.pipeName)
		return
	}
	posPos, err := dao_pipe.GetPosition(i.pipeName)
	if err != nil {
		return
	}
	if posPos == nil {
		posPos = &pipeline2.Position{
		}
	}
	pipe, err := pipeline.New(
		pipeline.OptionPipeline(pipeInfo),
		pipeline.OptionPosition(posPos),
	)
	if err != nil {
		return
	}
	insModel := &pipeline2.Instance{
		PipelineName: i.pipeName,
		NodeName:     i.nodeName,
		CreateTime:   time.Now(),
	}
	reg, err := register.New(
		register.WithKey(dao_register.PipeInstancePrefix()+"/"+i.pipeName),
		register.WithData(insModel),
	)
	i.pipeInfo = pipeInfo
	i.pipeIns = pipe
	i.pipeReg = reg
	return
}

func (i *instance) start() (err error) {
	if i.status == STATUS_RUN {
		return
	}
	i.status = STATUS_RUN
	defer func() {
		i.status = STATUS_STOP
	}()
	err = i.init()
	if err != nil {
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		cancel()
	}()
	i.cancel = cancel
	i.pipeReg.Run(ctx)
	i.pipeIns.Run(ctx)
	logrus.Info("pipeline instance start: ", i.pipeName)

	select {
	case <-ctx.Done():
		{
			return
		}
	case <-i.pipeIns.Context().Done():
		{
			return
		}
	case <-i.pipeReg.Context().Done():
		{
			return
		}
	}
}

func (i *instance) stop() {
	if i.status == STATUS_STOP {
		return
	}
	defer func() {
		i.status = STATUS_STOP
	}()
	i.cancel()
	logrus.Info("pipeline instance stop: ", i.pipeName)
	return
}
