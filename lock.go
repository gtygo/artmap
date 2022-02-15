package artmap

import (
	"runtime"
	"sync/atomic"
)

const (
	OBSOLETE_FLAG        = uint64(1)
	LOCKED_FLAG          = uint64(2)
	OBSOLETE_LOCKED_FLAG = uint64(3)

	spin_count = int8(64)
)

func readLockOrRestart(n *n) (uint64, bool) {
	v := awaitNodeUnLocked(n)
	if isObsolete(v) {
		return 0, false
	}
	return v, true
}

func readUnLockOrRestart(n *n, version uint64) bool {
	if n == nil {
		return true
	}
	if version != atomic.LoadUint64(&n.version) {
		return false
	}
	return true
}

func readUnLockWithNodeOrRestart(n *n, version uint64, uln *n) bool {
	if n == nil {
		return true
	}
	if version != atomic.LoadUint64(&n.version) {
		writeUnLock(uln)
		return false
	}
	return true
}

func upgradeToWriteLockOrRestart(n *n, version uint64) bool {
	if n == nil {
		return true
	}
	return atomic.CompareAndSwapUint64(&n.version, version, (version)+LOCKED_FLAG)
}

func upgradeToWriteLockWithNodeOrRestart(n *n, version uint64, uln *n) bool {
	if n == nil {
		return true
	}
	if !atomic.CompareAndSwapUint64(&n.version, version, (version)+LOCKED_FLAG) {
		writeUnLock(uln)
		return false
	}
	return true
}

func writeLockOrRestart(n *n) bool {
	var (
		v  uint64
		ok bool
	)
	for {
		v, ok = readLockOrRestart(n)
		if !ok {
			return false
		}
		if upgradeToWriteLockOrRestart(n, v) {
			break
		}
	}
	return false
}

func writeUnLock(n *n) {
	if n == nil {
		return
	}
	atomic.AddUint64(&n.version, LOCKED_FLAG)
}

func writeUnLockObsolete(n *n) {
	if n == nil {
		return
	}
	atomic.AddUint64(&n.version, OBSOLETE_LOCKED_FLAG)
}

func awaitNodeUnLocked(n *n) uint64 {
	v := atomic.LoadUint64(&n.version)
	c := spin_count
	for (v & LOCKED_FLAG) == LOCKED_FLAG { //spin lock
		if c <= 0 {
			runtime.Gosched()
			c = spin_count
		}
		c--
		v = atomic.LoadUint64(&n.version)
	}
	return v
}

func isObsolete(version uint64) bool {
	return (version & OBSOLETE_FLAG) == OBSOLETE_FLAG
}
