package sorted_set_test

import (
	"math/rand"
	"testing"

	"github.com/lleo/go-functional-collections/sorted"
	"github.com/lleo/go-functional-collections/sorted_set"
)

func buildKeys(numKeys, numKeysXtra int) ([]sorted.Key, []sorted.Key) {
	var keys = make([]sorted.Key, numKeys+numKeysXtra)

	for i := 0; i < numKeys+numKeysXtra; i++ {
		keys[i] = sorted.IntKey(i)
	}

	//randomize keys
	for i := len(keys) - 1; i > 0; i-- {
		var j = rand.Intn(i + 1)
		keys[i], keys[j] = keys[j], keys[i]
	}

	var xtra = keys[len(keys)-numKeysXtra:]
	keys = keys[:len(keys)-numKeysXtra]

	return keys, xtra
}

func buildSet(keys []sorted.Key) *sorted_set.Set {
	var s = sorted_set.New()
	for _, key := range keys {
		s = s.Set(key)
	}
	return s
}

const NumKeys10 = 1 * 10
const NumKeys100 = 1 * 100
const NumKeys1M = 1 * 1000
const NumKeys10M = 10 * 1000
const NumKeys100M = 100 * 1000
const NumKeys1MM = 1 * 1000 * 1000
const NumKeys10MM = 10 * 1000 * 1000

const NumKeysExtra10 = 2 * (NumKeys10 / 10)
const NumKeysExtra100 = 2 * (NumKeys100 / 10)
const NumKeysExtra1M = 2 * (NumKeys1M / 10)
const NumKeysExtra10M = 2 * (NumKeys10M / 10)
const NumKeysExtra100M = 20 * (NumKeys100M / 10)
const NumKeysExtra1MM = 20 * (NumKeys1MM / 10)
const NumKeysExtra10MM = 20 * (NumKeys10MM / 10)

var SSet10 *sorted_set.Set
var SSet100 *sorted_set.Set
var SSet1M *sorted_set.Set
var SSet10M *sorted_set.Set
var SSet100M *sorted_set.Set
var SSet1MM *sorted_set.Set
var SSet10MM *sorted_set.Set

var XtraKeys10 []sorted.Key
var XtraKeys100 []sorted.Key
var XtraKeys1M []sorted.Key
var XtraKeys10M []sorted.Key
var XtraKeys100M []sorted.Key
var XtraKeys1MM []sorted.Key
var XtraKeys10MM []sorted.Key

func BenchmarkSetOne10(b *testing.B) {
	//log.Printf("BenchmarkSetOne10: called b.N=%d\n", b.N)
	var keys, XtraKeys10 []sorted.Key
	if SSet10 == nil || XtraKeys10 == nil {
		//log.Println("Generating Sset10 & XtraKeys10")
		keys, XtraKeys10 = buildKeys(NumKeys10, NumKeysExtra10)
		SSet10 = buildSet(keys)
	}
	var s = SSet10
	var xtraKeys = XtraKeys10
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var i = rand.Int() % NumKeysExtra10
		_ = s.Set(xtraKeys[i])
	}
}

func BenchmarkSetOne100(b *testing.B) {
	//log.Printf("BenchmarkSetOne100: called b.N=%d\n", b.N)
	var keys, XtraKeys100 []sorted.Key
	if SSet100 == nil || XtraKeys100 == nil {
		//log.Println("Generating Sset100 & XtraKeys100")
		keys, XtraKeys100 = buildKeys(NumKeys100, NumKeysExtra100)
		SSet100 = buildSet(keys)
	}
	var s = SSet100
	var xtraKeys = XtraKeys100
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var i = rand.Int() % NumKeysExtra100
		_ = s.Set(xtraKeys[i])
	}
}

func BenchmarkSetOne1M(b *testing.B) {
	//log.Printf("BenchmarkSetOne1M: called b.N=%d\n", b.N)
	var keys, XtraKeys1M []sorted.Key
	if SSet1M == nil || XtraKeys1M == nil {
		//log.Println("Generating Sset1M & XtraKeys1M")
		keys, XtraKeys1M = buildKeys(NumKeys1M, NumKeysExtra1M)
		SSet1M = buildSet(keys)
	}
	var s = SSet1M
	var xtraKeys = XtraKeys1M
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var i = rand.Int() % NumKeysExtra1M
		_ = s.Set(xtraKeys[i])
	}
}

func BenchmarkSetOne10M(b *testing.B) {
	//log.Printf("BenchmarkSetOne10M: called b.N=%d\n", b.N)
	var keys, XtraKeys10M []sorted.Key
	if SSet10M == nil || XtraKeys10M == nil {
		//log.Println("Generating Sset10M & XtraKeys10M")
		keys, XtraKeys10M = buildKeys(NumKeys10M, NumKeysExtra10M)
		SSet10M = buildSet(keys)
	}
	var s = SSet10M
	var xtraKeys = XtraKeys10M
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var i = rand.Int() % NumKeysExtra10M
		_ = s.Set(xtraKeys[i])
	}
}

func BenchmarkSetOne100M(b *testing.B) {
	//log.Printf("BenchmarkSetOne100M: called b.N=%d\n", b.N)
	var keys, XtraKeys100M []sorted.Key
	if SSet100M == nil || XtraKeys100M == nil {
		//log.Println("Generating Sset100M & XtraKeys100M")
		keys, XtraKeys100M = buildKeys(NumKeys100M, NumKeysExtra100M)
		SSet100M = buildSet(keys)
	}
	var s = SSet100M
	var xtraKeys = XtraKeys100M
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var i = rand.Int() % NumKeysExtra100M
		_ = s.Set(xtraKeys[i])
	}
}

func BenchmarkSetOne1MM(b *testing.B) {
	//log.Printf("BenchmarkSetOne1MM: called b.N=%d\n", b.N)
	var keys, XtraKeys1MM []sorted.Key
	if SSet1MM == nil || XtraKeys1MM == nil {
		//log.Println("Generating Sset1MM & XtraKeys1MM")
		keys, XtraKeys1MM = buildKeys(NumKeys1MM, NumKeysExtra1MM)
		SSet1MM = buildSet(keys)
	}
	var s = SSet1MM
	var xtraKeys = XtraKeys1MM
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var i = rand.Int() % NumKeysExtra1MM
		_ = s.Set(xtraKeys[i])
	}
}

func BenchmarkSetOne10MM(b *testing.B) {
	//log.Printf("BenchmarkSetOne10MM: called b.N=%d\n", b.N)
	var keys, XtraKeys10MM []sorted.Key
	if SSet10MM == nil || XtraKeys10MM == nil {
		//log.Println("Generating Sset10MM & XtraKeys10MM")
		keys, XtraKeys10MM = buildKeys(NumKeys10MM, NumKeysExtra10MM)
		SSet10MM = buildSet(keys)
	}
	var s = SSet10MM
	var xtraKeys = XtraKeys10MM
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		//START HERE
		var i = rand.Int() % NumKeysExtra10MM
		_ = s.Set(xtraKeys[i])
	}
}
