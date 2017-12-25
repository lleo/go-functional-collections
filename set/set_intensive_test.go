package set_test

import (
	"log"
	"testing"
	"time"

	"github.com/lleo/go-functional-collections/set"
)

var sizeBig = 1000000

func TestIntensiveButildSetBig(t *testing.T) {
	var s = set.New()

	var keyStr = "a"
	for i := 0; i < sizeBig; i++ {
		s = s.Set(StringKey(keyStr))
		keyStr = Inc(keyStr)
	}

	//log.Println("Big(%d) large Set\n", sizeBig, s.TreeString(""))
}

// 32: 1st level collisions "a"&"ae", "b"&"af", "aa"&"e", "f"&"ab", "ac"&"g"
// 10,000: 2nd level collisions "gug","crr","akc","ert","dri","fkp","ipv"
// 10,000: 3rd level collisions "ktx","qk"

func TestIntensiveDestroySetBig(t *testing.T) {
	var s = set.New()
	var keys = make([]string, sizeBig)

	var keyStr = "a"
	for i := 0; i < sizeBig; i++ {
		keys[i] = keyStr
		s = s.Set(StringKey(keyStr))
		keyStr = Inc(keyStr)
	}

	//destroy keys
	var deleted bool
	for _, k := range keys {
		s, deleted = s.Remove(StringKey(k))
		if !deleted {
			t.Fatalf("Failed to delete key=%q\n", k)
		}
	}

	if s.NumEntries() != 0 {
		t.Fatal("Failed to empty Set")
	}
}

//findAndRemove is just here to demonstrate how slow array O(n) remove is versus
//HAMT O(log16(n)) remove is.
func findAndRemove(k set.SetKey, keys *[]set.SetKey) bool {
	for i := 0; i < len(*keys); i++ {
		if k.Equals((*keys)[i]) {
			(*keys)[i] = (*keys)[len(*keys)-1]
			(*keys) = (*keys)[:len(*keys)-1]
			//log.Printf("findAndRemove: found i=%d; k=%s\n", i, k)
			return true
		}
	}
	return false
}

func TestIntensiveIterBig(t *testing.T) {
	var keys = buildKeys(sizeBig)
	var s = buildSet(keys)

	var start = time.Now()
	var numRemoved int
	var it = s.Iter()
	for k := it.Next(); k != nil; k = it.Next() {
		var found bool
		s, found = s.Remove(k)
		//found = findAndRemove(k, &keys) //between 900 & 2700 times slower
		if !found {
			t.Fatalf("Failed to find k=%s", k)
		}
		//log.Printf("removed k=%s", k)
		numRemoved++
		if numRemoved%10000 == 0 {
			var timediff = time.Since(start)
			var rate = 10000 * 1000000 / float64(timediff) //millisec
			var numLeft = s.NumEntries()
			//var numLeft = len(keys)
			log.Printf("found numRemoved=%d; numLeft=%d; rate=%.3g 1/ms",
				numRemoved, numLeft, rate)
			start = time.Now()
		}
	}
}
