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

func genIntKeyVals(n int) []KeyVal {
	var kvs = make([]KeyVal, n)

	for i := 0; i < n; i++ {
		var x = (i + 1) * 10
		var k = sorted_map.IntKey(x)
		var v = x
		kvs[i] = KeyVal{k, v}
	}

	return kvs
}

func randomizeKeyVals(kvs []KeyVal) []KeyVal {
	var randKvs = make([]KeyVal, len(kvs))
	copy(randKvs, kvs)
	//var randKvs = kvs

	//From: https://en.wikipedia.org/wiki/Fisher%E2%80%93Yates_shuffle#The_modern_algorithm
	for i := len(randKvs) - 1; i > 0; i-- {
		var j = rand.Intn(i + 1)
		randKvs[i], randKvs[j] = randKvs[j], randKvs[i]
	}

	return randKvs
}

func extractRandKvs(nXtract int, kvs []KeyVal) ([]KeyVal, []KeyVal) {
	if nXtract >= len(kvs) {
		log.Panicf("nXtract,%d >= len(kvs),%d", nXtract, len(kvs))
	}
	var xtractKvs = make([]KeyVal, 0, nXtract)

	for len(xtractKvs) < nXtract {
		var i = rand.Int() % len(kvs)
		xtractKvs = append(xtractKvs, kvs[i])
		kvs = kvs[:i+copy(kvs[i:], kvs[i+1:])] //remove kvs[i] from kvs
	}

	return xtractKvs, kvs
}

func genMap(kvs []KeyVal) *sorted_map.Map {
	var m = sorted_map.New()
	for _, kv := range kvs {
		m = m.Put(kv.Key, kv.Val)
	}
	return m
}

const NumKvsXtra = 1000
const NumKvs1M = 1 * 1000
const NumKvs10M = 10 * 1000
const NumKvs100M = 100 * 1000
const NumKvs1MM = 1 * 1000 * 1000
const NumKvs10MM = 10 * 1000 * 1000

var SMap1M *sorted_map.Map
var SMap10M *sorted_map.Map
var SMap100M *sorted_map.Map
var SMap1MM *sorted_map.Map
var SMap10MM *sorted_map.Map

var XtraKvs1M []KeyVal
var XtraKvs10M []KeyVal
var XtraKvs100M []KeyVal
var XtraKvs1MM []KeyVal
var XtraKvs10MM []KeyVal

func BenchmarkPutOne1M(b *testing.B) {
	log.Printf("BenchmarkPutOne1M: called b.N=%d\n", b.N)
	if SMap1M == nil || XtraKvs1M == nil {
		log.Println("Generating Smap1M & XtraKvs1M")
		var kvs = genIntKeyVals(NumKvs1M + NumKvsXtra)
		XtraKvs1M, kvs = extractRandKvs(NumKvsXtra, kvs)
		SMap1M = genMap(kvs)
	}
	var m = SMap1M
	var xtraKvs = XtraKvs1M
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var i = rand.Int() % NumKvsXtra
		_ = m.Put(xtraKvs[i].Key, xtraKvs[i].Val)
	}
}

func BenchmarkPutOne10M(b *testing.B) {
	log.Printf("BenchmarkPutOne10M: called b.N=%d\n", b.N)
	if SMap10M == nil || XtraKvs10M == nil {
		log.Println("Generating Smap10M & XtraKvs10M")
		var kvs = genIntKeyVals(NumKvs10M + NumKvsXtra)
		XtraKvs10M, kvs = extractRandKvs(NumKvsXtra, kvs)
		SMap10M = genMap(kvs)
	}
	var m = SMap10M
	var xtraKvs = XtraKvs10M
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var i = rand.Int() % NumKvsXtra
		_ = m.Put(xtraKvs[i].Key, xtraKvs[i].Val)
	}
}

func BenchmarkPutOne100M(b *testing.B) {
	log.Printf("BenchmarkPutOne100M: called b.N=%d\n", b.N)
	if SMap100M == nil || XtraKvs100M == nil {
		log.Println("Generating Smap100M & XtraKvs100M")
		var kvs = genIntKeyVals(NumKvs100M + NumKvsXtra)
		XtraKvs100M, kvs = extractRandKvs(NumKvsXtra, kvs)
		SMap100M = genMap(kvs)
	}
	var m = SMap100M
	var xtraKvs = XtraKvs100M
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var i = rand.Int() % NumKvsXtra
		_ = m.Put(xtraKvs[i].Key, xtraKvs[i].Val)
	}
}

func BenchmarkPutOne1MM(b *testing.B) {
	log.Printf("BenchmarkPutOne1MM: called b.N=%d\n", b.N)
	if SMap1MM == nil || XtraKvs1MM == nil {
		log.Println("Generating Smap1MM & XtraKvs1MM")
		var kvs = genIntKeyVals(NumKvs1MM + NumKvsXtra)
		XtraKvs1MM, kvs = extractRandKvs(NumKvsXtra, kvs)
		SMap1MM = genMap(kvs)
	}
	var m = SMap1MM
	var xtraKvs = XtraKvs1MM
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var i = rand.Int() % NumKvsXtra
		_ = m.Put(xtraKvs[i].Key, xtraKvs[i].Val)
	}
}

func BenchmarkPutOne10MM(b *testing.B) {
	log.Printf("BenchmarkPutOne10MM: called b.N=%d\n", b.N)
	if SMap10MM == nil || XtraKvs10MM == nil {
		log.Println("Generating Smap10MM & XtraKvs10MM")
		var kvs = genIntKeyVals(NumKvs10MM + NumKvsXtra)
		XtraKvs10MM, kvs = extractRandKvs(NumKvsXtra, kvs)
		SMap10MM = genMap(kvs)
	}
	var m = SMap10MM
	var xtraKvs = XtraKvs10MM
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var i = rand.Int() % NumKvsXtra
		_ = m.Put(xtraKvs[i].Key, xtraKvs[i].Val)
	}
}
