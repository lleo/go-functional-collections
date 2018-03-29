package fmap_test

import (
	"testing"

	"github.com/lleo/go-functional-collections/fmap"
	"github.com/lleo/go-functional-collections/key"
)

var sizeBig = 1000000

func TestIntensiveButildMapBig(t *testing.T) {
	var m = fmap.New()

	var s = "a"
	for i := 0; i < sizeBig; i++ {
		m = m.Put(key.Str(s), i)
		s = Inc(s)
	}
}

// 32: 1st level collisions "a"&"ae", "b"&"af", "aa"&"e", "f"&"ab", "ac"&"g"
// 10,000: 2nd level collisions "gug","crr","akc","ert","dri","fkp","ipv"
// 10,000: 3rd level collisions "ktx","qk"

type strVal struct {
	Str string
	Val interface{}
}

func TestIntensiveDestroyMapBig(t *testing.T) {
	var m = fmap.New()
	var data = make([]KeyVal, sizeBig)

	var s = "a"
	for i := 0; i < sizeBig; i++ {
		var k = key.Str(s)
		data[i] = KeyVal{k, i}
		m = m.Put(k, i)
		s = Inc(s)
	}

	// destroy data
	var val interface{}
	var deleted bool
	for _, kv := range data {
		var k = kv.Key
		var v = kv.Val
		m, val, deleted = m.Remove(k)
		if !deleted {
			t.Fatalf("Failed to delete k=%q v=%d k.Hash()=%s\n", k, v, k.Hash())
		}
		if val != v {
			t.Fatalf("For key=%q, Value stored val=%d != expected v=%d\n",
				k, val, v)
		}
	}

	if m.NumEntries() != 0 {
		t.Fatal("Failed to empty Map")
	}
}

// findAndRemove is just here to demonstrate how slow array O(n) remove is
// versus HAMT O(log16(n)) remove is.
func findAndRemove(k key.Hash, kvs *[]KeyVal) bool {
	for i := 0; i < len(*kvs); i++ {
		if k.Equals((*kvs)[i].Key) {
			// BTW this is the fast non-order preserving element deletion
			(*kvs)[i] = (*kvs)[len(*kvs)-1]
			(*kvs) = (*kvs)[:len(*kvs)-1]
			//log.Printf("findAndRemove: found i=%d; k=%s\n", i, k)
			return true
		}
	}
	return false
}

func TestIntensiveIterBig(t *testing.T) {
	//var sizeBig = 100000 //20secs for linear remove; over 1/100 for persistent HAMT
	var kvs = buildKvs(sizeBig)

	var s = buildMap(kvs)
	//var start = time.Now()
	//var numRemoved int
	var it = s.Iter()
	for k, _ := it.Next(); k != nil; k, _ = it.Next() {
		var found bool
		s, _, found = s.Remove(k)
		if !found {
			t.Fatalf("Failed to find k=%s", k)
		}
		//numRemoved++
		//if numRemoved%10000 == 0 {
		//	var timediff = time.Since(start)
		//	//var rate = 10000 * 1000000 / float64(timediff) // millisec
		//	var rate = float64(timediff) / 10000 // nanosecond per removed entry
		//	var numLeft = s.NumEntries()
		//	//var numLeft = len(kvs)
		//	//log.Printf("found numRemoved=%d; numLeft=%d; rate=%.3g rm/ms",
		//	log.Printf("found numRemoved=%d; numLeft=%d; rate=%.0f ns/rm\n",
		//		numRemoved, numLeft, rate)
		//	start = time.Now()
		//}
	}

	////takes ~10^3 times longer
	//var kvs0 = make([]KeyVal, sizeBig)
	//copy(kvs0, kvs)
	//for _, kv := range kvs0 {
	//	var k = kv.Key
	//	var found = findAndRemove(k, &kvs)
	//	if !found {
	//		t.Fatalf("Failed to find k=%s", k)
	//	}
	//}
}
