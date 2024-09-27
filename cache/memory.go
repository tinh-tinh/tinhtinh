package cache

import (
	"slices"
	"sync"
	"time"

	"github.com/tinh-tinh/tinhtinh/utils"
)

type item struct {
	v interface{}
	e uint32
}

type Memory struct {
	max  int
	ttl  time.Duration
	data map[string]item
	sync.RWMutex
}

type MemoryOptions struct {
	Max int
	Ttl time.Duration
}

func NewInMemory(opt MemoryOptions) Store {
	store := &Memory{
		data: make(map[string]item),
		ttl:  opt.Ttl,
		max:  opt.Max,
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

func (m *Memory) Set(key string, val interface{}, ttl ...time.Duration) {
	var exp uint32
	if len(ttl) > 0 {
		exp = uint32(ttl[0].Seconds()) + utils.Timestamp()
	} else {
		exp = uint32(m.ttl.Seconds()) + utils.Timestamp()
	}
	i := item{e: exp, v: val}
	for m.Count()+1 > m.max {
		m.removeOldEle()
	}
	m.Lock()
	m.data[key] = i
	m.Unlock()
}

func (m *Memory) Keys() []string {
	keys := make([]string, 0, len(m.data))
	for k := range m.data {
		keys = append(keys, k)
	}
	return keys
}

func (m *Memory) Count() int {
	return len(m.Keys())
}

func (m *Memory) removeOldEle() {
	smallest := ^uint32(0)
	var key string
	for k, v := range m.data {
		if v.e < smallest {
			smallest = v.e
			key = k
		}
	}

	m.Delete(key)
}

func (m *Memory) Has(key string) bool {
	return slices.Contains(m.Keys(), key)
}

func (m *Memory) Delete(key string) {
	m.Lock()
	delete(m.data, key)
	m.Unlock()
}

func (m *Memory) Clear() {
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
