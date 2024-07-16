package table

import (
	"cmp"
	"sync"
)

type dataType interface {
	cmp.Ordered
}

type cmpMap[K, D dataType] struct {
	lock sync.RWMutex

	body map[K]D
}

func (m *cmpMap[K, D]) Set(key K, value D) {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.body[key] = value
}

func (m *cmpMap[K, D]) Get(key K) (D, bool) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	value, exist := m.body[key]
	return value, exist
}

func (m *cmpMap[K, D]) Inc(key K, count D) {
	m.lock.Lock()
	defer m.lock.Unlock()

	value, exist := m.body[key]
	if exist {
		m.body[key] = count + value
	} else {
		m.body[key] = count
	}
}

func (m *cmpMap[K, D]) Del(key K) {
	m.lock.Lock()
	defer m.lock.Unlock()
	delete(m.body, key)
}

func (m *cmpMap[K, D]) Each(f func(key K, value D)) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	for key, value := range m.body {
		f(key, value)
	}

}

type sliceType interface {
	[]byte | []string | []int
}

type sliceMap[K dataType, D sliceType] struct {
	lock sync.RWMutex

	body map[K]D
}

func (m *sliceMap[K, D]) Set(key K, value D) {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.body[key] = value
}

func (m *sliceMap[K, D]) Get(key K) (D, bool) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	value, exist := m.body[key]
	return value, exist
}

func (m *sliceMap[K, D]) Each(f func(key K, value D)) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	for key, value := range m.body {
		f(key, value)
	}

}
