package hash

import (
	"hash/fnv"
	"log"
	"testing"
)

func TestCalcHash(t *testing.T) {
	var key = "a"
	var v = CalcHash([]byte(key))
	//log.Println(v.String())
	if hashSize == 32 {
		var h = fnv.New32()
		h.Write([]byte(key))
		var val = h.Sum32()
		if uint32(v) != val {
			t.Fatalf("v,%d != val,%d", v, val)
		}
	} else if hashSize == 64 {
		var h = fnv.New64()
		h.Write([]byte(key))
		var val = h.Sum64()
		if uint64(v) != val {
			t.Fatalf("v,%d != val,%d", v, val)
		}
	} else {
		t.Fatalf("unknown hashSize,%d", hashSize)
	}
}

func TestValString(t *testing.T) {
	var key = "a"
	var v = CalcHash([]byte(key))

	var expected string
	if hashSize == 32 {
		expected = "/14/07/13/05/12/00/05/00"
	} else if hashSize == 64 {
		expected = "/14/11/07/11/01/00/06/08/12/04/13/11/03/06/15/10"
	} else {
		t.Fatalf("unknown hashSize,%d", hashSize)
	}

	var got = v.String()
	if got != expected {
		log.Printf("TestValString: failed: ")
		t.Fatalf("got %q expected %q", got, expected)
	}
}
