package watcher

import (
	"encoding/json"
	"fmt"

	"github.com/jin06/binlogo/v2/pkg/store/model"
	"github.com/jin06/binlogo/v2/pkg/store/model/node"
	"github.com/jin06/binlogo/v2/pkg/store/model/pipeline"
	"github.com/sirupsen/logrus"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func WrapStrHandler() Handler {
	return func(e *clientv3.Event) (ev *Event, err error) {
		ev = &Event{}
		var m string
		ev.Event = e
		ev.Data = m
		if e.Type == mvccpb.DELETE {

		} else {
			err = json.Unmarshal(e.Kv.Value, &m)
			if err != nil {
				logrus.Error(err)
				return
			}
		}
		return
	}
}

func WrapSchedulerBinding() Handler {
	return func(e *clientv3.Event) (ev *Event, err error) {
		ev = &Event{}
		m := &model.PipelineBind{}
		ev.Event = e
		ev.Data = m
		if e.Type != mvccpb.DELETE {
			err = m.Unmarshal(e.Kv.Value)
			if err != nil {
				logrus.Error(err)
				return
			}
		}
		return
	}
}

func WrapPipeline(prefix string, pipeName string) Handler {
	return func(e *clientv3.Event) (ev *Event, err error) {
		ev = &Event{}
		m := &pipeline.Pipeline{}
		ev.Event = e
		ev.Data = m
		if e.Type == mvccpb.DELETE {
			if pipeName == "" {
				_, err = fmt.Sscanf(string(e.Kv.Key), prefix+"/%s", &m.Name)
				if err != nil {
					logrus.Error(err)
					return
				}
			} else {
				m.Name = pipeName
			}
		} else {
			err = m.Unmarshal(e.Kv.Value)
			if err != nil {
				logrus.Error(err)
				return
			}
		}
		return
	}
}

func WrapNodeStatus(prefix string, name string) Handler {
	return func(e *clientv3.Event) (ev *Event, err error) {
		ev = &Event{}
		m := &node.Status{}
		ev.Event = e
		ev.Data = m
		if e.Type == mvccpb.DELETE {
			if name == "" {
				_, err = fmt.Sscanf(string(e.Kv.Key), prefix+"/%s", &m.NodeName)
				if err != nil {
					return
				}
			} else {
				m.NodeName = name
			}
		} else {
			err = json.Unmarshal(e.Kv.Value, m)
			if err != nil {
				return
			}
		}
		return
	}
}

func WrapNode(prefix string, nodeName string) Handler {
	return func(e *clientv3.Event) (ev *Event, err error) {
		ev = &Event{}
		m := &node.Node{}
		ev.Event = e
		ev.Data = m
		if e.Type == mvccpb.DELETE {
			if nodeName == "" {
				_, err = fmt.Sscanf(string(e.Kv.Key), prefix+"/%s", &m.Name)
				if err != nil {
					return
				}
			} else {
				m.Name = nodeName
			}
		} else {
			err = m.Unmarshal(e.Kv.Value)
			if err != nil {
				return
			}
		}
		return
	}
}

func WrapInstance(prefix string, pipeName string) Handler {
	return func(e *clientv3.Event) (ev *Event, err error) {
		ev = &Event{}
		m := &pipeline.Instance{}
		ev.Event = e
		ev.Data = m
		if e.Type == mvccpb.DELETE {
			if pipeName == "" {
				_, err = fmt.Sscanf(string(e.Kv.Key), prefix+"/%s", &m.PipelineName)
				if err != nil {
					logrus.Error(err)
					return
				}
			} else {
				m.PipelineName = pipeName
			}
		} else {
			err = json.Unmarshal(e.Kv.Value, &m)
			if err != nil {
				logrus.Error(err)
				return
			}
		}
		return
	}
}
