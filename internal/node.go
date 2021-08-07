package internal

import "unsafe"

const maxPrefixLen = 8

type n struct {
	typ        uint8
	numChild   uint8
	prefixLen  uint32
	version    uint64
	prefixLeaf uint64
	prefix     [maxPrefixLen]byte
}

type n4 struct {
	n     n
	keys  [4]byte
	child [4]unsafe.Pointer
}

type n16 struct {
	n     n
	keys  [16]byte
	child [16]unsafe.Pointer
}

type n48 struct {
	n     n
	keys  [256]byte
	child [48]unsafe.Pointer
}

type n256 struct {
	n     n
	child [256]unsafe.Pointer
}
