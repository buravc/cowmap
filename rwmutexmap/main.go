package rwmutexmap

import (
	"sync"
)

type Rwmutexmap struct {
	innerMap map[interface{}]interface{}
	sync.RWMutex
}

func New() *Rwmutexmap {
	innerMap := make(map[interface{}]interface{})
	rwmutex := &Rwmutexmap{innerMap, sync.RWMutex{}}
	return rwmutex
}

func (h *Rwmutexmap) Set(k interface{}, v interface{}) {
	h.Lock()
	h.innerMap[k] = v
	h.Unlock()
}

func (h *Rwmutexmap) Get(k interface{}) interface{} {
	h.RLock()
	val := h.innerMap[k]
	h.RUnlock()
	return val
}
