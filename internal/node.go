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

func (cn *n) search(key []byte, depth int, pn *n, pv uint64) (interface{}, bool) {
	var version uint64

CUR:
	if !rLock(cn, &version) {
		return nil, false
	}

	if !rUnLock(cn, pv) {
		return nil, false
	}

	if cn.checkPrefix(key, depth) != min(int(cn.prefixLen), maxPrefixLen) {
		if !rUnLock(cn, version) {
			return nil, false
		}
		return nil, true
	}

	depth += int(cn.prefixLen)
	if depth == len(key) {
		l := (* leaf)(atomic.LoadPointer(&cn.prefixLeaf))
		var v interface{}
		if l != nil && l.compareKey(key) {
			v = l.value
		}
		if !rUnLock(cn, version) {
			return nil, false
		}
		return v, true
	}

	if depth > len(key) {
		return nil, rUnLock(cn, version)
	}

	locator := cn.findChild(key[depth])

	var nextNode *n

	if locator != nil {
		nextNode = (*n)(atomic.LoadPointer(locator))
	}

	if nextNode == nil {
		if !rUnLock(cn, version) {
			return nil, false
		}
		return nil, true
	}

	if !rUnLock(cn, version) {
		return nil, false
	}

	if cn.typ == typeLeaf {
		l := (*leaf)(unsafe.Pointer(nextNode))

		if !rUnLock(cn, version) {
			return nil, false
		}

		if l.compareKey(key) {
			if !rUnLock(cn, version) {
				return nil, false
			}
			//Get success
			return l.value, true
		}
		return nil, true
	}

	depth++
	pn = cn
	pv = version
	cn = nextNode
	goto CUR
}

func (cn *n) insert(key []byte, value interface{}, depth int, pn *n, pv uint64, locator *unsafe.Pointer) bool {
	var version uint64
CUR:
	if !rLock(cn, &version) {
		return false
	}
	prefixLen, comKey, ok := cn.prefixMismatch(key, depth, pn, version, pv)
	if !ok {
		return false
	}

	if cn.prefixLen != uint32(prefixLen) {
		if !upgrade(pn, &pv) {
			return false
		}
		if !upgradeWithUnlock(cn, &version, pn) {
			return false
		}
		cn.commonInsert(key, comKey, value, depth, prefixLen, locator)
		unlock(cn)
		unlock(pn)
		return true
	}

	depth += int(cn.prefixLen)

	if depth == len(key) {
		upgrade(cn, &version)
		rUnLock()



	}

}

type n4 struct {
	n
	keys  [4]byte
	child [4]unsafe.Pointer
}

func (n *n4) addChild(key byte, child unsafe.Pointer) {
	i := 0
	for ; i < int(n.numChild); i++ {
		if key < n.keys[i] {
			break
		}
	}
	//back to next byte
	copy(n.keys[i+1:], n.keys[i:])
	copy(n.child[i+1:], n.child[i:])
	n.keys[i] = key
	atomic.StorePointer(&n.child[i], child)
	n.numChild++
}

func makeN4() *n4 {
	n := new(n4)
	n.typ = typeN4
	return n
}

type n16 struct {
	n
	keys  [16]byte
	child [16]unsafe.Pointer
}

func makeN16() *n16 {
	n := new(n16)
	n.typ = typeN16
	return n
}

type n48 struct {
	n
	keys  [256]byte
	child [48]unsafe.Pointer
}

func makeN48() *n48 {
	n := new(n48)
	n.typ = typeN48
	return n
}

type n256 struct {
	n
	child [256]unsafe.Pointer
}

func makeN256() *n256 {
	n := new(n256)
	n.typ = typeN256
	return n
}

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
