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
		var k = kv.Key
		var v = kv.Val

		//log.Printf("TestCompBuildMap: calling m.Put(%s, %v)\n", k, v)
		//log.Println("==================================")
		//log.Printf("Map m =\n%s", m.TreeString())
		m = m.Put(k, v)

		var err = m.Valid()
		if err != nil {
			t.Fatalf("Invalid Tree. err=%s\n", err)
		}
	}

	//log.Printf("Map m =\n%s", m.TreeString())
	//log.Printf("Map m =\n%s", m.String())

	var i int
	var fn = func(k0 sorted_map.MapKey, v0 interface{}) bool {
		var k1 = inOrderKvs[i].Key
		var v1 = inOrderKvs[i].Val
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
	//var ns = int64(time.Now().Nanosecond())
	//log.Printf("ns = %d\n", ns)
	//rand.Seed(ns)
	var buildKvs = randomizeKeyVals(kvs)
	var destroyKvs = randomizeKeyVals(kvs)

	log.Printf("buildKvs = %v\n", buildKvs)

	var m = sorted_map.New()
	for i, kv := range buildKvs {
		var added bool
		m, added = m.Store(kv.Key, kv.Val)
		if !added {
			t.Fatal("Attempted to Store(%s, %v) but added=%v\n",
				kv.Key, kv.Val, added)
		}
		var err = m.Valid()
		if err != nil {
			t.Fatalf("INVALID TREE Store: i=%d; kv.Key=%s; err=%s\n",
				i, kv.Key, err)
		}
		//log.Printf("\n%s", m.TreeString())
		log.Println("*************************************")
	}

	log.Printf("AFTER ALL MAP BUILDING: Map m=\n%s", m.TreeString())
	log.Printf("m.NumEntries() = %d\n", m.NumEntries())
	log.Printf("m = %s\n", m.String())

	log.Printf("destroyKvs = %v\n", destroyKvs)

	var shouldHaveKvs = make([]KeyVal, len(destroyKvs))
	copy(shouldHaveKvs, destroyKvs)

	for i, kv := range destroyKvs {
		//var origM = m
		//var dupOrigM = m.Dup()

		log.Printf("********* Removing kv.Key=%s; i=%d *********\n", kv.Key, i)

		log.Printf("BEFORE REMOVE: Map m=\n%s", m.TreeString())

		var val interface{}
		var found bool
		m, val, found = m.Remove(kv.Key)

		log.Printf("AFTER REMOVE Map m=\n%s", m.TreeString())

		if !found {
			t.Fatalf("Remove: i=%d; kv.Key=%s not found!\n", i, kv.Key)
		}

		if val != kv.Val {
			t.Fatalf("Remove: i=%d; kv.Key=%s; val=%d != expected kv.Val=%d\n",
				i, kv.Key, val, kv.Val)
		}

		var err = m.Valid()
		if err != nil {
			log.Printf("Map m=\n%s", m.TreeString())
			t.Fatalf("INVALID TREE Remove: i=%d; kv.Key=%s; err=%s\n",
				i, kv.Key, err)
		}

		shouldHaveKvs = shouldHaveKvs[1:] //take the first elt off ala range
		for _, kv0 := range shouldHaveKvs {
			var val, found = m.Load(kv0.Key)
			if !found {
				t.Fatalf("Remove: for key=%s: "+
					"failed to find shouldHave key=%s", kv.Key, kv0.Key)
			}
			if val != kv0.Val {
				t.Fatalf("Remove: found val,%v != expected val,%v",
					val, kv0.Val)
			}
		}
		//if !origM.Equiv(dupOrigM) {
		//	t.Fatal("the original map was modified during Remove(%s).", kv.Key)
		//}
	}

	if m.NumEntries() != 0 {
		t.Fatal("m.NumEntries() != 0")
	}
}
