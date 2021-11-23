package manager_event

import (
	"context"
	"github.com/jin06/binlogo/app/server/node/manager"
	"github.com/jin06/binlogo/pkg/store/dao/dao_event"
	"github.com/jin06/binlogo/pkg/store/dao/dao_node"
	"github.com/jin06/binlogo/pkg/store/dao/dao_pipe"
	"github.com/sirupsen/logrus"
	"strconv"
	"time"
)

// Manager is event manager
type Manager struct {
	cancel context.CancelFunc
	*manager.Base
}

// New returns a new Manager
func New() (m *Manager) {
	m = &Manager{
		nil, &manager.Base{Status: manager.STOP},
	}
	return
}

// Run start working when become a leader
func (m *Manager) Run(ctx context.Context) (err error) {
	m.Lock()
	defer m.Unlock()
	if m.Status == manager.START {
		return
	}
	myCtx, cancel := context.WithCancel(ctx)
	m.cancel = cancel
	m.Status = manager.START
	go func() {
		defer func() {
			m.Lock()
			m.Unlock()
			m.Status = manager.STOP
		}()
		go m.cleanHistoryEvent()
		for {
			select {
			case <-myCtx.Done():
				{
					return
				}
			case <-ctx.Done():
				{
					return
				}
			case <-time.Tick(time.Hour * 60):
				{
					err1 := m.cleanHistoryEvent()
					if err1 != nil {
						logrus.Errorln("Clean event error: ", err1)
					}
				}
			}
		}
	}()
	return
}

func (m *Manager) cleanHistoryEvent() (err error) {
	pipes, err := dao_pipe.AllPipelines()
	if err != nil {
		return
	}
	nodes, err := dao_node.AllNodes()

	if err != nil {
		return
	}
	deleteTime := time.Now().Add(-time.Hour * 24 * 3).UnixNano()
	for _, pipe := range pipes {
		for _, n := range nodes {
			from := n.Name + ".0"
			end := n.Name + "." + strconv.FormatInt(deleteTime, 10)
			prefix := dao_event.PipelinePrefix() + "/" + pipe.Name + "/"
			_, err1 := dao_event.DeleteRange(prefix+from, prefix+end)
			if err1 != nil {
				logrus.Errorln("Clean history events error: ", err1)
			}
			prefix = dao_event.NodePrefix() + "/" + n.Name + "/"
			_, err2 := dao_event.DeleteRange(prefix+from, prefix+end)
			if err2 != nil {
				logrus.Errorln("Clean history events error: ", err2)
			}
		}
	}
	return
}

// Stop stop working when lost leader
func (m *Manager) Stop() {
	m.Lock()
	defer m.Unlock()
	if m.Status == manager.STOP {
		return
	}
	//if m.cancel != nil {
	//	m.cancel()
	//}
	m.cancel()
	m.Status = manager.STOP
	return
}
