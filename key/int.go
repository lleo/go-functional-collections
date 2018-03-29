package key

import (
	"log"
	"strconv"

	"github.com/lleo/go-functional-collections/hash"
)

type Int int

const uintSizeBits = 32 << (^uint(0) >> 32 & 1)

//const uintSizeBytes = 4 << (^uint(0) >> 32 & 1)

// Int2ByteSlice is a function that gets initialized at runtime based on the
// size of a 'int'.
//
// If the size of an 'int' is neither 4bytes nor 8bytes the function will panic
// if called.
var Int2ByteSlice func(int) []byte

func init() {
	switch uintSizeBits {
	case 32:
		Int2ByteSlice = func(i int) []byte {
			return []byte{
				byte(i),
				byte((uint(i) >> 8)),
				byte((uint(i) >> 16)),
				byte((uint(i) >> 24)),
			}
		}
	case 64:
		Int2ByteSlice = func(i int) []byte {
			return []byte{
				byte(i),
				byte((uint(i) >> 8)),
				byte((uint(i) >> 16)),
				byte((uint(i) >> 24)),
				byte((uint(i) >> 32)),
				byte((uint(i) >> 40)),
				byte((uint(i) >> 48)),
				byte((uint(i) >> 56)),
			}
		}
	default:
		Int2ByteSlice = func(i int) []byte {
			log.Panicf("Int2ByteSlice not defined for uintSizeBits=%d\n",
				uintSizeBits)
			return nil
		}
	}
}

func (ik Int) Less(okey Sort) bool {
	var oik, ok = okey.(Int)
	if !ok {
		panic("okey is not a key.Int")
		//return false
	}
	return ik < oik
}

// Hash calculates the hash.Val of the Int receiver every time it is called.
func (ik Int) Hash() hash.Val {
	var ib = Int2ByteSlice(int(ik))
	return hash.Calculate(ib)
}

// Equals determines if the given Key is equivalent, by value, to the receiver.
func (ik Int) Equals(okey Hash) bool {
	var oik, ok = okey.(Int)
	if !ok {
		return false
	}
	return ik == oik
}

// String returns a string representation of the receiver.
func (ik Int) String() string {
	return strconv.FormatInt(int64(ik), 10)
}
