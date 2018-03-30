package set_test

import (
	"sort"
	"testing"

	"github.com/lleo/go-functional-collections/key"
	"github.com/lleo/go-functional-collections/set"
)

func TestBasicBuildSimpleSet(t *testing.T) {
	var s = set.New()
	s = s.
		Set(key.Str("a")).
		Set(key.Str("b")).
		Set(key.Str("c"))

	if !s.IsSet(key.Str("a")) {
		t.Fatal("s.IsSet(\"a\") not true")
	}

	if !s.IsSet(key.Str("b")) {
		t.Fatal("s.IsSet(\"b\") not true")
	}

	if !s.IsSet(key.Str("c")) {
		t.Fatal("s.IsSet(\"c\") not true")
	}
}

func TestBasicIsSet(t *testing.T) {
	var s = set.New()
	s = s.Set(key.Str("a"))

	var found = s.IsSet(key.Str("a"))
	if !found {
		t.Fatal("failed to s.IsSet(\"a\")")
	}

	found = s.IsSet(key.Str("b"))
	if found {
		t.Fatal("WTF! \"b\" found")
	}
}

func TestBasicAdd(t *testing.T) {
	var s = set.New()

	var added bool
	s, added = s.Add(key.Str("a"))

	if !added {
		t.Fatal("added for s.Add(\"a\") is false")
	} else {
		if !s.IsSet(key.Str("a")) {
			t.Fatal("s.Add(\"a\") was added, but s.IsSet(\"a\") not true")
		}
	}

	s, added = s.Add(key.Str("a"))
	if added {
		t.Fatal("added == true for second s.Add(\"a\")")
	} else {
		if !s.IsSet(key.Str("a")) {
			t.Fatal("s.Add(\"a\") was added, but s.IsSet(\"a\") not true")
		}
	}
}

func TestBasicUnset(t *testing.T) {
	var s = set.New()

	if s.NumEntries() != 0 {
		t.Fatal("s.NumEntries() != 0")
	}

	s = s.
		Set(key.Str("a")).
		Set(key.Str("b")).
		Set(key.Str("c"))

	if s.NumEntries() != 3 {
		t.Fatal("s.NumEntries() != 3")
	}

	s = s.Unset(key.Str("b"))

	if s.NumEntries() != 2 {
		t.Fatal("s.NumEntries() != 2")
	}

	if s.IsSet(key.Str("b")) {
		t.Fatal("found \"b\" after s.Unset(\"b\")")
	}

	s = s.Unset(key.Str("a")).Unset(key.Str("c"))

	if s.NumEntries() != 0 {
		t.Fatal("s.NumEntries() != 0")
	}
}

func TestBasicRemove(t *testing.T) {
	var s = set.New()

	s = s.
		Set(key.Str("a")).
		Set(key.Str("b")).
		Set(key.Str("c"))

	if s.NumEntries() != 3 {
		t.Fatal("s.NumEntries() != 3")
	}

	var found bool
	var key key.Str

	key = "d"
	s, found = s.Remove(key)
	if found {
		t.Fatalf("found key=%#v that does not exist.", key)
	}

	key = "b"
	s, found = s.Remove(key)
	if !found {
		t.Fatalf("failed to find & remove key=%#v", key)
	}

	s, found = s.Remove(key)
	if found {
		t.Fatalf("found key=%#v entry for key just Removed;", key)
	}
}

func TestBasicRange(t *testing.T) {
	var s = set.New().
		Set(key.Str("a")).
		Set(key.Str("b")).
		Set(key.Str("c"))

	var keys = make([]key.Hash, s.NumEntries())
	var i int
	s.Range(func(k key.Hash) bool {
		keys[i] = k
		i++
		return true
	})
	sort.Slice(keys, func(i, j int) bool {
		ki := keys[i].(key.Str)
		kj := keys[j].(key.Str)
		return string(ki) < string(kj)
	})
	var str = "a"
	for _, k := range keys {
		var sk = key.Str(str)
		if !k.Equals(sk) {
			t.Fatalf("k,%s != sk,%s", s.String(), sk.String())
		}
		str = Inc(str)
	}
}

func TestBasicString(t *testing.T) {
	var s = set.New()
	s = s.
		Set(key.Str("a")).
		Set(key.Str("b")).
		Set(key.Str("c"))

	var str = s.String()
	//log.Printf("s.String()=%s\n", str)

	var expectedStr = "Set{\"c\",\"b\",\"a\"}"
	if str != expectedStr {
		t.Fatalf("str,%q != expectedStr,%q", str, expectedStr)
	}
}

func TestBasicCount(t *testing.T) {
	//var s = set.New().
	//	Set(key.Str("a")).
	//	Set(key.Str("b")).
	//	Set(key.Str("c"))

	var s = set.New()
	var str = "a"
	for i := 0; i < 100000; i++ {
		s = s.Set(key.Str(str))
		str = Inc(str)
	}

	if s.NumEntries() != s.Count() {
		t.Fatalf("s.NumEntries(),%d != s.Count(),%d", s.NumEntries(), s.Count())
	}
}

func TestBasicNewFromList(t *testing.T) {
	var keys = buildKeys(100)

	var s = set.NewFromList(keys)

	for _, k := range keys {
		var isSet = s.IsSet(k)
		if !isSet {
			t.Fatalf("key=%s is not set", k)
		}
	}
}

func TestBasicBulkInsert(t *testing.T) {
	var keys = buildKeys(100)
	var s = set.NewFromList(keys[:50])

	s = s.BulkInsert(keys[50:])

	if s.NumEntries() != 100 {
		t.Fatalf("s.NumEntries(),%d != 100", s.NumEntries())
	}

	for _, k := range keys {
		var isSet = s.IsSet(k)
		if !isSet {
			t.Fatalf("key=%s is not set", k)
		}
	}
}

func TestBasicBulkInsertConflict(t *testing.T) {
	var keys = buildKeys(100)
	var s = set.NewFromList(keys[:70])
	var s0 = s.BulkInsert(keys[50:])

	if s0.NumEntries() != 100 {
		t.Fatalf("s0.NumEntries(),%d != 100", s0.NumEntries())
	}

	for _, k := range keys {
		var isSet = s0.IsSet(k)
		if !isSet {
			t.Fatalf("key=%s is not set", k)
		}
	}
}

func TestBasicMerge(t *testing.T) {
	var keys = buildKeys(100)

	var s0 = set.NewFromList(keys[:50])
	var s1 = set.NewFromList(keys[50:])
	var s = s0.Merge(s1)

	for _, k := range keys {
		var isSet = s.IsSet(k)
		if !isSet {
			t.Fatalf("k=%s is not set", k)
		}
	}
}

func TestBasicMergeConflict(t *testing.T) {
	var keys = buildKeys(100)

	var s0 = set.NewFromList(keys[:60])
	var s1 = set.NewFromList(keys[50:])
	var s = s0.Merge(s1)

	for _, k := range keys {
		var isSet = s.IsSet(k)
		if !isSet {
			t.Fatalf("k=%s is not set", k)
		}
	}
}

func TestBasicBulkDelete(t *testing.T) {
	var keys = buildKeys(100)

	var origSet = set.NewFromList(keys)
	var copySet = origSet.DeepCopy()

	if !origSet.Equiv(copySet) {
		t.Fatal("origSet != copySet")
	}

	var s *set.Set
	var notFound []key.Hash
	s, notFound = origSet.BulkDelete(keys[50:])

	if !origSet.Equiv(copySet) {
		t.Fatal("origSet != copySet after BulkDelete")
	}

	if len(notFound) != 0 {
		t.Fatalf("len(notFound),%d != 0", len(notFound))
	}

	if numEnts := s.NumEntries(); numEnts != 50 {
		t.Fatalf("s.NumEntries(),%d != 50", numEnts)
	}

	if count := s.Count(); count != 50 {
		t.Fatalf("s.Count(),%d != 50", count)
	}

	for _, k := range keys[:50] {
		var isSet = s.IsSet(k)
		if !isSet {
			t.Fatalf("k=%s is not set", k)
		}
	}
}

func isMember(k key.Hash, keys []key.Hash) bool {
	for _, key := range keys {
		if k.Equals(key) {
			return true
		}
	}
	return false
}

func TestBasicBulkDeleteNotFound(t *testing.T) {
	var keys = buildKeys(100)

	var s = set.NewFromList(keys[:70])

	//log.Println(s.TreeString(""))

	var notFound []key.Hash
	s, notFound = s.BulkDelete(keys)

	if len(notFound) != 30 {
		t.Fatalf("len(notFound),%d != 30", len(notFound))
	}

	if s.NumEntries() != 0 {
		t.Fatalf("s.NumEntries(),%d != 0", s.NumEntries())
	}

	if count := s.Count(); count != 0 {
		t.Fatalf("s.Count(),%d != 0", count)
	}

	// slice keys[70:] not added to s *Map.
	for _, k := range keys[70:] {
		if !isMember(k, notFound) {
			t.Fatalf("expected to find k=%s in notFound", k)
		}
	}
}

func TestBasicDifference(t *testing.T) {
	var tot = 10
	var big, sml = tot * 6 / 10, tot * 4 / 10
	var keys = buildKeys(tot)
	var setA = set.NewFromList(keys[:big])
	var copySetA = setA.DeepCopy()
	var setB = set.NewFromList(keys[sml:])
	var copySetB = setB.DeepCopy()

	var diffSet = setA.Difference(setB)
	var diffKeys = diffSet.Keys()
	sort.Slice(diffKeys, func(i, j int) bool {
		var ki = diffKeys[i].(key.Str)
		var kj = diffKeys[j].(key.Str)
		return string(ki) < string(kj)
	})
	for i, k := range diffKeys {
		var expectedK = keys[i]
		//log.Printf("i=%d; k=%q ?= expectedK=%q",
		//	i, k.String(), expectedK.String())
		if k.String() != expectedK.String() {
			t.Fatalf("%q != %q", k.String(), expectedK.String())
		}
	}
	var diffLen = len(diffKeys)
	var expectedLen = len(keys[:big]) - len(keys[sml:big])
	//log.Printf("(len(keys[:%d]),%d-len(keys[%d:%d]),%d),%d ?= len(diffKeys),%d",
	//	big, len(keys[:big]), sml, big, len(keys[sml:big]), expectedLen, diffLen)
	if diffLen != expectedLen {
		t.Fatalf("diffLen,%d != expectedLen,%d", diffLen, expectedLen)
	}

	if !copySetA.Equiv(setA) {
		t.Fatal("!copySetA.Equiv(setA)")
	}
	if !copySetB.Equiv(setB) {
		t.Fatal("!copySetB.Equiv(setB)")
	}
}

//func TestBasicDifference2(t *testing.T) {
//	var tot = 10
//	var big, sml = tot * 6 / 10, tot * 4 / 10
//	var keys = buildKeys(tot)
//	var setA = set.NewFromList(keys[:big])
//	var copySetA = setA.DeepCopy()
//	var setB = set.NewFromList(keys[sml:])
//	var copySetB = setB.DeepCopy()
//
//	var diffSet = setA.Difference2(setB)
//	var diffKeys = diffSet.Keys()
//	sort.Slice(diffKeys, func(i, j int) bool {
//		var ki = diffKeys[i].(key.Str)
//		var kj = diffKeys[j].(key.Str)
//		return string(ki) < string(kj)
//	})
//	for i, k := range diffKeys {
//		var expectedK = keys[i]
//		//log.Printf("i=%d; k=%q ?= expectedK=%q",
//		//	i, k.String(), expectedK.String())
//		if k.String() != expectedK.String() {
//			t.Fatalf("%q != %q", k.String(), expectedK.String())
//		}
//	}
//	var diffLen = len(diffKeys)
//	var expectedLen = len(keys[:big]) - len(keys[sml:big])
//	//log.Printf("(len(keys[:%d]),%d-len(keys[%d:%d]),%d),%d ?= len(diffKeys),%d",
//	//	big, len(keys[:big]), sml, big, len(keys[sml:big]), expectedLen, diffLen)
//	if diffLen != expectedLen {
//		t.Fatalf("diffLen,%d != expectedLen,%d", diffLen, expectedLen)
//	}
//
//	if !copySetA.Equiv(setA) {
//		t.Fatal("!copySetA.Equiv(setA)")
//	}
//	if !copySetB.Equiv(setB) {
//		t.Fatal("!copySetB.Equiv(setB)")
//	}
//}
