package fmap_test

import (
	"log"
	"math/rand"
	"testing"

	"github.com/lleo/go-functional-collections/fmap"
	"github.com/lleo/go-functional-collections/key"
)

func buildKvs2(numMapKvs, numKvsXtra int) ([]KeyVal, []KeyVal) {
	var kvs = make([]KeyVal, numMapKvs+numKvsXtra)

	var s = "a"
	for i := 0; i < numMapKvs+numKvsXtra; i++ {
		kvs[i] = KeyVal{key.Str(s), i}
		s = Inc(s)
	}

	// randomize kvs
	// https://en.wikipedia.org/wiki/Fisher–Yates_shuffle#The_modern_algorithm
	for i := len(kvs) - 1; i > 0; i-- {
		var j = rand.Intn(i + 1)
		kvs[i], kvs[j] = kvs[j], kvs[i]
	}

	var xtra = kvs[len(kvs)-numKvsXtra:]
	kvs = kvs[:len(kvs)-numKvsXtra]

	return kvs, xtra
}

const NumKvs10 = 1 * 10
const NumKvs100 = 1 * 100
const NumKvs1M = 1 * 1000
const NumKvs10M = 10 * 1000
const NumKvs100M = 100 * 1000
const NumKvs1MM = 1 * 1000 * 1000
const NumKvs10MM = 10 * 1000 * 1000
const NumKvs100MM = 100 * 1000 * 1000

const NumKvsExtra10 = 2 * (NumKvs10 / 10)
const NumKvsExtra100 = 2 * (NumKvs100 / 10)
const NumKvsExtra1M = 2 * (NumKvs1M / 10)
const NumKvsExtra10M = 2 * (NumKvs10M / 10)
const NumKvsExtra100M = 20 * (NumKvs100M / 10)
const NumKvsExtra1MM = 20 * (NumKvs1MM / 10)
const NumKvsExtra10MM = 20 * (NumKvs10MM / 10)
const NumKvsExtra100MM = 20 * (NumKvs100MM / 10)

var FMap10 *fmap.Map
var FMap100 *fmap.Map
var FMap1M *fmap.Map
var FMap10M *fmap.Map
var FMap100M *fmap.Map
var FMap1MM *fmap.Map
var FMap10MM *fmap.Map
var FMap100MM *fmap.Map

var XtraKvs10 []KeyVal
var XtraKvs100 []KeyVal
var XtraKvs1M []KeyVal
var XtraKvs10M []KeyVal
var XtraKvs100M []KeyVal
var XtraKvs1MM []KeyVal
var XtraKvs10MM []KeyVal
var XtraKvs100MM []KeyVal

func BenchmarkPutOne10(b *testing.B) {
	var xtra = XtraKvs10
	var m = FMap10
	if m == nil {
		var kvs []KeyVal
		kvs, xtra = buildKvs2(NumKvs10, NumKvsExtra10)
		m = buildMap(kvs)
		XtraKvs10 = xtra
		FMap10 = m
		log.Println("BenchmarkPutOne10: built XtraKvs10 & FMap10.")
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var j = rand.Int() % len(xtra)
		var kv = xtra[j]
		_ = m.Put(kv.Key, kv.Val)
	}
}

func BenchmarkPutOne100(b *testing.B) {
	log.Printf("BenchmarkPutOne100: b.N=%d\n", b.N)
	var xtra = XtraKvs100
	var m = FMap100
	if m == nil {
		var kvs []KeyVal
		kvs, xtra = buildKvs2(NumKvs100, NumKvsExtra100)
		m = buildMap(kvs)
		XtraKvs100 = xtra
		FMap100 = m
		log.Println("BenchmarkPutOne100: built XtraKvs100 & FMap100.")
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var j = rand.Int() % len(xtra)
		var kv = xtra[j]
		_ = m.Put(kv.Key, kv.Val)
	}
}

func BenchmarkPutOne1M(b *testing.B) {
	log.Printf("BenchmarkPutOne1M: b.N=%d\n", b.N)
	var xtra = XtraKvs1M
	var m = FMap1M
	if m == nil {
		var kvs []KeyVal
		kvs, xtra = buildKvs2(NumKvs1M, NumKvsExtra1M)
		m = buildMap(kvs)
		XtraKvs1M = xtra
		FMap1M = m
		log.Println("BenchmarkPutOne1M: built XtraKvs1M & FMap1M.")
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var j = rand.Int() % len(xtra)
		var kv = xtra[j]
		_ = m.Put(kv.Key, kv.Val)
	}
}

func BenchmarkPutOne10M(b *testing.B) {
	log.Printf("BenchmarkPutOne10M: b.N=%d\n", b.N)
	var xtra = XtraKvs10M
	var m = FMap10M
	if m == nil {
		var kvs []KeyVal
		kvs, xtra = buildKvs2(NumKvs10M, NumKvsExtra10M)
		m = buildMap(kvs)
		XtraKvs10M = xtra
		FMap10M = m
		log.Println("BenchmarkPutOne10M: built XtraKvs10M & FMap10M.")
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var j = rand.Int() % len(xtra)
		var kv = xtra[j]
		_ = m.Put(kv.Key, kv.Val)
	}
}

func BenchmarkPutOne100M(b *testing.B) {
	log.Printf("BenchmarkPutOne100M: b.N=%d\n", b.N)
	var xtra = XtraKvs100M
	var m = FMap100M
	if m == nil {
		var kvs []KeyVal
		kvs, xtra = buildKvs2(NumKvs100M, NumKvsExtra100M)
		m = buildMap(kvs)
		XtraKvs100M = xtra
		FMap100M = m
		log.Println("BenchmarkPutOne100M: built XtraKvs100M & FMap100M.")
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var j = rand.Int() % len(xtra)
		var kv = xtra[j]
		_ = m.Put(kv.Key, kv.Val)
	}
}

func BenchmarkPutOne1MM(b *testing.B) {
	log.Printf("BenchmarkPutOne1MM: b.N=%d\n", b.N)
	var xtra = XtraKvs1MM
	var m = FMap1MM
	if m == nil {
		var kvs []KeyVal
		kvs, xtra = buildKvs2(NumKvs1MM, NumKvsExtra1MM)
		m = buildMap(kvs)
		XtraKvs1MM = xtra
		FMap1MM = m
		log.Println("BenchmarkPutOne1MM: built XtraKvs1MM & FMap1MM.")
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var j = rand.Int() % len(xtra)
		var kv = xtra[j]
		_ = m.Put(kv.Key, kv.Val)
	}
}

func BenchmarkPutOne10MM(b *testing.B) {
	log.Printf("BenchmarkPutOne10MM: b.N=%d\n", b.N)
	var xtra = XtraKvs10MM
	var m = FMap10MM
	if m == nil {
		var kvs []KeyVal
		kvs, xtra = buildKvs2(NumKvs10MM, NumKvsExtra10MM)
		m = buildMap(kvs)
		XtraKvs10MM = xtra
		FMap10MM = m
		log.Println("BenchmarkPutOne10MM: built XtraKvs10MM & FMap10MM.")
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var j = rand.Int() % len(xtra)
		var kv = xtra[j]
		_ = m.Put(kv.Key, kv.Val)
	}
}

func xBenchmarkPutOne100MM(b *testing.B) {
	log.Printf("BenchmarkPutOne100MM: b.N=%d\n", b.N)
	var xtra = XtraKvs100MM
	var m = FMap100MM
	if m == nil {
		var kvs []KeyVal
		kvs, xtra = buildKvs2(NumKvs100MM, NumKvsExtra100MM)
		m = buildMap(kvs)
		XtraKvs100MM = xtra
		FMap100MM = m
		log.Println("BenchmarkPutOne100MM: built XtraKvs100MM & FMap100MM.")
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var j = rand.Int() % len(xtra)
		var kv = xtra[j]
		_ = m.Put(kv.Key, kv.Val)
	}
}

func BenchmarkIterNext10(b *testing.B) {
	log.Printf("BenchmarkIterNext10: b.N=%d\n", b.N)
	var xtra = XtraKvs10
	var m = FMap10
	if m == nil {
		var kvs []KeyVal
		kvs, xtra = buildKvs2(NumKvs10, NumKvsExtra10)
		m = buildMap(kvs)
		XtraKvs10 = xtra
		FMap10 = m
		log.Println("BenchmarkPutOne10: built XtraKvs10 & FMap10.")
	}
	b.ResetTimer()

	var it = m.Iter()
	for i := 0; i < b.N; i++ {
		var k, _ = it.Next()
		if k == nil {
			it = m.Iter()
			k, _ = it.Next()
		}
	}
}

func BenchmarkIterNext100(b *testing.B) {
	log.Printf("BenchmarkIterNext100: b.N=%d\n", b.N)
	var xtra = XtraKvs100
	var m = FMap100
	if m == nil {
		var kvs []KeyVal
		kvs, xtra = buildKvs2(NumKvs100, NumKvsExtra100)
		m = buildMap(kvs)
		XtraKvs100 = xtra
		FMap100 = m
		log.Println("BenchmarkPutOne100: built XtraKvs100 & FMap100.")
	}
	b.ResetTimer()

	var it = m.Iter()
	for i := 0; i < b.N; i++ {
		var k, _ = it.Next()
		if k == nil {
			it = m.Iter()
			k, _ = it.Next()
		}
	}
}

func BenchmarkIterNext1M(b *testing.B) {
	log.Printf("BenchmarkIterNext1M: b.N=%d\n", b.N)
	var xtra = XtraKvs1M
	var m = FMap1M
	if m == nil {
		var kvs []KeyVal
		kvs, xtra = buildKvs2(NumKvs1M, NumKvsExtra1M)
		m = buildMap(kvs)
		XtraKvs1M = xtra
		FMap1M = m
		log.Println("BenchmarkPutOne1M: built XtraKvs1M & FMap1M.")
	}
	b.ResetTimer()

	var it = m.Iter()
	for i := 0; i < b.N; i++ {
		var k, _ = it.Next()
		if k == nil {
			it = m.Iter()
			k, _ = it.Next()
		}
	}
}

func BenchmarkIterNext10M(b *testing.B) {
	log.Printf("BenchmarkIterNext10M: b.N=%d\n", b.N)
	var xtra = XtraKvs10M
	var m = FMap10M
	if m == nil {
		var kvs []KeyVal
		kvs, xtra = buildKvs2(NumKvs10M, NumKvsExtra10M)
		m = buildMap(kvs)
		XtraKvs10M = xtra
		FMap10M = m
		log.Println("BenchmarkPutOne10M: built XtraKvs10M & FMap10M.")
	}
	b.ResetTimer()

	var it = m.Iter()
	for i := 0; i < b.N; i++ {
		var k, _ = it.Next()
		if k == nil {
			it = m.Iter()
			k, _ = it.Next()
		}
	}
}

func BenchmarkIterNext100M(b *testing.B) {
	log.Printf("BenchmarkIterNext100M: b.N=%d\n", b.N)
	var xtra = XtraKvs100M
	var m = FMap100M
	if m == nil {
		var kvs []KeyVal
		kvs, xtra = buildKvs2(NumKvs100M, NumKvsExtra100M)
		m = buildMap(kvs)
		XtraKvs100M = xtra
		FMap100M = m
		log.Println("BenchmarkPutOne100M: built XtraKvs100M & FMap100M.")
	}
	b.ResetTimer()

	var it = m.Iter()
	for i := 0; i < b.N; i++ {
		var k, _ = it.Next()
		if k == nil {
			it = m.Iter()
			k, _ = it.Next()
		}
	}
}

func BenchmarkIterNext1MM(b *testing.B) {
	log.Printf("BenchmarkIterNext1MM: b.N=%d\n", b.N)
	var xtra = XtraKvs1MM
	var m = FMap1MM
	if m == nil {
		var kvs []KeyVal
		kvs, xtra = buildKvs2(NumKvs1MM, NumKvsExtra1MM)
		m = buildMap(kvs)
		XtraKvs1MM = xtra
		FMap1MM = m
		log.Println("BenchmarkPutOne1MM: built XtraKvs1MM & FMap1MM.")
	}
	b.ResetTimer()

	var it = m.Iter()
	for i := 0; i < b.N; i++ {
		var k, _ = it.Next()
		if k == nil {
			it = m.Iter()
			k, _ = it.Next()
		}
	}
}

func BenchmarkIterNext10MM(b *testing.B) {
	log.Printf("BenchmarkIterNext10MM: b.N=%d\n", b.N)
	var xtra = XtraKvs10MM
	var m = FMap10MM
	if m == nil {
		var kvs []KeyVal
		kvs, xtra = buildKvs2(NumKvs10MM, NumKvsExtra10MM)
		m = buildMap(kvs)
		XtraKvs10MM = xtra
		FMap10MM = m
		log.Println("BenchmarkPutOne10MM: built XtraKvs10MM & FMap10MM.")
	}
	b.ResetTimer()

	var it = m.Iter()
	for i := 0; i < b.N; i++ {
		var k, _ = it.Next()
		if k == nil {
			it = m.Iter()
			k, _ = it.Next()
		}
	}
}

func xBenchmarkIterNext100MM(b *testing.B) {
	log.Printf("BenchmarkIterNext100MM: b.N=%d\n", b.N)
	var xtra = XtraKvs100MM
	var m = FMap100MM
	if m == nil {
		var kvs []KeyVal
		kvs, xtra = buildKvs2(NumKvs100MM, NumKvsExtra100MM)
		m = buildMap(kvs)
		XtraKvs100MM = xtra
		FMap100MM = m
		log.Println("BenchmarkPutOne100MM: built XtraKvs100MM & FMap100MM.")
	}
	b.ResetTimer()

	var it = m.Iter()
	for i := 0; i < b.N; i++ {
		var k, _ = it.Next()
		if k == nil {
			it = m.Iter()
			k, _ = it.Next()
		}
	}
}

func BenchmarkBuildFunctional(b *testing.B) {
	log.Printf("BenchmarkBuildFunctional: b.N=%d\n", b.N)
	var kvs = buildKvs(b.N)
	var m = fmap.New()
	b.ResetTimer()
	for _, kv := range kvs {
		var k, v = kv.Key, kv.Val
		m = m.Put(k, v)
	}
}

func BenchmarkBuildNewFromList(b *testing.B) {
	log.Printf("BenchmarkBuildNewFromList: b.N=%d\n", b.N)
	var kvs = buildKvs(b.N)
	b.ResetTimer()
	_ = fmap.NewFromList(kvs)
}

func BenchmarkBuildBulkInsert(b *testing.B) {
	log.Printf("BenchmarkBuildBulkInsert: b.N=%d\n", b.N)
	var kvs = buildKvs(b.N)
	var m = fmap.New()
	b.ResetTimer()
	_ = m.BulkInsert(kvs, fmap.KeepOrigVal)
}
