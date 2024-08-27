package cache

import (
	"sync"
	"time"

	"github.com/tinh-tinh/tinhtinh/utils"
)

type item struct {
	v any
	e uint32
}

type Memory struct {
	data map[string]item
	sync.RWMutex
}

func NewInMemory() *Memory {
	store := &Memory{
		data: make(map[string]item),
	}
	utils.StartTimeStampUpdater()
	go store.gc(1 * time.Second)
	return store
}

func (m *Memory) Get(key string) interface{} {
	m.RLock()
	v, ok := m.data[key]
	m.RUnlock()

	if !ok || v.e != 0 && v.e <= utils.Timestamp() {
		return nil
	}
	return v.v
}

func (m *Memory) Set(key string, val interface{}, ttl time.Duration) {
	var exp uint32
	if ttl > 0 {
		exp = uint32(ttl.Seconds()) + utils.Timestamp()
	}
	i := item{e: exp, v: val}
	m.Lock()
	m.data[key] = i
	m.Unlock()
}

func (m *Memory) Delete(key string) {
	m.Lock()
	delete(m.data, key)
	m.Unlock()
}

func (m *Memory) Reset() {
	md := make(map[string]item)
	m.Lock()
	m.data = md
	m.Unlock()
}

func (m *Memory) gc(sleep time.Duration) {
	ticker := time.NewTimer(sleep)
	defer ticker.Stop()
	var expired []string

	for range ticker.C {
		ts := utils.Timestamp()
		expired = expired[:0]
		m.RLock()
		for key, v := range m.data {
			if v.e != 0 && v.e <= ts {
				expired = append(expired, key)
			}
		}
		m.RUnlock()
		m.Lock()

		for i := range expired {
			v := m.data[expired[i]]
			if v.e != 0 && v.e <= ts {
				delete(m.data, expired[i])
			}
		}
		m.Unlock()
	}
}
