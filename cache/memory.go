package cache

import (
	"slices"
	"sync"
	"time"

	"github.com/tinh-tinh/tinhtinh/utils"
)

type item[V any] struct {
	v V
	e uint32
}

type Memory[K string, V any] struct {
	max  int
	ttl  time.Duration
	data map[K]item[V]
	sync.RWMutex
}

type MemoryOptions struct {
	Max int
	Ttl time.Duration
}

func NewInMemory[K string, V interface{}](opt MemoryOptions) Store[K, V] {
	store := &Memory[K, V]{
		data: make(map[K]item[V]),
		ttl:  opt.Ttl,
		max:  opt.Max,
	}
	utils.StartTimeStampUpdater()
	go store.gc(1 * time.Second)
	return store
}

func (m *Memory[K, V]) Get(key K) *V {
	m.RLock()
	v, ok := m.data[key]
	m.RUnlock()

	if !ok || v.e != 0 && v.e <= utils.Timestamp() {
		return nil
	}
	return &v.v
}

func (m *Memory[K, V]) Set(key K, val V, ttl ...time.Duration) {
	var exp uint32
	if len(ttl) > 0 {
		exp = uint32(ttl[0].Seconds()) + utils.Timestamp()
	} else {
		exp = uint32(m.ttl.Seconds()) + utils.Timestamp()
	}
	i := item[V]{e: exp, v: val}
	for m.Count()+1 >= m.max {
		m.removeOldEle()
	}
	m.Lock()
	m.data[key] = i
	m.Unlock()
}

func (m *Memory[K, V]) Keys() []K {
	keys := make([]K, 0, len(m.data))
	for k := range m.data {
		keys = append(keys, k)
	}
	return keys
}

func (m *Memory[K, V]) Count() int {
	return len(m.Keys())
}

func (m *Memory[K, V]) removeOldEle() {
	smallest := ^uint32(0)
	var key K
	for k, v := range m.data {
		if v.e < smallest {
			smallest = v.e
			key = k
		}
	}

	m.Delete(key)
}

func (m *Memory[K, V]) Has(key K) bool {
	return slices.Contains(m.Keys(), key)
}

func (m *Memory[K, V]) Delete(key K) {
	m.Lock()
	delete(m.data, key)
	m.Unlock()
}

func (m *Memory[K, V]) Clear() {
	md := make(map[K]item[V])
	m.Lock()
	m.data = md
	m.Unlock()
}

func (m *Memory[K, V]) gc(sleep time.Duration) {
	ticker := time.NewTimer(sleep)
	defer ticker.Stop()
	var expired []K

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
