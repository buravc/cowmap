package cowmap

import (
	"sync/atomic"
	"unsafe"
)

type Cowmap struct {
	innerMap unsafe.Pointer
}

func New() *Cowmap {
	innerMap := new(map[interface{}]interface{})
	cowmap := &Cowmap{unsafe.Pointer(&innerMap)}
	return cowmap
}

func (h *Cowmap) Set(k interface{}, v interface{}) {
	copyMap := make(map[interface{}]interface{})
	innerMap := *(*map[interface{}]interface{})(atomic.LoadPointer(&h.innerMap))

	for k, v := range innerMap {
		copyMap[k] = v
	}
	copyMap[k] = v

	cptr := unsafe.Pointer(&copyMap)

	atomic.StorePointer(&h.innerMap, cptr)
}

func (h *Cowmap) Get(k interface{}) interface{} {
	innerMap := *(*map[interface{}]interface{})(atomic.LoadPointer(&h.innerMap))
	return innerMap[k]
}

func (h *Cowmap) Iterate(f func(k interface{}, v interface{})) {
	innerMap := *(*map[interface{}]interface{})(atomic.LoadPointer(&h.innerMap))
	for key, val := range innerMap {
		f(key, val)
	}
}
