package fmap_test

import (
	"log"
	"testing"

	"github.com/lleo/go-functional-collections/fmap"
	"github.com/lleo/go-functional-collections/hash"
)

func TestBasicButildSimpleMap(t *testing.T) {
	var m = fmap.New()
	m = m.
		Put(hash.StringKey("a"), 1).
		Put(hash.StringKey("b"), 2).
		Put(hash.StringKey("c"), 3)

	if m.Get(hash.StringKey("a")) != 1 {
		t.Fatal("m.Get(\"a\") != 1")
	}

	if m.Get(hash.StringKey("b")) != 2 {
		t.Fatal("m.Get(\"b\") != 2")
	}

	if m.Get(hash.StringKey("c")) != 3 {
		t.Fatal("m.Get(\"c\") != 3")
	}
}

func TestBasicLoad(t *testing.T) {
	var m = fmap.New()
	m = m.Put(hash.StringKey("a"), nil)

	var val interface{}
	var found bool

	val, found = m.Load(hash.StringKey("a"))
	if !found {
		t.Fatal("failed to m.Load(\"a\")")
	} else {
		if val != nil {
			t.Fatal("val != nil")
		}
	}

	val, found = m.Load(hash.StringKey("b"))
	if found {
		t.Fatal("WTF! \"b\" found")
	} else {
		if val != nil {
			t.Fatal("WTF! val!=nil for !found \"b\" ")
		}
	}
}

func TestBasicLoadOrStore(t *testing.T) {
	var m = fmap.New()
	m = m.Put(hash.StringKey("a"), 1)

	var val interface{}
	var loaded bool

	m, val, loaded = m.LoadOrStore(hash.StringKey("b"), 2)
	if loaded {
		t.Fatal("failed to store (!loaded) (\"b\",2)")
	} else {
		if val != nil {
			t.Fatal("previous value val!=nil for store of m.LoadOrStore(\"b\", 2) call")
		}
	}

	m, val, loaded = m.LoadOrStore(hash.StringKey("a"), 3)
	if !loaded {
		t.Fatal("failed to load m.LoadOrStore(\"b\", 3)")
	} else {
		if val != 1 {
			t.Fatal("val != 1 for m.LoadOrStore(\"b\", 3) call")
		}
	}

	val = m.Get(hash.StringKey("a"))
	if val != 1 {
		t.Fatalf("val != 1 prior call to m.LoadOrStore(\"b\", 3) changed val=%d", val)
	}
}

func TestBasicStore(t *testing.T) {
	var m = fmap.New()

	var added bool
	m, added = m.Store(hash.StringKey("a"), 1)

	if !added {
		t.Fatal("added for m.Store(\"a\", 1)")
	} else {
		if m.Get(hash.StringKey("a")) != 1 {
			t.Fatal("m.Get(\"a\") != 1")
		}
	}

	m, added = m.Store(hash.StringKey("a"), 2)
	if added {
		t.Fatal("added == true for second m.Store(\"a\", 2)")
	} else {
		if m.Get(hash.StringKey("a")) != 2 {
			t.Fatal("m.Get(\"a\") != 2")
		}
	}
}

func TestBasicDelete(t *testing.T) {
	var m = fmap.New()

	if m.NumEntries() != 0 {
		t.Fatal("m.NumEntries() != 0")
	}

	m = m.
		Put(hash.StringKey("a"), 1).
		Put(hash.StringKey("b"), 2).
		Put(hash.StringKey("c"), 3)

	if m.NumEntries() != 3 {
		t.Fatal("m.NumEntries() != 3")
	}

	m = m.Del(hash.StringKey("b"))

	if m.NumEntries() != 2 {
		t.Fatal("m.NumEntries() != 2")
	}

	var _, found = m.Load(hash.StringKey("b"))
	if found {
		t.Fatal("found \"b\" after m.Del(\"b\")")
	}

	m = m.Del(hash.StringKey("a")).Del(hash.StringKey("c"))

	if m.NumEntries() != 0 {
		t.Fatal("m.NumEntries() != 0")
	}
}

func TestBasicRemove(t *testing.T) {
	var m = fmap.New()

	m = m.
		Put(hash.StringKey("a"), 1).
		Put(hash.StringKey("b"), 2).
		Put(hash.StringKey("c"), 3)

	if m.NumEntries() != 3 {
		t.Fatal("m.NumEntries() != 3")
	}

	var val interface{}
	var found bool
	var key hash.Key

	key = hash.StringKey("d")
	m, val, found = m.Remove(key)
	if found {
		t.Fatalf("found val=%#v for key=%#v that does not exist.", val, key)
	}

	key = hash.StringKey("b")
	m, val, found = m.Remove(key)
	if !found {
		t.Fatalf("failed to find & remove key=%#v", key)
	} else if val != 2 {
		t.Fatalf("found key=%#v, but val=%#v was not the expected val=2",
			key, val)
	}

	m, val, found = m.Remove(key)
	if found {
		t.Fatalf("found key=%#v entry for key just Removed; val=%#v", key, val)
	}
}

func TestBasicString(t *testing.T) {
	var m = fmap.New()
	m = m.
		Put(hash.StringKey("a"), 1).
		Put(hash.StringKey("b"), 2).
		Put(hash.StringKey("c"), 3)

	var str = m.String()
	log.Printf("m.String()=%s\n", str)

	var expectedStr = "Map{\"c\":3,\"b\":2,\"a\":1}"
	if str != expectedStr {
		t.Fatalf("str,%q != expectedStr,%q", str, expectedStr)
	}
}

func TestBasicNewFromList(t *testing.T) {
	var kvs = buildKvs(100)

	var m = fmap.NewFromList(kvs)

	for _, kv := range kvs {
		var k, v = kv.Key, kv.Val
		var val, found = m.Load(k)
		if !found {
			t.Fatalf("failed to find key=%s", k)
		}
		if val != v {
			t.Fatalf("val,%d != v,%d", val, v)
		}
	}
}

//FIXME: need more; test conflict
func TestBasicBulkInsert(t *testing.T) {
	var kvs = buildKvs(100)
	var m = fmap.NewFromList(kvs[:50])

	m = m.BulkInsert(kvs[50:], fmap.KeepOrigVal)

	if m.NumEntries() != 100 {
		t.Fatalf("m.NumEntries(),%d != 100", m.NumEntries())
	}

	for _, kv := range kvs {
		var k, v = kv.Key, kv.Val
		var val, found = m.Load(k)
		if !found {
			t.Fatalf("failed to find key=%s", k)
		}
		if val != v {
			t.Fatalf("val,%d != v,%d", val, v)
		}
	}
}

//FIXME: need more; test not found
func TestBasicBulkDelete(t *testing.T) {
	var kvs = buildKvs(100)
	var keys = make([]hash.Key, len(kvs))
	for i, kv := range kvs {
		keys[i] = kv.Key
	}
	var m = fmap.NewFromList(kvs)

	var notFound []hash.Key
	m, notFound = m.BulkDelete(keys[50:])

	if len(notFound) != 0 {
		t.Fatalf("len(notFound),%d != 0", len(notFound))
	}

	if m.NumEntries() != 50 {
		t.Fatalf("m.NumEntries(),%d != 50", m.NumEntries())
	}

	for _, kv := range kvs[:50] {
		var k, v = kv.Key, kv.Val
		var val, found = m.Load(k)
		if !found {
			t.Fatalf("k=%s not found", k)
		}
		if v != val {
			t.Fatalf("v,%d != val,%d", v, val)
		}
	}
}

//FIXME: need more; test conflict
func TestBasicMerge(t *testing.T) {
	var kvs = buildKvs(100)

	var m0 = fmap.NewFromList(kvs[:50])
	var m1 = fmap.NewFromList(kvs[50:])
	var m = m0.Merge(m1, fmap.KeepOrigVal)

	for _, kv := range kvs {
		var k, v = kv.Key, kv.Val
		var val, found = m.Load(k)
		if !found {
			t.Fatalf("k=%s not found", k)
		}
		if val != v {
			t.Fatalf("val,%d != v,%d", val, v)
		}
	}
}
