package set_test

import (
	"testing"

	"github.com/lleo/go-functional-collections/set"
)

var size1MM = 1000000

func Test_Intensive_ButildSet1MM(t *testing.T) {
	var m = set.New()

	var s = "a"
	for i := 0; i < size1MM; i++ {
		m = m.Set(StringKey(s))
		s = Inc(s)
	}

	//log.Println("1MM large Set\n", m.LongString(""))
}

// 32: 1st level collisions "a"&"ae", "b"&"af", "aa"&"e", "f"&"ab", "ac"&"g"
// 10,000: 2nd level collisions "gug","crr","akc","ert","dri","fkp","ipv"
// 10,000: 3rd level collisions "ktx","qk"

func Test_Intensive_DestroySet1MM(t *testing.T) {
	var m = set.New()
	var keys = make([]string, size1MM)

	var s = "a"
	for i := 0; i < size1MM; i++ {
		keys[i] = s
		m = m.Set(StringKey(s))
		s = Inc(s)
	}

	//destroy keys
	var deleted bool
	for _, k := range keys {
		m, deleted = m.Remove(StringKey(k))
		if !deleted {
			t.Fatalf("Failed to delete key=%q\n", k)
		}
	}

	if m.NumEntries() != 0 {
		t.Fatal("Failed to empty Set")
	}
}
