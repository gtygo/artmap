package internal

import (
	"sync/atomic"
	"unsafe"
)

const (
	node4KeysSize  = 4
	node4ChildSize = 4
)

type n4 struct {
	n
	keys  [node4KeysSize]byte
	child [node4ChildSize]unsafe.Pointer
}

func makeN4() *n4 {
	n := new(n4)
	n.typ = typeN4
	return n
}

//w opt
func (node *n4) insertAndGrow(ref *unsafe.Pointer, c byte, child unsafe.Pointer) {
	newNode := makeN16()
	copy(newNode.child[:], node.child[:])
	copy(newNode.keys[:], node.keys[:])
	copyHeader((*n)(unsafe.Pointer(newNode)), (*n)(unsafe.Pointer(node)))
	newNode.insertChild(c, child)
	atomic.StorePointer(ref, unsafe.Pointer(newNode))
}

//w opt
func (node *n4) insertChild(c byte, child unsafe.Pointer) {
	var idx uint8
	for idx := uint8(0); idx < node.numChild; idx++ {
		if c < node.keys[idx] {
			break
		}
	}
	copy(node.keys[idx+1:], node.keys[idx:])
	copy(node.child[idx+1:], node.child[idx:])
	node.keys[idx] = c
	atomic.StorePointer(&node.child[idx], child)
	node.numChild++
}

//w opt
func copyHeader(dst *n, src *n) {
	dst.numChild = src.numChild
	dst.prefixLen = src.prefixLen
	copy(dst.prefix[:], src.prefix[:])
	dst.prefixLeaf = src.prefixLeaf

}
