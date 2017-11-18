package sorted_map_test

import (
	"log"
	"testing"

	"github.com/lleo/go-functional-collections/sorted_map"
)

var gMap *sorted_map.Map

func Test_Basic_BuildSimpleMap(t *testing.T) {
	var gMap = sorted_map.New()

	log.Println("insert #1 100,0")
	gMap = gMap.Put(IntKey(100), 0)
	//log.Printf("gMap=\n%s\n", gMap.TreeString())

	log.Println("insert #2 50,0")
	gMap = gMap.Put(IntKey(50), 0)
	//log.Printf("gMap.NumEntries()=%d; nodes=\n%s\n", gMap.NumEntries(), gMap.TreeString())

	log.Println("insert #3 70,0")
	//gMap = gMap.Put(IntKey(70), 0)
	gMap = gMap.Put(IntKey(30), 0)
	//log.Printf("gMap=\n%s\n", gMap.TreeString())

	log.Println("insert #3 40,0")
	//gMap = gMap.Put(IntKey(40), 0)
	gMap = gMap.Put(IntKey(40), 0)
	//log.Printf("gMap=\n%s\n", gMap.TreeString())

	log.Println("insert #4 150,0")
	gMap = gMap.Put(IntKey(150), 0)
	//log.Printf("gMap=\n%s\n", gMap.TreeString())

	log.Println("insert #5 200,0")
	gMap = gMap.Put(IntKey(200), 0)
	//log.Printf("gMap=\n%s\n", gMap.TreeString())

	log.Println("insert #6 250,0")
	gMap = gMap.Put(IntKey(250), 0)
	//log.Printf("gMap=\n%s\n", gMap.TreeString())

	log.Println("insert #7 300,0")
	gMap = gMap.Put(IntKey(300), 0)
	//log.Printf("gMap=\n%s\n", gMap.TreeString())

	log.Println("insert #8 350,0")
	gMap = gMap.Put(IntKey(350), 0)
	log.Printf("gMap=\n%s\n", gMap.TreeString())

	log.Println("replace #8 350,10")
	gMap = gMap.Put(IntKey(350), 10)
	log.Printf("gMap=\n%s\n", gMap.TreeString())
	log.Printf("keyVals = %q\n", gMap)
	//gMap.String = "{30: 0, 50: 400, 10: 0, 15: 0, 20: 0, 25: 0, 30: 0, 35: 10}"
}

func Test_Basic_Del_OnlyRoot(t *testing.T) {
	var m = sorted_map.New()

	m = m.Put(IntKey(10), 0)
	m = m.Del(IntKey(10))

	if m.NumEntries() != 0 {
		t.Fatal("m.NumEntries() != 0")
	}
}

func Test_Basic_Del_TwoNodes0(t *testing.T) {
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

func Test_Basic_Del_TwoNodes(t *testing.T) {
	var m, keys = buildMap(2)

	var expectedNumEntries = uint(len(keys))
	for i := 0; i < 2; i++ {
		var k = keys[i]
		m = m.Del(k)
		expectedNumEntries--

		if m.NumEntries() != expectedNumEntries {
			t.Fatal("m.Del(%s) failed.")
		}
	}

	log.Println("m =", m)
}
