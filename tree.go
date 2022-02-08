package artmap

import (
	"sync/atomic"
	"unsafe"
)

type tree struct {
	count uint64
	root  unsafe.Pointer
}

func New() *tree {
	return &tree{
		count: 0,
		root:  unsafe.Pointer(MakeN4()),
	}
}

func (t *tree) Get(key []byte) (interface{}, bool) {
	for {
		n := (*n)(atomic.LoadPointer(&t.root))
		v, ok := n.LookupOpt(key, 0, nil, 0)
		if ok {
			if v != nil {
				return v, true
			}
			return nil, false
		}
	}
}

func (t *tree) Set(key []byte, value interface{}) {
	for {
		n := (*n)(atomic.LoadPointer(&t.root))
		if n.InsertOpt(key, value, 0, nil, 0, &t.root) {
			atomic.AddUint64(&t.count, 1)
			return
		}
	}
}

func (t *tree) Remove(key []byte) {

}

func (t *tree) Pop(key []byte) (interface{}, bool) {
	return nil, false
}

func (t *tree) Count() uint64 {
	return atomic.LoadUint64(&t.count)
}

func (t *tree) Clear() {

}

func (t *tree) Has() {

}

func (t *tree) IsEmpty() {

}
