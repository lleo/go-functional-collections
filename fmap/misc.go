// +build go1.9

package fmap

import (
	"math/bits"
)

const uintSize = 32 << (^uint(0) >> 32 & 1)

func pow2GreaterThan(n uint) int {
	//return 1 << (1 + bits.LeadingZeros(0) - bits.LeadingZeros(size))
	return 1 << uint(1+uintSize-bits.LeadingZeros(n))
}
