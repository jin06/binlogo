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
		closing: make(chan struct{}),
		closed:  make(chan struct{}),
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
	isClosed     bool
	closing      chan struct{}
	closed       chan struct{}
	closeOnce    sync.Once
	completeOnce sync.Once
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
	defer r.CompleteClose()
	defer r.Close()
	logrus.WithField("key", r.registerKey).Info("register run")
	defer func() {
		if re := recover(); re != nil {
			logrus.Errorln("register panic, ", re)
		}
		if err != nil {
			logrus.WithField("err", err.Error()).Infoln("register process quit unexpectedly")
		}
		logrus.WithField("key", r.registerKey).Info("register stopped")
	}()
	// if err = r.init(); err != nil {
	// 	return
	// }
	if err = r.reg(ctx); err != nil {
		logrus.Errorln(err)
		return
	}
	// defer func() {
	// 	errR := r.revoke(stx)
	// 	if errR != nil {
	// 		logrus.Errorln(errR)
	// 	}
	// 	logrus.Errorln("Register end: ", r.registerKey)
	// }()

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
				ok, err := r.watch(ctx)
				if err != nil {
					watchErrCount++
					if watchErrCount >= 3 {
						logrus.Errorln("Pipeline instance watch failed")
						return err
					}
				}
				watchErrCount = 0
				if !ok {
					logrus.WithField("Register Key", r.registerKey).Debug("Register exit, watch none")
					return err
				}
			}
		case <-keepTicker.C:
			{
				if err := r.keepOnce(ctx); err != nil {
					keepErrCount++
					if keepErrCount >= 3 {
						logrus.Errorln("Pipeline instance lease failed")
						return err
					}
				}
			}
		}
	}
}

func (r *Register) reg(ctx context.Context) (err error) {
	return dao.RegisterInstance(ctx, r.registerData, r.ttl)
}

func (r *Register) keepOnce(ctx context.Context) error {
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

func (r *Register) Close() error {
	r.closeOnce.Do(func() {
		r.revoke(context.Background())
		close(r.closing)
	})
	return nil
}

func (r *Register) CompleteClose() {
	r.completeOnce.Do(func() {
		close(r.closed)
	})
}

func (r *Register) Closed() chan struct{} {
	return r.closed
}

func (r *Register) IsClosed() bool {
	return r.isClosed
}
