package sorted_map_test

import (
	"log"
	"testing"
)

func TestBasicLoadOrStoreTree0(t *testing.T) {
	var m0 = mkmap(
		mknod(20, Black,
			mknod(10, Red, nil, nil),
			mknod(30, Red, nil, nil)))

	var origM = m0
	var dupM = m0.Dup()

	//log.Printf("Before LoadOrStore: m0=\n%s", m0.TreeString())

	var m1, val, found = m0.LoadOrStore(IntKey(10), 10)

	//log.Printf("After LoadOrStore: m1=\n%s", m1.TreeString())

	if !found {
		t.Fatal("key not found; it was expected it to be found.")
	}

	if val != 10 {
		t.Fatal("val,%d != expected val,%v,", val, m0.Root().Ln().Val())
	}

	if m1 != m0 {
		t.Fatal("current tree not same as orig tree.")
	}

	if !origM.Equiv(dupM) {
		t.Fatal("TestBasicPutCase1: orig Map and duplicate of orig Map are not identical.")
	}
}

func TestBasicLoadOrStoreTree1(t *testing.T) {
	var m0 = mkmap(
		mknod(20, Black,
			nil,
			mknod(30, Red, nil, nil)))

	var origM = m0
	var dupM = m0.Dup()

	//log.Printf("Before LoadOrStore: m0=\n%s", m0.TreeString())

	var m1, val, found = m0.LoadOrStore(IntKey(10), 10)

	//log.Printf("After LoadOrStore: m1=\n%s", m1.TreeString())

	if found {
		t.Fatal("key was found; it was not expected to be found.")
	}

	if val != nil {
		t.Fatal("val,%d != expected val,%v,", val, nil)
	}

	if m1.NumEntries() != 3 {
		t.Fatal("m1.NumEntries() != 3")
	}

	if m1 == m0 {
		t.Fatal("current tree is the same as orig tree.")
	}

	if !origM.Equiv(dupM) {
		t.Fatal("TestBasicPutCase1: orig Map and duplicate of orig Map are not identical.")
	}
}

func TestBasicLoadOrStoreTree2(t *testing.T) {
	var m0 = mkmap(
		mknod(10, Black,
			nil,
			mknod(30, Red, nil, nil)))

	var origM = m0
	var dupM = m0.Dup()

	//log.Printf("Before LoadOrStore: m0=\n%s", m0.TreeString())

	var m1, val, found = m0.LoadOrStore(IntKey(20), 20)

	//log.Printf("After LoadOrStore: m1=\n%s", m1.TreeString())

	if found {
		t.Fatal("key was found; it was not expected to be found.")
	}

	if val != nil {
		t.Fatal("val,%d != expected val,%v,", val, nil)
	}

	if m1.NumEntries() != 3 {
		t.Fatal("m1.NumEntries() != 3")
	}

	if m1 == m0 {
		t.Fatal("current tree is the same as orig tree.")
	}

	if !origM.Equiv(dupM) {
		t.Fatal("TestBasicPutCase1: orig Map and duplicate of orig Map are not identical.")
	}
}

func TestBasicLoadOrStoreTree3(t *testing.T) {
	var m0 = mkmap(
		mknod(60, Black,
			mknod(20, Red,
				mknod(10, Black, nil, nil),
				mknod(40, Black,
					nil,
					mknod(50, Red, nil, nil))),
			mknod(80, Black,
				mknod(70, Red, nil, nil),
				mknod(90, Red, nil, nil))))

	if err := m0.Valid(); err != nil {
		t.Fatal("m0 is invalid; err=%s", err)
	}

	var origM = m0
	var dupM = m0.Dup()

	//log.Printf("Before LoadOrStore: m0=\n%s", m0.TreeString())

	var m1, val, found = m0.LoadOrStore(IntKey(30), 30)

	//log.Printf("After LoadOrStore: m1=\n%s", m1.TreeString())

	if found {
		t.Fatal("key was found; it was not expected to be found.")
	}

	if val != nil {
		t.Fatal("val,%d != expected val,%v,", val, nil)
	}

	if m1.NumEntries() != 9 {
		t.Fatal("m1.NumEntries() != 9")
	}

	if m1 == m0 {
		t.Fatal("current tree is the same as orig tree.")
	}

	if !origM.Equiv(dupM) {
		t.Fatal("TestBasicPutCase1: orig Map and duplicate of orig Map are not identical.")
	}
}

func TestBasicPutCase1(t *testing.T) {
	var m = mkmap(nil)

	var origM = m
	var dupM = m.Dup()

	//log.Printf("BEFORE Put m =\n%s", m.TreeString())

	m = m.Put(IntKey(10), 10)

	//log.Printf("AFTER Put m =\n%s", m.TreeString())

	if m.NumEntries() != 1 {
		t.Fatal("m.NumEntries() != 1")
	}

	if err := m.Valid(); err != nil {
		t.Fatalf("map not valid; err=%s", err)
	}

	if !origM.Equiv(dupM) {
		t.Fatal("TestBasicPutCase1: orig Map and duplicate of orig Map are not identical.")
	}
}

func TestBasicPutCase2(t *testing.T) {
	var m = mkmap(
		mknod(20, Black,
			mknod(10, Red, nil, nil),
			nil))

	var origM = m
	var dupM = m.Dup()

	//log.Printf("BEFORE Put m =\n%s", m.TreeString())

	m = m.Put(IntKey(30), 30)

	//log.Printf("AFTER Put m =\n%s", m.TreeString())

	if m.NumEntries() != 3 {
		t.Fatal("m.NumEntries() != 1")
	}

	if err := m.Valid(); err != nil {
		t.Fatalf("map not valid; err=%s", err)
	}

	if !origM.Equiv(dupM) {
		t.Fatal("TestBasicPutCase2: orig Map and duplicate of orig Map are not identical.")
	}
}

func TestBasicPutCase3(t *testing.T) {
	//insert order 10, 20, 50, 40, 30, 60
	var m = mkmap(
		mknod(20, Black,
			mknod(10, Black, nil, nil),
			mknod(40, Black,
				mknod(30, Red, nil, nil),
				mknod(50, Red, nil, nil),
			),
		))

	var origM = m
	var dupM = m.Dup()

	//log.Printf("BEFORE Put m =\n%s", m.TreeString())

	m = m.Put(IntKey(60), 60)

	//log.Printf("AFTER Put m =\n%s", m.TreeString())

	if m.NumEntries() != 6 {
		t.Fatal("m.NumEntries() != 6")
	}

	if err := m.Valid(); err != nil {
		t.Fatalf("map not valid; err=%s", err)
	}

	if !origM.Equiv(dupM) {
		t.Fatal("TestBasicPutCase3: orig Map and duplicate of orig Map are not identical.")
	}
}

func TestBasicPutCase4(t *testing.T) {
	//var m = mkmap(
	//	mknod(7940, Black,
	//		mknod(4930, Black,
	//			nil,
	//			mknod(7100, Red, nil, nil)),
	//		mknod(8090, Black,
	//			nil,
	//			mknod(10050, Red, nil, nil)),
	//	))
	//insert order 50, 20, 60, 40, 70 ???
	var m = mkmap(
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

	//log.Printf("BEFORE Put m =\n%s", m.TreeString())

	//m = m.Put(IntKey(5310), 5310)
	m = m.Put(IntKey(30), 30)

	//log.Printf("AFTER Put m =\n%s", m.TreeString())

	if m.NumEntries() != 6 {
		t.Fatal("m.NumEntries() != 6")
	}

	if err := m.Valid(); err != nil {
		t.Fatalf("map not valid; err=%s", err)
	}

	if !origM.Equiv(dupM) {
		t.Fatal("TestBasicPutCase4: orig Map and duplicate of orig Map are not identical.")
	}
}

func TestBasicDelCase1Tree0(t *testing.T) {
	var m0 = mkmap(
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
	var m0 = mkmap(
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

func TestBasicDelCase1Tree2(t *testing.T) {
	var m0 = mkmap(
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
	var m = mkmap(
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
	var m = mkmap(
		mknod(20, Black,
			mknod(10, Red, nil, nil),
			nil,
		))

	//log.Printf("BEFORE REMOVE: Map m=\n%s", m.TreeString())

	m = m.Del(IntKey(10))

	//log.Printf("AFTER REMOVE Map m=\n%s", m.TreeString())

	if m.NumEntries() != 1 {
		t.Fatalf("m.NumEntries(),%d != 1", m.NumEntries())
	}

	if !m.Root().IsBlack() {
		t.Fatal("!m.Root().IsBlack()")
	}

	if m.Root().Ln() != nil {
		t.Fatal("m.Root().Ln() != nil")
	}

	if m.Root().Rn() != nil {
		t.Fatal("m.Root().Rn() != nil")
	}
}

func TestBasicDelCase3Tree0(t *testing.T) {
	var m = mkmap(
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
	var m = mkmap(
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

	//log.Printf("origMapStr0 =\n%s", origMapStr0)

	m = m.Del(IntKey(30))

	if m.NumEntries() != 5 {
		t.Fatalf("m.NumEntries(),%d != 5", m.NumEntries())
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

func TestBasicDelTwoChildTree0(t *testing.T) {
	var m = mkmap(
		mknod(40, Black,
			mknod(20, Black,
				mknod(10, Red, nil, nil),
				mknod(30, Red, nil, nil)),
			mknod(70, Red,
				mknod(50, Black, nil, nil),
				mknod(80, Black, nil, nil))))

	var shouldHaveKvs = []KeyVal{
		{IntKey(10), 10},
		{IntKey(30), 30},
		{IntKey(40), 40},
		{IntKey(50), 50},
		{IntKey(70), 70},
		{IntKey(80), 80},
	}

	var origM = m
	var dupOrigM = m.Dup()
	var origMapStr0 = m.TreeString()

	//log.Printf("BEFORE DEL m = \n%s", m.TreeString())

	m = m.Del(IntKey(20))

	//log.Printf("AFTER DEL m = \n%s", m.TreeString())

	if err := m.Valid(); err != nil {
		t.Fatalf("INVALID TREE AFTER Del(IntKey(20)); err=%s\n", err)
	}

	if m.NumEntries() != 6 {
		t.Fatalf("m.NumEntries(),%d != 6", m.NumEntries())
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

//deleteCase4
func TestBasicDelTwoChildTree1(t *testing.T) {
	var m = mkmap(
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

	//log.Printf("BEFORE DEL m = \n%s", m.TreeString())

	m = m.Del(IntKey(70))

	//log.Printf("AFTER DEL m = \n%s", m.TreeString())

	if err := m.Valid(); err != nil {
		t.Fatalf("INVALID TREE AFTER Del(IntKey(70)); err=%s\n", err)
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

func TestBasicDelTwoChildTree2(t *testing.T) {
	var m = mkmap(
		mknod(40, Black,
			mknod(10, Black, nil, nil),
			mknod(70, Red,
				mknod(50, Black, nil, nil),
				mknod(80, Black, nil, nil))))

	//shouldHave after Del(70)
	var shouldHaveKvs = []KeyVal{
		{IntKey(10), 10},
		{IntKey(50), 50},
		{IntKey(70), 70},
		{IntKey(80), 80},
	}

	var origM = m
	var dupOrigM = m.Dup()
	var origMapStr0 = m.TreeString()

	//log.Printf("BEFORE DEL m = \n%s", m.TreeString())

	m = m.Del(IntKey(40))

	//log.Printf("AFTER DEL m = \n%s", m.TreeString())

	if err := m.Valid(); err != nil {
		t.Fatalf("INVALID TREE AFTER Del(IntKey(40)); err=%s\n", err)
	}

	if m.NumEntries() != 4 {
		t.Fatalf("m.NumEntries(),%d != 4", m.NumEntries())
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

func TestBasicDelTwoChildTree3(t *testing.T) {
	var m = mkmap(
		mknod(50, Black,
			mknod(20, Red,
				mknod(10, Black, nil, nil),
				mknod(40, Black,
					mknod(30, Red, nil, nil),
					nil)),
			mknod(80, Black, nil, nil)))

	//shouldHave after Del(20)
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

	//log.Printf("BEFORE DEL m = \n%s", m.TreeString())

	m = m.Del(IntKey(20))

	//log.Printf("AFTER DEL m = \n%s", m.TreeString())

	if err := m.Valid(); err != nil {
		t.Fatalf("INVALID TREE AFTER Del(IntKey(20)); err=%s\n", err)
	}

	if m.NumEntries() != 5 {
		t.Fatalf("m.NumEntries(),%d != 5", m.NumEntries())
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

func TestBasicDelTwoChildTree4(t *testing.T) {
	var m = mkmap(
		mknod(60, Black,
			mknod(20, Black,
				mknod(10, Black, nil, nil),
				mknod(40, Red,
					mknod(30, Black, nil, nil),
					mknod(50, Black, nil, nil))),
			mknod(80, Black,
				mknod(70, Black, nil, nil),
				mknod(90, Black, nil, nil))))

	if err := m.Valid(); err != nil {
		t.Fatalf("INVALID TREE; err=%s\n", err)
	}

	//shouldHave after Del(80)
	var shouldHaveKvs = []KeyVal{
		{IntKey(10), 10},
		{IntKey(20), 20},
		{IntKey(30), 30},
		{IntKey(40), 40},
		{IntKey(50), 50},
		{IntKey(60), 60},
		{IntKey(70), 70},
		//{IntKey(80), 80},
		{IntKey(90), 90},
	}

	var origM = m
	var dupOrigM = m.Dup()
	var origMapStr0 = m.TreeString()

	//log.Printf("BEFORE DEL m = \n%s", m.TreeString())

	m = m.Del(IntKey(80))

	//log.Printf("AFTER DEL m = \n%s", m.TreeString())

	if err := m.Valid(); err != nil {
		t.Fatalf("INVALID TREE AFTER Del(IntKey(80)); err=%s\n", err)
	}

	if m.NumEntries() != 8 {
		t.Fatalf("m.NumEntries(),%d != 8", m.NumEntries())
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

//func TestBasicRange(t *testing.T) {
//	var m = mkmap()
//}
