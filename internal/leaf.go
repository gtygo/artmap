package internal

import (
	"bytes"
	"unsafe"
)

type leaf struct {
	typ   uint8
	key   []byte
	value interface{}
}

func makeLeaf(k []byte, v interface{}) *leaf {
	l := &leaf{
		key:   k,
		value: v,
	}
	return l
}

func (l *leaf) compareKey(k []byte) bool {
	return len(k) == len(l.key) && bytes.Compare(k, l.key) == 0
}

func (l *leaf) insertExpandLeaf(key []byte, value interface{}, depth int, loc *unsafe.Pointer) {

}
