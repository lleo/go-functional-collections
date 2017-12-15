package sorted_map_test

import (
	"log"
	"math/rand"
	"testing"

	"github.com/lleo/go-functional-collections/sorted_map"
)

type KeyVal struct {
	Key sorted_map.MapKey
	Val interface{}
}

func buildKvs(numMapKvs, numKvsXtra int) ([]KeyVal, []KeyVal) {
	var kvs = make([]KeyVal, numMapKvs+numKvsXtra)

	for i := 0; i < numMapKvs+numKvsXtra; i++ {
		kvs[i] = KeyVal{sorted_map.IntKey(i), i}
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

func buildMap(kvs []KeyVal) *sorted_map.Map {
	log.Printf("buildMap: len(kvs)=%d;\n", len(kvs))
	var m = sorted_map.New()
	for _, kv := range kvs {
		m = m.Put(kv.Key, kv.Val)
	}
	return m
}

const NumKvs10 = 1 * 10
const NumKvs100 = 1 * 100
const NumKvs1M = 1 * 1000
const NumKvs10M = 10 * 1000
const NumKvs100M = 100 * 1000
const NumKvs1MM = 1 * 1000 * 1000
const NumKvs10MM = 10 * 1000 * 1000

const NumKvsExtra10 = 2 * (NumKvs10 / 10)
const NumKvsExtra100 = 2 * (NumKvs100 / 10)
const NumKvsExtra1M = 2 * (NumKvs1M / 10)
const NumKvsExtra10M = 2 * (NumKvs10M / 10)
const NumKvsExtra100M = 20 * (NumKvs100M / 10)
const NumKvsExtra1MM = 20 * (NumKvs1MM / 10)
const NumKvsExtra10MM = 20 * (NumKvs10MM / 10)

var SMap10 *sorted_map.Map
var SMap100 *sorted_map.Map
var SMap1M *sorted_map.Map
var SMap10M *sorted_map.Map
var SMap100M *sorted_map.Map
var SMap1MM *sorted_map.Map
var SMap10MM *sorted_map.Map

var XtraKvs10 []KeyVal
var XtraKvs100 []KeyVal
var XtraKvs1M []KeyVal
var XtraKvs10M []KeyVal
var XtraKvs100M []KeyVal
var XtraKvs1MM []KeyVal
var XtraKvs10MM []KeyVal

func BenchmarkPutOne10(b *testing.B) {
	var xtra = XtraKvs10
	var m = SMap10
	if m == nil {
		var kvs []KeyVal
		kvs, xtra = buildKvs(NumKvs10, NumKvsExtra10)
		m = buildMap(kvs)
		XtraKvs10 = xtra
		SMap10 = m
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
	var m = SMap100
	if m == nil {
		var kvs []KeyVal
		kvs, xtra = buildKvs(NumKvs100, NumKvsExtra100)
		m = buildMap(kvs)
		XtraKvs100 = xtra
		SMap100 = m
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
	var m = SMap1M
	if m == nil {
		var kvs []KeyVal
		kvs, xtra = buildKvs(NumKvs1M, NumKvsExtra1M)
		m = buildMap(kvs)
		XtraKvs1M = xtra
		SMap1M = m
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
	var m = SMap10M
	if m == nil {
		var kvs []KeyVal
		kvs, xtra = buildKvs(NumKvs10M, NumKvsExtra10M)
		m = buildMap(kvs)
		XtraKvs10M = xtra
		SMap10M = m
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
	var m = SMap100M
	if m == nil {
		var kvs []KeyVal
		kvs, xtra = buildKvs(NumKvs100M, NumKvsExtra100M)
		m = buildMap(kvs)
		XtraKvs100M = xtra
		SMap100M = m
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
	var m = SMap1MM
	if m == nil {
		var kvs []KeyVal
		kvs, xtra = buildKvs(NumKvs1MM, NumKvsExtra1MM)
		m = buildMap(kvs)
		XtraKvs1MM = xtra
		SMap1MM = m
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
	var m = SMap10MM
	if m == nil {
		var kvs []KeyVal
		kvs, xtra = buildKvs(NumKvs10MM, NumKvsExtra10MM)
		m = buildMap(kvs)
		XtraKvs10MM = xtra
		SMap10MM = m
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var j = rand.Int() % len(xtra)
		var kv = xtra[j]
		_ = m.Put(kv.Key, kv.Val)
	}
}
