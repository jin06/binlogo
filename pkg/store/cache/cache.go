package cache

import (
	"encoding/json"
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
	data  map[string]map[string][]byte
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
		map[string]map[string][]byte{},
		map[string]*menu{},
		sync.RWMutex{},
	}

	for _, v := range validMenus {
		cm.data[v] = map[string][]byte{}
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
	var newData map[string]map[string][]byte
	if err := json.NewDecoder(serialized).Decode(newData); err != nil {
		return err
	}

	c.Lock()
	defer c.Unlock()
	c.data = newData

	return nil
}

func (c *CacheManager) Set(m string, key string, val []byte) error {

	c.menus[m].Lock()
	defer c.menus[m].Unlock()
	c.data[m][key] = val
	return nil
}

func (c *CacheManager) Get(menu string, key string) []byte {
	if _, ok := c.menus[menu]; !ok {
		return nil
	}

	c.menus[menu].RLock()
	defer c.menus[menu].RUnlock()
	return c.data[menu][key]
}

func (c *CacheManager) SetEntry(e Entry) (err error) {
	return c.Set(e.Menu, e.Key, e.Value)
}

func (c *CacheManager) GetEntry(menu string, key string) (e *Entry, err error) {
	b := c.Get(menu, key)
	e = &Entry{}
	err = json.Unmarshal(b, e)
	return
}

func (c *CacheManager) check(m string, key string) (ok bool, err error) {
	if _, ok = c.menus[m]; !ok {
		return
	}
	return
}
