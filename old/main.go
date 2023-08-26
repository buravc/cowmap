package old

import (
	"sync/atomic"
	"unsafe"
)

type cowmap struct {
	innerMap **map[interface{}]interface{}
}

func New() *cowmap {
	innerMap := new(map[interface{}]interface{})
	cowmap := &cowmap{&innerMap}
	return cowmap
}

func (h *cowmap) Set(k interface{}, v interface{}) {
	copyMap := make(map[interface{}]interface{})
	for k, v := range **h.innerMap {
		copyMap[k] = v
	}

	copyMap[k] = v

	pointerOfInnerMap := unsafe.Pointer(h.innerMap)
	pointerOfCopyMap := unsafe.Pointer(&copyMap)

	atomic.SwapPointer((*unsafe.Pointer)(pointerOfInnerMap), pointerOfCopyMap)
}

func (h *cowmap) Get(k interface{}) interface{} {
	return (**h.innerMap)[k]
}
