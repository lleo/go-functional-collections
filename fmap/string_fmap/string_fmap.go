package string_keyed_fmap

import (
	"fmt"
	"strings"

	"github.com/lleo/go-functional-collections/fmap"
	"github.com/lleo/go-functional-collections/hash"
)

type StringKeyedMap fmap.Map

func New() *StringKeyedMap {
	var m = new(StringKeyedMap)
	return m
}

func (m *StringKeyedMap) Get(k string) interface{} {
	return (*fmap.Map)(m).Get(hash.StringKey(k))
}

func (m *StringKeyedMap) Load(k string) (interface{}, bool) {
	return (*fmap.Map)(m).Load(hash.StringKey(k))
}

func (m *StringKeyedMap) LoadOrStore(k string, v interface{}) (
	*StringKeyedMap, interface{}, bool,
) {
	var nm, val, found = (*fmap.Map)(m).LoadOrStore(hash.StringKey(k), v)
	return (*StringKeyedMap)(nm), val, found
}

func (m *StringKeyedMap) Put(k string, v interface{}) *StringKeyedMap {
	return (*StringKeyedMap)((*fmap.Map)(m).Put(hash.StringKey(k), v))
}

func (m *StringKeyedMap) Store(k string, v interface{}) (
	*StringKeyedMap, bool,
) {
	var nm, added = (*fmap.Map)(m).Store(hash.StringKey(k), v)
	return (*StringKeyedMap)(nm), added
}

func (m *StringKeyedMap) Del(k string) *StringKeyedMap {
	return (*StringKeyedMap)((*fmap.Map)(m).Del(hash.StringKey(k)))
}

func (m *StringKeyedMap) Remove(k string) (
	*StringKeyedMap, interface{}, bool,
) {
	var nm, val, removed = (*fmap.Map)(m).Remove(hash.StringKey(k))
	return (*StringKeyedMap)(nm), val, removed
}

type Iter fmap.Iter

func (it *Iter) Next() (string, interface{}) {
	var k, val = (*fmap.Iter)(it).Next()
	var sk hash.StringKey
	var s string

	if k == nil {
		s = ""
	} else {
		sk = k.(hash.StringKey)
		s = string(sk)

	}
	return s, val
}

func (m *StringKeyedMap) Iter() *Iter {
	return (*Iter)((*fmap.Map)(m).Iter())
}

func (m *StringKeyedMap) Range(f func(string, interface{}) bool) {
	var fn = func(mk hash.Key, v interface{}) bool {
		var sk = mk.(hash.StringKey)
		return f(string(sk), v)
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

	var it = m.Iter()
	for k, v := it.Next(); k != ""; k, v = it.Next() {
		ents[i] = fmt.Sprintf("%q:%#v", k, v)
		i++
	}

	//m.Range(func(k hash.Key, v interface{}) bool {
	//	ents[i] = fmt.Sprintf("%q:%#v", k, v)
	//	i++
	//	return true
	//})

	return "StringKeyedMap{" + strings.Join(ents, ",") + "}"
}
