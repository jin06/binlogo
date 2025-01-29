package manager_pipe

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/jin06/binlogo/v2/app/pipeline/pipeline"
	"github.com/jin06/binlogo/v2/pkg/event"
	"github.com/jin06/binlogo/v2/pkg/register"
	"github.com/jin06/binlogo/v2/pkg/store/dao"
	"github.com/jin06/binlogo/v2/pkg/store/model"
	mPipe "github.com/jin06/binlogo/v2/pkg/store/model/pipeline"
	"github.com/sirupsen/logrus"
)

type instance struct {
	pipeName     string
	nodeName     string
	pipeIns      *pipeline.Pipeline
	pipeInfo     *mPipe.Pipeline
	pipeReg      *register.Register
	mutex        sync.Mutex
	startTime    time.Time
	manager      *Manager
	started      chan struct{}
	closing      chan struct{}
	closed       chan struct{}
	closeOnce    sync.Once
	completeOnce sync.Once

	//stopped   chan struct{}
}

func newInstance(pipeName string, nodeName string, manager *Manager) *instance {
	ins := &instance{
		pipeName:  pipeName,
		nodeName:  nodeName,
		mutex:     sync.Mutex{},
		startTime: time.Time{},
		closing:   make(chan struct{}),
		started:   make(chan struct{}),
		closed:    make(chan struct{}),
		manager:   manager,
	}
	return ins
}

func (i *instance) init(ctx context.Context) (err error) {
	pipeInfo, err := dao.GetPipeline(context.TODO(), i.pipeName)
	if err != nil {
		return
	}
	if pipeInfo == nil {
		return errors.New("no pipeline: " + i.pipeName)
	}
	posPos, err := dao.GetPosition(ctx, i.pipeName)
	if err != nil {
		return
	}
	if posPos == nil {
		posPos = &mPipe.Position{}
	}
	pipe, err := pipeline.New(
		pipeline.OptionPipeline(pipeInfo),
		pipeline.OptionPosition(posPos),
	)
	if err != nil {
		return
	}
	insModel := &mPipe.Instance{
		PipelineName: i.pipeName,
		NodeName:     i.nodeName,
		CreateTime:   time.Now(),
	}
	reg := register.New(
		register.WithKey(dao.PipeInstancePrefix()+"/"+i.pipeName),
		register.WithData(insModel),
	)
	i.pipeInfo = pipeInfo
	i.pipeIns = pipe
	i.pipeReg = reg
	return
}

func (i *instance) start(ctx context.Context) (err error) {
	defer i.CompleteClose()
	defer i.Close()
	i.startTime = time.Now()
	logrus.Infof("Pipeline instance start: %s", i.pipeName)
	defer func() {
		if r := recover(); r != nil {
			logrus.Errorf("instance run panic: %v", r)
		}
		logrus.Infof("pipeline instance stopped: %s", i.pipeName)
		if err != nil {
			event.Event(model.NewErrorPipeline(i.pipeName, "Pipeline instance stopped error: "+err.Error()))
		}
		event.Event(model.NewInfoPipeline(i.pipeName, "Pipeline instance stopped"))
	}()
	if err = i.init(ctx); err != nil {
		return
	}
	go func() {
		i.pipeReg.Run(ctx)
		i.Close()
	}()
	go func() {
		i.pipeIns.Run(ctx)
		i.Close()
	}()
	event.Event(model.NewInfoPipeline(i.pipeName, "Pipeline instance start success"))
	close(i.started)

	select {
	case <-ctx.Done():
		return
	case <-i.closing:
		return
	case <-i.pipeReg.Closed():
		return
	case <-i.pipeIns.Closed():
		return
	}
}

func (i *instance) CompleteClose() {
	i.completeOnce.Do(func() {
		<-i.pipeReg.Closed()
		<-i.pipeIns.Closed()
		i.manager.Remove(i.pipeName)
		close(i.closed)
	})
}

func (i *instance) Close() {
	i.closeOnce.Do(func() {
		i.pipeReg.Close()
		i.pipeIns.Close()
		close(i.closing)
	})
}

// StartTime returns instance start time
func (i *instance) StartTime() time.Time {
	return i.startTime
}

func (i *instance) Closed() chan struct{} {
	return i.closed
}
