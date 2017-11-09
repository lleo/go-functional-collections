package fmap_test

import (
	"log"
	"testing"
	"time"

	"github.com/lleo/go-functional-collections/fmap"
)

func buildKeys(numKeys int) []StringKey {
	var keys = make([]StringKey, numKeys)
	var s = "a"
	for i := 0; i < numKeys; i++ {
		keys[i] = StringKey(s)
		s = Inc(s)
	}
	return keys
}

var NumKey10MM = BaseMap10MMSize + (50 * 1000)
var Keys10MM []StringKey

var BaseMap10xSize = 10
var BaseMap10x *fmap.Map

func Benchmark_BaseMap10x_PutOne(b *testing.B) {
	log.Printf("Benchmark_BaseMap10x_PutOne: b.N=%d\n", b.N)

	var keys = Keys10MM
	if keys == nil {
		log.Printf("Making Keys10MM=%d...\n", NumKey10MM)

		var startKeys = time.Now()
		keys = buildKeys(NumKey10MM)
		var tookKeys = time.Since(startKeys)

		log.Printf("build keys took       -> %s\n", tookKeys)

		Keys10MM = keys

		log.Printf("Made Keys10MM; last key = Keys10MM[%d]=%s\n",
			len(keys), keys[len(keys)-1])
	}

	var m = BaseMap10x
	if m == nil {
		log.Println("Making Map...")

		var startMap = time.Now()
		m = fmap.New()
		for i := 0; i < BaseMap10xSize; i++ {
			m = m.Put(keys[i], i)
		}
		var tookMap = time.Since(startMap)

		log.Printf("build Map took        -> %s\n", tookMap)

		BaseMap10x = m
	}

	b.ResetTimer()

	var keyStart = BaseMap10xSize
	keys = keys[keyStart : keyStart+2000]

	for i := 0; i < b.N; i++ {
		var j = i % 2000
		_ = m.Put(keys[j], i)
	}
}

var BaseMap10Size = 10
var BaseMap10 *fmap.Map

func Benchmark_BaseMap10_PutOne(b *testing.B) {
	log.Printf("Benchmark_BaseMap10_PutOne: b.N=%d\n", b.N)

	var keys = Keys10MM
	if keys == nil {
		log.Printf("Making Keys10MM=%d...\n", NumKey10MM)

		var startKeys = time.Now()
		keys = buildKeys(NumKey10MM)
		var tookKeys = time.Since(startKeys)

		log.Printf("build keys took       -> %s\n", tookKeys)

		Keys10MM = keys

		log.Printf("Made Keys10MM; last key = Keys10MM[%d]=%s\n",
			len(keys), keys[len(keys)-1])
	}

	var m = BaseMap10
	if m == nil {
		log.Println("Making Map...")

		var startMap = time.Now()
		m = fmap.New()
		for i := 0; i < BaseMap10Size; i++ {
			m = m.Put(keys[i], i)
		}
		var tookMap = time.Since(startMap)

		log.Printf("build Map took        -> %s\n", tookMap)

		BaseMap10 = m
	}

	b.ResetTimer()

	var keyStart = BaseMap10Size
	keys = keys[keyStart : keyStart+2000]

	for i := 0; i < b.N; i++ {
		var j = i % 2000
		_ = m.Put(keys[j], i)
	}
}

var BaseMap100Size = 100
var BaseMap100 *fmap.Map

func Benchmark_BaseMap100_PutOne(b *testing.B) {
	log.Printf("Benchmark_BaseMap100_PutOne: b.N=%d\n", b.N)

	var keys = Keys10MM
	if keys == nil {
		log.Printf("Making Keys10MM=%d...\n", NumKey10MM)

		var startKeys = time.Now()
		keys = buildKeys(NumKey10MM)
		var tookKeys = time.Since(startKeys)

		log.Printf("build keys took       -> %s\n", tookKeys)

		Keys10MM = keys

		log.Printf("Made Keys10MM; last key = Keys10MM[%d]=%s\n",
			len(keys), keys[len(keys)-1])
	}

	var m = BaseMap100
	if m == nil {
		log.Println("Making Map...")

		var startMap = time.Now()
		m = fmap.New()
		for i := 0; i < BaseMap100Size; i++ {
			m = m.Put(keys[i], i)
		}
		var tookMap = time.Since(startMap)

		log.Printf("build Map took        -> %s\n", tookMap)

		BaseMap100 = m
	}

	b.ResetTimer()

	var keyStart = BaseMap100Size
	keys = keys[keyStart : keyStart+2000]

	for i := 0; i < b.N; i++ {
		var j = i % 2000
		_ = m.Put(keys[j], i)
	}
}

var BaseMap1MSize = 1 * 1000
var BaseMap1M *fmap.Map

func Benchmark_BaseMap1M_PutOne(b *testing.B) {
	log.Printf("Benchmark_BaseMap1M_PutOne: b.N=%d\n", b.N)

	var keys = Keys10MM
	if keys == nil {
		log.Printf("Making Keys10MM=%d...\n", NumKey10MM)

		var startKeys = time.Now()
		keys = buildKeys(NumKey10MM)
		var tookKeys = time.Since(startKeys)

		log.Printf("build keys took       -> %s\n", tookKeys)

		Keys10MM = keys

		log.Printf("Made Keys10MM; last key = Keys10MM[%d]=%s\n",
			len(keys), keys[len(keys)-1])
	}

	var m = BaseMap1M
	if m == nil {
		log.Println("Making Map...")

		var startMap = time.Now()
		m = fmap.New()
		for i := 0; i < BaseMap1MSize; i++ {
			m = m.Put(keys[i], i)
		}
		var tookMap = time.Since(startMap)

		log.Printf("build Map took        -> %s\n", tookMap)

		BaseMap1M = m
	}

	b.ResetTimer()

	var keyStart = BaseMap1MSize
	keys = keys[keyStart : keyStart+2000]

	for i := 0; i < b.N; i++ {
		var j = i % 2000
		_ = m.Put(keys[j], i)
	}
}

var BaseMap10MSize = 10 * 1000
var BaseMap10M *fmap.Map

func Benchmark_BaseMap10M_PutOne(b *testing.B) {
	log.Printf("Benchmark_BaseMap10M_PutOne: b.N=%d\n", b.N)

	var keys = Keys10MM
	if keys == nil {
		log.Printf("Making Keys10MM=%d...\n", NumKey10MM)

		var startKeys = time.Now()
		keys = buildKeys(NumKey10MM)
		var tookKeys = time.Since(startKeys)

		log.Printf("build keys took       -> %s\n", tookKeys)

		Keys10MM = keys

		log.Printf("Made Keys10MM; last key = Keys10MM[%d]=%s\n",
			len(keys), keys[len(keys)-1])
	}

	var m = BaseMap10M
	if m == nil {
		log.Println("Making Map...")

		var startMap = time.Now()
		m = fmap.New()
		for i := 0; i < BaseMap10MSize; i++ {
			m = m.Put(keys[i], i)
		}
		var tookMap = time.Since(startMap)

		log.Printf("build Map took        -> %s\n", tookMap)

		BaseMap10M = m
	}

	b.ResetTimer()

	var keyStart = BaseMap10MSize
	keys = keys[keyStart : keyStart+2000]

	for i := 0; i < b.N; i++ {
		var j = i % 2000
		_ = m.Put(keys[j], i)
	}
}

var BaseMap100MSize = 100 * 1000
var BaseMap100M *fmap.Map

func Benchmark_BaseMap100M_PutOne(b *testing.B) {
	log.Printf("Benchmark_BaseMap100M_PutOne: b.N=%d\n", b.N)

	var keys = Keys10MM
	if keys == nil {
		log.Printf("Making Keys10MM=%d...\n", NumKey10MM)

		var startKeys = time.Now()
		keys = buildKeys(NumKey10MM)
		var tookKeys = time.Since(startKeys)

		log.Printf("build keys took       -> %s\n", tookKeys)

		Keys10MM = keys

		log.Printf("Made Keys10MM; last key = Keys10MM[%d]=%s\n",
			len(keys), keys[len(keys)-1])
	}

	var m = BaseMap100M
	if m == nil {
		log.Println("Making Map...")

		var startMap = time.Now()
		m = fmap.New()
		for i := 0; i < BaseMap100MSize; i++ {
			m = m.Put(keys[i], i)
		}
		var tookMap = time.Since(startMap)

		log.Printf("build Map took        -> %s\n", tookMap)

		BaseMap100M = m
	}

	b.ResetTimer()

	var keyStart = BaseMap100MSize
	keys = keys[keyStart : keyStart+2000]

	for i := 0; i < b.N; i++ {
		var j = i % 2000
		_ = m.Put(keys[j], i)
	}
}

var BaseMap1MMSize = 1 * 1000 * 1000
var BaseMap1MM *fmap.Map

func Benchmark_BaseMap1MM_PutOne(b *testing.B) {
	log.Printf("Benchmark_BaseMap1MM_PutOne: b.N=%d\n", b.N)

	var keys = Keys10MM
	if keys == nil {
		log.Printf("Making Keys10MM=%d...\n",
			NumKey10MM)

		var startKeys = time.Now()
		keys = buildKeys(NumKey10MM)
		var tookKeys = time.Since(startKeys)

		log.Printf("build keys took       -> %s\n", tookKeys)

		Keys10MM = keys

		log.Printf("Made Keys10MM; last key = Keys10MM[%d]=%s\n",
			len(keys), keys[len(keys)-1])
	}

	var m = BaseMap1MM
	if m == nil {
		log.Println("Making Map...")

		var startMap = time.Now()
		m = fmap.New()
		for i := 0; i < BaseMap1MMSize; i++ {
			m = m.Put(keys[i], i)
		}
		var tookMap = time.Since(startMap)

		log.Printf("build Map took        -> %s\n", tookMap)

		BaseMap1MM = m
	}

	b.ResetTimer()

	var keyStart = BaseMap1MMSize
	keys = keys[keyStart : keyStart+2000]

	for i := 0; i < b.N; i++ {
		var j = i % 2000
		_ = m.Put(keys[j], i)
	}
}

var BaseMap10MMSize = 10 * 1000 * 1000
var BaseMap10MM *fmap.Map

func Benchmark_BaseMap10MM_PutOne(b *testing.B) {
	log.Printf("Benchmark_BaseMap10MM_PutOne: b.N=%d\n", b.N)

	var keys = Keys10MM
	if keys == nil {
		log.Printf("Making Keys10MM=%d...\n", NumKey10MM)

		var startKeys = time.Now()
		keys = buildKeys(NumKey10MM)
		var tookKeys = time.Since(startKeys)

		log.Printf("build keys took       -> %s\n", tookKeys)

		Keys10MM = keys

		log.Printf("Made Keys10MM; last key = Keys10MM[%d]=%s\n",
			len(keys), keys[len(keys)-1])
	}

	var m = BaseMap10MM
	if m == nil {
		log.Println("Making Map...")
		var startMap = time.Now()
		m = fmap.New()
		for i := 0; i < BaseMap10MMSize; i++ {
			m = m.Put(keys[i], i)
		}
		var tookMap = time.Since(startMap)

		log.Printf("build Map took        -> %s\n", tookMap)

		BaseMap10MM = m
	}

	b.ResetTimer()

	var keyStart = BaseMap10MMSize
	keys = keys[keyStart : keyStart+2000]

	for i := 0; i < b.N; i++ {
		var j = i % 2000
		_ = m.Put(keys[j], i)
	}
}
