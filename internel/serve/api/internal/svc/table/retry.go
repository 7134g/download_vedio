package table

import "sync"

var errorRetry sync.Map

func IncErrCount(key string) {
	count, exist := errorRetry.Load(key)
	if exist {
		errorRetry.Store(key, count.(uint)+1)
	} else {
		errorRetry.Store(key, uint(1))
	}
}

func GetErrCount(key string) uint {
	value, ok := errorRetry.Load(key)
	if ok {
		return value.(uint)
	} else {
		return 0
	}
}
