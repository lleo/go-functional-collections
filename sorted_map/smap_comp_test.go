package sorted_map_test

import (
	"log"
	"testing"

	"github.com/lleo/go-functional-collections/sorted_map"
)

func TestCompBuildMap(t *testing.T) {
	var numKeys = 1024 //tested upto 10240
	var inOrderKvs = genIntKeyVals(numKeys)
	//rand.Seed(int64(time.Now().Nanosecond()))
	var kvs = randomizeKeyVals(inOrderKvs)

	//log.Println("kvs[:10] =", kvs[:10])
	var m = sorted_map.New()

	for _, kv := range kvs {
		var k = kv.Key
		var v = kv.Val

		m = m.Put(k, v)

		var err = m.Valid()
		if err != nil {
			t.Fatalf("Invalid Tree. err=%s\n", err)
		}
	}

	//log.Printf("Map m =\n%s", m.TreeString())
	//log.Printf("Map m =\n%s", m.String())

	var i int
	var fn = func(k0 sorted_map.MapKey, v0 interface{}) bool {
		var k1 = inOrderKvs[i].Key
		var v1 = inOrderKvs[i].Val
		if k0.Less(k1) || k1.Less(k0) { //k0 != k1
			t.Fatalf("InOrder keys: i=%d; found k0=%s not the expected k1=%s\n",
				i, k0, k1)
		}
		if v0 != v1 {
			t.Fatalf("InOrder vals: i=%d; found v0=%d not the expected v1=%d\n",
				i, v0, v1)
		}
		i++
		return true
	}
	m.Range(fn)
}

func TestCompDestroyMap(t *testing.T) {
	var numKeys = 1024 //tested upto 10240
	var kvs = genIntKeyVals(numKeys)
	//rand.Seed(int64(time.Now().Nanosecond()))
	var buildKvs = randomizeKeyVals(kvs)
	var destroyKvs = randomizeKeyVals(kvs)

	//log.Printf("buildKvs = %v\n", buildKvs)

	var m = sorted_map.New()
	for i, kv := range buildKvs {
		m = m.Put(kv.Key, kv.Val)
		var err = m.Valid()
		if err != nil {
			t.Fatalf("INVALID TREE Store: i=%d; kv.Key=%s; err=%s\n",
				i, kv.Key, err)
		}
		//log.Printf("Map m =\n%s", m.TreeString())
	}

	//log.Printf("AFTER ALL MAP BUILDING: Map m=\n%s", m.TreeString())
	//log.Printf("m = %s\n", m.String())

	//log.Printf("destroyKvs = %v\n", destroyKvs)

	var shouldHaveKvs = make([]KeyVal, len(destroyKvs))
	copy(shouldHaveKvs, destroyKvs)

	for i, kv := range destroyKvs {
		//var origM = m
		//var dupOrigM = m.Dup()

		//log.Printf("*******Removing kv.Key=%s; i=%d \n", kv.Key, i)

		//log.Printf("BEFORE REMOVE: Map m=\n%s", m.TreeString())

		var val interface{}
		var found bool
		m, val, found = m.Remove(kv.Key)

		//log.Printf("AFTER REMOVE Map m=\n%s", m.TreeString())

		if !found {
			t.Fatalf("Remove: i=%d; kv.Key=%s not found!\n", i, kv.Key)
		}

		if val != kv.Val {
			t.Fatalf("Remove: i=%d; kv.Key=%s; val=%d != expected kv.Val=%d\n",
				i, kv.Key, val, kv.Val)
		}

		var err = m.Valid()
		if err != nil {
			log.Printf("Map m=\n%s", m.TreeString())
			t.Fatalf("INVALID TREE Remove: i=%d; kv.Key=%s; err=%s\n",
				i, kv.Key, err)
		}

		shouldHaveKvs = shouldHaveKvs[1:] //take the first elt off ala range
		for _, kv0 := range shouldHaveKvs {
			var val, found = m.Load(kv0.Key)
			if !found {
				t.Fatalf("Remove: i=%d; for key=%s: "+
					"failed to find shouldHave key=%s", i, kv.Key, kv0.Key)
			}
			if val != kv0.Val {
				t.Fatalf("Remove: found val,%v != expected val,%v",
					val, kv0.Val)
			}
		}
		//if !origM.Equiv(dupOrigM) {
		//	t.Fatal("the original map was modified during Remove(%s).", kv.Key)
		//}
	}

	if m.NumEntries() != 0 {
		t.Fatal("m.NumEntries() != 0")
	}
}

func TestCompRangeForwAll(t *testing.T) {
	var numKeys = 1024 //tested upto 10240
	var inOrderKvs = genIntKeyVals(numKeys)
	//rand.Seed(int64(time.Now().Nanosecond()))
	var buildKvs = randomizeKeyVals(inOrderKvs)

	var m = sorted_map.New()
	for i, kv := range buildKvs {
		m = m.Put(kv.Key, kv.Val)
		var err = m.Valid()
		if err != nil {
			log.Printf("INVALID TREE Map m =\n%s", m.TreeString())
			t.Fatalf("INVALID TREE Store: i=%d; kv.Key=%s; err=%s\n",
				i, kv.Key, err)
		}
	}

	var i int
	var fn = func(k0 sorted_map.MapKey, v0 interface{}) bool {
		//log.Printf("i=%d; k0=%s; v0=%d;", i, k0, v0)
		var k1 = inOrderKvs[i].Key
		var v1 = inOrderKvs[i].Val
		//log.Printf("i=%d; k1=%s; v1=%d;", i, k1, v1)
		if k0.Less(k1) || k1.Less(k0) {
			t.Fatalf("InOrder keys: i=%d; found k0=%s not the expected k1=%s",
				i, k0, k1)
		}
		if v0 != v1 {
			t.Fatalf("InOrder vals: i=%d; found v0=%d not the expected v1=%d",
				i, v0, v1)
		}
		i++
		return true
	}
	m.Range(fn)
	//m.RangeLimit(sorted_map.Inf(-1), sorted_map.Inf(1), fn)
}

func TestCompRangeForwBeg(t *testing.T) {
	var numKeys = 1024 //tested upto 10240
	var inOrderKvs = genIntKeyVals(numKeys)
	//rand.Seed(int64(time.Now().Nanosecond()))
	var buildKvs = randomizeKeyVals(inOrderKvs)

	//log.Printf("buildKvs = %v\n", buildKvs)

	var m = sorted_map.New()
	for i, kv := range buildKvs {
		m = m.Put(kv.Key, kv.Val)
		var err = m.Valid()
		if err != nil {
			log.Printf("INVALID TREE Map m =\n%s", m.TreeString())
			t.Fatalf("INVALID TREE Store: i=%d; kv.Key=%s; err=%s\n",
				i, kv.Key, err)
		}
	}

	var eltOffset = 13
	var startElt = eltOffset
	var i = startElt - 1 //index starts at zero
	var fn = func(k0 sorted_map.MapKey, v0 interface{}) bool {
		//log.Printf("i=%d; k0=%s; v0=%d;", i, k0, v0)
		var k1 = inOrderKvs[i].Key
		var v1 = inOrderKvs[i].Val
		//log.Printf("i=%d; k1=%s; v1=%d;", i, k1, v1)
		if k0.Less(k1) || k1.Less(k0) {
			t.Fatalf("InOrder keys: i=%d; found k0=%s not the expected k1=%s",
				i, k0, k1)
		}
		if v0 != v1 {
			t.Fatalf("InOrder vals: i=%d; found v0=%d not the expected v1=%d",
				i, v0, v1)
		}
		i++
		return true
	}
	m.RangeLimit(IntKey(eltOffset*10), sorted_map.Inf(1), fn)
	//m.RangeLimit(IntKey(130), sorted_map.Inf(1), fn)
}

func TestCompRangeForwEnd(t *testing.T) {
	var numKeys = 1024 //tested upto 10240
	var inOrderKvs = genIntKeyVals(numKeys)
	//rand.Seed(int64(time.Now().Nanosecond()))
	var buildKvs = randomizeKeyVals(inOrderKvs)
	//var destroyKvs = randomizeKeyVals(kvs)

	//log.Printf("buildKvs = %v\n", buildKvs)

	var m = sorted_map.New()
	for i, kv := range buildKvs {
		m = m.Put(kv.Key, kv.Val)
		var err = m.Valid()
		if err != nil {
			t.Fatalf("INVALID TREE Store: i=%d; kv.Key=%s; err=%s\n",
				i, kv.Key, err)
		}
		//log.Printf("Map m =\n%s", m.TreeString())
	}

	var eltOffset = 13
	var startElt = 1
	var i = startElt - 1 //index starts at zero
	var fn = func(k0 sorted_map.MapKey, v0 interface{}) bool {
		//log.Printf("i=%d; k0=%s; v0=%d;", i, k0, v0)
		var k1 = inOrderKvs[i].Key
		var v1 = inOrderKvs[i].Val
		//log.Printf("i=%d; k1=%s; v1=%d;", i, k1, v1)
		if k0.Less(k1) || k1.Less(k0) {
			t.Fatalf("InOrder keys: i=%d; found k0=%s not the expected k1=%s",
				i, k0, k1)
		}
		if v0 != v1 {
			t.Fatalf("InOrder vals: i=%d; found v0=%d not the expected v1=%d",
				i, v0, v1)
		}
		i++
		return true
	}
	m.RangeLimit(sorted_map.Inf(-1), IntKey((numKeys-eltOffset)*10), fn)
	//m.RangeLimit(sorted_map.Inf(-1), IntKey(10110), fn)
}

func TestCompRangeForwBoth(t *testing.T) {
	var numKeys = 1024 //tested upto 10240
	var inOrderKvs = genIntKeyVals(numKeys)
	//rand.Seed(int64(time.Now().Nanosecond()))
	var buildKvs = randomizeKeyVals(inOrderKvs)
	//var destroyKvs = randomizeKeyVals(kvs)

	//log.Printf("buildKvs = %v\n", buildKvs)

	var m = sorted_map.New()
	for i, kv := range buildKvs {
		m = m.Put(kv.Key, kv.Val)
		var err = m.Valid()
		if err != nil {
			t.Fatalf("INVALID TREE Store: i=%d; kv.Key=%s; err=%s\n",
				i, kv.Key, err)
		}
		//log.Printf("Map m =\n%s", m.TreeString())
	}

	var eltOffset = 13
	var startElt = eltOffset
	var i = startElt - 1 //index starts at zero
	var fn = func(k0 sorted_map.MapKey, v0 interface{}) bool {
		//log.Printf("i=%d; k0=%s; v0=%d;", i, k0, v0)
		var k1 = inOrderKvs[i].Key
		var v1 = inOrderKvs[i].Val
		//log.Printf("i=%d; k1=%s; v1=%d;", i, k1, v1)
		if k0.Less(k1) || k1.Less(k0) {
			t.Fatalf("InOrder keys: i=%d; found k0=%s not the expected k1=%s",
				i, k0, k1)
		}
		if v0 != v1 {
			t.Fatalf("InOrder vals: i=%d; found v0=%d not the expected v1=%d",
				i, v0, v1)
		}
		i++
		return true
	}
	m.RangeLimit(IntKey(startElt*10), IntKey((numKeys-eltOffset)*10), fn)
	//m.RangeLimit(IntKey(130), IntKey(10110), fn)
}

func TestCompRangeRevAll(t *testing.T) {
	var numKeys = 1025 //tested upto 10240
	var inOrderKvs = genIntKeyVals(numKeys)
	//rand.Seed(int64(time.Now().Nanosecond()))
	var buildKvs = randomizeKeyVals(inOrderKvs)
	//var destroyKvs = randomizeKeyVals(kvs)

	//log.Printf("buildKvs = %v\n", buildKvs)

	var m = sorted_map.New()
	for i, kv := range buildKvs {
		m = m.Put(kv.Key, kv.Val)
		var err = m.Valid()
		if err != nil {
			t.Fatalf("INVALID TREE Store: i=%d; kv.Key=%s; err=%s\n",
				i, kv.Key, err)
		}
		//log.Printf("Map m =\n%s", m.TreeString())
	}

	var i = numKeys - 1 //index starts at zero
	var fn = func(k0 sorted_map.MapKey, v0 interface{}) bool {
		log.Printf("i=%d; k0=%s; v0=%d;", i, k0, v0)
		var k1 = inOrderKvs[i].Key
		var v1 = inOrderKvs[i].Val
		log.Printf("i=%d; k1=%s; v1=%d;", i, k1, v1)
		if k0.Less(k1) || k1.Less(k0) {
			t.Fatalf("InOrder keys: i=%d; found k0=%s not the expected k1=%s",
				i, k0, k1)
		}
		if v0 != v1 {
			t.Fatalf("InOrder vals: i=%d; found v0=%d not the expected v1=%d",
				i, v0, v1)
		}
		i--
		return true
	}
	m.RangeLimit(sorted_map.Inf(1), sorted_map.Inf(-1), fn)
	//m.RangeLimit(IntKey(130), IntKey(10110), fn)
}

func TestCompRangeRevBeg(t *testing.T) {
	var numKeys = 1024 //tested upto 10240
	var inOrderKvs = genIntKeyVals(numKeys)
	//rand.Seed(int64(time.Now().Nanosecond()))
	var buildKvs = randomizeKeyVals(inOrderKvs)

	//log.Printf("buildKvs = %v\n", buildKvs)

	var m = sorted_map.New()
	for i, kv := range buildKvs {
		m = m.Put(kv.Key, kv.Val)
		var err = m.Valid()
		if err != nil {
			log.Printf("INVALID TREE Map m =\n%s", m.TreeString())
			t.Fatalf("INVALID TREE Store: i=%d; kv.Key=%s; err=%s\n",
				i, kv.Key, err)
		}
	}

	var eltOffset = 13
	var startElt = numKeys - eltOffset
	var i = startElt - 1 //index starts at zero
	var fn = func(k0 sorted_map.MapKey, v0 interface{}) bool {
		//log.Printf("i=%d; k0=%s; v0=%d;", i, k0, v0)
		var k1 = inOrderKvs[i].Key
		var v1 = inOrderKvs[i].Val
		//log.Printf("i=%d; k1=%s; v1=%d;", i, k1, v1)
		if k0.Less(k1) || k1.Less(k0) {
			t.Fatalf("InOrder keys: i=%d; found k0=%s not the expected k1=%s",
				i, k0, k1)
		}
		if v0 != v1 {
			t.Fatalf("InOrder vals: i=%d; found v0=%d not the expected v1=%d",
				i, v0, v1)
		}
		i--
		return true
	}
	m.RangeLimit(IntKey((numKeys-eltOffset)*10), sorted_map.Inf(-1), fn)
}

func TestCompRangeRevEnd(t *testing.T) {
	var numKeys = 1024 //tested upto 10240
	var inOrderKvs = genIntKeyVals(numKeys)
	//rand.Seed(int64(time.Now().Nanosecond()))
	var buildKvs = randomizeKeyVals(inOrderKvs)
	//var destroyKvs = randomizeKeyVals(kvs)

	//log.Printf("buildKvs = %v\n", buildKvs)

	var m = sorted_map.New()
	for i, kv := range buildKvs {
		m = m.Put(kv.Key, kv.Val)
		var err = m.Valid()
		if err != nil {
			t.Fatalf("INVALID TREE Store: i=%d; kv.Key=%s; err=%s\n",
				i, kv.Key, err)
		}
		//log.Printf("Map m =\n%s", m.TreeString())
	}

	var eltOffset = 13
	var startElt = numKeys //- eltOffset
	var i = startElt - 1   //index starts at zero
	var fn = func(k0 sorted_map.MapKey, v0 interface{}) bool {
		//log.Printf("i=%d; k0=%s; v0=%d;", i, k0, v0)
		var k1 = inOrderKvs[i].Key
		var v1 = inOrderKvs[i].Val
		//log.Printf("i=%d; k1=%s; v1=%d;", i, k1, v1)
		if k0.Less(k1) || k1.Less(k0) {
			t.Fatalf("InOrder keys: i=%d; found k0=%s not the expected k1=%s",
				i, k0, k1)
		}
		if v0 != v1 {
			t.Fatalf("InOrder vals: i=%d; found v0=%d not the expected v1=%d",
				i, v0, v1)
		}
		i--
		return true
	}
	m.RangeLimit(sorted_map.Inf(1), IntKey(eltOffset*10), fn)
}

func TestCompRangeRevBoth(t *testing.T) {
	var numKeys = 1024 //tested upto 10240
	var inOrderKvs = genIntKeyVals(numKeys)
	//rand.Seed(int64(time.Now().Nanosecond()))
	var buildKvs = randomizeKeyVals(inOrderKvs)
	//var destroyKvs = randomizeKeyVals(kvs)

	//log.Printf("buildKvs = %v\n", buildKvs)

	var m = sorted_map.New()
	for i, kv := range buildKvs {
		m = m.Put(kv.Key, kv.Val)
		var err = m.Valid()
		if err != nil {
			t.Fatalf("INVALID TREE Store: i=%d; kv.Key=%s; err=%s\n",
				i, kv.Key, err)
		}
		//log.Printf("Map m =\n%s", m.TreeString())
	}

	var eltOffset = 13
	var startElt = numKeys - eltOffset
	var i = startElt - 1 //index starts at zero
	var fn = func(k0 sorted_map.MapKey, v0 interface{}) bool {
		//log.Printf("i=%d; k0=%s; v0=%d;", i, k0, v0)
		var k1 = inOrderKvs[i].Key
		var v1 = inOrderKvs[i].Val
		//log.Printf("i=%d; k1=%s; v1=%d;", i, k1, v1)
		if k0.Less(k1) || k1.Less(k0) {
			t.Fatalf("InOrder keys: i=%d; found k0=%s not the expected k1=%s",
				i, k0, k1)
		}
		if v0 != v1 {
			t.Fatalf("InOrder vals: i=%d; found v0=%d not the expected v1=%d",
				i, v0, v1)
		}
		i--
		return true
	}
	m.RangeLimit(IntKey(startElt*10), IntKey((numKeys-eltOffset)*10), fn)
}
