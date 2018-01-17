package fmap_test

import (
	"log"
	"math/rand"
	"testing"

	"github.com/lleo/go-functional-collections/fmap"
	"github.com/lleo/go-functional-collections/hash"
)

func buildKvs2(numMapKvs, numKvsXtra int) ([]keyVal, []keyVal) {
	var kvs = make([]keyVal, numMapKvs+numKvsXtra)

	var s = "a"
	for i := 0; i < numMapKvs+numKvsXtra; i++ {
		kvs[i] = keyVal{hash.StringKey(s), i}
		s = Inc(s)
	}

	//randomize kvs
	//https://en.wikipedia.org/wiki/Fisher–Yates_shuffle#The_modern_algorithm
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

var XtraKvs10 []keyVal
var XtraKvs100 []keyVal
var XtraKvs1M []keyVal
var XtraKvs10M []keyVal
var XtraKvs100M []keyVal
var XtraKvs1MM []keyVal
var XtraKvs10MM []keyVal
var XtraKvs100MM []keyVal

func BenchmarkPutOne10(b *testing.B) {
	var xtra = XtraKvs10
	var m = FMap10
	if m == nil {
		var kvs []keyVal
		kvs, xtra = buildKvs2(NumKvs10, NumKvsExtra10)
		m = buildMap(kvs)
		XtraKvs10 = xtra
		FMap10 = m
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
		var kvs []keyVal
		kvs, xtra = buildKvs2(NumKvs100, NumKvsExtra100)
		m = buildMap(kvs)
		XtraKvs100 = xtra
		FMap100 = m
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
		var kvs []keyVal
		kvs, xtra = buildKvs2(NumKvs1M, NumKvsExtra1M)
		m = buildMap(kvs)
		XtraKvs1M = xtra
		FMap1M = m
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
		var kvs []keyVal
		kvs, xtra = buildKvs2(NumKvs10M, NumKvsExtra10M)
		m = buildMap(kvs)
		XtraKvs10M = xtra
		FMap10M = m
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
		var kvs []keyVal
		kvs, xtra = buildKvs2(NumKvs100M, NumKvsExtra100M)
		m = buildMap(kvs)
		XtraKvs100M = xtra
		FMap100M = m
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
		var kvs []keyVal
		kvs, xtra = buildKvs2(NumKvs1MM, NumKvsExtra1MM)
		m = buildMap(kvs)
		XtraKvs1MM = xtra
		FMap1MM = m
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
		var kvs []keyVal
		kvs, xtra = buildKvs2(NumKvs10MM, NumKvsExtra10MM)
		m = buildMap(kvs)
		XtraKvs10MM = xtra
		FMap10MM = m
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var j = rand.Int() % len(xtra)
		var kv = xtra[j]
		_ = m.Put(kv.Key, kv.Val)
	}
}

func BenchmarkPutOne100MM(b *testing.B) {
	log.Printf("BenchmarkPutOne100MM: b.N=%d\n", b.N)
	var xtra = XtraKvs100MM
	var m = FMap100MM
	if m == nil {
		var kvs []keyVal
		kvs, xtra = buildKvs2(NumKvs100MM, NumKvsExtra100MM)
		m = buildMap(kvs)
		XtraKvs100MM = xtra
		FMap100MM = m
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var j = rand.Int() % len(xtra)
		var kv = xtra[j]
		_ = m.Put(kv.Key, kv.Val)
	}
}
