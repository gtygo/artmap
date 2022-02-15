package artmap

import (
	"fmt"
	"github.com/orcaman/concurrent-map"
	"sync"
	"testing"
)

var benchSize = 1

func BenchmarkInsertArtMap(b *testing.B) {
	m := New()
	dataSet1 := make([][]byte, benchSize)

	for i := 0; i < benchSize; i++ {
		dataSet1 = append(dataSet1, []byte(fmt.Sprintf("test:%d", i)))
	}
	for n := 0; n < b.N; n++ {
		for i := 0; i < benchSize; i++ {
			m.Set(dataSet1[i], 1)
			m.Get(dataSet1[i])
		}
	}
}

func BenchmarkInsertMutexMap(b *testing.B) {
	m := make(map[string]int, benchSize)
	lock := sync.Mutex{}
	dataSet1 := make([]string, benchSize)

	for i := 0; i < benchSize; i++ {
		dataSet1 = append(dataSet1, fmt.Sprintf("test:%d", i))
	}
	for n := 0; n < b.N; n++ {
		for i := 0; i < benchSize; i++ {
			lock.Lock()
			m[dataSet1[i]] = 1
			lock.Unlock()
		}
	}
}

func BenchmarkInsertSyncMap(b *testing.B) {
	m := sync.Map{}
	dataSet1 := make([]string, benchSize)

	for i := 0; i < benchSize; i++ {
		dataSet1 = append(dataSet1, fmt.Sprintf("test:%d", i))
	}
	for n := 0; n < b.N; n++ {
		for i := 0; i < benchSize; i++ {
			m.Store(dataSet1[i], 1)
			m.Load(dataSet1[i])
		}
	}
}

func BenchmarkInsertConcurrentMap(b *testing.B) {
	conmap := cmap.New()

	dataSet1 := make([]string, benchSize)

	for i := 0; i < benchSize; i++ {
		dataSet1 = append(dataSet1, fmt.Sprintf("test:%d", i))
	}
	for n := 0; n < b.N; n++ {
		for i := 0; i < benchSize; i++ {
			conmap.Set(dataSet1[i], 1)
			conmap.Get(dataSet1[i])
		}
	}
}
