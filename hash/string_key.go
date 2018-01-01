package hash

type StringKey string

func (sk StringKey) Hash() Val {
	return CalcHash([]byte(sk))
}

func (sk StringKey) Equals(okey Key) bool {
	var osk, ok = okey.(StringKey)
	if !ok {
		return false
	}
	return sk == osk
}

func (sk StringKey) String() string {
	return string(sk)
}
