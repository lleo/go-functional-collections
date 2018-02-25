// +build !go1.9

package set

const uintSize = 32 << (^uint(0) >> 32 & 1)

func pow2GreaterThan(n uint) int {
	var m uint = 2
	for i := 1; i < uintSize; i++ {
		m = m << 1
		if m > n {
			return m
		}
	}
	return m
}

// topBit returns the position of the first bit set. If no bits are set, then it
// returns 0. So positions are 1 to uintSize (32 or 64).
func topBit(n uint) uint {
	var i uint
	for i = 0; i < uintSize; i++ { // for i = 1 to 32
		if (n >> i) == 0 {
			break
		}
	}
	return i
	//for i := uint(0); i < uintSize; i++ {
	//	if (n >> i) == 0 {
	//		return i
	//	}
	//}
	//return uintSize
}
