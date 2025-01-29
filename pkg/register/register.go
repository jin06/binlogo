package register

import (
	"context"
	"sync"
	"time"

	"github.com/jin06/binlogo/v2/pkg/store/dao"
	"github.com/jin06/binlogo/v2/pkg/store/model/pipeline"
	"github.com/sirupsen/logrus"
)

// New returns a new *Register
func New(opts ...Option) (r *Register) {
	r = &Register{
		ttl:     time.Second * 500,
		stopped: make(chan struct{}),
	}
	//r.client = etcdclient.Default()
	for _, v := range opts {
		v(r)
	}
	return
}

// Register Encapsulation of etcd register
type Register struct {
	ttl          time.Duration
	registerKey  string
	registerData *pipeline.Instance
	stopped      chan struct{}
	closeOnce    sync.Once
}

// func (r *Register) init() (err error) {
// 	logrus.Info("init register")
// 	if r.client, err = etcdclient.New(); err != nil {
// 		return
// 	}
// 	return
// }

// Run start register
func (r *Register) Run(ctx context.Context) (err error) {
	logrus.WithField("key", r.registerKey).Info("register run")
	defer func() {
		if re := recover(); re != nil {
			logrus.Errorln("register panic, ", re)
		}
		if err != nil {
			logrus.WithField("err", err.Error()).Infoln("register process quit unexpectedly")
		}
		logrus.WithField("key", r.registerKey).Info("register stopped")
		r.close()
	}()
	// if err = r.init(); err != nil {
	// 	return
	// }
	stx, cancel := context.WithCancel(ctx)
	defer cancel()
	if err = r.reg(stx); err != nil {
		logrus.Errorln(err)
		return
	}
	defer func() {
		errR := r.revoke(stx)
		if errR != nil {
			logrus.Errorln(errR)
		}
		logrus.Errorln("Register end: ", r.registerKey)
	}()

	keepTicker := time.NewTicker(time.Second)
	keepErrCount := 0
	defer keepTicker.Stop()
	watchTicker := time.NewTicker(time.Second)
	watchErrCount := 0
	defer watchTicker.Stop()
	for {
		select {
		case <-ctx.Done():
			{
				return
			}
		case <-watchTicker.C:
			{
				ok, err := r.watch(stx)
				if err != nil {
					watchErrCount++
					if watchErrCount >= 3 {
						logrus.Errorln("Pipeline instance watch failed")
						return err
					}
				}
				watchErrCount = 0
				if !ok {
					return err
				}
			}
		case <-keepTicker.C:
			{
				ok, err := r.keepOnce(stx)
				if err != nil {
					keepErrCount++
					if keepErrCount >= 3 {
						logrus.Errorln("Pipeline instance lease failed")
						return err
					}
				}
				keepErrCount = 0
				if !ok {
					return err
				}
			}
		}
	}
}

func (r *Register) reg(ctx context.Context) (err error) {
	return dao.RegisterInstance(ctx, r.registerData, r.ttl)
}

func (r *Register) keepOnce(ctx context.Context) (bool, error) {
	return dao.LeaseInstance(ctx, r.registerData.PipelineName, r.ttl)
}

func (r *Register) revoke(ctx context.Context) (err error) {
	return dao.UnRegisterInstance(ctx, r.registerData.PipelineName, r.registerData.NodeName)
}

func (r *Register) watch(ctx context.Context) (ok bool, err error) {
	ins, err := dao.GetInstance(ctx, r.registerData.PipelineName)
	if err != nil {
		return false, err
	}
	if ins.NodeName == r.registerData.NodeName {
		ok = true
	} else {
		ok = false
	}

	return
}

func (r *Register) close() error {
	r.closeOnce.Do(func() {
		close(r.stopped)
	})
	return nil
}

func (r *Register) Stopped() chan struct{} {
	return r.stopped
}
