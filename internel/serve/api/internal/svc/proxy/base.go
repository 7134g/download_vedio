package proxy

import (
	"cmp"
	"sync"
	"time"
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

type hostTable struct {
	lock sync.RWMutex

	reqTimeMap map[string]int64
	resTimeMap map[string]int64
	urlToBody  map[string][]byte
}

func newHostTable() hostTable {
	return hostTable{
		lock:       sync.RWMutex{},
		reqTimeMap: make(map[string]int64),
		resTimeMap: make(map[string]int64),
		urlToBody:  make(map[string][]byte),
	}
}

func (m *hostTable) AddBody(url string, body []byte) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	m.urlToBody[url] = body
}

func (m *hostTable) AddReqUrl(url string) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	m.reqTimeMap[url] = time.Now().Unix()
}

func (m *hostTable) AddResUrl(url string) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	m.resTimeMap[url] = time.Now().Unix()
}

func (m *hostTable) Find(s message) [][]byte {
	m.lock.RLock()
	defer m.lock.RUnlock()

	reqClock := m.reqTimeMap[s.source]
	var urls = make([]string, 0)
	for k, v := range m.reqTimeMap {
		if reqClock-v > 0 && reqClock-v < 10 {
			urls = append(urls, k)
		}
	}

	var bs = make([][]byte, 0)
	for _, url := range urls {
		body, exist := m.urlToBody[url]
		if !exist {
			_, ok := m.reqTimeMap[url]
			if ok {
				// 等待全部页面响应
				s.sleep = time.Second
				sourceChan <- s
				return nil
			}
			continue
		}
		bs = append(bs, body)
	}

	return bs
}
