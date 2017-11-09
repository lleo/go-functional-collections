package fmap_test

import (
	"testing"

	"github.com/lleo/go-functional-collections/fmap"
)

var size1MM = 1000000

func Test_Intensive_ButildMap1MM(t *testing.T) {
	var m = fmap.New()

	var s = "a"
	for i := 0; i < size1MM; i++ {
		m = m.Put(StringKey(s), i)
		s = Inc(s)
	}

	//log.Println("1MM large Map\n", m.LongString(""))
}

// 32: 1st level collisions "a"&"ae", "b"&"af", "aa"&"e", "f"&"ab", "ac"&"g"
// 10,000: 2nd level collisions "gug","crr","akc","ert","dri","fkp","ipv"
// 10,000: 3rd level collisions "ktx","qk"

func Test_Intensive_DestroyMap1MM(t *testing.T) {
	var m = fmap.New()
	var data = make(map[string]int, size1MM)

	var s = "a"
	for i := 0; i < size1MM; i++ {
		data[s] = i
		m = m.Put(StringKey(s), i)
		s = Inc(s)
	}

	//destroy data
	var val interface{}
	var deleted bool
	for k, v := range data {
		m, val, deleted = m.Remove(StringKey(k))
		if !deleted {
			t.Fatalf("Failed to delete key=%q\n", k)
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
