package internal

import (
	"runtime"
	"sync/atomic"
)

const d_flag = uint64(1)
const l_flag = uint64(1) << 1
const l_and_d_flag = uint64(1)<<2 - 1
const spin_count=int8(64)

func rLock(n *n, version *uint64) bool {
	v:=waitUnLock(n)
	if (v&d_flag)==d_flag{
		*version=0
		return false
	}
	*version=v
	return false
}

func rUnLock(n *n,version *uint64)bool{
	if n==nil{
		return true
	}
	return *version==atomic.LoadUint64(&n.version)
}

func lock(n *n)bool {
	v := uint64(0)
	for {
		if !rLock(n, &v) {
			return false
		}
		if upgrade(n,&v){
			break;
		}
	}
	return false
}

func unlock(n *n) {
	if n == nil {
		return
	}
	atomic.AddUint64(&n.version, l_flag)
}

func upgrade(n *n ,version *uint64)bool{
	if n ==nil{
		return true
	}
	return atomic.CompareAndSwapUint64(&n.version,*version,(*version)+l_flag)
}

func waitUnLock(n *n)uint64{
	v:=atomic.LoadUint64(&n.version)
	c:=spin_count
	for(v & l_flag) == l_flag {
		if c<=0{
			runtime.Gosched()
			c=spin_count
		}
		c--
		v=atomic.LoadUint64(&n.version)
	}
	return v
}


