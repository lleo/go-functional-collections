package hash

// StringKey defines a type of 'string' that hash the methods that satisfy the
// hash.Key interface.
type StringKey string

// Hash calculates the hash.Val if the StringKey receiver every time it is
// called.
func (sk StringKey) Hash() Val {
	return CalcHash([]byte(sk))
}

// Equals determines if the given Key is equivalent, by value, to the receiver.
func (sk StringKey) Equals(okey Key) bool {
	var osk, ok = okey.(StringKey)
	if !ok {
		return false
	}
	return sk == osk
}

// String returns a string representation of the receiver.
func (sk StringKey) String() string {
	return string(sk)
}
