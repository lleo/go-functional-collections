package sorted_map_test

import (
	"log"
	"testing"

	"github.com/lleo/go-functional-collections/sorted_map"
)

func Test_Basic_BuildMap(t *testing.T) {
	var m = sorted_map.New()

	log.Println("insert #1 100,0")
	m = m.Put(IntKey(100), 0)
	//log.Printf("m=\n%s\n", m.TreeString())

	log.Println("insert #2 50,0")
	m = m.Put(IntKey(50), 0)
	//log.Printf("m.NumEntries()=%d; nodes=\n%s\n", m.NumEntries(), m.TreeString())

	log.Println("insert #3 70,0")
	//m = m.Put(IntKey(70), 0)
	m = m.Put(IntKey(30), 0)
	//log.Printf("m=\n%s\n", m.TreeString())

	log.Println("insert #3 40,0")
	//m = m.Put(IntKey(40), 0)
	m = m.Put(IntKey(40), 0)
	//log.Printf("m=\n%s\n", m.TreeString())

	log.Println("insert #4 150,0")
	m = m.Put(IntKey(150), 0)
	//log.Printf("m=\n%s\n", m.TreeString())

	log.Println("insert #5 200,0")
	m = m.Put(IntKey(200), 0)
	//log.Printf("m=\n%s\n", m.TreeString())

	log.Println("insert #6 250,0")
	m = m.Put(IntKey(250), 0)
	//log.Printf("m=\n%s\n", m.TreeString())

	log.Println("insert #7 300,0")
	m = m.Put(IntKey(300), 0)
	//log.Printf("m=\n%s\n", m.TreeString())

	log.Println("insert #8 350,0")
	m = m.Put(IntKey(350), 0)
	log.Printf("m=\n%s\n", m.TreeString())

	log.Println("replace #8 350,10")
	m = m.Put(IntKey(350), 10)
	log.Printf("m=\n%s\n", m.TreeString())
	log.Printf("keyVals = %q\n", m)
	//m.String = "{30: 0, 50: 400, 10: 0, 15: 0, 20: 0, 25: 0, 30: 0, 35: 10}"
}

func Test_Basic_Del_OnlyRoot(t *testing.T) {
	var m = sorted_map.New()

	m = m.Put(IntKey(10), 0)
	m = m.Del(IntKey(10))

	if m.NumEntries() != 0 {
		t.Fatal("m.NumEntries() != 0")
	}
}

func tTest_Basic_Del_TwoNodes0(t *testing.T) {
	var m = sorted_map.New()

	m = m.Put(IntKey(100), 10)
	m = m.Put(IntKey(150), 15)

	log.Println("Test_Basic_Del_TwoNodes: *********** calling m.Del(IntKey(100)) ***********")
	m = m.Del(IntKey(100))

	if m.NumEntries() != 1 {
		t.Fatal("m.NumEntries() != 1")
	}

	var val100 = m.Get(IntKey(100))
	if val100 != nil {
		t.Fatal("m.Get(IntKey(100)) != nil")
	}

	var val150 = m.Get(IntKey(150))
	if val150 != 15 {
		t.Fatalf("m.Get(IntKey(150)),%d != 15", val150)
	}

	log.Println("m =", m)
}

func buildMap(n int) (*sorted_map.Map, []sorted_map.MapKey) {
	var m = sorted_map.New()
	var keys = make([]sorted_map.MapKey, n)
	for i := 0; i < n; i++ {
		var k = IntKey(i * 100)
		var v interface{} = k + 1
		keys[i] = k
		m = m.Put(k, v)
	}
	return m, keys
}

func Test_Basic_Del_LoHi_TwoNodes(t *testing.T) {
	var m, keys = buildMap(2)

	log.Printf("JUST BUILT: m=\n%s", m.TreeString())

	var expectedNumEntries = uint(len(keys))
	for i := 0; i < 2; i++ {
		var k = keys[i]
		m = m.Del(k)
		expectedNumEntries--

		if m.NumEntries() != expectedNumEntries {
			t.Fatal("m.Del(%s) failed.")
		}
	}

	log.Printf("COMPLETELY DELETED m =\n%s", m.TreeString())
}

func Test_Basic_Del_HiLo_TwoNodes(t *testing.T) {
	var m, keys = buildMap(2)

	log.Printf("JUST BUILT: m=\n%s", m.TreeString())

	var expectedNumEntries = uint(len(keys))
	for i := 1; i >= 0; i-- {
		var k = keys[i]
		m = m.Del(k)
		expectedNumEntries--

		if m.NumEntries() != expectedNumEntries {
			t.Fatal("m.Del(%s) failed.")
		}
	}

	log.Printf("COMPLETELY DELETED m =\n%s", m.TreeString())
}

const Black = sorted_map.Black
const Red = sorted_map.Red

var mknod = sorted_map.MakeNode
var mkmap = sorted_map.MakeMap

func Test_Del_Case1(t *testing.T) {
	var m = mkmap(3,
		mknod(IntKey(20), 20, Black,
			nil,
			mknod(IntKey(30), 30, Red, nil, nil),
		))

	m = m.Del(IntKey(30))
	if m.NumEntries() != 2 {
		t.Fatalf("m.NumEntries(),%d != 2", m.NumEntries())
	}

	if !m.Root().IsBlack() {
		t.Fatal("!m.Root().IsBlack()")
	}

	if m.Root().Ln() != nil {
		t.Fatal("m.Root().Rn() != nil")
	}

	if m.Root().Rn() != nil {
		t.Fatal("m.Root().Ln() != nil")
	}
}

func Test_Del_Case3(t *testing.T) {
	var m = mkmap(3,
		mknod(IntKey(20), 20, Black,
			mknod(IntKey(10), 10, Black, nil, nil),
			mknod(IntKey(30), 30, Black, nil, nil),
		))

	m = m.Del(IntKey(30))
	if m.NumEntries() != 2 {
		t.Fatalf("m.NumEntries(),%d != 2", m.NumEntries())
	}

	if !m.Root().IsBlack() {
		t.Fatal("!m.Root().IsBlack()")
	}

	if !m.Root().Ln().IsRed() {
		t.Fatal("!m.Root().Ln().IsRed()")
	}

	if m.Root().Rn() != nil {
		t.Fatal("m.Root().Rn() != nil")
	}
}

//func Test_Del_Case1(t *testing.T) {
//	var m = mkmap(3,
//		mknod(20, Black
//			Black,
//			mknod(10, Black,
//				mknod(5, Red, nil, nil),
//				mknod(15, Red, nil, nil)),
//			mknod(30, Black,
//				mknod(25, Red, nil, nil),
//				mknod(35, Red, nil, nil)),
//		))
//
//	m.Del(IntKey(25))
//	if m.NumEntries() != 6 {
//		t.Fatal("m.NumEntries() != 2")
//	}
//
//	if m.Root().IsRed() {
//		t.Fatal("m.Root().IsRed()")
//	}
//
//	if m.Root().Ln().IsRed() {
//		t.Fatal("m.Root().Ln().IsRed()")
//	}
//
//	//FIXME...
//
//	if m.Root().Rn() != nil {
//		t.Fatal("m.Root().Rn() != nil")
//	}
//}
