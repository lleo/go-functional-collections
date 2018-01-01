package hash

import (
	"strconv"
)

type IntKey int

const uintSize = 4 << (^uint(0) >> 32 & 1) // 32 or 64

func (ik IntKey) Hash() Val {
	var ib []byte
	var i = int(ik)
	if uintSize == 4 {
		ib = []byte{
			byte(uint(i) & 0xff),
			byte((uint(i) & 0xff00) >> 8),
			byte((uint(i) & 0xff0000) >> 16),
			byte((uint(i) & 0xff000000) >> 24),
		}
	} else if uintSize == 8 {
		ib = []byte{
			byte(uint(i) & 0xff),
			byte((uint(i) & 0xff00) >> 8),
			byte((uint(i) & 0xff0000) >> 16),
			byte((uint(i) & 0xff000000) >> 24),
			byte((uint(i) & 0xff00000000) >> 32),
			byte((uint(i) & 0xff0000000000) >> 40),
			byte((uint(i) & 0xff000000000000) >> 48),
			byte((uint(i) & 0xff00000000000000) >> 56),
		}
	} else {
		panic("invalid int uintSize")
	}
	return CalcHash(ib)
}

func (ik IntKey) Equals(okey Key) bool {
	var oik, ok = okey.(IntKey)
	if !ok {
		return false
	}
	return ik == oik
}

func (ik IntKey) String() string {
	return strconv.FormatInt(int64(ik), 10)
}
