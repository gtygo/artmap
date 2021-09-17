package internal

import (
	"bytes"
	"sync/atomic"
	"unsafe"
)

const (
	minN16  = 4
	minN48  = 13
	minN256 = 38

	maxPrefixLen = 8
)

const (
	typeN4 = iota
	typeN16
	typeN48
	typeN256
	typeLeaf
)

type n struct {
	typ        uint8
	numChild   uint8
	prefixLen  uint32
	version    uint64
	prefixLeaf unsafe.Pointer
	prefix     [maxPrefixLen]byte
}

func (n *n)get(key []byte,depth int,p *n,pv uint64)(interface{},bool){
	var version uint64

RECUR:

	if !rLock(n,&version){
		return nil,false
	}

	if !rUnLock(n,pv){
		return nil,false
	}

	if n.checkPrefix(key,depth)!=min(int(n.prefixLen),maxPrefixLen){
		if !rUnLock(n,version){
			return nil,false
		}
		return nil,true
	}

	depth+=int(n.prefixLen)
	if depth==len(key){
		l:=(* leaf)(atomic.LoadPointer(&n.prefixLeaf))
		var v interface{}
		if l!=nil && l.compareKey(key) {
			v =l.value
		}
		if !rUnLock(n,version){
			return nil,false
		}
		return v,true
	}

	if depth>len(key) {
		return nil,rUnLock(n,version)
	}

	locator:=



}





type n4 struct {
	n
	keys  [4]byte
	child [4]unsafe.Pointer
}

func makeN4()*n4{
	n:=new(n4)
	n.typ=typeN4
	return n
}

type n16 struct {
	n
	keys  [16]byte
	child [16]unsafe.Pointer
}

func makeN16()*n16{
	n:=new(n16)
	n.typ=typeN16
	return n
}

type n48 struct {
	n
	keys  [256]byte
	child [48]unsafe.Pointer
}

func makeN48()*n48{
	n:=new(n48)
	n.typ=typeN48
	return n
}

type n256 struct {
	n
	child [256]unsafe.Pointer
}

func makeN256()*n256{
	n:=new(n256)
	n.typ=typeN256
	return n
}


type leaf struct {
	typ   uint8
	key   []byte
	value interface{}
}

func makeLeaf(k []byte, v interface{}) *leaf {
	l := &leaf{
		typ:typeLeaf,
		key: k,
		value: v,
	}
	return l
}

func (l *leaf) compareKey(k []byte) bool {
	return len(k) == len(l.key) && bytes.Compare(k, l.key) == 0
}
