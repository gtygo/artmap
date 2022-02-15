package artmap

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMapCreation(t *testing.T) {
	m := New()
	assert.NotEmpty(t, m)
	assert.Equal(t, m.count, uint64(0))
}

func TestInsert(t *testing.T) {
	m := New()

	m.Set([]byte("foo"), []byte("1"))
	m.Set([]byte("bar"), []byte("1"))
	assert.Equal(t, m.count, uint64(2))

}

func TestGet(t *testing.T) {
	m := New()

	// Get a missing element.
	val, ok := m.Get([]byte("Money"))
	assert.Equal(t, ok, false)
	assert.Equal(t, val, nil)

	elephant := "elephant"
	m.Set([]byte("elephant"), elephant)

	// Retrieve inserted element.
	tmp, ok := m.Get([]byte("elephant"))
	assert.Equal(t, ok, true)

	v, ok := tmp.(string) // Type assertion.
	assert.Equal(t, ok, true)

	assert.Equal(t, v, "elephant")
}

func TestClear(t *testing.T) {

	m := New()
	m.Set([]byte("elephant"), "elephant")
	m.Clear()
	assert.Equal(t, m.Count(), uint64(0))
}

func TestInsert4To16(t *testing.T) {
	m := New()

	m.Set([]byte("a"), "1")
	m.Set([]byte("b"), "2")
	m.Set([]byte("c"), "3")
	m.Set([]byte("d"), "4")
	m.Set([]byte("e"), "5")
	assert.Equal(t, m.count, uint64(5))

	if v, ok := m.Get([]byte("a")); ok {
		ans := v.(string)
		assert.Equal(t, ans, "1")
	}

	if v, ok := m.Get([]byte("e")); ok {
		ans := v.(string)
		assert.Equal(t, ans, "5")
	}

	if v, ok := m.Get([]byte("b")); ok {
		ans := v.(string)
		assert.Equal(t, ans, "2")
	}
}

func TestInsert16To48(t *testing.T) {
	m := New()

	for i := 16; i >= 0; i-- {
		k := fmt.Sprintf("%c", 'b'+i)
		m.Set([]byte(k), k)
		v, ok := m.Get([]byte(k))
		assert.Equal(t, ok, true)
		assert.Equal(t, v, k)
	}
	assert.Equal(t, m.count, uint64(17))
	m.Set([]byte("a"), "a")
	for i := 0; i < 18; i++ {
		k := fmt.Sprintf("%c", 'a'+i)
		v, ok := m.Get([]byte(k))
		assert.Equal(t, ok, true)
		assert.Equal(t, v, k)
	}
}

func TestInsert48To256(t *testing.T) {
	m := New()

	for i := 0; i < 52; i++ {
		if i == 51 {
			print("a")
		}
		k := fmt.Sprintf("%c", '!'+i)
		m.Set([]byte(k), k)
		for j := 0; j < i; j++ {
			kk := fmt.Sprintf("%c", '!'+j)
			v, ok := m.Get([]byte(kk))
			assert.Equal(t, ok, true)
			assert.Equal(t, v, kk)
		}

	}
	for i := 0; i < 52; i++ {
		k := fmt.Sprintf("%c", '!'+i)
		v, ok := m.Get([]byte(k))
		assert.Equal(t, ok, true)
		assert.Equal(t, v, k)
	}
}

func TestInsertMultiKey(t *testing.T) {
	m := New()
	for i := 0; i < 100000; i++ {
		k := fmt.Sprintf("%d", i)
		m.Set([]byte(k), k)
		v, ok := m.Get([]byte(k))
		assert.Equal(t, ok, true)
		assert.Equal(t, v, k)
	}
	for i := 0; i < 100000; i++ {
		k := fmt.Sprintf("%d", i)
		v, ok := m.Get([]byte(k))
		assert.Equal(t, ok, true)
		assert.Equal(t, v, k)
	}

	assert.Equal(t, m.count, uint64(100000))
}


