package sorted_map_test

import (
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
	for i := len(kvs) - 1; i > 0; i-- {
		var j = rand.Intn(i + 1)
		kvs[i], kvs[j] = kvs[j], kvs[i]
	}

	var xtra = kvs[len(kvs)-numKvsXtra:]
	kvs = kvs[:len(kvs)-numKvsXtra]

	return kvs, xtra
}

func buildMap(kvs []KeyVal) *sorted_map.Map {
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
	//log.Printf("BenchmarkPutOne10: called b.N=%d\n", b.N)
	var kvs, XtraKvs10 []KeyVal
	if SMap10 == nil || XtraKvs10 == nil {
		//log.Println("Generating Smap10 & XtraKvs10")
		kvs, XtraKvs10 = buildKvs(NumKvs10, NumKvsExtra10)
		SMap10 = buildMap(kvs)
	}
	var m = SMap10
	var xtraKvs = XtraKvs10
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var i = rand.Int() % NumKvsExtra10
		_ = m.Put(xtraKvs[i].Key, xtraKvs[i].Val)
	}
}

func BenchmarkPutOne100(b *testing.B) {
	//log.Printf("BenchmarkPutOne100: called b.N=%d\n", b.N)
	var kvs, XtraKvs100 []KeyVal
	if SMap100 == nil || XtraKvs100 == nil {
		//log.Println("Generating Smap100 & XtraKvs100")
		kvs, XtraKvs100 = buildKvs(NumKvs100, NumKvsExtra100)
		SMap100 = buildMap(kvs)
	}
	var m = SMap100
	var xtraKvs = XtraKvs100
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var i = rand.Int() % NumKvsExtra100
		_ = m.Put(xtraKvs[i].Key, xtraKvs[i].Val)
	}
}

func BenchmarkPutOne1M(b *testing.B) {
	//log.Printf("BenchmarkPutOne1M: called b.N=%d\n", b.N)
	var kvs, XtraKvs1M []KeyVal
	if SMap1M == nil || XtraKvs1M == nil {
		//log.Println("Generating Smap1M & XtraKvs1M")
		kvs, XtraKvs1M = buildKvs(NumKvs1M, NumKvsExtra1M)
		SMap1M = buildMap(kvs)
	}
	var m = SMap1M
	var xtraKvs = XtraKvs1M
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var i = rand.Int() % NumKvsExtra1M
		_ = m.Put(xtraKvs[i].Key, xtraKvs[i].Val)
	}
}

func BenchmarkPutOne10M(b *testing.B) {
	//log.Printf("BenchmarkPutOne10M: called b.N=%d\n", b.N)
	var kvs, XtraKvs10M []KeyVal
	if SMap10M == nil || XtraKvs10M == nil {
		//log.Println("Generating Smap10M & XtraKvs10M")
		kvs, XtraKvs10M = buildKvs(NumKvs10M, NumKvsExtra10M)
		SMap10M = buildMap(kvs)
	}
	var m = SMap10M
	var xtraKvs = XtraKvs10M
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var i = rand.Int() % NumKvsExtra10M
		_ = m.Put(xtraKvs[i].Key, xtraKvs[i].Val)
	}
}

func BenchmarkPutOne100M(b *testing.B) {
	//log.Printf("BenchmarkPutOne100M: called b.N=%d\n", b.N)
	var kvs, XtraKvs100M []KeyVal
	if SMap100M == nil || XtraKvs100M == nil {
		//log.Println("Generating Smap100M & XtraKvs100M")
		kvs, XtraKvs100M = buildKvs(NumKvs100M, NumKvsExtra100M)
		SMap100M = buildMap(kvs)
	}
	var m = SMap100M
	var xtraKvs = XtraKvs100M
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var i = rand.Int() % NumKvsExtra100M
		_ = m.Put(xtraKvs[i].Key, xtraKvs[i].Val)
	}
}

func BenchmarkPutOne1MM(b *testing.B) {
	//log.Printf("BenchmarkPutOne1MM: called b.N=%d\n", b.N)
	var kvs, XtraKvs1MM []KeyVal
	if SMap1MM == nil || XtraKvs1MM == nil {
		//log.Println("Generating Smap1MM & XtraKvs1MM")
		kvs, XtraKvs1MM = buildKvs(NumKvs1MM, NumKvsExtra1MM)
		SMap1MM = buildMap(kvs)
	}
	var m = SMap1MM
	var xtraKvs = XtraKvs1MM
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var i = rand.Int() % NumKvsExtra1MM
		_ = m.Put(xtraKvs[i].Key, xtraKvs[i].Val)
	}
}

func BenchmarkPutOne10MM(b *testing.B) {
	//log.Printf("BenchmarkPutOne10MM: called b.N=%d\n", b.N)
	var kvs, XtraKvs10MM []KeyVal
	if SMap10MM == nil || XtraKvs10MM == nil {
		//log.Println("Generating Smap10MM & XtraKvs10MM")
		kvs, XtraKvs10MM = buildKvs(NumKvs10MM, NumKvsExtra10MM)
		SMap10MM = buildMap(kvs)
	}
	var m = SMap10MM
	var xtraKvs = XtraKvs10MM
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		//START HERE
		var i = rand.Int() % NumKvsExtra10MM
		_ = m.Put(xtraKvs[i].Key, xtraKvs[i].Val)
	}
}
