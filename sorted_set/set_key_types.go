package sorted_set

import "strconv"

type StringKey string

func (sk StringKey) Less(o SetKey) bool {
	var osk, ok = o.(StringKey)
	if !ok {
		panic("o is not a StringKey")
	}
	return sk < osk
}

func (sk StringKey) String() string {
	return string(sk)
}

type IntKey int

func (ik IntKey) Less(o SetKey) bool {
	var oik, ok = o.(IntKey)
	if !ok {
		panic("o is not a IntKey")
	}
	return ik < oik
}

func (ik IntKey) String() string {
	return strconv.Itoa(int(ik))
}
