package artmap

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"runtime"
	"sync"
	"testing"
)

var wg sync.WaitGroup

const letterBytes = "~!@#$%^&*()_+-=1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// go test tree_concurrent_test.go node*.go lock.go core.go tree.go leaf.go -cpuprofile=cpu.out

func TestConcurrentInsert1(t *testing.T) {
	runtime.GOMAXPROCS(8)

	for j := 0; j < 1000; j++ {
		m := New()
		dataSet := GenerateKey(10000, 8)
		threadN := 16
		for i := 0; i < threadN; i++ {
			wg.Add(1)
			go InsertHelper(dataSet, m)
		}
		wg.Wait()
		GetHelper(dataSet, m, t)
		assert.Equal(t, m.Count(), uint64(10000))

		runtime.GC()
	}
}

func TestConcurrentInsert2(t *testing.T) {
	runtime.GOMAXPROCS(8)

	for j := 0; j < 1; j++ {
		m := New()
		dataSet := GenerateKey(10000, 8)
		threadN := 16
		for i := 0; i < threadN; i++ {
			wg.Add(1)
			go InsertAndGetHelper(dataSet, m, t)
		}
		wg.Wait()
		GetHelper(dataSet, m, t)
		assert.Equal(t, m.Count(), uint64(10000))
		runtime.GC()
	}
}

func InsertHelper(keys [][]byte, tree *Tree) {
	for _, key := range keys {
		tree.Set(key, key)
	}
	wg.Done()
}

func InsertAndGetHelper(keys [][]byte, tree *Tree, t assert.TestingT) {
	for _, key := range keys {
		tree.Set(key, key)
		v, ok := tree.Get(key)
		assert.Equal(t, v, key)
		assert.Equal(t, ok, true)
	}
	wg.Done()
}
func GetHelper(keys [][]byte, tree *Tree, t assert.TestingT) {
	for _, key := range keys {
		v, ok := tree.Get(key)
		assert.Equal(t, v, key)
		assert.Equal(t, ok, true)
	}
}

func RandStringBytes(n int) []byte {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return b
}

func GenerateKey(count int, keySize int) [][]byte {
	data := make([][]byte, 0, count)

	for i := 0; i < count; i++ {
		data = append(data, []byte(fmt.Sprintf("test_key_id:%d", i)))
		//data = append(data, RandStringBytes(keySize))
	}
	return data
}
