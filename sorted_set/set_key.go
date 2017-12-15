package sorted_set

type SetKey interface {
	Less(SetKey) bool
	String() string
}

// nInf is a SetKey for negative infinity
type nInf struct{}

func (nInf) Less(SetKey) bool {
	return true
}

func (nInf) String() string {
	return "nInf"
}

// pInf is a SetKey for positive infinity
type pInf struct{}

func (pInf) Less(SetKey) bool {
	return false
}

func (pInf) String() string {
	return "pInf"
}

var (
	ninf = nInf{}
	pinf = pInf{}
)

//InfKey() if passed a non-negative iteger it will return a key that is greater
//than any oter key, other wise (for a negetive integer) it will return a key
//that is less than any other key.
func InfKey(sign int) SetKey {
	if sign < 0 {
		return ninf
	}
	return pinf
}

func less(x, y SetKey) bool {
	if x == pinf || y == ninf {
		return false
	}
	if x == ninf || y == pinf {
		return true
	}
	return x.Less(y)
}

func cmp(x, y SetKey) int {
	if less(x, y) {
		return -1
	} else if less(y, x) {
		return 1
	}
	return 0
}
