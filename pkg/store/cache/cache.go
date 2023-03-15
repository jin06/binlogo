package cache

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"sync"
)

var defaultCM *CacheManager

func Default() *CacheManager {
	if defaultCM == nil {
		defaultCM = NewCacheManager()
	}
	return defaultCM
}

type CacheManager struct {
	data  map[string]map[string]Entry
	menus map[string]*menu
	sync.RWMutex
}

type menu struct {
	sync.RWMutex
}

type Entry struct {
	Menu     string
	Key      string
	Value    []byte
	Revision uint64
	Delete   bool
}

var validMenus = []string{
	Pipeline, Node, Event, Scheduler,
}

const Pipeline = "pipeline"
const Node = "node"
const Event = "event"
const Scheduler = "scheduler"

func NewCacheManager() *CacheManager {
	cm := &CacheManager{
		map[string]map[string]Entry{},
		map[string]*menu{},
		sync.RWMutex{},
	}

	for _, v := range validMenus {
		cm.data[v] = map[string]Entry{}
		cm.menus[v] = &menu{
			sync.RWMutex{},
		}
	}
	return cm
}

// Marshal serializes cache data
func (c *CacheManager) Marshal() ([]byte, error) {
	c.RLock()
	defer c.RUnlock()
	dataBytes, err := json.Marshal(c.data)
	return dataBytes, err
}

// UnMarshal deserializes cache data
func (c *CacheManager) UnMarshal(serialized io.ReadCloser) error {
	var newData map[string]map[string]Entry
	if err := json.NewDecoder(serialized).Decode(newData); err != nil {
		return err
	}

	c.Lock()
	defer c.Unlock()
	c.data = newData

	return nil
}

func (c *CacheManager) GetEntry(menu string, key string) (e *Entry) {
	c.menus[menu].RLock()
	defer c.menus[menu].RUnlock()
	if !c.check(menu, key) {
		return
	}
	e = &Entry{}
	*e = c.data[menu][key]
	return
}

func (c *CacheManager) GetEntries(menu string) (list *map[string]Entry) {
	c.menus[menu].RLock()
	defer c.menus[menu].RUnlock()
	_, ok := c.data[menu]
	if !ok {
		return
	}
	list = &map[string]Entry{}

	*list = c.data[menu]
	return
}

func (c *CacheManager) SetEntry(e Entry) (err error) {
	c.menus[e.Menu].Lock()
	defer c.menus[e.Menu].Unlock()
	if old, has := c.data[e.Menu][e.Key]; has {
		if e.Revision <= old.Revision {
			err = errors.New(fmt.Sprintf("Revision conflict, current revision is %d, request revision is %d. Menu: %s, Key: %s", old.Revision, e.Revision, e.Menu, e.Key))
			return
		}
	}
	c.data[e.Menu][e.Key] = e
	return
}

func (c *CacheManager) check(m string, key string) (ok bool) {
	if _, ok = c.data[m]; !ok {
		return
	}
	_, ok = c.data[m][key]
	return

}
