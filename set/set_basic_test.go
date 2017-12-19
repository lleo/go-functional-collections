package set_test

import (
	"testing"

	"github.com/lleo/go-functional-collections/set"
)

func Test_Basic_ButildSimpleSet(t *testing.T) {
	var s = set.New()
	s = s.
		Set(StringKey("a")).
		Set(StringKey("b")).
		Set(StringKey("c"))

	if !s.IsSet(StringKey("a")) {
		t.Fatal("s.IsSet(\"a\") not true")
	}

	if !s.IsSet(StringKey("b")) {
		t.Fatal("s.IsSet(\"b\") not true")
	}

	if !s.IsSet(StringKey("c")) {
		t.Fatal("s.IsSet(\"c\") not true")
	}
}

func Test_Basic_IsSet(t *testing.T) {
	var s = set.New()
	s = s.Set(StringKey("a"))

	var found = s.IsSet(StringKey("a"))
	if !found {
		t.Fatal("failed to s.IsSet(\"a\")")
	}

	found = s.IsSet(StringKey("b"))
	if found {
		t.Fatal("WTF! \"b\" found")
	}
}

func Test_Basic_Add(t *testing.T) {
	var s = set.New()

	var added bool
	s, added = s.Add(StringKey("a"))

	if !added {
		t.Fatal("added for s.Add(\"a\") is false")
	} else {
		if !s.IsSet(StringKey("a")) {
			t.Fatal("s.Add(\"a\") was added, but s.IsSet(\"a\") not true")
		}
	}

	s, added = s.Add(StringKey("a"))
	if added {
		t.Fatal("added == true for second s.Add(\"a\")")
	} else {
		if !s.IsSet(StringKey("a")) {
			t.Fatal("s.Add(\"a\") was added, but s.IsSet(\"a\") not true")
		}
	}
}

func Test_Basic_Unset(t *testing.T) {
	var s = set.New()

	if s.NumEntries() != 0 {
		t.Fatal("s.NumEntries() != 0")
	}

	s = s.
		Set(StringKey("a")).
		Set(StringKey("b")).
		Set(StringKey("c"))

	if s.NumEntries() != 3 {
		t.Fatal("s.NumEntries() != 3")
	}

	s = s.Unset(StringKey("b"))

	if s.NumEntries() != 2 {
		t.Fatal("s.NumEntries() != 2")
	}

	if s.IsSet(StringKey("b")) {
		t.Fatal("found \"b\" after s.Unset(\"b\")")
	}

	s = s.Unset(StringKey("a")).Unset(StringKey("c"))

	if s.NumEntries() != 0 {
		t.Fatal("s.NumEntries() != 0")
	}
}

func Test_Basic_Remove(t *testing.T) {
	var s = set.New()

	s = s.
		Set(StringKey("a")).
		Set(StringKey("b")).
		Set(StringKey("c"))

	if s.NumEntries() != 3 {
		t.Fatal("s.NumEntries() != 3")
	}

	var found bool
	var key set.SetKey

	key = StringKey("d")
	s, found = s.Remove(key)
	if found {
		t.Fatalf("found key=%#v that does not exist.", key)
	}

	key = StringKey("b")
	s, found = s.Remove(key)
	if !found {
		t.Fatalf("failed to find & remove key=%#v", key)
	}

	s, found = s.Remove(key)
	if found {
		t.Fatalf("found key=%#v entry for key just Removed;", key)
	}
}
