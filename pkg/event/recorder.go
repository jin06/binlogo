package event

import (
	"context"
	"github.com/golang/groupcache/lru"
	"github.com/jin06/binlogo/pkg/store/dao/dao_event"
	"github.com/jin06/binlogo/pkg/store/model/event"
	"github.com/spf13/viper"
	"strings"
	"time"
)

type Recorder struct {
	incomeChan chan *event.Event
	flushChan  chan *event.Event
	cache      *lru.Cache
	nodeName   string
	flushMap   map[string]*event.Event
}

func New() (*Recorder, error) {
	r := &Recorder{}
	r.incomeChan = make(chan *event.Event, 16384)
	r.flushChan = make(chan *event.Event, 4096)
	r.cache = lru.New(4096)
	r.nodeName = viper.GetString("node.name")
	r.flushMap = map[string]*event.Event{}
	return r, nil
}

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
	for {
		select {
		case <-ctx.Done():
			{
				return
			}
		case e := <-r.flushChan:
			{
				r.flushMap[e.Key] = e
				r.flush(false)
			}
		case <-time.Tick(time.Second * 10):
			{
				r.flush(true)
			}
		}
	}
}

func (r Recorder) flush(force bool) {
	if force || len(r.flushMap) >= 100 {
		for _, v := range r.flushMap {
			dao_event.Update(v)
		}
		r.flushMap = map[string]*event.Event{}
	}
}

func (r Recorder) Event(e *event.Event) {
	r.incomeChan <- e
	return
}

func (r Recorder) dispatch(new *event.Event) {
	aggKey := aggregatorKey(new)
	oldInter, ok := r.cache.Get(aggKey)
	if ok {
		if old, is := oldInter.(*event.Event); is {
			if time.Now().Unix()-old.FirstTime.Unix() < int64(time.Minute*5) {
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

func (r Recorder) update(e *event.Event) {
	r.flushChan <- e
}

func aggregatorKey(e *event.Event) string {
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
