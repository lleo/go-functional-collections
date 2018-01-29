// +build !go1.9

package fmap

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
