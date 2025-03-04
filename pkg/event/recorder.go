package event

import (
	"context"
	"strings"
	"time"

	"github.com/golang/groupcache/lru"
	"github.com/jin06/binlogo/v2/configs"
	"github.com/jin06/binlogo/v2/pkg/store/dao"
	"github.com/jin06/binlogo/v2/pkg/store/model"
	"github.com/sirupsen/logrus"
)

// Recorder record event.
// Aggregate events and delete old events regularly
type Recorder struct {
	incomeChan chan *model.Event
	flushChan  chan *model.Event
	cache      *lru.Cache
	nodeName   string
	flushMap   map[string]*model.Event
}

// New Returns a new Recorder
func New() (*Recorder, error) {
	r := &Recorder{}
	r.incomeChan = make(chan *model.Event, 16384)
	r.flushChan = make(chan *model.Event, 4096)
	r.cache = lru.New(4096)
	r.nodeName = configs.Default.NodeName
	r.flushMap = map[string]*model.Event{}
	return r, nil
}

// Loop start event record loop
func (r Recorder) Loop(ctx context.Context) {
	go r._dispatch(ctx)
	go r._send(ctx)
}

func (r Recorder) _dispatch(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			{
				return
			}
		case e := <-r.incomeChan:
			{
				r.dispatch(e)
			}
		}
	}
}

func (r Recorder) _send(ctx context.Context) {
	ticker := time.NewTicker(time.Second * 10)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			{
				return
			}
		case e := <-r.flushChan:
			{
				r.flushMap[e.K] = e
				r.flush(false)
			}
		case <-ticker.C:
			{
				r.flush(true)
			}
		}
	}
}

func (r Recorder) flush(force bool) {
	if force || len(r.flushMap) >= 100 {
		for _, v := range r.flushMap {
			err := dao.UpdateEvent(context.Background(), v)
			if err != nil {
				logrus.Errorln("Write event to etcd failed: ", err)
			}
		}
		r.flushMap = map[string]*model.Event{}
	}
}

// Event pass event to income chan
func (r Recorder) Event(e *model.Event) {
	r.incomeChan <- e
}

func (r Recorder) dispatch(new *model.Event) {
	aggKey := aggregatorKey(new)
	oldInter, ok := r.cache.Get(aggKey)
	if ok {
		if old, is := oldInter.(*model.Event); is {
			if isExceedTime(time.Now(), old.FirstTime) {
				old.Count = old.Count + 1
				old.LastTime = time.Now()
				r.cache.Add(aggKey, old)
				r.update(old)
				return
			} else {
				r.cache.Remove(aggKey)
			}
		}
	}
	new.Count = 1
	//new.Key = key.GetKey()
	new.FirstTime = time.Now()
	new.LastTime = time.Now()
	new.NodeName = r.nodeName
	r.cache.Add(aggKey, new)
	r.update(new)
}

func (r Recorder) update(e *model.Event) {
	r.flushChan <- e
}

func aggregatorKey(e *model.Event) string {
	return strings.Join([]string{
		string(e.Type),
		string(e.ResourceType),
		e.ResourceName,
		e.Message,
		e.NodeName,
	},
		".",
	)
}

func isExceedTime(newTime time.Time, oldTime time.Time) bool {
	exceedTime := time.Minute * 5
	return oldTime.Add(exceedTime).Before(newTime)
}
