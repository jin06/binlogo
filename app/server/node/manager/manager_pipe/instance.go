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
	pipeName  string
	nodeName  string
	pipeIns   *pipeline.Pipeline
	pipeInfo  *mPipe.Pipeline
	pipeReg   *register.Register
	cancel    context.CancelFunc
	mutex     sync.Mutex
	startTime time.Time
	stopping  chan struct{}
	stopOnce  sync.Once
	started   chan struct{}
	stopped   chan struct{}
	exit      bool
	//stopped   chan struct{}
}

func newInstance(pipeName string, nodeName string) *instance {
	ins := &instance{
		pipeName:  pipeName,
		nodeName:  nodeName,
		mutex:     sync.Mutex{},
		startTime: time.Time{},
		stopping:  make(chan struct{}),
		started:   make(chan struct{}),
		stopped:   make(chan struct{}),
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
		close(i.stopped)
		i.exit = true
	}()
	if err = i.init(ctx); err != nil {
		return
	}
	stx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		i.pipeReg.Run(stx)
		i.stop()
	}()
	go func() {
		i.pipeIns.Run(stx)
		i.stop()
	}()
	event.Event(model.NewInfoPipeline(i.pipeName, "Pipeline instance start success"))
	close(i.started)

	select {
	case <-ctx.Done():
		return
	case <-i.stopping:
		return
	}
}

func (i *instance) stop() {
	i.stopOnce.Do(func() {
		i.startTime = time.Time{}
		close(i.stopping)
	})
}

// StartTime returns instance start time
func (i *instance) StartTime() time.Time {
	return i.startTime
}
