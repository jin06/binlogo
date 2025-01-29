package manager_pipe

import (
	"context"
	"sync"
	"time"

	"github.com/jin06/binlogo/v2/pkg/store/dao"
	"github.com/jin06/binlogo/v2/pkg/store/model"
	"github.com/jin06/binlogo/v2/pkg/store/model/node"
	"github.com/sirupsen/logrus"
)

// Manager is pipeline instance manager.
// If there is a pipeline bound to the current node, run an instance
// will also stop inbound instance on the current node
type Manager struct {
	mapping      map[string]bool
	mappingIns   map[string]*instance
	node         *node.Node
	mutex        sync.Mutex
	closing      chan struct{}
	closeOnce    sync.Once
	closed       chan struct{}
	completeOnce sync.Once
}

// New returns a new Manager
func New(n *node.Node) (m *Manager) {
	m = &Manager{
		mapping:    map[string]bool{},
		mappingIns: map[string]*instance{},
		// mutex:      sync.Mutex{},
		closed:  make(chan struct{}),
		closing: make(chan struct{}),
	}
	m.node = n
	logrus.Debug("New pipeline manager ")
	logrus.Debug("NodeName:", n.Name)
	return
}

// Run start woking
func (m *Manager) Run(ctx context.Context) {
	defer m.CompleteClose()
	defer m.Close()
	defer func() {
		if r := recover(); r != nil {
			logrus.Errorln("pipeline manager panic, ", r)
			panic(r)
		}
	}()
	if err := m.scanPipelines(nil); err != nil {
		logrus.Error(err)
	}
	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()
	for {
		select {
		case <-m.closing:
			return
		case <-ctx.Done():
			{
				return
			}
		case <-ticker.C:
			{
				if serr := m.scanPipelines(nil); serr != nil {
					logrus.Error(serr)
				}
				m.dispatch(ctx)
			}
		}
	}
}

// scanPipelines scan pipeline bind, find pipelines that should run in this node
func (m *Manager) scanPipelines(pb *model.PipelineBind) (err error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if pb == nil {
		pb, err = dao.GetPipelineBind(context.Background())
		if err != nil {
			return
		}
	}
	for pipeName, nodeName := range pb.Bindings {
		if nodeName == m.node.Name {
			m.mapping[pipeName] = true
			continue
		}
		if nodeName != m.node.Name {
			m.mapping[pipeName] = false
		}
	}

	for pName := range m.mapping {
		if _, ok := pb.Bindings[pName]; !ok {
			m.mapping[pName] = false
		}
	}
	return
}

func (m *Manager) dispatch(ctx context.Context) {
	for pName, shouldRun := range m.mapping {
		// _, isExist := m.mappingIns[pName]
		ins := m.Get(pName)
		if shouldRun && ins == nil {
			m.Add(ctx, pName)
			go func() {
				err := m.mappingIns[pName].start(ctx)
				if err != nil {
					logrus.WithError(err).Error("Pipeline instance with error")
				}
			}()
		}
		if !shouldRun && ins != nil {
			m.mappingIns[pName].Close()
			<-m.mappingIns[pName].Closed()
			// m.Remove(pName)
		}
	}
}

func (m *Manager) Close() {
	m.closeOnce.Do(func() {
		close(m.closing)
	})
}

func (m *Manager) Closed() chan struct{} {
	return m.closed
}

func (m *Manager) CompleteClose() {
	m.completeOnce.Do(func() {
		close(m.closed)
	})
}

func (m *Manager) Get(pipeName string) *instance {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if ins, ok := m.mappingIns[pipeName]; ok {
		return ins
	}
	return nil
}

func (m *Manager) Add(ctx context.Context, pipeName string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	ins := newInstance(pipeName, m.node.Name, m)

	m.mappingIns[pipeName] = ins
}

func (m *Manager) Remove(pipeName string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	delete(m.mappingIns, pipeName)
}
