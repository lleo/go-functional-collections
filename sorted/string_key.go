package sorted

// StringKey is a type wrapper of string values that implements the sorted.Key
// interface.
type StringKey string

// Less returns true if passed another StringKey that it based on a string value
// that is less than the string value of the receiver. String comparison is
// used.
//
// If the sorted.Key value that is passed to the Less method not an IntKey type,
// Less will panic.
func (sk StringKey) Less(o Key) bool {
	var osk, ok = o.(StringKey)
	if !ok {
		panic("o is not a StringKey")
	}
	return sk < osk
}

func (sk StringKey) String() string {
	return string(sk)
}
