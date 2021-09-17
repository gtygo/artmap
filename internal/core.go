package internal

import "unsafe"

func (n *n) checkPrefix(key []byte, depth int) int {
	if n.prefixLen == 0 {
		return 0
	}

	maxCmp := min(min(int(n.prefixLen), maxPrefixLen), len(key)-depth)

	var idx int
	for idx = 0; idx < maxCmp; idx++ {
		if n.prefix[idx] != key[depth+idx] {
			return idx
		}
	}
	return idx
}

func (n *n) findChild(c byte) *unsafe.Pointer {

	switch n.typ {
	case typeN4:
		{
			node:=(* n4)(unsafe.Pointer(n))
			//todo: Accelerated traversal using SIMD
			for i:=0;i<int(node.numChild);i++{
				if node.keys[i]==c{
					return &node.child[i]
				}
			}
			break
		}
	case typeN16:
		{
		    node:=(*n16)(unsafe.Pointer(n))

		    for i:=0;i<int(node.numChild);i++{
		    	if node.keys[i]==c{
		    		return &node.child[i]
				}
			}
			break
		}
	case typeN48:
		{
			node:=(*n48)(unsafe.Pointer(n))
			i:=node.keys[c]
			if i>0 {
				return &node.child[i]
			}
			break
		}
	case typeN256:
		{

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
