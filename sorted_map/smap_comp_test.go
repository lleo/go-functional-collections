package sorted_map_test

import (
	"log"
	"testing"

	"github.com/lleo/go-functional-collections/sorted_map"
)

func TestCompBuildMap(t *testing.T) {
	var numKeys = 64 //tested upto 10240
	var inOrderKvs = genIntKeyVals(numKeys)
	var kvs = randomizeKeyVals(inOrderKvs)

	//log.Println("kvs[:10] =", kvs[:10])
	var m = sorted_map.New()

	for _, kv := range kvs {
		var k = kv.key
		var v = kv.val

		//log.Printf("TestCompBuildMap: calling m.Put(%s, %v)\n", k, v)
		//log.Println("==================================")
		//log.Printf("Map m =\n%s", m.TreeString())
		m = m.Put(k, v)

		if !m.Valid() {
			t.Fatal("Invalid Tree.")
		}
	}

	//log.Printf("Map m =\n%s", m.TreeString())
	//log.Printf("Map m =\n%s", m.String())

	var i int
	var fn = func(k0 sorted_map.MapKey, v0 interface{}) bool {
		var k1 = inOrderKvs[i].key
		var v1 = inOrderKvs[i].val
		if k0.Less(k1) || k1.Less(k0) { //k0 != k1
			t.Fatalf("InOrder keys: i=%d; found k0=%s not the expected k1=%s\n",
				i, k0, k1)
		}
		if v0 != v1 {
			t.Fatalf("InOrder vals: i=%d; found v0=%d not the expected v1=%d\n",
				i, v0, v1)
		}
		i++
		return true
	}
	m.Range(fn)
}

func TestCompDestroyMap(t *testing.T) {
	var numKeys = 8
	var kvs = genIntKeyVals(numKeys)
	var buildKvs = randomizeKeyVals(kvs)
	var destroyKvs = randomizeKeyVals(kvs)

	log.Printf("buildKvs = %v\n", buildKvs)

	var m = sorted_map.New()
	for i, kv := range buildKvs {
		var added bool
		m, added = m.Store(kv.key, kv.val)
		if !added {
			t.Fatal("Attempted to Store(%s, %v) but added=%v\n",
				kv.key, kv.val, added)
		}
		if !m.Valid() {
			t.Fatalf("!!! INVALID TREE !!! Store: i=%d; kv.key=%s;\n",
				i, kv.key)
		}
		//log.Printf("\n%s", m.TreeString())
		log.Println("*************************************")
	}

	log.Printf("AFTER ALL MAP BUILDING: Map m=\n%s", m.TreeString())
	log.Printf("m.NumEntries() = %d\n", m.NumEntries())
	log.Printf("m = %s\n", m.String())

	log.Printf("destroyKvs = %v\n", destroyKvs)

	log.Printf("BEFORE REMOVE: Map m=\n%s", m.TreeString())

	for i, kv := range destroyKvs {
		//var origM = m
		//var dupOrigM = m.Dup()

		log.Printf("********* Removing kv.key=%s; i=%d *********\n", kv.key, i)

		var val interface{}
		var found bool
		m, val, found = m.Remove(kv.key)

		log.Printf("AFTER REMOVE Map m=\n%s", m.TreeString())

		if !found {
			t.Fatalf("Remove: i=%d; kv.key=%s not found!\n", i, kv.key)
		}

		if val != kv.val {
			t.Fatalf("Remove: i=%d; kv.key=%s; val=%d != expected kv.val=%d\n",
				i, kv.key, val, kv.val)
		}

		if !m.Valid() {
			log.Printf("Map m=\n%s", m.TreeString())
			t.Fatalf("!!! INVALID TREE !!! Remove: i=%d; kv.key=%s;\n",
				i, kv.key)
		}

		//if !origM.Equiv(dupOrigM) {
		//	t.Fatal("the original map was modified during Remove(%s).", kv.key)
		//}
	}

	if m.NumEntries() != 0 {
		t.Fatal("m.NumEntries() != 0")
	}
}
