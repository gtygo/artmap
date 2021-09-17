package internal

import (
	"sync/atomic"
	"unsafe"
)

type Art struct {
	count int
	root  unsafe.Pointer
}

func NewTree() *Art {
	return &Art{
		count: 0,
		root:  unsafe.Pointer(makeN4()),
	}
}

func (t *Art) Get(key []byte)(interface{},bool){
	for {
		n:=atomic.LoadPointer(&t.root)

		v:=




	}
}







