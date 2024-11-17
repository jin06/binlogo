package manager_event

import (
	"context"
	"strconv"
	"sync"
	"time"

	"github.com/jin06/binlogo/v2/app/server/node/manager"
	"github.com/jin06/binlogo/v2/pkg/store/dao"
	"github.com/jin06/binlogo/v2/pkg/store/dao/dao_pipe"
	"github.com/sirupsen/logrus"
)

// Manager is event manager
type Manager struct {
	*manager.Base
	Exit     bool
	stopOnce sync.Once
	stopping chan struct{}
}

// New returns a new Manager
func New() (m *Manager) {
	m = &Manager{
		Base:     &manager.Base{Status: manager.STOP},
		Exit:     false,
		stopping: make(chan struct{}),
	}
	return
}

// Run start working when become a leader
func (m *Manager) Run(ctx context.Context) (err error) {
	defer func() {
		m.Exit = true
	}()

	ticker := time.NewTicker(time.Hour * 60)
	defer ticker.Stop()
	for {
		select {
		case <-m.stopping:
			return
		case <-ctx.Done():
			{
				return
			}
		case <-ticker.C:
			{
				err1 := m.cleanHistoryEvent()
				if err1 != nil {
					logrus.Errorln("Clean event error: ", err1)
				}
			}
		}
	}
}

func (m *Manager) cleanHistoryEvent() (err error) {
	pipes, err := dao_pipe.AllPipelines(context.Background())
	if err != nil {
		return
	}
	nodes, err := dao.AllNodes()

	if err != nil {
		return
	}
	deleteTime := time.Now().Add(-time.Hour * 24 * 3).UnixNano()
	for _, pipe := range pipes {
		for _, n := range nodes {
			from := n.Name + ".0"
			end := n.Name + "." + strconv.FormatInt(deleteTime, 10)
			prefix := dao.PipelinePrefix() + "/" + pipe.Name + "/"

			_, err1 := dao.DeleteRangeEvent(context.Background(), prefix+from, prefix+end)
			if err1 != nil {
				logrus.Errorln("Clean history events error: ", err1)
			}
			prefix = dao.EventNodePrefix() + "/" + n.Name + "/"
			_, err2 := dao.DeleteRangeEvent(context.Background(), prefix+from, prefix+end)
			if err2 != nil {
				logrus.Errorln("Clean history events error: ", err2)
			}
		}
	}
	return
}

// Stop stop working when lost leader
func (m *Manager) Stop() {
	m.stop()
}

func (m *Manager) stop() {
	m.stopOnce.Do(func() {
		close(m.stopping)
	})
}
