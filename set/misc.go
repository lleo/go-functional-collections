// +build go1.9

package set

import (
	"math/bits"
)

const uintSize = 32 << (^uint(0) >> 32 & 1)

func pow2GreaterThan(n uint) int {
	//return 1 << (1 + bits.LeadingZeros(0) - bits.LeadingZeros(size))
	return 1 << uint(1+uintSize-bits.LeadingZeros(n))
}

// topBit returns the position of the first bit set. If no bits are set, then it
// returns 0. So positions are 1 to uintSize (32 or 64).
func topBit(n uint) uint {
	return uintSize - uint(bits.LeadingZeros(n))
}
