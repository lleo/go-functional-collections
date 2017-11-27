package sorted_map_test

import (
	"log"
	"testing"

	"github.com/lleo/go-functional-collections/sorted_map"
)

const Black = sorted_map.Black
const Red = sorted_map.Red

var mknod = sorted_map.MakeNode
var mkmap = sorted_map.MakeMap

func TestBasicPutCase1(t *testing.T) {
	t.Fatal("not implemented")
}

func TestBasicPutCase2(t *testing.T) {
	t.Fatal("not implemented")
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
		t.Fatal("TestBasicPutCase4: orig Map and duplicate of orig Map are not identical.")
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

	log.Printf("TestBasicPutCase4: m=\n%s", m.Root().ToString(-1))

	if m.NumEntries() != 6 {
		t.Fatal("m.NumEntries() != 6")
	}

	if !origM.Equiv(dupM) {
		t.Fatal("TestBasicPutCase4: orig Map and duplicate of orig Map are not identical.")
	}
}

func TestBasicDelCase1(t *testing.T) {
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

func TestBasicDelCase3(t *testing.T) {
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

//func TestBasicDelCase1(t *testing.T) {
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
