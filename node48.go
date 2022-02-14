package artmap

import (
	"sync/atomic"
	"unsafe"
)

const (
	node48KeysSize  = 256
	node48ChildSize = 48
)

type n48 struct {
	n
	keys  [node48KeysSize]byte
	child [node48ChildSize]unsafe.Pointer
}

func makeN48() *n48 {
	n := new(n48)
	n.typ = typeN48
	return n
}

//w opt
func (node *n48) insertAndGrow(ref *unsafe.Pointer, c byte, child unsafe.Pointer) {
	newNode := makeN256()
	node.copyData(newNode)
	copyHeader((*n)(unsafe.Pointer(newNode)), (*n)(unsafe.Pointer(node)))
	newNode.insertChild(c, child)
	atomic.StorePointer(ref, unsafe.Pointer(newNode))
}

//w opt
func (node *n48) insertChild(c byte, child unsafe.Pointer) {
	pos := 0
	for node.child[pos] != nil {
		pos++
	}
	atomic.StorePointer(&node.child[pos], child)
	node.keys[c] = uint8(pos + 1)
	node.numChild++
}

func (node *n48) copyData(newNode *n256) {
	for i := 0; i < node48KeysSize; i++ {
		if node.keys[i] != 0 {
			newNode.child[i] = node.child[node.keys[i]-1]
		}
	}
}
