package fmap

import (
	"fmt"
	"strings"

	"github.com/lleo/go-functional-collections/key/hash"
)

// bitmapShift is 3 because we are using uint8 as the base bitmap type.
const bitmapShift uint = 3

// bitmapSize is the number of uint8 needed to cover hash.IndexLimit bits.
const bitmapSize uint = (hash.IndexLimit + (1 << bitmapShift) - 1) /
	(1 << bitmapShift)

type bitmap [bitmapSize]uint8

// bitFmtStr cannot be a constant if I use a library function to convert
// hash.IndexLimit to a string, so I am going to hardcode it manually.
//const bitsFmtStr = fmt.Sprintf("%%0%db", 1<<bitmapShift)
//const bitsFmtStr = "%0" + strconv.Itoa(1<<bitmapShift) + "b"
//const bitsFmtStr = "%032b" //for bitmap being array of uint32 aka bitmapShift=5
const bitsFmtStr = "%08b" //for bitmap being array of uint8 aka bitmapShift=3

const byteMask = (1 << bitmapShift) - 1

func (bm *bitmap) String() string {
	// Show all bits in bitmap because hash.IndexLimit is a multiple of the
	// bitmap base type.
	var strs = make([]string, bitmapSize)
	//var fmtStr = fmt.Sprintf("%%0%db", 1<<bitmapShift)
	for i := uint(0); i < bitmapSize; i++ {
		strs[i] = fmt.Sprintf(bitsFmtStr, bm[i])
	}

	return strings.Join(strs, " ")
}

// IsSet returns a bool indicating whether the i'th bit is 1(true) or 0(false).
func (bm *bitmap) isSet(idx uint) bool {
	var nth = idx >> bitmapShift
	var bit = idx & byteMask

	return (bm[nth] & (1 << bit)) > 0
}

// Set marks the i'th bit 1.
func (bm *bitmap) set(idx uint) *bitmap {
	var nth = idx >> bitmapShift
	var bit = idx & byteMask

	bm[nth] |= 1 << bit

	return bm
}

// Unset marks the i'th bit to 0.
func (bm *bitmap) unset(idx uint) *bitmap {
	var nth = idx >> bitmapShift
	var bit = idx & byteMask

	//if bm[nth]&(1<<bit) > 0 {
	//	bm[nth] &^= 1 << bit
	//}
	bm[nth] &^= 1 << bit

	return bm
}

// Count returns the numbers of bits set below the i'th bit.
func (bm *bitmap) count(idx uint) uint {
	var nth = idx >> bitmapShift
	var bit = idx & byteMask

	var count uint
	for i := uint(0); i < nth; i++ {
		count += bitCount8(bm[i])
	}

	count += bitCount8(bm[nth] & ((1 << bit) - 1))

	return count
}
