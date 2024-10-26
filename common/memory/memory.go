package memory

import (
	"slices"
	"sync"
	"time"

	"github.com/tinh-tinh/tinhtinh/common"
	"github.com/tinh-tinh/tinhtinh/common/era"
)

type item struct {
	v interface{}
	e uint32
}

type Store struct {
	max  int
	ttl  time.Duration
	data map[string]item
	sync.RWMutex
}

type Options struct {
	Max int
	Ttl time.Duration
}

func New(opt Options) *Store {
	if opt.Max == 0 {
		opt.Max = common.MaxInt
	}
	store := &Store{
		data: make(map[string]item),
		ttl:  opt.Ttl,
		max:  opt.Max,
	}
	era.StartTimeStampUpdater()
	go store.gc(1 * time.Second)
	return store
}

func (m *Store) Get(key string) interface{} {
	m.RLock()
	v, ok := m.data[key]
	m.RUnlock()

	if !ok || v.e != 0 && v.e <= era.Timestamp() {
		return nil
	}
	return v.v
}

func (m *Store) Set(key string, val interface{}, ttl ...time.Duration) {
	var exp uint32
	if len(ttl) > 0 {
		exp = uint32(ttl[0].Seconds()) + era.Timestamp()
	} else {
		exp = uint32(m.ttl.Seconds()) + era.Timestamp()
	}
	i := item{e: exp, v: val}
	for m.Count()+1 > m.max {
		m.removeOldEle()
	}
	m.Lock()
	m.data[key] = i
	m.Unlock()
}

func (m *Store) Keys() []string {
	keys := make([]string, 0, len(m.data))
	for k := range m.data {
		keys = append(keys, k)
	}
	return keys
}

func (m *Store) Count() int {
	return len(m.Keys())
}

func (m *Store) removeOldEle() {
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

func (m *Store) Has(key string) bool {
	return slices.Contains(m.Keys(), key)
}

func (m *Store) Delete(key string) {
	m.Lock()
	delete(m.data, key)
	m.Unlock()
}

func (m *Store) Clear() {
	md := make(map[string]item)
	m.Lock()
	m.data = md
	m.Unlock()
}

func (m *Store) gc(sleep time.Duration) {
	ticker := time.NewTimer(sleep)
	defer ticker.Stop()
	var expired []string

	for range ticker.C {
		ts := era.Timestamp()
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
