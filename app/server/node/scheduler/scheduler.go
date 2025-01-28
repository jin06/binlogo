package scheduler

import (
	"context"
	"sync"
	"time"

	"github.com/jin06/binlogo/v2/pkg/event"
	"github.com/jin06/binlogo/v2/pkg/store/dao"
	"github.com/jin06/binlogo/v2/pkg/store/model"
	"github.com/jin06/binlogo/v2/pkg/store/model/pipeline"
	"github.com/sirupsen/logrus"
)

// Scheduler schedule the pipeline.
// allocate the pipeline to the appropriate node for operation
type Scheduler struct {
	lock     sync.Mutex
	stopOnce sync.Once
	stopping chan struct{}
	Exit     bool
}

// New returns a new Scheduler
func New() (s *Scheduler) {
	s = &Scheduler{
		stopping: make(chan struct{}),
	}
	return
}

// Run start working
func (s *Scheduler) Run(ctx context.Context) (err error) {
	logrus.Info("Pipeline Scheduler run")
	s.schedule(ctx)
	s.Exit = true
	logrus.Info("Pipeline Scheduler stop")
	return
}

func (s *Scheduler) schedule(ctx context.Context) {
	// stx, cancel := context.WithCancel(ctx)
	// defer cancel()
	// key := fmt.Sprintf("%s/%s", dao.SchedulerPrefix(), "pipeline_bind")
	// w, err := watcher.New(watcher.WithHandler(watcher.WrapSchedulerBinding()), watcher.WithKey(key))
	// if err != nil {
	// 	return
	// }
	// defer w.Close()
	// waCh, err := w.WatchEtcd(stx)
	// if err != nil {
	// return
	// }
	schedulePipes(ctx, nil)
	ticker := time.NewTicker(time.Second * 1)
	defer ticker.Stop()
	for {
		select {
		case <-s.stopping:
			return
		case <-ctx.Done():
			return
		// case ev, whOK := <-waCh:
		// 	if !whOK {
		// 		return
		// 	}
		// 	{
		// 		if ev.Event.Type == mvccpb.PUT {
		// 			if val, ok := ev.Data.(*model.PipelineBind); ok {
		// 				schedulePipes(ctx, val)
		// 			}
		// 		}
		// 	}
		case <-ticker.C:
			schedulePipes(ctx, nil)
		}
	}
}

// Stop stop working
func (s *Scheduler) Stop() {
	s.stop()
}

func (s *Scheduler) stop() {
	s.stopOnce.Do(func() {
		close(s.stopping)
	})
}

func scheduleOne(ctx context.Context, p *pipeline.Pipeline) (err error) {
	a := newAlgorithm(withAlgPipe(p))
	err = a.cal()
	if err != nil {
		return
	}
	if _, err = dao.UpdatePipelineBind(ctx, p.Name, a.bestNode.Name); err != nil {
		return
	}
	return
}

// noScheduledPipelines find pipelines that not be scheduled
func noScheduledPipelines(pb *model.PipelineBind) (pipes []*pipeline.Pipeline, err error) {
	if pb == nil {
		pb, err = dao.GetPipelineBind(context.Background())
		if err != nil {
			return
		}
	}
	pipes = []*pipeline.Pipeline{}
	for k, v := range pb.Bindings {
		if v == "" {
			pipes = append(pipes, &pipeline.Pipeline{Name: k})
		}
	}
	return
}

func schedulePipes(ctx context.Context, pb *model.PipelineBind) {
	pipes, err := noScheduledPipelines(pb)
	if err != nil {
		return
	}
	for _, v := range pipes {
		er := scheduleOne(ctx, v)
		if er != nil {
			event.EventErrorPipeline(v.Name, "SCHEDULING: "+er.Error())
		} else {
			event.EventInfoPipeline(v.Name, "SCHEDULING SUCCESS")
		}
	}
}
