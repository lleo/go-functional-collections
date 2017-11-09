package string_keyed_map

import (
	"fmt"
	"strings"

	"github.com/lleo/go-functional-collections/fmap"
)

type StringKeyedMap fmap.Map

func New() *StringKeyedMap {
	var m = new(StringKeyedMap)
	return m
}

func (m *StringKeyedMap) Get(s string) interface{} {
	return (*fmap.Map)(m).Get(StringKey(s))
}

func (m *StringKeyedMap) Load(s string) (interface{}, bool) {
	return (*fmap.Map)(m).Load(StringKey(s))
}

func (m *StringKeyedMap) LoadOrStore(
	s string,
	v interface{},
) (*StringKeyedMap, interface{}, bool) {
	var nm, val, found = (*fmap.Map)(m).LoadOrStore(StringKey(s), v)
	return (*StringKeyedMap)(nm), val, found
}

func (m *StringKeyedMap) Put(s string, v interface{}) *StringKeyedMap {
	return (*StringKeyedMap)((*fmap.Map)(m).Put(StringKey(s), v))
}

func (m *StringKeyedMap) Store(
	s string,
	v interface{},
) (*StringKeyedMap, bool) {
	var nm, added = (*fmap.Map)(m).Store(StringKey(s), v)
	return (*StringKeyedMap)(nm), added
}

func (m *StringKeyedMap) Del(s string) *StringKeyedMap {
	return (*StringKeyedMap)((*fmap.Map)(m).Del(StringKey(s)))
}

func (m *StringKeyedMap) Delete(s string) *StringKeyedMap {
	return (*StringKeyedMap)((*fmap.Map)(m).Delete(StringKey(s)))
}

func (m *StringKeyedMap) Remove(s string) (
	*StringKeyedMap,
	interface{},
	bool,
) {
	var nm, val, removed = (*fmap.Map)(m).Remove(StringKey(s))
	return (*StringKeyedMap)(nm), val, removed
}

func (m *StringKeyedMap) Range(f func(string, interface{}) bool) {
	var fn = func(mk fmap.MapKey, v interface{}) bool {
		var sk = mk.(StringKey)
		var s = string(sk)
		return f(s, v)
	}
	(*fmap.Map)(m).Range(fn)
	return
}

func (m *StringKeyedMap) NumEntries() uint {
	return (*fmap.Map)(m).NumEntries()
}

func (m *StringKeyedMap) String() string {
	var ents = make([]string, m.NumEntries())
	var i int = 0
	m.Range(func(k string, v interface{}) bool {
		ents[i] = fmt.Sprintf("%q:%#v", k, v)
		i++
		return true
	})
	return "[" + strings.Join(ents, ",") + "]"
}

func (m *StringKeyedMap) LongString(indent string) string {
	return (*fmap.Map)(m).LongString(indent)
}

func (m *StringKeyedMap) Stats() *fmap.Stats {
	return (*fmap.Map)(m).Stats()
}
