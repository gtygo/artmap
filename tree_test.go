package artmap

import "testing"

func TestMapCreation(t *testing.T) {
	m := New()
	if m == nil {
		t.Error("map is null.")
	}

	if m.Count() != 0 {
		t.Error("new map should be empty.")
	}
}

func TestInsert(t *testing.T) {
	m := New()

	m.Set([]byte("foo"), []byte("1"))
	m.Set([]byte("bar"), []byte("1"))

	if m.Count() != 2 {
		t.Error("map should contain exactly two elements.")
	}
}

func TestInsert4To16(t *testing.T) {
	m := New()

	m.Set([]byte("a"), []byte("1"))
	m.Set([]byte("b"), []byte("2"))
	m.Set([]byte("c"), []byte("3"))
	m.Set([]byte("d"), []byte("4"))
	m.Set([]byte("e"), []byte("5"))

	if m.Count() != 5 {
		t.Error("map should contain exactly two elements.")
	}
	if v, ok := m.Get([]byte("a")); ok {
		ans := v.(string)
		if ans != "1" {
			t.Error("map has data.")
		}
	}

	if v, ok := m.Get([]byte("e")); ok {
		ans := v.(string)
		if ans != "5" {
			t.Error("map has data.")
		}
	}
	if v, ok := m.Get([]byte("d")); ok {
		ans := v.(string)
		if ans != "4" {
			t.Error("map has data.")
		}
	}
	if v, ok := m.Get([]byte("c")); ok {
		ans := v.(string)
		if ans != "3" {
			t.Error("map has data.")
		}
	}
	if v, ok := m.Get([]byte("b")); ok {
		ans := v.(string)
		if ans != "2" {
			t.Error("map has data.")
		}
	}

}

func TestGet(t *testing.T) {
	m := New()

	// Get a missing element.
	val, ok := m.Get([]byte("Money"))

	if ok == true {
		t.Error("ok should be false when item is missing from map.")
	}

	if val != nil {
		t.Error("Missing values should return as null.")
	}

	elephant := "elephant"
	m.Set([]byte("elephant"), elephant)

	// Retrieve inserted element.
	tmp, ok := m.Get([]byte("elephant"))
	if ok == false {
		t.Error("ok should be true for item stored within the map.")
	}

	elephant, ok = tmp.(string) // Type assertion.
	if !ok {
		t.Error("expecting an element, not null.")
	}

	if elephant != "elephant" {
		t.Error("item was modified.")
	}
}
