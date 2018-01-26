package fmap_test

import (
	"log"
	"testing"
	"time"

	"github.com/lleo/go-functional-collections/fmap"
	"github.com/lleo/go-functional-collections/hash"
)

var sizeBig = 1000000

func TestIntensiveButildMapBig(t *testing.T) {
	var m = fmap.New()

	var s = "a"
	for i := 0; i < sizeBig; i++ {
		m = m.Put(hash.StringKey(s), i)
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
	var data = make([]keyVal, sizeBig)

	var s = "a"
	for i := 0; i < sizeBig; i++ {
		var k = hash.StringKey(s)
		data[i] = keyVal{k, i}
		m = m.Put(k, i)
		s = Inc(s)
	}

	//destroy data
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

//findAndRemove is just here to demonstrate how slow array O(n) remove is versus
//HAMT O(log16(n)) remove is.
func findAndRemove(k hash.Key, kvs *[]keyVal) bool {
	for i := 0; i < len(*kvs); i++ {
		if k.Equals((*kvs)[i].Key) {
			//BTW this is the fast non-order preserving element deletion
			(*kvs)[i] = (*kvs)[len(*kvs)-1]
			(*kvs) = (*kvs)[:len(*kvs)-1]
			//log.Printf("findAndRemove: found i=%d; k=%s\n", i, k)
			return true
		}
	}
	return false
}

func TestIntensiveIterBig(t *testing.T) {
	var kvs = buildKvs(sizeBig)
	var s = buildMap(kvs)

	var start = time.Now()
	var numRemoved int
	var it = s.Iter()
	for k, v := it.Next(); k != nil; k, v = it.Next() {
		var found bool
		var val interface{}
		s, val, found = s.Remove(k)
		//found = findAndRemove(k, &kvs) //between 900 & 2700 times slower
		if !found {
			t.Fatalf("Failed to find k=%s", k)
		}
		if v != val {
			t.Fatalf("Found val,%v != expected v,%v for key=%s;", val, v, k)
		}
		//log.Printf("removed k=%s", k)
		numRemoved++
		if numRemoved%10000 == 0 {
			var timediff = time.Since(start)
			var rate = 10000 * 1000000 / float64(timediff) //millisec
			var numLeft = s.NumEntries()
			//var numLeft = len(kvs)
			log.Printf("found numRemoved=%d; numLeft=%d; rate=%.3g 1/ms",
				numRemoved, numLeft, rate)
			start = time.Now()
		}
	}
}
