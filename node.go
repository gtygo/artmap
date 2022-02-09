package artmap

import (
	"sync/atomic"
	"unsafe"
)

const (
	typeN4 = iota
	typeN16
	typeN48
	typeN256
	typeLeaf

	maxPrefixLen = 8
)

type n struct {
	typ        uint8
	numChild   uint8
	prefixLen  uint32
	version    uint64
	prefixLeaf unsafe.Pointer
	prefix     [maxPrefixLen]byte
}

func (cn *n) LookupOpt(key []byte, depth int, pn *n, pv uint64) (interface{}, bool) {
	var (
		cv    uint64
		rFlag bool
	)

CUR:

	if cv, rFlag = readLockOrRestart(cn); !rFlag {
		return nil, false
	}
	if pn != nil {
		if !readUnLockOrRestart(pn, pv) {
			return nil, false
		}
	}
	if cn.checkPrefix(key, depth) != min(int(cn.prefixLen), maxPrefixLen) {
		if !readUnLockOrRestart(cn, cv) {
			return nil, false
		}
		return nil, true
	}

	depth += int(cn.prefixLen)
	if depth == len(key) {
		l := (*leaf)(atomic.LoadPointer(&cn.prefixLeaf))
		var v interface{}
		if l != nil && l.compareKey(key) {
			v = l.value
		}
		if !readUnLockOrRestart(cn, cv) {
			return nil, false
		}
		return v, true
	}

	if depth > len(key) {
		return nil, readUnLockOrRestart(cn, cv)
	}

	locator := cn.findChild(key[depth])

	var nextNode *n

	if locator != nil {
		nextNode = (*n)(*locator)
	}

	if nextNode == nil {
		if !readUnLockOrRestart(cn, cv) {
			return nil, false
		}
		return nil, true
	}

	if !readUnLockOrRestart(cn, cv) {
		return nil, false
	}

	if nextNode.typ == typeLeaf {
		l := (*leaf)(unsafe.Pointer(nextNode))

		if !readUnLockOrRestart(cn, cv) {
			return nil, false
		}

		if l.compareKey(key) {
			if !readUnLockOrRestart(cn, cv) {
				return nil, false
			}
			//Get success
			return l.value, true
		}
		return nil, true
	}

	depth++
	pn = cn
	pv = cv
	cn = nextNode
	goto CUR
}

func (cn *n) InsertOpt(key []byte, value interface{}, depth int, pn *n, pv uint64, locator *unsafe.Pointer) bool {
	var (
		cv    uint64
		rFlag bool
	)
CUR:
	if cv, rFlag = readLockOrRestart(cn); !rFlag {
		return false
	}
	prefixLen, comKey, ok := cn.prefixMismatch(key, depth, pn, cv, pv)
	if !ok {
		return false
	}

	if cn.prefixLen != uint32(prefixLen) {
		if !upgradeToWriteLockOrRestart(pn, pv) {
			return false
		}
		if !upgradeToWriteLockWithNodeOrRestart(cn, cv, pn) {
			return false
		}
		cn.insertSplitPrefix(key, comKey, value, depth, prefixLen, locator)
		writeUnLock(cn)
		writeUnLock(pn)
		return true
	}

	depth += int(cn.prefixLen)

	if depth == len(key) {
		//lock current
		if !upgradeToWriteLockOrRestart(cn, cv) {
			return false
		}
		//release father
		if !readUnLockWithNodeOrRestart(pn, pv, cn) {
			return false
		}
		cn.updatePrefixLeaf(key, value)
		writeUnLock(cn)
		return true

	}
	loc := cn.findChild(key[depth])
	var nextNode unsafe.Pointer = nil

	if loc != nil && (*loc) != nil {
		nextNode = *loc
	}
	if nextNode == nil {
		if cn.isFull() {
			if !upgradeToWriteLockOrRestart(pn, pv) {
				return false
			}
			if !upgradeToWriteLockWithNodeOrRestart(cn, cv, pn) {
				return false
			}
			//Capacity Expansion
			cn.insertAndGrow(locator, key[depth], unsafe.Pointer(makeLeaf(key, value)))
			writeUnLockObsolete(cn)
			writeUnLock(pn)
		} else {
			if !upgradeToWriteLockOrRestart(cn, cv) {
				return false
			}
			if !readUnLockWithNodeOrRestart(pn, pv, cn) {
				return false
			}
			cn.insert(key[depth], unsafe.Pointer(makeLeaf(key, value)))
			writeUnLock(cn)

		}
		return true
	}
	if pn != nil {
		if !readUnLockOrRestart(pn, pv) {
			return false
		}
	}

	if ((*n)(nextNode)).typ == typeLeaf {
		if !upgradeToWriteLockOrRestart(cn, cv) {
			return false
		}
		(*leaf)(nextNode).insertExpandLeaf(key, value, depth+1, loc)
		writeUnLock(cn)
		return true
	}
	depth += 1
	pn = cn
	pv = cv
	locator = loc
	cn = (*n)(nextNode)
	goto CUR
}

func (cn *n) updatePrefixLeaf(key []byte, value interface{}) {
	if cn.prefixLeaf == nil {
		atomic.StorePointer(&cn.prefixLeaf, unsafe.Pointer(makeLeaf(key, value)))
	} else {
		l := (*leaf)(cn.prefixLeaf)
		l.value = value
	}
	//check leaf type
}

func (cn *n) isFull() bool {
	switch cn.typ {
	case typeN4:
		{
			return cn.numChild == node4ChildSize
		}
	case typeN16:
		{
			return cn.numChild == node16ChildSize

		}
	case typeN48:
		{
			return cn.numChild == node48ChildSize
		}
	case typeN256:
		{
			return cn.numChild == node256ChildSize
		}
	}
	return true
}

func (cn *n) insertAndGrow(ref *unsafe.Pointer, c byte, child unsafe.Pointer) {
	switch cn.typ {
	case typeN4:
		{
			((*n4)(unsafe.Pointer(cn))).insertAndGrow(ref, c, child)
		}
	case typeN16:
		{
			((*n16)(unsafe.Pointer(cn))).insertAndGrow(ref, c, child)
		}
	case typeN48:
		{
			((*n48)(unsafe.Pointer(cn))).insertAndGrow(ref, c, child)
		}
	case typeN256:
		{
			((*n256)(unsafe.Pointer(cn))).insertAndGrow(c, child)
		}
	}
}

func (cn *n) insert(c byte, child unsafe.Pointer) {
	switch cn.typ {
	case typeN4:
		{
			((*n4)(unsafe.Pointer(cn))).insertChild(c, child)
		}
	case typeN16:
		{
			((*n16)(unsafe.Pointer(cn))).insertChild(c, child)
		}
	case typeN48:
		{
			((*n48)(unsafe.Pointer(cn))).insertChild(c, child)
		}
	case typeN256:
		{
			((*n256)(unsafe.Pointer(cn))).insertChild(c, child)
		}
	}
}

func (cn *n) findChild(c byte) *unsafe.Pointer {

	switch cn.typ {
	case typeN4:
		{
			node := (*n4)(unsafe.Pointer(cn))
			for i := 0; i < int(node.numChild); i++ {
				if node.keys[i] == c {
					return &node.child[i]
				}
			}
			break
		}
	case typeN16:
		{
			node := (*n16)(unsafe.Pointer(cn))
			i, j := 0, int(node.numChild)
			for i < j {
				h := int(uint(i+j) >> 1)
				if node.keys[h] > c {
					j = h
				} else if node.keys[h] < c {
					i = h + 1
				} else if node.keys[h] == c {
					return &node.child[h]
				}
			}
			break
		}
	case typeN48:
		{
			node := (*n48)(unsafe.Pointer(cn))
			i := node.keys[c]
			if i > 0 {
				return &node.child[i-1]
			}

		}
	case typeN256:
		{
			node := (*n256)(unsafe.Pointer(cn))
			if node.child[(int)(c)] != nil {
				return &node.child[(int)(c)]
			}
			break
		}
	}
	return nil

}
