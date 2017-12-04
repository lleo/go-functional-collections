package sorted_map_test

import (
	"log"
	"testing"
)

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
		mknod(20, Black,
			mknod(10, Red, nil, nil),
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
		mknod(20, Black,
			mknod(10, Black, nil, nil),
			mknod(40, Black,
				mknod(30, Red, nil, nil),
				mknod(50, Red, nil, nil),
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
	//	mknod(7940, Black,
	//		mknod(4930, Black,
	//			nil,
	//			mknod(7100, Red, nil, nil)),
	//		mknod(8090, Black,
	//			nil,
	//			mknod(10050, Red, nil, nil)),
	//	))
	//insert order 50, 20, 60, 40, 70 ???
	var m = mkmap(5,
		mknod(50, Black,
			mknod(20, Black,
				nil,
				mknod(40, Red, nil, nil)),
			mknod(60, Black,
				nil,
				mknod(70, Red, nil, nil)),
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
		mknod(10, Black, nil, nil))

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
		mknod(10, Black, nil, nil))

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
		mknod(10, Black,
			nil,
			mknod(20, Red, nil, nil),
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
		mknod(20, Black,
			mknod(10, Red, nil, nil),
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
		mknod(20, Black,
			nil,
			mknod(30, Red, nil, nil),
		))

	//log.Printf("BEFORE REMOVE: Map m=\n%s", m.TreeString())

	m = m.Del(IntKey(30))

	//log.Printf("AFTER REMOVE Map m=\n%s", m.TreeString())

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
		mknod(20, Black,
			mknod(10, Red, nil, nil),
			nil,
		))

	log.Printf("BEFORE REMOVE: Map m=\n%s", m.TreeString())

	m = m.Del(IntKey(10))

	log.Printf("AFTER REMOVE Map m=\n%s", m.TreeString())

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
		mknod(20, Black,
			mknod(10, Black, nil, nil),
			mknod(30, Black, nil, nil),
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
		mknod(20, Black,
			mknod(10, Black, nil, nil),
			mknod(40, Red,
				mknod(30, Black, nil, nil),
				mknod(50, Black,
					nil,
					mknod(60, Red, nil, nil)))))

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

func TestBasicDelTermTree0(t *testing.T) {
	var m = mkmap(7,
		mknod(40, Black,
			mknod(20, Black,
				mknod(10, Red, nil, nil),
				mknod(30, Red, nil, nil)),
			mknod(70, Red,
				mknod(50, Black, nil, nil),
				mknod(80, Black, nil, nil))))

	var origM = m
	var dupOrigM = m.Dup()
	var origMapStr0 = m.TreeString()

	log.Printf("BEFORE DEL m = \n%s", m.TreeString())

	m = m.Del(IntKey(20))

	log.Printf("AFTER DEL m = \n%s", m.TreeString())

	var err = m.Valid()
	if err != nil {
		t.Fatalf("INVALID TREE AFTER Del(IntKey(20)); err=%s\n", err)
	}

	if m.NumEntries() != 6 {
		t.Fatal("m.NumEntries(),%d != 6", m.NumEntries())
	}

	var origMapStr1 = origM.TreeString()
	if origMapStr0 != origMapStr1 {
		log.Printf("origMapStr0 != origMapStr1:\n"+
			"origMapStr0=\n%s\norigMapStr1=\n%s", origMapStr0, origMapStr1)
	}

	if !origM.Equiv(dupOrigM) {
		t.Fatal("TestBasicDelTwoChildrenCase2: " +
			"orig Map and duplicate of orig Map are not identical.")
	}
}

//deleteCase4
func TestBasicDelTermTree1(t *testing.T) {
	var m = mkmap(6,
		mknod(40, Black,
			mknod(10, Black,
				nil,
				mknod(30, Red, nil, nil)),
			mknod(70, Red,
				mknod(50, Black, nil, nil),
				mknod(80, Black, nil, nil))))

	//shouldHave after Del(70)
	var shouldHaveKvs = []KeyVal{
		{IntKey(10), 10},
		{IntKey(30), 30},
		{IntKey(40), 40},
		{IntKey(50), 50},
		{IntKey(80), 80},
	}

	var origM = m
	var dupOrigM = m.Dup()
	var origMapStr0 = m.TreeString()

	log.Printf("BEFORE DEL m = \n%s", m.TreeString())

	m = m.Del(IntKey(70))

	log.Printf("AFTER DEL m = \n%s", m.TreeString())

	var err = m.Valid()
	if err != nil {
		t.Fatalf("INVALID TREE AFTER Del(IntKey(20)); err=%s\n", err)
	}

	if m.NumEntries() != 5 {
		t.Fatal("m.NumEntries(),%d != 5", m.NumEntries())
	}

	for _, kv := range shouldHaveKvs {
		var val, found = m.Load(kv.Key)
		if !found {
			log.Printf("failed to find shouldHave key=%s in Tree=\n%s",
				kv.Key, m.TreeString())
			t.Fatalf("Failed to find shouldHave kv.Key=%s", kv.Key)
		}
		if val != kv.Val {
			t.Fatalf("found val,%v != expected val,%v", val, kv.Val)
		}
	}

	var origMapStr1 = origM.TreeString()
	if origMapStr0 != origMapStr1 {
		log.Printf("origMapStr0 != origMapStr1:\n"+
			"origMapStr0=\n%s\norigMapStr1=\n%s", origMapStr0, origMapStr1)
	}

	if !origM.Equiv(dupOrigM) {
		t.Fatal("TestBasicDelTwoChildrenCase2: " +
			"orig Map and duplicate of orig Map are not identical.")
	}
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
