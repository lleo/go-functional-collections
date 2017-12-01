package sorted_map_test

import (
	"log"
	"testing"

	"github.com/lleo/go-functional-collections/sorted_map"
)

const Black = sorted_map.Black
const Red = sorted_map.Red

var mkmap = sorted_map.MakeMap

//var mknod = sorted_map.MakeNode

func mknod(i int, ln, rn *sorted_map.Node) *sorted_map.Node {
	return sorted_map.MakeNode(IntKey(i), i, ln, rn)
}

func TestBasicPutCase1(t *testing.T) {
	var m = mkmap(0, nil)

	var origM = m
	var dupM = m.Dup()

	m = m.Put(IntKey(10), 10)

	if m.NumEntries() != 1 {
		t.Fatal("m.NumEntries() != 1")
	}

	if !origM.Equiv(dupM) {
		t.Fatal("TestBasicPutCase1: orig Map and duplicate of orig Map are not identical.")
	}
}

func TestBasicPutCase2(t *testing.T) {
	var m = mkmap(2,
		mknod(IntKey(20), 20, Black,
			mknod(IntKey(10), 10, Red, nil, nil),
			nil))

	var origM = m
	var dupM = m.Dup()

	m = m.Put(IntKey(30), 30)

	if m.NumEntries() != 3 {
		t.Fatal("m.NumEntries() != 1")
	}

	if !origM.Equiv(dupM) {
		t.Fatal("TestBasicPutCase2: orig Map and duplicate of orig Map are not identical.")
	}
}

func TestBasicPutCase3(t *testing.T) {
	//insert order 10, 20, 50, 40, 30, 60
	var m = mkmap(5,
		mknod(IntKey(20), 20, Black,
			mknod(IntKey(10), 10, Black, nil, nil),
			mknod(IntKey(40), 40, Black,
				mknod(IntKey(30), 30, Red, nil, nil),
				mknod(IntKey(50), 50, Red, nil, nil),
			),
		))

	var origM = m
	var dupM = m.Dup()

	m = m.Put(IntKey(60), 60)

	if m.NumEntries() != 6 {
		t.Fatal("m.NumEntries() != 6")
	}

	log.Printf("\n%s", m.TreeString())

	if !origM.Equiv(dupM) {
		t.Fatal("TestBasicPutCase3: orig Map and duplicate of orig Map are not identical.")
	}
}

func TestBasicPutCase4(t *testing.T) {
	//var m = mkmap(5,
	//	mknod(IntKey(7940), 7940, Black,
	//		mknod(IntKey(4930), 4930, Black,
	//			nil,
	//			mknod(IntKey(7100), 7100, Red, nil, nil)),
	//		mknod(IntKey(8090), 8090, Black,
	//			nil,
	//			mknod(IntKey(10050), 10050, Red, nil, nil)),
	//	))
	//insert order 50, 20, 60, 40, 70 ???
	var m = mkmap(5,
		mknod(IntKey(50), 50, Black,
			mknod(IntKey(20), 20, Black,
				nil,
				mknod(IntKey(40), 40, Red, nil, nil)),
			mknod(IntKey(60), 60, Black,
				nil,
				mknod(IntKey(70), 70, Red, nil, nil)),
		))

	var origM = m      //copy the pointer
	var dupM = m.Dup() //copy the value

	//m = m.Put(IntKey(5310), 5310)
	m = m.Put(IntKey(30), 30)

	log.Printf("TestBasicPutCase4: m=\n%s", m.Root().TreeString())

	if m.NumEntries() != 6 {
		t.Fatal("m.NumEntries() != 6")
	}

	if !origM.Equiv(dupM) {
		t.Fatal("TestBasicPutCase4: orig Map and duplicate of orig Map are not identical.")
	}
}

func TestBasicDelCase1Tree0(t *testing.T) {
	var m0 = mkmap(1,
		mknod(IntKey(10), 10, Black, nil, nil))

	var then = m0.TreeString()
	//var dupM0 = m0.Dup()

	var m1 = m0.Del(IntKey(10))

	if m1.NumEntries() != 0 {
		t.Fatal("m.NumEntries() != 0")
	}

	var now = m0.TreeString()
	if then != now {
		log.Printf("origninal tree changeed:\nTHEN: %s\nNOW: %s",
			then, now)
		t.Fatal("The original tree changed.")
	}

	//if !m0.Equiv(dupM0) {
	//	t.Fatal("The original tree changed.")
	//}
}

func TestBasicDelCase1Tree1(t *testing.T) {
	var m0 = mkmap(1,
		mknod(IntKey(10), 10, Black, nil, nil))

	var then = m0.TreeString()
	//var dupM0 = m0.Dup()

	var m1 = m0.Del(IntKey(10))

	if m1.NumEntries() != 0 {
		t.Fatal("m.NumEntries() != 0")
	}

	var now = m0.TreeString()
	if then != now {
		log.Printf("origninal tree changeed:\nTHEN: %s\nNOW: %s",
			then, now)
		t.Fatal("The original tree changed.")
	}

	//if !m0.Equiv(dupM0) {
	//	t.Fatal("The original tree changed.")
	//}
}

func TestBasicDelCase1Tree2(t *testing.T) {
	var m0 = mkmap(2,
		mknod(IntKey(10), 10, Black,
			nil,
			mknod(IntKey(20), 20, Red, nil, nil),
		))

	var then = m0.TreeString()

	var m1 = m0.Del(IntKey(10))

	if m1.NumEntries() != 1 {
		t.Fatal("m.NumEntries() != 1")
	}

	var now = m0.TreeString()
	if then != now {
		log.Printf("origninal tree changeed:\nTHEN: %s\nNOW: %s",
			then, now)
		t.Fatal("The original tree changed.")
	}
}

func TestBasicDelCase1Tree3(t *testing.T) {
	var m0 = mkmap(2,
		mknod(IntKey(20), 20, Black,
			mknod(IntKey(10), 10, Red, nil, nil),
			nil,
		))

	var then = m0.TreeString()

	var m1 = m0.Del(IntKey(20))

	if m1.NumEntries() != 1 {
		t.Fatal("m.NumEntries() != 1")
	}

	var now = m0.TreeString()
	if then != now {
		log.Printf("origninal tree changeed:\nTHEN: %s\nNOW: %s",
			then, now)
		t.Fatal("The original tree changed.")
	}
}

// DeleteCase1 is exhaustively tested.

func TestBasicDelCase2Tree0(t *testing.T) {
	var m = mkmap(2,
		mknod(IntKey(20), 20, Black,
			nil,
			mknod(IntKey(30), 30, Red, nil, nil),
		))

	m = m.Del(IntKey(30))

	if m.NumEntries() != 1 {
		t.Fatalf("m.NumEntries(),%d != 1", m.NumEntries())
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

func TestBasicDelCase2Tree1(t *testing.T) {
	var m = mkmap(2,
		mknod(IntKey(20), 20, Black,
			mknod(IntKey(10), 30, Red, nil, nil),
			nil,
		))

	m = m.Del(IntKey(10))

	if m.NumEntries() != 1 {
		t.Fatalf("m.NumEntries(),%d != 1", m.NumEntries())
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

func TestBasicDelCase3Tree0(t *testing.T) {
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

func TestBasicDelCase6Tree0(t *testing.T) {
	var m = mkmap(6,
		mknod(IntKey(20), 20, Black,
			mknod(IntKey(10), 10, Black, nil, nil),
			mknod(IntKey(40), 40, Red,
				mknod(IntKey(30), 30, Black, nil, nil),
				mknod(IntKey(50), 50, Black,
					nil,
					mknod(IntKey(60), 60, Red, nil, nil)))))

	var origM = m
	var dupOrigM = m.Dup()
	var origMapStr0 = m.TreeString()

	log.Printf("origMapStr0 =\n%s", origMapStr0)

	m = m.Del(IntKey(30))
	if m.NumEntries() != 5 {
		t.Fatal("m.NumEntries() != 5")
	}

	var origMapStr1 = origM.TreeString()
	if origMapStr0 != origMapStr1 {
		log.Printf("origMapStr0 != origMapStr1:\n"+
			"origMapStr0=\n%s\norigMapStr1=\n%s", origMapStr0, origMapStr1)
	}

	if !origM.Equiv(dupOrigM) {
		t.Fatal("TestBasicPutCase4: orig Map and duplicate of orig Map are not identical.")
	}
}

func TestBasicDelTwoChildrenCase2(t *testing.T) {
	var m = mkmap(7,
		mknod(IntKey()))
}

//func TestBasicDelCase1Deeper(t *testing.T) {
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
