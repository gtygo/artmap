package artmap

import (
	"sync/atomic"
	"unsafe"
)

const (
	node256ChildSize = uint8(255)
)

type n256 struct {
	n
	child [node256ChildSize]unsafe.Pointer
}

func makeN256() *n256 {
	n := new(n256)
	n.typ = typeN256
	return n
}
//w opt
func (node *n256) insertChild(c byte, child unsafe.Pointer) {
	node.numChild++
	atomic.StorePointer(&node.child[c], child)
}
