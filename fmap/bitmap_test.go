package fmap

import (
	"testing"

	"github.com/lleo/go-functional-collections/hash"
)

func TestBitmapIsSet(t *testing.T) {
	var bm = bitmap{2}
	//log.Printf("TestBitmapIsSet: bm=%s", bm.String())
	if !bm.isSet(1) {
		t.Fatal("bm.isSet(1) returned false")
	}
}

func TestBitmapSet(t *testing.T) {
	var bm bitmap
	bm.set(1).set(9)

	//log.Printf("TestBitmapSet: bm=%s", bm.String())

	if !bm.isSet(1) {
		t.Fatalf("bm.isSet(1),%t", bm.isSet(1))
	}
	if !bm.isSet(9) {
		t.Fatalf("bm.isSet(9),%t", bm.isSet(9))
	}
}

func TestBitmapUnset(t *testing.T) {
	var bm bitmap
	for i := uint(0); i < hash.IndexLimit; i++ {
		bm.set(i)
	}
	//log.Printf("TestBitmapUnset: bm=%s", bm.String())
	for i := uint(0); i < hash.IndexLimit; i++ {
		if i%2 == 0 {
			bm.unset(i)
		}
	}
	//log.Printf("TestBitmapUnset: bm=%s", bm.String())
	var isSet bool
	for i := uint(0); i < hash.IndexLimit; i++ {
		if bm.isSet(i) != isSet {
			t.Fatalf("bm.isSet(i,%d),%t != isSet,%t", i, bm.isSet(i), isSet)
		}
		isSet = !isSet
	}
}

func TestBitmapCount(t *testing.T) {
	var bm bitmap
	bm.set(0).
		set(1).
		set(2).
		set(9)
	//log.Printf("TestBitmapCount: bm=%s", bm.String())
	var idx uint
	if idx = 0; bm.count(idx) != 0 {
		t.Fatalf("idx=%d; bm.count(idx),%d != 0", idx, bm.count(idx))
	}
	if idx = 1; bm.count(idx) != 1 {
		t.Fatalf("idx=%d; bm.count(idx),%d != 1", idx, bm.count(idx))
	}
	if idx = 2; bm.count(idx) != 2 {
		t.Fatalf("idx=%d; bm.count(idx),%d != 2", idx, bm.count(idx))
	}
	if idx = 3; bm.count(idx) != 3 {
		t.Fatalf("idx=%d; bm.count(idx),%d != 3", idx, bm.count(idx))
	}
	if idx = 4; bm.count(idx) != 3 {
		t.Fatalf("idx=%d; bm.count(idx),%d != 3", idx, bm.count(idx))
	}
	if idx = 5; bm.count(idx) != 3 {
		t.Fatalf("idx=%d; bm.count(idx),%d != 3", idx, bm.count(idx))
	}
	if idx = 9; bm.count(idx) != 3 {
		t.Fatalf("idx=%d; bm.count(idx),%d != 3", idx, bm.count(idx))
	}
	if idx = 10; bm.count(idx) != 4 {
		t.Fatalf("idx=%d; bm.count(idx),%d != 4", idx, bm.count(idx))
	}
	if idx = hash.MaxIndex; bm.count(idx) != 4 {
		t.Fatalf("idx=%d; bm.count(idx),%d != 4", idx, bm.count(idx))
	}
}
