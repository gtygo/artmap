package artmap

import (
	"sync/atomic"
	"unsafe"
)

type Tree struct {
	count uint64
	root  unsafe.Pointer
}

func New() *Tree {
	return &Tree{
		count: 0,
		root:  unsafe.Pointer(makeN4()),
	}
}

func (t *Tree) Get(key []byte) (interface{}, bool) {
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

func (t *Tree) Set(key []byte, value interface{}) {
	for {
		n := (*n)(atomic.LoadPointer(&t.root))
		r, a := n.InsertOpt(key, value, 0, nil, 0, &t.root);
		if r {
			if a == 1 {
				atomic.AddUint64(&t.count, 1)
			}
			return
		}
	}
}

func (t *Tree) Remove(key []byte) {

}

func (t *Tree) Pop(key []byte) (interface{}, bool) {
	return nil, false
}

func (t *Tree) Count() uint64 {
	return atomic.LoadUint64(&t.count)
}

func (t *Tree) Clear() {
	atomic.StorePointer(&t.root, unsafe.Pointer(makeN4()))
	atomic.StoreUint64(&t.count, 0)
	//don't runtime.GC
}

func (t *Tree) Has() {

}

func (t *Tree) IsEmpty() {

}
