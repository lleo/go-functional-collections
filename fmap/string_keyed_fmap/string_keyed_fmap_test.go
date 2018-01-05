package string_keyed_fmap_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/lleo/go-functional-collections/fmap/string_keyed_fmap"
)

func TestButildSimpleMap(t *testing.T) {
	var m = string_keyed_fmap.New()
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
	var m = string_keyed_fmap.New()
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
	var m = string_keyed_fmap.New()
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
	var m = string_keyed_fmap.New()

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

func TestDel(t *testing.T) {
	var m = string_keyed_fmap.New()

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
	var m = string_keyed_fmap.New()

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

	var found bool
	var val interface{}
	m, val, found = m.Remove("b")

	if !found {
		t.Fatal("m.Remove(\"b\") not found")
	}

	if val != 2 {
		t.Fatal("val != 2")
	}

	if m.NumEntries() != 2 {
		t.Fatal("m.NumEntries() != 2")
	}

	_, found = m.Load("b")
	if found {
		t.Fatal("found \"b\" after m.Del(\"b\")")
	}
}

func TestIter(t *testing.T) {
	var m = string_keyed_fmap.New()

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

	var it = m.Iter()
	var k, v = it.Next()
	if k == "" {
		t.Fatal("k == \"\"")
	}
	if k != "c" {
		t.Fatal("k != \"c\"")
	}
	if v != 3 {
		t.Fatal("v != 3")
	}

	k, v = it.Next()
	if k == "" {
		t.Fatal("k == \"\"")
	}
	if k != "b" {
		t.Fatal("k != \"b\"")
	}
	if v != 2 {
		t.Fatal("v != 2")
	}

	k, v = it.Next()
	if k == "" {
		t.Fatal("k == \"\"")
	}
	if k != "a" {
		t.Fatal("k != \"a\"")
	}
	if v != 1 {
		t.Fatal("v != 1")
	}

	k, _ = it.Next()
	if k != "" {
		t.Fatal("k != \"\"")
	}
}

func TestRange(t *testing.T) {
	var m = string_keyed_fmap.New().
		Put("a", 1).
		Put("b", 2).
		Put("c", 3)

	var ents = make([]string, m.NumEntries())

	var i int = 0
	m.Range(func(k string, v interface{}) bool {
		ents[i] = fmt.Sprintf("%q:%#v", k, v)
		i++
		return true
	})

	var str = strings.Join(ents, ",")
	var expected_str = "\"c\":3,\"b\":2,\"a\":1"
	if str != expected_str {
		t.Fatalf("str,%s != expected_str,%s", str, expected_str)
	}
}

func TestString(t *testing.T) {
	var m = string_keyed_fmap.New().
		Put("a", 1).
		Put("b", 2).
		Put("c", 3)

	var str = m.String()
	var expected_str = "StringKeyedMap{\"c\":3,\"b\":2,\"a\":1}"
	if str != expected_str {
		t.Fatalf("str,%s != expected_str,%s", str, expected_str)
	}
}
