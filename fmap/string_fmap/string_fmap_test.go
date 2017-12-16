package string_keyed_map_test

import (
	"testing"

	"github.com/lleo/go-functional-collections/fmap/string_keyed_map"
)

func TestButildSimpleMap(t *testing.T) {
	var m = string_keyed_map.New()
	m = m.
		Put("a", 1).
		Put("b", 2).
		Put("c", 3)

	if m.Get("a") != 1 {
		t.Fatal("m.Get(\"a\") != 1")
	}

	if m.Get("b") != 2 {
		t.Fatal("m.Get(\"b\") != 2")
	}

	if m.Get("c") != 3 {
		t.Fatal("m.Get(\"c\") != 3")
	}
}

func TestLoad(t *testing.T) {
	var m = string_keyed_map.New()
	m = m.Put("a", nil)

	var val interface{}
	var found bool

	val, found = m.Load("a")
	if !found {
		t.Fatal("failed to m.Load(\"a\")")
	} else {
		if val != nil {
			t.Fatal("val != nil")
		}
	}

	val, found = m.Load("b")
	if found {
		t.Fatal("WTF! \"b\" found")
	} else {
		if val != nil {
			t.Fatal("WTF! val!=nil for !found \"b\" ")
		}
	}
}

func TestLoadOrStore(t *testing.T) {
	var m = string_keyed_map.New()
	m = m.Put("a", 1)

	var val interface{}
	var loaded bool

	m, val, loaded = m.LoadOrStore("b", 2)
	if loaded {
		t.Fatal("failed to store (!loaded) (\"b\",2)")
	} else {
		if val != nil {
			t.Fatal("previous value val!=nil for store of m.LoadOrStore(\"b\", 2) call")
		}
	}

	m, val, loaded = m.LoadOrStore("a", 3)
	if !loaded {
		t.Fatal("failed to load m.LoadOrStore(\"b\", 3)")
	} else {
		if val != 1 {
			t.Fatal("val != 1 for m.LoadOrStore(\"b\", 3) call")
		}
	}

	val = m.Get("a")
	if val != 1 {
		t.Fatalf("val != 1 prior call to m.LoadOrStore(\"b\", 3) changed val=%d", val)
	}
}

func TestStore(t *testing.T) {
	var m = string_keyed_map.New()

	var added bool
	m, added = m.Store("a", 1)

	if !added {
		t.Fatal("added for m.Store(\"a\", 1)")
	} else {
		if m.Get("a") != 1 {
			t.Fatal("m.Get(\"a\") != 1")
		}
	}

	m, added = m.Store("a", 2)
	if added {
		t.Fatal("added == true for second m.Store(\"a\", 2)")
	} else {
		if m.Get("a") != 2 {
			t.Fatal("m.Get(\"a\") != 2")
		}
	}
}

func TestDelete(t *testing.T) {
	var m = string_keyed_map.New()

	if m.NumEntries() != 0 {
		t.Fatal("m.NumEntries() != 0")
	}

	m = m.
		Put("a", 1).
		Put("b", 2).
		Put("c", 3)

	if m.NumEntries() != 3 {
		t.Fatal("m.NumEntries() != 3")
	}

	m = m.Del("b")

	if m.NumEntries() != 2 {
		t.Fatal("m.NumEntries() != 2")
	}

	var _, found = m.Load("b")
	if found {
		t.Fatal("found \"b\" after m.Del(\"b\")")
	}

	m = m.Del("a").Del("c")

	if m.NumEntries() != 0 {
		t.Fatal("m.NumEntries() != 0")
	}
}

func TestRemove(t *testing.T) {

}
