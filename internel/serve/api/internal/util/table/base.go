package table

import (
	"cmp"
	"sync"
)

type dataType interface {
	cmp.Ordered
}

type cmpMap[D dataType] struct {
	lock sync.RWMutex

	body map[string]D
}

func (m *cmpMap[D]) Set(key string, value D) {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.body[key] = value
}

func (m *cmpMap[D]) Get(key string) (D, bool) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	value, exist := m.body[key]
	return value, exist
}

type sliceType interface {
	[]byte | []string | []int
}

type sliceMap[D sliceType] struct {
	lock sync.RWMutex

	body map[string]D
}

func (m *sliceMap[D]) Set(key string, value D) {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.body[key] = value
}

func (m *sliceMap[D]) Get(key string) (D, bool) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	value, exist := m.body[key]
	return value, exist
}
