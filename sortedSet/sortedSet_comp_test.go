package sorted_set

import (
	"log"
	"testing"

	"github.com/lleo/go-functional-collections/sorted"
)

func TestCompBuildSet(t *testing.T) {
	var numKeys = 1024 //tested upto 10240
	var inOrderKeys = buildKeys(numKeys)
	//rand.Seed(int64(time.Now().Nanosecond()))
	var keys = randomizeKeys(inOrderKeys)

	//log.Println("keys[:10] =", keys[:10])
	var s = buildSet(keys)

	for _, key := range keys {
		s = s.Set(key)

		var err = s.valid()
		if err != nil {
			t.Fatalf("Invalid Tree. err=%s\n", err)
		}
	}

	//log.Printf("Set s =\n%s", s.treeString())
	//log.Printf("Set s =\n%s", s.String())

	var i int
	var fn = func(k0 sorted.Key) bool {
		var k1 = inOrderKeys[i]
		if sorted.Cmp(k0, k1) != 0 { //k0 != k1
			t.Fatalf("InOrder keys: i=%d; found k0=%s not the expected k1=%s\n",
				i, k0, k1)
		}
		i++
		return true
	}
	s.Range(fn)
}

func TestCompDestroySet(t *testing.T) {
	var numKeys = 1024 //tested upto 10240
	var keys = buildKeys(numKeys)
	//rand.Seed(int64(time.Now().Nanosecond()))
	var buildKeys = randomizeKeys(keys)
	var destroyKeys = randomizeKeys(keys)

	//log.Printf("buildKeys = %v\n", buildKeys)

	var s = buildSet(buildKeys)

	//log.Printf("AFTER ALL SET BUILDING: Set s=\n%s", s.treeString())
	//log.Printf("s = %s\n", s.String())

	//log.Printf("destroyKeys = %v\n", destroyKeys)

	var shouldHaveKeys = make([]sorted.Key, len(destroyKeys))
	copy(shouldHaveKeys, destroyKeys)

	for i, key := range destroyKeys {
		//var origM = s
		//var dupOrigM = s.dup()

		//log.Printf("*******Removing key=%s; i=%d \n", key, i)

		//log.Printf("BEFORE REMOVE: Set s=\n%s", s.treeString())

		var found bool
		s, found = s.Remove(key)

		//log.Printf("AFTER REMOVE Set s=\n%s", s.treeString())

		if !found {
			t.Fatalf("Remove: i=%d; key=%s not found!\n", i, key)
		}

		var err = s.valid()
		if err != nil {
			log.Printf("Set s=\n%s", s.treeString())
			t.Fatalf("INVALID TREE Remove: i=%d; key=%s; err=%s\n",
				i, key, err)
		}

		shouldHaveKeys = shouldHaveKeys[1:] //take the first elt off ala range
		for _, key0 := range shouldHaveKeys {
			var isSet = s.IsSet(key0)
			if !isSet {
				t.Fatalf("Remove: i=%d; for key=%s: "+
					"failed to find shouldHave key0=%s", i, key, key0)
			}
		}
		//if !origM.equiv(dupOrigM) {
		//	t.Fatal("the original set was modified during Remove(%s).", key)
		//}
	}

	if s.NumEntries() != 0 {
		t.Fatal("s.NumEntries() != 0")
	}
}

func TestCompRangeForwAll(t *testing.T) {
	var numKeys = 1024 //tested upto 10240
	var inOrderKeys = buildKeys(numKeys)
	//rand.Seed(int64(time.Now().Nanosecond()))
	var buildKeys = randomizeKeys(inOrderKeys)

	var s = buildSet(buildKeys)

	var i int
	var fn = func(k0 sorted.Key) bool {
		//log.Printf("i=%d; k0=%s;", i, k0)
		var k1 = inOrderKeys[i]
		//log.Printf("i=%d; k1=%s;", i, k1)
		if sorted.Cmp(k0, k1) != 0 {
			t.Fatalf("InOrder keys: i=%d; found k0=%s not the expected k1=%s",
				i, k0, k1)
		}
		i++
		return true
	}
	s.Range(fn)
}

func TestCompRangeForwBeg(t *testing.T) {
	var numKeys = 1024 //tested upto 10240
	var inOrderKeys = buildKeys(numKeys)
	//rand.Seed(int64(time.Now().Nanosecond()))
	var buildKeys = randomizeKeys(inOrderKeys)

	//log.Printf("buildKeys = %v\n", buildKeys)

	var s = buildSet(buildKeys)

	var eltOffset = 13
	var startElt = eltOffset
	var i = startElt - 1 //index starts at zero
	var fn = func(k0 sorted.Key) bool {
		//log.Printf("i=%d; k0=%s;", i, k0)
		var k1 = inOrderKeys[i]
		//log.Printf("i=%d; k1=%s;", i, k1)
		if sorted.Cmp(k0, k1) != 0 {
			t.Fatalf("InOrder keys: i=%d; found k0=%s not the expected k1=%s",
				i, k0, k1)
		}
		i++
		return true
	}
	s.RangeLimit(sorted.IntKey(eltOffset*10), sorted.InfKey(1), fn)
}

func TestCompRangeForwEnd(t *testing.T) {
	var numKeys = 1024 //tested upto 10240
	var inOrderKeys = buildKeys(numKeys)
	//rand.Seed(int64(time.Now().Nanosecond()))
	var buildKeys = randomizeKeys(inOrderKeys)

	//log.Printf("buildKeys = %v\n", buildKeys)

	var s = buildSet(buildKeys)

	var eltOffset = 13
	var startElt = 1
	var i = startElt - 1 //index starts at zero
	var fn = func(k0 sorted.Key) bool {
		//log.Printf("i=%d; k0=%s; v0=%d;", i, k0)
		var k1 = inOrderKeys[i]
		//log.Printf("i=%d; k1=%s; v1=%d;", i, k1)
		if sorted.Cmp(k0, k1) != 0 {
			t.Fatalf("InOrder keys: i=%d; found k0=%s not the expected k1=%s",
				i, k0, k1)
		}
		i++
		return true
	}
	s.RangeLimit(sorted.InfKey(-1), sorted.IntKey((numKeys-eltOffset)*10), fn)
}

func TestCompRangeForwBoth(t *testing.T) {
	var numKeys = 1024 //tested upto 10240
	var inOrderKeys = buildKeys(numKeys)
	//rand.Seed(int64(time.Now().Nanosecond()))
	var buildKeys = randomizeKeys(inOrderKeys)

	//log.Printf("buildKeys = %v\n", buildKeys)

	var s = buildSet(buildKeys)

	var eltOffset = 13
	var startElt = eltOffset
	var i = startElt - 1 //index starts at zero
	var fn = func(k0 sorted.Key) bool {
		//log.Printf("i=%d; k0=%s; v0=%d;", i, k0)
		var k1 = inOrderKeys[i]
		//log.Printf("i=%d; k1=%s; v1=%d;", i, k1)
		if sorted.Cmp(k0, k1) != 0 {
			t.Fatalf("InOrder keys: i=%d; found k0=%s not the expected k1=%s",
				i, k0, k1)
		}
		i++
		return true
	}
	s.RangeLimit(sorted.IntKey(startElt*10),
		sorted.IntKey((numKeys-eltOffset)*10), fn)
	//s.RangeLimit(sorted.IntKey(130), sorted.IntKey(10110), fn)
}

func TestCompRangeRevAll(t *testing.T) {
	var numKeys = 1025 //tested upto 10240
	var inOrderKeys = buildKeys(numKeys)
	//rand.Seed(int64(time.Now().Nanosecond()))
	var buildKeys = randomizeKeys(inOrderKeys)

	//log.Printf("buildKeys = %v\n", buildKeys)

	var s = buildSet(buildKeys)

	var i = numKeys - 1 //index starts at zero
	var fn = func(k0 sorted.Key) bool {
		var k1 = inOrderKeys[i]
		//log.Printf("i=%d; k0=%s; v0=%d;", i, k0)
		//log.Printf("i=%d; k1=%s; v1=%d;", i, k1)
		if sorted.Cmp(k0, k1) != 0 {
			t.Fatalf("InOrder keys: i=%d; found k0=%s not the expected k1=%s",
				i, k0, k1)
		}
		i--
		return true
	}
	s.RangeLimit(sorted.InfKey(1), sorted.InfKey(-1), fn)
}

func TestCompRangeRevBeg(t *testing.T) {
	var numKeys = 1024 //tested upto 10240
	var inOrderKeys = buildKeys(numKeys)
	//rand.Seed(int64(time.Now().Nanosecond()))
	var buildKeys = randomizeKeys(inOrderKeys)

	//log.Printf("buildKeys = %v\n", buildKeys)

	var s = buildSet(buildKeys)

	var eltOffset = 13
	var startElt = numKeys - eltOffset
	var i = startElt - 1 //index starts at zero
	var fn = func(k0 sorted.Key) bool {
		var k1 = inOrderKeys[i]
		//log.Printf("i=%d; k0=%s; v0=%d;", i, k0)
		//log.Printf("i=%d; k1=%s; v1=%d;", i, k1)
		if sorted.Cmp(k0, k1) != 0 {
			t.Fatalf("InOrder keys: i=%d; found k0=%s not the expected k1=%s",
				i, k0, k1)
		}
		i--
		return true
	}
	s.RangeLimit(sorted.IntKey((numKeys-eltOffset)*10), sorted.InfKey(-1), fn)
}

func TestCompRangeRevEnd(t *testing.T) {
	var numKeys = 1024 //tested upto 10240
	var inOrderKeys = buildKeys(numKeys)
	//rand.Seed(int64(time.Now().Nanosecond()))
	var buildKeys = randomizeKeys(inOrderKeys)

	//log.Printf("buildKeys = %v\n", buildKeys)

	var s = buildSet(buildKeys)

	var eltOffset = 13
	var startElt = numKeys //- eltOffset
	var i = startElt - 1   //index starts at zero
	var fn = func(k0 sorted.Key) bool {
		var k1 = inOrderKeys[i]
		//log.Printf("i=%d; k0=%s; v0=%d;", i, k0)
		//log.Printf("i=%d; k1=%s; v1=%d;", i, k1)
		if sorted.Cmp(k0, k1) != 0 {
			t.Fatalf("InOrder keys: i=%d; found k0=%s not the expected k1=%s",
				i, k0, k1)
		}
		i--
		return true
	}
	s.RangeLimit(sorted.InfKey(1), sorted.IntKey(eltOffset*10), fn)
}

func TestCompRangeRevBoth(t *testing.T) {
	var numKeys = 1024 //tested upto 10240
	var inOrderKeys = buildKeys(numKeys)
	//rand.Seed(int64(time.Now().Nanosecond()))
	var buildKeys = randomizeKeys(inOrderKeys)

	//log.Printf("buildKeys = %v\n", buildKeys)

	var s = buildSet(buildKeys)

	var eltOffset = 13
	var startElt = numKeys - eltOffset
	var i = startElt - 1 //index starts at zero
	var fn = func(k0 sorted.Key) bool {
		var k1 = inOrderKeys[i]
		//log.Printf("i=%d; k0=%s; v0=%d;", i, k0)
		//log.Printf("i=%d; k1=%s; v1=%d;", i, k1)
		if sorted.Cmp(k0, k1) != 0 {
			t.Fatalf("InOrder keys: i=%d; found k0=%s not the expected k1=%s",
				i, k0, k1)
		}
		i--
		return true
	}
	s.RangeLimit(sorted.IntKey(startElt*10),
		sorted.IntKey((numKeys-eltOffset)*10), fn)
}
