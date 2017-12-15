package fmap_test

import (
	"log"
	"math/rand"
	"testing"

	"github.com/lleo/go-functional-collections/fmap"
)

type KeyVal struct {
	Key fmap.MapKey
	Val interface{}
}

func buildKvs(numMapKvs, numKvsXtra int) ([]KeyVal, []KeyVal) {
	var kvs = make([]KeyVal, numMapKvs+numKvsXtra)

	var s = "a"
	for i := 0; i < numMapKvs+numKvsXtra; i++ {
		kvs[i] = KeyVal{StringKey(s), i}
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

func buildMap(kvs []KeyVal) *fmap.Map {
	log.Printf("buildMap: len(kvs)=%d;\n", len(kvs))
	var m = fmap.New()
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

var FMap10 *fmap.Map
var FMap100 *fmap.Map
var FMap1M *fmap.Map
var FMap10M *fmap.Map
var FMap100M *fmap.Map
var FMap1MM *fmap.Map
var FMap10MM *fmap.Map

var XtraKvs10 []KeyVal
var XtraKvs100 []KeyVal
var XtraKvs1M []KeyVal
var XtraKvs10M []KeyVal
var XtraKvs100M []KeyVal
var XtraKvs1MM []KeyVal
var XtraKvs10MM []KeyVal

func benchPutOne(xtra []KeyVal, m *fmap.Map, b *testing.B) {
	for i := 0; i < b.N; i++ {
		var j = rand.Int() % len(xtra)
		var kv = xtra[j]
		_ = m.Put(kv.Key, kv.Val)
	}
}

func BenchmarkPutOne10(b *testing.B) {
	var xtra = XtraKvs10
	var m = FMap10
	if m == nil {
		var kvs []KeyVal
		kvs, xtra = buildKvs(NumKvs10, NumKvsExtra10)
		m = buildMap(kvs)
		XtraKvs10 = xtra
		FMap10 = m
	}
	b.ResetTimer()
	benchPutOne(xtra, m, b)
}

func BenchmarkPutOne100(b *testing.B) {
	log.Printf("BenchmarkPutOne100: b.N=%d\n", b.N)
	var xtra = XtraKvs100
	var m = FMap100
	if m == nil {
		var kvs []KeyVal
		kvs, xtra = buildKvs(NumKvs100, NumKvsExtra100)
		m = buildMap(kvs)
		XtraKvs100 = xtra
		FMap100 = m
	}
	b.ResetTimer()
	benchPutOne(xtra, m, b)
}

func BenchmarkPutOne1M(b *testing.B) {
	log.Printf("BenchmarkPutOne1M: b.N=%d\n", b.N)
	var xtra = XtraKvs1M
	var m = FMap1M
	if m == nil {
		var kvs []KeyVal
		kvs, xtra = buildKvs(NumKvs1M, NumKvsExtra1M)
		m = buildMap(kvs)
		XtraKvs1M = xtra
		FMap1M = m
	}
	b.ResetTimer()
	benchPutOne(xtra, m, b)
}

func BenchmarkPutOne10M(b *testing.B) {
	log.Printf("BenchmarkPutOne10M: b.N=%d\n", b.N)
	var xtra = XtraKvs10M
	var m = FMap10M
	if m == nil {
		var kvs []KeyVal
		kvs, xtra = buildKvs(NumKvs10M, NumKvsExtra10M)
		m = buildMap(kvs)
		XtraKvs10M = xtra
		FMap10M = m
	}
	b.ResetTimer()
	benchPutOne(xtra, m, b)
}

func BenchmarkPutOne100M(b *testing.B) {
	log.Printf("BenchmarkPutOne100M: b.N=%d\n", b.N)
	var xtra = XtraKvs100M
	var m = FMap100M
	if m == nil {
		var kvs []KeyVal
		kvs, xtra = buildKvs(NumKvs100M, NumKvsExtra100M)
		m = buildMap(kvs)
		XtraKvs100M = xtra
		FMap100M = m
	}
	b.ResetTimer()
	benchPutOne(xtra, m, b)
}

func BenchmarkPutOne1MM(b *testing.B) {
	log.Printf("BenchmarkPutOne1MM: b.N=%d\n", b.N)
	var xtra = XtraKvs1MM
	var m = FMap1MM
	if m == nil {
		var kvs []KeyVal
		kvs, xtra = buildKvs(NumKvs1MM, NumKvsExtra1MM)
		m = buildMap(kvs)
		XtraKvs1MM = xtra
		FMap1MM = m
	}
	b.ResetTimer()
	benchPutOne(xtra, m, b)
}

func BenchmarkPutOne10MM(b *testing.B) {
	log.Printf("BenchmarkPutOne10MM: b.N=%d\n", b.N)
	var xtra = XtraKvs10MM
	var m = FMap10MM
	if m == nil {
		var kvs []KeyVal
		kvs, xtra = buildKvs(NumKvs10MM, NumKvsExtra10MM)
		m = buildMap(kvs)
		XtraKvs10MM = xtra
		FMap10MM = m
	}
	b.ResetTimer()
	benchPutOne(xtra, m, b)
}

//func buildKeys(numKeys int) []StringKey {
//	var keys = make([]StringKey, numKeys)
//	var s = "a"
//	for i := 0; i < numKeys; i++ {
//		keys[i] = StringKey(s)
//		s = Inc(s)
//	}
//	return keys
//}
//
//var NumKey10MM = BaseMap10MMSize + (50 * 1000)
//var Keys10MM []StringKey
//
//var BaseMap10xSize = 10
//var BaseMap10x *fmap.Map
//
//func Benchmark_BaseMap10x_PutOne(b *testing.B) {
//	log.Printf("Benchmark_BaseMap10x_PutOne: b.N=%d\n", b.N)
//
//	var keys = Keys10MM
//	if keys == nil {
//		log.Printf("Making Keys10MM=%d...\n", NumKey10MM)
//
//		var startKeys = time.Now()
//		keys = buildKeys(NumKey10MM)
//		var tookKeys = time.Since(startKeys)
//
//		log.Printf("build keys took       -> %s\n", tookKeys)
//
//		Keys10MM = keys
//
//		log.Printf("Made Keys10MM; last key = Keys10MM[%d]=%s\n",
//			len(keys), keys[len(keys)-1])
//	}
//
//	var m = BaseMap10x
//	if m == nil {
//		log.Println("Making Map...")
//
//		var startMap = time.Now()
//		m = fmap.New()
//		for i := 0; i < BaseMap10xSize; i++ {
//			m = m.Put(keys[i], i)
//		}
//		var tookMap = time.Since(startMap)
//
//		log.Printf("build Map took        -> %s\n", tookMap)
//
//		BaseMap10x = m
//	}
//
//	b.ResetTimer()
//
//	var keyStart = BaseMap10xSize
//	keys = keys[keyStart : keyStart+2000]
//
//	for i := 0; i < b.N; i++ {
//		var j = i % 2000
//		_ = m.Put(keys[j], i)
//	}
//}
//
//var BaseMap10Size = 10
//var BaseMap10 *fmap.Map
//
//func Benchmark_BaseMap10_PutOne(b *testing.B) {
//	log.Printf("Benchmark_BaseMap10_PutOne: b.N=%d\n", b.N)
//
//	var keys = Keys10MM
//	if keys == nil {
//		log.Printf("Making Keys10MM=%d...\n", NumKey10MM)
//
//		var startKeys = time.Now()
//		keys = buildKeys(NumKey10MM)
//		var tookKeys = time.Since(startKeys)
//
//		log.Printf("build keys took       -> %s\n", tookKeys)
//
//		Keys10MM = keys
//
//		log.Printf("Made Keys10MM; last key = Keys10MM[%d]=%s\n",
//			len(keys), keys[len(keys)-1])
//	}
//
//	var m = BaseMap10
//	if m == nil {
//		log.Println("Making Map...")
//
//		var startMap = time.Now()
//		m = fmap.New()
//		for i := 0; i < BaseMap10Size; i++ {
//			m = m.Put(keys[i], i)
//		}
//		var tookMap = time.Since(startMap)
//
//		log.Printf("build Map took        -> %s\n", tookMap)
//
//		BaseMap10 = m
//	}
//
//	b.ResetTimer()
//
//	var keyStart = BaseMap10Size
//	keys = keys[keyStart : keyStart+2000]
//
//	for i := 0; i < b.N; i++ {
//		var j = i % 2000
//		_ = m.Put(keys[j], i)
//	}
//}
//
//var BaseMap100Size = 100
//var BaseMap100 *fmap.Map
//
//func Benchmark_BaseMap100_PutOne(b *testing.B) {
//	log.Printf("Benchmark_BaseMap100_PutOne: b.N=%d\n", b.N)
//
//	var keys = Keys10MM
//	if keys == nil {
//		log.Printf("Making Keys10MM=%d...\n", NumKey10MM)
//
//		var startKeys = time.Now()
//		keys = buildKeys(NumKey10MM)
//		var tookKeys = time.Since(startKeys)
//
//		log.Printf("build keys took       -> %s\n", tookKeys)
//
//		Keys10MM = keys
//
//		log.Printf("Made Keys10MM; last key = Keys10MM[%d]=%s\n",
//			len(keys), keys[len(keys)-1])
//	}
//
//	var m = BaseMap100
//	if m == nil {
//		log.Println("Making Map...")
//
//		var startMap = time.Now()
//		m = fmap.New()
//		for i := 0; i < BaseMap100Size; i++ {
//			m = m.Put(keys[i], i)
//		}
//		var tookMap = time.Since(startMap)
//
//		log.Printf("build Map took        -> %s\n", tookMap)
//
//		BaseMap100 = m
//	}
//
//	b.ResetTimer()
//
//	var keyStart = BaseMap100Size
//	keys = keys[keyStart : keyStart+2000]
//
//	for i := 0; i < b.N; i++ {
//		var j = i % 2000
//		_ = m.Put(keys[j], i)
//	}
//}
//
//var BaseMap1MSize = 1 * 1000
//var BaseMap1M *fmap.Map
//
//func Benchmark_BaseMap1M_PutOne(b *testing.B) {
//	log.Printf("Benchmark_BaseMap1M_PutOne: b.N=%d\n", b.N)
//
//	var keys = Keys10MM
//	if keys == nil {
//		log.Printf("Making Keys10MM=%d...\n", NumKey10MM)
//
//		var startKeys = time.Now()
//		keys = buildKeys(NumKey10MM)
//		var tookKeys = time.Since(startKeys)
//
//		log.Printf("build keys took       -> %s\n", tookKeys)
//
//		Keys10MM = keys
//
//		log.Printf("Made Keys10MM; last key = Keys10MM[%d]=%s\n",
//			len(keys), keys[len(keys)-1])
//	}
//
//	var m = BaseMap1M
//	if m == nil {
//		log.Println("Making Map...")
//
//		var startMap = time.Now()
//		m = fmap.New()
//		for i := 0; i < BaseMap1MSize; i++ {
//			m = m.Put(keys[i], i)
//		}
//		var tookMap = time.Since(startMap)
//
//		log.Printf("build Map took        -> %s\n", tookMap)
//
//		BaseMap1M = m
//	}
//
//	b.ResetTimer()
//
//	var keyStart = BaseMap1MSize
//	keys = keys[keyStart : keyStart+2000]
//
//	for i := 0; i < b.N; i++ {
//		var j = i % 2000
//		_ = m.Put(keys[j], i)
//	}
//}
//
//var BaseMap10MSize = 10 * 1000
//var BaseMap10M *fmap.Map
//
//func Benchmark_BaseMap10M_PutOne(b *testing.B) {
//	log.Printf("Benchmark_BaseMap10M_PutOne: b.N=%d\n", b.N)
//
//	var keys = Keys10MM
//	if keys == nil {
//		log.Printf("Making Keys10MM=%d...\n", NumKey10MM)
//
//		var startKeys = time.Now()
//		keys = buildKeys(NumKey10MM)
//		var tookKeys = time.Since(startKeys)
//
//		log.Printf("build keys took       -> %s\n", tookKeys)
//
//		Keys10MM = keys
//
//		log.Printf("Made Keys10MM; last key = Keys10MM[%d]=%s\n",
//			len(keys), keys[len(keys)-1])
//	}
//
//	var m = BaseMap10M
//	if m == nil {
//		log.Println("Making Map...")
//
//		var startMap = time.Now()
//		m = fmap.New()
//		for i := 0; i < BaseMap10MSize; i++ {
//			m = m.Put(keys[i], i)
//		}
//		var tookMap = time.Since(startMap)
//
//		log.Printf("build Map took        -> %s\n", tookMap)
//
//		BaseMap10M = m
//	}
//
//	b.ResetTimer()
//
//	var keyStart = BaseMap10MSize
//	keys = keys[keyStart : keyStart+2000]
//
//	for i := 0; i < b.N; i++ {
//		var j = i % 2000
//		_ = m.Put(keys[j], i)
//	}
//}
//
//var BaseMap100MSize = 100 * 1000
//var BaseMap100M *fmap.Map
//
//func Benchmark_BaseMap100M_PutOne(b *testing.B) {
//	log.Printf("Benchmark_BaseMap100M_PutOne: b.N=%d\n", b.N)
//
//	var keys = Keys10MM
//	if keys == nil {
//		log.Printf("Making Keys10MM=%d...\n", NumKey10MM)
//
//		var startKeys = time.Now()
//		keys = buildKeys(NumKey10MM)
//		var tookKeys = time.Since(startKeys)
//
//		log.Printf("build keys took       -> %s\n", tookKeys)
//
//		Keys10MM = keys
//
//		log.Printf("Made Keys10MM; last key = Keys10MM[%d]=%s\n",
//			len(keys), keys[len(keys)-1])
//	}
//
//	var m = BaseMap100M
//	if m == nil {
//		log.Println("Making Map...")
//
//		var startMap = time.Now()
//		m = fmap.New()
//		for i := 0; i < BaseMap100MSize; i++ {
//			m = m.Put(keys[i], i)
//		}
//		var tookMap = time.Since(startMap)
//
//		log.Printf("build Map took        -> %s\n", tookMap)
//
//		BaseMap100M = m
//	}
//
//	b.ResetTimer()
//
//	var keyStart = BaseMap100MSize
//	keys = keys[keyStart : keyStart+2000]
//
//	for i := 0; i < b.N; i++ {
//		var j = i % 2000
//		_ = m.Put(keys[j], i)
//	}
//}
//
//var BaseMap1MMSize = 1 * 1000 * 1000
//var BaseMap1MM *fmap.Map
//
//func Benchmark_BaseMap1MM_PutOne(b *testing.B) {
//	log.Printf("Benchmark_BaseMap1MM_PutOne: b.N=%d\n", b.N)
//
//	var keys = Keys10MM
//	if keys == nil {
//		log.Printf("Making Keys10MM=%d...\n",
//			NumKey10MM)
//
//		var startKeys = time.Now()
//		keys = buildKeys(NumKey10MM)
//		var tookKeys = time.Since(startKeys)
//
//		log.Printf("build keys took       -> %s\n", tookKeys)
//
//		Keys10MM = keys
//
//		log.Printf("Made Keys10MM; last key = Keys10MM[%d]=%s\n",
//			len(keys), keys[len(keys)-1])
//	}
//
//	var m = BaseMap1MM
//	if m == nil {
//		log.Println("Making Map...")
//
//		var startMap = time.Now()
//		m = fmap.New()
//		for i := 0; i < BaseMap1MMSize; i++ {
//			m = m.Put(keys[i], i)
//		}
//		var tookMap = time.Since(startMap)
//
//		log.Printf("build Map took        -> %s\n", tookMap)
//
//		BaseMap1MM = m
//	}
//
//	b.ResetTimer()
//
//	var keyStart = BaseMap1MMSize
//	keys = keys[keyStart : keyStart+2000]
//
//	for i := 0; i < b.N; i++ {
//		var j = i % 2000
//		_ = m.Put(keys[j], i)
//	}
//}
//
//var BaseMap10MMSize = 10 * 1000 * 1000
//var BaseMap10MM *fmap.Map
//
//func Benchmark_BaseMap10MM_PutOne(b *testing.B) {
//	log.Printf("Benchmark_BaseMap10MM_PutOne: b.N=%d\n", b.N)
//
//	var keys = Keys10MM
//	if keys == nil {
//		log.Printf("Making Keys10MM=%d...\n", NumKey10MM)
//
//		var startKeys = time.Now()
//		keys = buildKeys(NumKey10MM)
//		var tookKeys = time.Since(startKeys)
//
//		log.Printf("build keys took       -> %s\n", tookKeys)
//
//		Keys10MM = keys
//
//		log.Printf("Made Keys10MM; last key = Keys10MM[%d]=%s\n",
//			len(keys), keys[len(keys)-1])
//	}
//
//	var m = BaseMap10MM
//	if m == nil {
//		log.Println("Making Map...")
//		var startMap = time.Now()
//		m = fmap.New()
//		for i := 0; i < BaseMap10MMSize; i++ {
//			m = m.Put(keys[i], i)
//		}
//		var tookMap = time.Since(startMap)
//
//		log.Printf("build Map took        -> %s\n", tookMap)
//
//		BaseMap10MM = m
//	}
//
//	b.ResetTimer()
//
//	var keyStart = BaseMap10MMSize
//	keys = keys[keyStart : keyStart+2000]
//
//	for i := 0; i < b.N; i++ {
//		var j = i % 2000
//		_ = m.Put(keys[j], i)
//	}
//}
//
