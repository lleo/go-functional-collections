package set_test

import (
	"log"
	"math/rand"
	"testing"

	"github.com/lleo/go-functional-collections/key"
	"github.com/lleo/go-functional-collections/set"
)

func buildKeysBench(numKeys, numKeysXtra int) ([]key.Hash, []key.Hash) {
	var keys = make([]key.Hash, numKeys+numKeysXtra)

	var s = "a"
	for i := 0; i < numKeys+numKeysXtra; i++ {
		keys[i] = key.Str(s)
		s = Inc(s)
	}

	//randomize keys
	//https://en.wikipedia.org/wiki/Fisher–Yates_shuffle#The_modern_algorithm
	for i := len(keys) - 1; i > 0; i-- {
		var j = rand.Intn(i + 1)
		keys[i], keys[j] = keys[j], keys[i]
	}

	var xtra = keys[len(keys)-numKeysXtra:]
	keys = keys[:len(keys)-numKeysXtra]

	return keys, xtra
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

var Set10 *set.Set
var Set100 *set.Set
var Set1M *set.Set
var Set10M *set.Set
var Set100M *set.Set
var Set1MM *set.Set
var Set10MM *set.Set

var XtraKeys10 []key.Hash
var XtraKeys100 []key.Hash
var XtraKeys1M []key.Hash
var XtraKeys10M []key.Hash
var XtraKeys100M []key.Hash
var XtraKeys1MM []key.Hash
var XtraKeys10MM []key.Hash

func BenchmarkSetOne10(b *testing.B) {
	var xtra = XtraKeys10
	var m = Set10
	if m == nil {
		var keys []key.Hash
		keys, xtra = buildKeysBench(NumKeys10, NumKeysExtra10)
		m = buildSet(keys)
		XtraKeys10 = xtra
		Set10 = m
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var j = rand.Int() % len(xtra)
		var key = xtra[j]
		_ = m.Set(key)
	}
}

func BenchmarkSetOne100(b *testing.B) {
	log.Printf("BenchmarkSetOne100: b.N=%d\n", b.N)
	var xtra = XtraKeys100
	var m = Set100
	if m == nil {
		var keys []key.Hash
		keys, xtra = buildKeysBench(NumKeys100, NumKeysExtra100)
		m = buildSet(keys)
		XtraKeys100 = xtra
		Set100 = m
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var j = rand.Int() % len(xtra)
		var key = xtra[j]
		_ = m.Set(key)
	}
}

func BenchmarkSetOne1M(b *testing.B) {
	log.Printf("BenchmarkSetOne1M: b.N=%d\n", b.N)
	var xtra = XtraKeys1M
	var m = Set1M
	if m == nil {
		var keys []key.Hash
		keys, xtra = buildKeysBench(NumKeys1M, NumKeysExtra1M)
		m = buildSet(keys)
		XtraKeys1M = xtra
		Set1M = m
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var j = rand.Int() % len(xtra)
		var key = xtra[j]
		_ = m.Set(key)
	}
}

func BenchmarkSetOne10M(b *testing.B) {
	log.Printf("BenchmarkSetOne10M: b.N=%d\n", b.N)
	var xtra = XtraKeys10M
	var m = Set10M
	if m == nil {
		var keys []key.Hash
		keys, xtra = buildKeysBench(NumKeys10M, NumKeysExtra10M)
		m = buildSet(keys)
		XtraKeys10M = xtra
		Set10M = m
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var j = rand.Int() % len(xtra)
		var key = xtra[j]
		_ = m.Set(key)
	}
}

func BenchmarkSetOne100M(b *testing.B) {
	log.Printf("BenchmarkSetOne100M: b.N=%d\n", b.N)
	var xtra = XtraKeys100M
	var m = Set100M
	if m == nil {
		var keys []key.Hash
		keys, xtra = buildKeysBench(NumKeys100M, NumKeysExtra100M)
		m = buildSet(keys)
		XtraKeys100M = xtra
		Set100M = m
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var j = rand.Int() % len(xtra)
		var key = xtra[j]
		_ = m.Set(key)
	}
}

func BenchmarkSetOne1MM(b *testing.B) {
	log.Printf("BenchmarkSetOne1MM: b.N=%d\n", b.N)
	var xtra = XtraKeys1MM
	var m = Set1MM
	if m == nil {
		var keys []key.Hash
		keys, xtra = buildKeysBench(NumKeys1MM, NumKeysExtra1MM)
		m = buildSet(keys)
		XtraKeys1MM = xtra
		Set1MM = m
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var j = rand.Int() % len(xtra)
		var key = xtra[j]
		_ = m.Set(key)
	}
}

func BenchmarkSetOne10MM(b *testing.B) {
	log.Printf("BenchmarkSetOne10MM: b.N=%d\n", b.N)
	var xtra = XtraKeys10MM
	var m = Set10MM
	if m == nil {
		var keys []key.Hash
		keys, xtra = buildKeysBench(NumKeys10MM, NumKeysExtra10MM)
		m = buildSet(keys)
		XtraKeys10MM = xtra
		Set10MM = m
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var j = rand.Int() % len(xtra)
		var key = xtra[j]
		_ = m.Set(key)
	}
}

var bulkSize = 10000
var bulkKeys []key.Hash
var bulkSet *set.Set

func BenchmarkBulkInsert(b *testing.B) {
	if bulkKeys == nil {
		log.Printf("BenchBulkDelete: building bulkKeys & bulkSet.")
		bulkKeys = buildKeys(bulkSize)
		//bulkSet = buildSet(bulkKeys)
	}
	b.ResetTimer()
	var s = set.New()
	for i := 0; i < b.N; i++ {
		s.BulkInsert(bulkKeys)
		s = set.New()
	}
}

func BenchmarkBulkInsert2(b *testing.B) {
	if bulkKeys == nil {
		log.Printf("BenchBulkDelete: building bulkKeys & bulkSet.")
		bulkKeys = buildKeys(bulkSize)
		//bulkSet = buildSet(bulkKeys)
	}
	b.ResetTimer()
	var s = set.New()
	for i := 0; i < b.N; i++ {
		//for j := 0; j < bulkSize; j++ {
		for j := 0; s.NumEntries() < bulkSize; j++ {
			s = s.Set(bulkKeys[j])
		}
		s = set.New()
	}
}

func BenchmarkBulkDelete(b *testing.B) {
	if bulkKeys == nil {
		log.Printf("BenchBulkDelete: building bulkKeys & bulkSet.")
		bulkKeys = buildKeys(bulkSize)
		bulkSet = buildSet(bulkKeys)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bulkSet.BulkDelete(bulkKeys)
	}
}

func BenchmarkBulkDelete2(b *testing.B) {
	if bulkKeys == nil {
		log.Printf("BenchBulkDelete2: building bulkKeys & bulkSet.")
		bulkKeys = buildKeys(bulkSize)
		bulkSet = buildSet(bulkKeys)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bulkSet.BulkDelete2(bulkKeys)
	}
}

func BenchmarkBulkDelete3(b *testing.B) {
	if bulkKeys == nil {
		log.Printf("BenchBulkDelete2: building bulkKeys & bulkSet.")
		bulkKeys = buildKeys(bulkSize)
		bulkSet = buildSet(bulkKeys)
	}
	b.ResetTimer()
	var s = bulkSet
	for i := 0; i < b.N; i++ {
		//for j := 0; j < len(bulkKeys); j++ {
		//	s, _ = s.Remove(bulkKeys[j])
		//}
		//for _, k := range bulkKeys {
		//	s = s.Unset(k)
		//}
		for j := 0; j < bulkSize; j++ {
			s = s.Unset(bulkKeys[j])
		}
		s = bulkSet
	}
}

var diffKeys []key.Hash
var setA, setB *set.Set
var origSetA, origSetB *set.Set
var tot = 100

var big, sml = tot * 6 / 10, tot * 4 / 10

//var big, sml = tot * 3 / 10, tot * 8 / 10
//var big, sml = tot * 8 / 10, tot * 7 / 10
//var big, sml = tot, tot * 9 / 10

//func init() {
//	diffKeys = buildKeys(tot)
//	setA = set.NewFromList(diffKeys[:big])
//	setB = set.NewFromList(diffKeys[sml:])
//}

func BenchmarkDifference(b *testing.B) {
	log.Println("b.N =", b.N)
	if diffKeys == nil {
		log.Printf("tot=%d; big=%d; sml=%d;\n", tot, big, sml)
		diffKeys = buildKeys(tot)
		setA = set.NewFromList(diffKeys[:big])
		setB = set.NewFromList(diffKeys[sml:])
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = setA.Difference(setB)
	}
}

//func BenchmarkDifference2(b *testing.B) {
//	if diffKeys == nil {
//		log.Printf("tot=%d; big=%d; sml=%d;\n", tot, big, sml)
//		diffKeys = buildKeys(tot)
//		setA = set.NewFromList(diffKeys[:big])
//		setB = set.NewFromList(diffKeys[sml:])
//	}
//	b.ResetTimer()
//
//	for i := 0; i < b.N; i++ {
//		_ = setA.Difference2(setB)
//	}
//}

//func BenchmarkDifference1(b *testing.B) {
//	if diffKeys == nil {
//		log.Printf("tot=%d; big=%d; sml=%d;\n", tot, big, sml)
//		diffKeys = buildKeys(tot)
//		setA = set.NewFromList(diffKeys[:big])
//		setB = set.NewFromList(diffKeys[sml:])
//	}
//	b.ResetTimer()
//
//	for i := 0; i < b.N; i++ {
//		_ = setA.Difference1(setB)
//	}
//}
//
//func BenchmarkDifference3(b *testing.B) {
//	if diffKeys == nil {
//		log.Printf("tot=%d; big=%d; sml=%d;\n", tot, big, sml)
//		diffKeys = buildKeys(tot)
//		setA = set.NewFromList(diffKeys[:big])
//		setB = set.NewFromList(diffKeys[sml:])
//	}
//	b.ResetTimer()
//
//	for i := 0; i < b.N; i++ {
//		_ = setA.Difference3(setB)
//	}
//}
