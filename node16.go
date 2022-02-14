package artmap

import (
	"sync/atomic"
	"unsafe"
)

const (
	node16KeysSize  = 16
	node16ChildSize = 16
)

type n16 struct {
	n
	keys  [node16KeysSize]byte
	child [node16ChildSize]unsafe.Pointer
}

func makeN16() *n16 {
	n := new(n16)
	n.typ = typeN16
	return n
}

//w opt
func (node *n16) insertAndGrow(ref *unsafe.Pointer, c byte, child unsafe.Pointer) {
	newNode := makeN48()
	node.copyData(newNode)
	copyHeader((*n)(unsafe.Pointer(newNode)), (*n)(unsafe.Pointer(node)))

	newNode.insertChild(c, child)

	atomic.StorePointer(ref, unsafe.Pointer(newNode))

}

//w opt
func (node *n16) insertChild(c byte, child unsafe.Pointer) {
	var idx uint8
	for idx = uint8(0); idx < node.numChild; idx++ {
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

func (node *n16) copyData(newNode *n48) {
	copy(newNode.child[:], node.child[:])
	for i := uint8(0); i < node.numChild; i++ {
		newNode.keys[node.keys[i]] = i + 1
	}
}
