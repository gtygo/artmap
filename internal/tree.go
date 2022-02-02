package internal

import (
	"sync/atomic"
	"unsafe"
)

type Art struct {
	count uint64
	root  unsafe.Pointer
}

func NewTree() *Art {
	return &Art{
		count: 0,
		root:  unsafe.Pointer(makeN4()),
	}
}

func (t *Art) Get(key []byte) (interface{}, bool) {
	for {
		n := (*n)(atomic.LoadPointer(&t.root))
		v, ok := n.search(key, 0, nil, 0)
		if ok {
			if v != nil {
				return v, true
			}
			return nil, false
		}
	}
}

func (t *Art) Set(key []byte, value interface{}) {








}

func (t *Art) Remove(key []byte) {

}

func (t *Art) Count() uint64 {
	return t.count
}
