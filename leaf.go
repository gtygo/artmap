package artmap

import (
	"bytes"
	"sync/atomic"
	"unsafe"
)

type leaf struct {
	typ   uint8
	key   []byte
	value interface{}
}

func makeLeaf(k []byte, v interface{}) *leaf {
	l := &leaf{
		typ:   typeLeaf,
		key:   k,
		value: v,
	}

	return l
}

func (l *leaf) compareKey(k []byte) bool {
	return len(k) == len(l.key) && bytes.Compare(k, l.key) == 0
}

func (l *leaf) insertExpandLeaf(key []byte, value interface{}, depth int, loc *unsafe.Pointer) {
	if l.compareKey(key) {
		l.value = value
	}

	prefixLen := min(len(key), len(l.key))
	var idx = 0
	for idx = depth; idx < prefixLen; idx++ {
		if l.key[idx] != key[idx] {
			break
		}
	}
	newNode := makeN4()
	newNode.prefixLen = uint32(idx - depth)
	copy(newNode.prefix[:], key[depth:])
	if idx == len(l.key) {
		newNode.prefixLeaf = unsafe.Pointer(l)

	} else {
		newNode.insertChild(l.key[idx], unsafe.Pointer(l))

	}
	if idx == len(key) {
		newNode.prefixLeaf = unsafe.Pointer(makeLeaf(key, value))
	} else {
		newNode.insertChild(key[idx], unsafe.Pointer(makeLeaf(key, value)))
	}
	atomic.StorePointer(loc, unsafe.Pointer(newNode))

}
