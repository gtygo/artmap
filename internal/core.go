package internal

import (
	"sync/atomic"
	"unsafe"
)

func (cn *n) checkPrefix(key []byte, depth int) int {
	if cn.prefixLen == 0 {
		return 0
	}

	maxCmp := min(min(int(cn.prefixLen), maxPrefixLen), len(key)-depth)

	var idx int
	for idx = 0; idx < maxCmp; idx++ {
		if cn.prefix[idx] != key[depth+idx] {
			return idx
		}
	}
	return idx
}



func (cn *n) prefixMismatch(key []byte, depth int, pn *n, v uint64, pv uint64) (int, []byte, bool) {
	if cn.prefixLen < maxPrefixLen {
		return cn.checkPrefix(key, depth), nil, true
	}
	var (
		completeKey []byte
		ok          bool
	)
	for {
		completeKey, ok = cn.findCompleteKey(v)
		if ok {
			break
		}

		if !readUnLockOrRestart(cn, v) || !readUnLockOrRestart(pn, pv) {
			return 0, nil, false
		}
	}
	i := depth
	minPrefixLen := min(len(key), depth+int(cn.prefixLen))

	for ; i < minPrefixLen; i++ {
		if key[i] != completeKey[i] {
			break
		}

	}
	return i - depth, completeKey, true
}

func (cn *n) commonInsert(key, comKey []byte, value interface{}, depth int, prefixLen int, nodeLoc *unsafe.Pointer) {
	nn4 := makeN4()
	tmpDepth := depth + prefixLen
	if len(key) == tmpDepth {
		nn4.prefixLeaf = unsafe.Pointer(makeLeaf(key, value))
	} else {
		nn4.addChild(key[depth], unsafe.Pointer(makeLeaf(key, value)))
	}
	nn4.prefixLen = uint32(prefixLen)
	copy(cn.prefix[:min(maxPrefixLen, prefixLen)], cn.prefix[:])

	if cn.prefixLen <= maxPrefixLen {
		nn4.addChild(cn.prefix[prefixLen], unsafe.Pointer(cn))
		cn.prefixLen -= uint32(prefixLen) + 1
		copy(cn.prefix[:min(maxPrefixLen, int(cn.prefixLen))], cn.prefix[prefixLen+1:])

	} else {
		nn4.addChild(comKey[depth+prefixLen], unsafe.Pointer(cn))
		cn.prefixLen -= uint32(prefixLen) + 1
		copy(cn.prefix[:min(maxPrefixLen, int(cn.prefixLen))], comKey[depth+prefixLen+1:])

	}

	atomic.StorePointer(nodeLoc, unsafe.Pointer(nn4))
}

func (cn *n) findCompleteKey(version uint64) ([]byte, bool) {
	//1. check prefix leaf

	prefixL := atomic.LoadPointer(&cn.prefixLeaf)
	if prefixL != nil {
		l := (*leaf)(prefixL)

		if !readUnLockOrRestart(cn, version) {
			return nil, false
		}
		return l.key, true
	}
	//2. check node leaf
	child := cn.findFirstChild()
	if !readUnLockOrRestart(cn, version) {
		return nil, false
	}

	if child.typ == typeLeaf {
		l := (*leaf)(unsafe.Pointer(child))
		key := l.key
		if !readUnLockOrRestart(cn, version) {
			return nil, false
		}
		return key, true
	}
	childVersion := uint64(0)
	if !readLockOrRestart(child, &childVersion) {
		return nil, false
	}
	return child.findCompleteKey(childVersion)
}

func (cn *n) findFirstChild() *n {
	switch cn.typ {
	case typeN4:
		{
			xn := (*n4)(unsafe.Pointer(cn))
			return (*n)(atomic.LoadPointer(&xn.child[0]))
		}
	case typeN16:
		{
			xn := (*n16)(unsafe.Pointer(cn))
			return (*n)(atomic.LoadPointer(&xn.child[0]))
		}
	case typeN48:
		{
			xn := (*n48)(unsafe.Pointer(cn))
			for i := 0; i < 256; i++ {
				idx := xn.keys[i]
				if idx != 0 {
					return (*n)(atomic.LoadPointer(&xn.child[idx-1]))
				}
			}
			return nil

		}
	case typeN256:
		{
			xn := (*n256)(unsafe.Pointer(cn))
			for i := 0; i < 256; i++ {
				if child := atomic.LoadPointer(&xn.child[i]); child != nil {
					return (*n)(child)
				}

			}
			return nil

		}

	}
	return nil

}

func min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}
