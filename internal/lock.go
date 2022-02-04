package internal

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

func readLockOrRestart(n *n, version *uint64) bool {
	v := awaitNodeUnLocked(n)
	if isObsolete(v) {
		*version = 0
		return false
	}
	*version = v
	return true
}

func readUnLockOrRestart(n *n, version uint64) bool {
	if n == nil {
		return true
	}
	return version == atomic.LoadUint64(&n.version)
}

func readUnLockWithNodeOrRestart(n *n, version uint64, uln *n) bool {
	if n == nil {
		return true
	}
	writeUnLock(uln)
	return version == atomic.LoadUint64(&n.version)
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
	v := uint64(0)
	for {
		if !readLockOrRestart(n, &v) {
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

func setLockedBit(version uint64) uint64 {
	return version + LOCKED_FLAG
}

func isObsolete(version uint64) bool {
	return (version & OBSOLETE_FLAG) == OBSOLETE_FLAG
}
