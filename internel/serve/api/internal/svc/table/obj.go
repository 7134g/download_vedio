package table

import "sync"

var ProxyCatchUrl = cmpMap[string, int]{
	lock: sync.RWMutex{},
	body: make(map[string]int),
}

// ProxyCatchHtmlTitle 用于获取title
var ProxyCatchHtmlTitle = cmpMap[string, string]{
	lock: sync.RWMutex{},
	body: make(map[string]string),
}
