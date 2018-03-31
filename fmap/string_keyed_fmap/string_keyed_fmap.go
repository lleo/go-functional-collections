package string_keyed_fmap

import (
	"fmt"
	"strings"

	"github.com/lleo/go-functional-collections/fmap"
	"github.com/lleo/go-functional-collections/key"
)

// StringKeyedMap is a wrapper of the fmap.Map structure that stores/returns
// plain strings as the key in the key/value mappings.
type StringKeyedMap fmap.Map

// New return a properly initialize pointer to a StringKeyedMap struct.
func New() *StringKeyedMap {
	//var m = new(StringKeyedMap) //DOESN'T WORK; there is no root table init.
	var m = fmap.New()
	return (*StringKeyedMap)(m)
}

// Get loads the value stored for the given key. If the key doesn't exist in the
// Map a nil is returned. If you need to store nil values and want to
// distinguish between a found existing mapping of the key to nil and a
// non-existent mapping for the key, you must use the Load method.
//
// The key is a plain string value. This method will panic if the key is the
// empty string.
func (m *StringKeyedMap) Get(k string) interface{} {
	if k == "" {
		panic("key is empty string")
	}
	return (*fmap.Map)(m).Get(key.Str(k))
}

// Load retrieves the value related to the string key in the Map data structure.
// It also return a bool to indicate the value was found. This allows you to
// store nil values in the Map data structure and distinguish between a found
// nil key/value mapping and a non-existant key/value mapping.
//
// The key is a plain string value. This method will panic if the key is the
// empty string.
func (m *StringKeyedMap) Load(k string) (interface{}, bool) {
	if k == "" {
		panic("key is empty string")
	}
	return (*fmap.Map)(m).Load(key.Str(k))
}

// LoadOrStore returns the existing value for the key if present. Otherwise,
// it stores a new key/value mapping and returns the given value. The loaded
// result is true if the key/value was loaded, false if a new key/value mapping
// was created. Lastly, if an existing key/value mapping was loaded then the
// returned map is the original *Map, if the a new key/value mapping was
// created returned *Map is a new persistent *Map.
//
// The key is a plain string value. This method will panic if the key is the
// empty string.
func (m *StringKeyedMap) LoadOrStore(k string, v interface{}) (
	*StringKeyedMap, interface{}, bool,
) {
	if k == "" {
		panic("key is empty string")
	}
	var nm, val, found = (*fmap.Map)(m).LoadOrStore(key.Str(k), v)
	return (*StringKeyedMap)(nm), val, found
}

// Put stores a new key/value mapping. It returns a new persistent *Map data
// structure.
//
// The key is a plain string value. This method will panic if the key is the
// empty string.
func (m *StringKeyedMap) Put(k string, v interface{}) *StringKeyedMap {
	if k == "" {
		panic("key is empty string")
	}
	return (*StringKeyedMap)((*fmap.Map)(m).Put(key.Str(k), v))
}

// Store stores a new key/value mapping. It returns a new persistent
// *Map data structure and a bool indicating if a new pair was added (true)
// or if the value merely replaced a prior value (false). Regardless of
// whether a new key/value mapping was created or mearly replaced, a new
// *Map is created.
//
// The key is a plain string value. This method will panic if the key is the
// empty string.
func (m *StringKeyedMap) Store(k string, v interface{}) (
	*StringKeyedMap, bool,
) {
	if k == "" {
		panic("key is empty string")
	}
	var nm, added = (*fmap.Map)(m).Store(key.Str(k), v)
	return (*StringKeyedMap)(nm), added
}

// Del deletes any entry with the given key, but does not indicate if the key
// existed or not. However, if the key did not exist the returned *Map will be
// the original *Map.
//
// The key is a plain string value. This method will panic if the key is the
// empty string.
func (m *StringKeyedMap) Del(k string) *StringKeyedMap {
	if k == "" {
		panic("key is empty string")
	}
	return (*StringKeyedMap)((*fmap.Map)(m).Del(key.Str(k)))
}

// Remove deletes any key/value mapping for the given key. It returns a
// *Map data structure, the possible value that was stored for that key,
// and a boolean idicating if the key was found and deleted. If the key didn't
// exist, then the value is set nil, and the original *Map is returned.
//
// The key is a plain string value. This method will panic if the key is the
// empty string.
func (m *StringKeyedMap) Remove(k string) (
	*StringKeyedMap, interface{}, bool,
) {
	if k == "" {
		panic("key is empty string")
	}
	var nm, val, removed = (*fmap.Map)(m).Remove(key.Str(k))
	return (*StringKeyedMap)(nm), val, removed
}

// Iter is a wrapper around fmap.Iter
type Iter fmap.Iter

// Next returns each sucessive key/value mapping in the *Map. When all enrties
// have been returned it will return an empty string as the key.
func (it *Iter) Next() (string, interface{}) {
	var k, val = (*fmap.Iter)(it).Next()
	var sk key.Str
	var s string

	if k == nil {
		s = ""
	} else {
		sk = k.(key.Str)
		s = string(sk)

	}
	return s, val
}

// Iter returns a *Iter structure. You can call the Next() method on the *Iter
// structure sucessively until it return a nil key value, to walk the key/value
// mappings in the Map data structure. This is safe under any usage of the *Map
// because the Map is immutable.
func (m *StringKeyedMap) Iter() *Iter {
	return (*Iter)((*fmap.Map)(m).Iter())
}

// Range applies the given function for every key/value mapping in the *Map
// data structure. Given that the *Map is immutable there is no danger with
// concurrent use of the *Map while the Range method is executing.
func (m *StringKeyedMap) Range(f func(string, interface{}) bool) {
	var fn = func(k key.Hash, v interface{}) bool {
		var sk key.Str
		var s string
		if k == nil {
			s = ""
		} else {
			sk = k.(key.Str)
			s = string(sk)
		}
		return f(s, v)
	}
	(*fmap.Map)(m).Range(fn)
	return
}

// NumEntries() returns the number of key/value entries in the *Map. This
// operation is O(1), because a current count of the number of entries is
// maintained at the top level of the *Map data structure, so walking the data
// structure is not required to get the current count of key/value entries.
func (m *StringKeyedMap) NumEntries() int {
	return (*fmap.Map)(m).NumEntries()
}

// String prints a string list all the key/value mappings in the *Map. It is
// intended to be simmilar to fmt.Printf("%#v") of a golang builtin map.
func (m *StringKeyedMap) String() string {
	var ents = make([]string, m.NumEntries())
	var i int = 0

	var it = m.Iter()
	for k, v := it.Next(); k != ""; k, v = it.Next() {
		ents[i] = fmt.Sprintf("%q:%#v", k, v)
		i++
	}

	//m.Range(func(k key.Hash, v interface{}) bool {
	//	ents[i] = fmt.Sprintf("%q:%#v", k, v)
	//	i++
	//	return true
	//})

	return "StringKeyedMap{" + strings.Join(ents, ",") + "}"
}
