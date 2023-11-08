package manager_pipe

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/jin06/binlogo/app/pipeline/pipeline"
	"github.com/jin06/binlogo/pkg/event"
	"github.com/jin06/binlogo/pkg/register"
	"github.com/jin06/binlogo/pkg/store/dao/dao_pipe"
	"github.com/jin06/binlogo/pkg/store/dao/dao_register"
	event2 "github.com/jin06/binlogo/pkg/store/model/event"
	pipeline2 "github.com/jin06/binlogo/pkg/store/model/pipeline"
	"github.com/sirupsen/logrus"
)

type instance struct {
	pipeName  string
	nodeName  string
	pipeIns   *pipeline.Pipeline
	pipeInfo  *pipeline2.Pipeline
	pipeReg   *register.Register
	cancel    context.CancelFunc
	mutex     sync.Mutex
	startTime time.Time
	stopping  chan struct{}
	stopOnce  sync.Once
	started   chan struct{}
	stopped   chan struct{}
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
		posPos = &pipeline2.Position{}
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
	reg := register.New(
		register.WithKey(dao_register.PipeInstancePrefix()+"/"+i.pipeName),
		register.WithData(insModel),
	)
	i.pipeInfo = pipeInfo
	i.pipeIns = pipe
	i.pipeReg = reg
	return
}

func (i *instance) start(ctx context.Context) (err error) {
	i.startTime = time.Now()
	defer func() {
		if r := recover(); r != nil {
			logrus.Errorln("instance run panic, ", r)
		}
		if err != nil {
			event.Event(event2.NewErrorPipeline(i.pipeName, "Pipeline instance stopped error: "+err.Error()))
		}
		event.Event(event2.NewInfoPipeline(i.pipeName, "Pipeline instance stopped"))
		close(i.stopped)
	}()
	if err = i.init(); err != nil {
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
	logrus.Info("pipeline instance start: ", i.pipeName)
	event.Event(event2.NewInfoPipeline(i.pipeName, "Pipeline instance start success"))
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
		logrus.Info("pipeline instance stop: ", i.pipeName)
	})
}

// StartTime returns instance start time
func (i *instance) StartTime() time.Time {
	return i.startTime
}
