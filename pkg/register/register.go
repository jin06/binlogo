package register

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/jin06/binlogo/v2/pkg/store/dao"
	"github.com/jin06/binlogo/v2/pkg/store/model/pipeline"
	"github.com/sirupsen/logrus"
)

// New returns a new *Register
func New(opts ...Option) (r *Register) {
	r = &Register{
		ttl:     time.Second * 5,
		closing: make(chan struct{}),
		closed:  make(chan struct{}),
	}
	//r.client = etcdclient.Default()
	for _, v := range opts {
		v(r)
	}
	r.log = logrus.WithField("RegisterKey", r.registerKey)

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
	log          *logrus.Entry
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
	// logrus.WithField("key", r.registerKey).Info("Register run")
	r.log.Info("Register run")
	defer func() {
		if re := recover(); re != nil {
			logrus.Errorln("register panic, ", re)
		}
		if err != nil {
			logrus.WithField("err", err.Error()).Infoln("register process quit unexpectedly")
		}
		// logrus.WithField("key", r.registerKey).Info("register stopped")
		r.log.Info("Register end")
	}()

	if err = r.reg(ctx); err != nil {
		r.log.Errorln(err)
		return
	}

	keepTicker := time.NewTicker(time.Second)
	defer keepTicker.Stop()
	watchTicker := time.NewTicker(time.Second)
	defer watchTicker.Stop()
	for {
		select {
		case <-ctx.Done():
			{
				return
			}
		case <-watchTicker.C:
			{
				err := r.watch(ctx)
				if err != nil {
					r.log.WithError(err).Errorln("Register watch error")
					return err
				}
			}
		case <-keepTicker.C:
			{
				if err := r.keepOnce(ctx); err != nil {
					r.log.WithError(err).Errorln("Pipeline instance lease failed")
					return err
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
	r.log.Debug("Register revoke")
	return dao.UnRegisterInstance(ctx, r.registerData.PipelineName, r.registerData.NodeName)
}

func (r *Register) watch(ctx context.Context) error {
	ins, err := dao.GetInstance(ctx, r.registerData.PipelineName)
	if err != nil {
		return err
	}
	if ins.NodeName != r.registerData.NodeName {
		return errors.New("register node is not expected")
	}
	return nil
}

func (r *Register) Close() error {
	r.closeOnce.Do(func() {
		logrus.Debug("Register exit")
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
