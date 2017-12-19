// +build go1.9

package set

import (
	"math/bits"
)

func bitCount8(n uint8) uint {
	return uint(bits.OnesCount8(n))
}
