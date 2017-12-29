package sorted_map

import (
	"log"
	"testing"
)

func TestBasicGetPos(t *testing.T) {
	var m = mkmap(
		mknod(20, black,
			mknod(10, red, nil, nil),
			mknod(30, red, nil, nil)))

	var val = m.Get(IntKey(20))

	if val != 20 {
		t.Fatal("m.Get(IntKey(20)) did not return a val==20")
	}
}

func TestBasicGetNeg(t *testing.T) {
	var m = mkmap(
		mknod(20, black,
			mknod(10, red, nil, nil),
			mknod(30, red, nil, nil)))

	var val = m.Get(IntKey(40))

	if val != nil {
		t.Fatal("m.Get(IntKey(40)) did not return a val==nil")
	}
}

func TestBasicLoadOrStoreTree0(t *testing.T) {
	var m0 = mkmap(
		mknod(20, black,
			mknod(10, red, nil, nil),
			mknod(30, red, nil, nil)))

	var origM = m0
	var dupM = m0.dup()

	//log.Printf("Before LoadOrStore: m0=\n%s", m0.treeString())

	var m1, val, found = m0.LoadOrStore(IntKey(10), 10)

	//log.Printf("After LoadOrStore: m1=\n%s", m1.treeString())

	if !found {
		t.Fatal("key not found; it was expected it to be found.")
	}

	if val != 10 {
		t.Fatal("val,%d != expected val,%v,", val, m0.root.ln.val)
	}

	if m1 != m0 {
		t.Fatal("current tree not same as orig tree.")
	}

	if !origM.equiv(dupM) {
		t.Fatal("TestBasicLoadOrStoreTree0: orig Map and duplicate of orig Map are not identical.")
	}
}

func TestBasicLoadOrStoreTree1(t *testing.T) {
	var m0 = mkmap(
		mknod(20, black,
			nil,
			mknod(30, red, nil, nil)))

	var origM = m0
	var dupM = m0.dup()

	//log.Printf("Before LoadOrStore: m0=\n%s", m0.treeString())

	var m1, val, found = m0.LoadOrStore(IntKey(10), 10)

	//log.Printf("After LoadOrStore: m1=\n%s", m1.treeString())

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

	if !origM.equiv(dupM) {
		t.Fatal("TestBasicLoadOrStoreTree1: orig Map and duplicate of orig Map are not identical.")
	}
}

func TestBasicLoadOrStoreTree2(t *testing.T) {
	var m0 = mkmap(
		mknod(10, black,
			nil,
			mknod(30, red, nil, nil)))

	var origM = m0
	var dupM = m0.dup()

	//log.Printf("Before LoadOrStore: m0=\n%s", m0.treeString())

	var m1, val, found = m0.LoadOrStore(IntKey(20), 20)

	//log.Printf("After LoadOrStore: m1=\n%s", m1.treeString())

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

	if !origM.equiv(dupM) {
		t.Fatal("TestBasicLoadOrStoreTree2: orig Map and duplicate of orig Map are not identical.")
	}
}

func TestBasicLoadOrStoreTree3(t *testing.T) {
	var m0 = mkmap(
		mknod(60, black,
			mknod(20, red,
				mknod(10, black, nil, nil),
				mknod(40, black,
					nil,
					mknod(50, red, nil, nil))),
			mknod(80, black,
				mknod(70, red, nil, nil),
				mknod(90, red, nil, nil))))

	if err := m0.valid(); err != nil {
		t.Fatal("m0 is invalid; err=%s", err)
	}

	var origM = m0
	var dupM = m0.dup()

	//log.Printf("Before LoadOrStore: m0=\n%s", m0.treeString())

	var m1, val, found = m0.LoadOrStore(IntKey(30), 30)

	//log.Printf("After LoadOrStore: m1=\n%s", m1.treeString())

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

	if !origM.equiv(dupM) {
		t.Fatal("TestBasicLoadOrStoreTree3: orig Map and duplicate of orig Map are not identical.")
	}
}

func TestBasicStoreReplace(t *testing.T) {
	var m0 = mkmap(
		mknod(20, black,
			mknod(10, red, nil, nil),
			mknod(30, red, nil, nil)))

	var origM0 = m0
	var dupM0 = m0.dup()

	var m1, added = m0.Store(IntKey(30), 31)

	if added {
		t.Fatal("Store added new entry when it should not")
	}

	var val = m1.Get(IntKey(30))
	if val != 31 {
		t.Fatal("new map did not return value of 31 for a lookup of IntKey(31)")
	}

	if !origM0.equiv(dupM0) {
		t.Fatal("TestBasicStoreReplace: orig Map and duplicate of orig Map are not identical.")
	}
}

func TestBasicPutCase1(t *testing.T) {
	var m = mkmap(nil)

	var origM = m
	var dupM = m.dup()

	//log.Printf("BEFORE Put m =\n%s", m.treeString())

	m = m.Put(IntKey(10), 10)

	//log.Printf("AFTER Put m =\n%s", m.treeString())

	if m.NumEntries() != 1 {
		t.Fatal("m.NumEntries() != 1")
	}

	if err := m.valid(); err != nil {
		t.Fatalf("map not valid; err=%s", err)
	}

	if !origM.equiv(dupM) {
		t.Fatal("TestBasicPutCase1: orig Map and duplicate of orig Map are not identical.")
	}
}

func TestBasicPutCase2(t *testing.T) {
	var m = mkmap(
		mknod(20, black,
			mknod(10, red, nil, nil),
			nil))

	var origM = m
	var dupM = m.dup()

	//log.Printf("BEFORE Put m =\n%s", m.treeString())

	m = m.Put(IntKey(30), 30)

	//log.Printf("AFTER Put m =\n%s", m.treeString())

	if m.NumEntries() != 3 {
		t.Fatal("m.NumEntries() != 1")
	}

	if err := m.valid(); err != nil {
		t.Fatalf("map not valid; err=%s", err)
	}

	if !origM.equiv(dupM) {
		t.Fatal("TestBasicPutCase2: orig Map and duplicate of orig Map are not identical.")
	}
}

func TestBasicPutCase3(t *testing.T) {
	//insert order 10, 20, 50, 40, 30, 60
	var m = mkmap(
		mknod(20, black,
			mknod(10, black, nil, nil),
			mknod(40, black,
				mknod(30, red, nil, nil),
				mknod(50, red, nil, nil),
			),
		))

	var origM = m
	var dupM = m.dup()

	//log.Printf("BEFORE Put m =\n%s", m.treeString())

	m = m.Put(IntKey(60), 60)

	//log.Printf("AFTER Put m =\n%s", m.treeString())

	if m.NumEntries() != 6 {
		t.Fatal("m.NumEntries() != 6")
	}

	if err := m.valid(); err != nil {
		t.Fatalf("map not valid; err=%s", err)
	}

	if !origM.equiv(dupM) {
		t.Fatal("TestBasicPutCase3: orig Map and duplicate of orig Map are not identical.")
	}
}

func TestBasicPutCase4(t *testing.T) {
	//var m = mkmap(
	//	mknod(7940, black,
	//		mknod(4930, black,
	//			nil,
	//			mknod(7100, red, nil, nil)),
	//		mknod(8090, black,
	//			nil,
	//			mknod(10050, red, nil, nil)),
	//	))
	//insert order 50, 20, 60, 40, 70 ???
	var m = mkmap(
		mknod(50, black,
			mknod(20, black,
				nil,
				mknod(40, red, nil, nil)),
			mknod(60, black,
				nil,
				mknod(70, red, nil, nil)),
		))

	var origM = m      //copy the pointer
	var dupM = m.dup() //copy the value

	//log.Printf("BEFORE Put m =\n%s", m.treeString())

	//m = m.Put(IntKey(5310), 5310)
	m = m.Put(IntKey(30), 30)

	//log.Printf("AFTER Put m =\n%s", m.treeString())

	if m.NumEntries() != 6 {
		t.Fatal("m.NumEntries() != 6")
	}

	if err := m.valid(); err != nil {
		t.Fatalf("map not valid; err=%s", err)
	}

	if !origM.equiv(dupM) {
		t.Fatal("TestBasicPutCase4: orig Map and duplicate of orig Map are not identical.")
	}
}

func TestBasicRemoveNeg(t *testing.T) {
	var m0 = mkmap(
		mknod(20, black,
			mknod(10, red, nil, nil),
			mknod(30, red, nil, nil)))

	var m1, val, found = m0.Remove(IntKey(40))

	if found {
		t.Fatal("found a key that does not exist")
	}

	if val != nil {
		t.Fatal("val != nil")
	}

	if m1 != m0 {
		t.Fatal("returned map not the same as the original map")
	}
}

func TestBasicDelCase1Tree0(t *testing.T) {
	var m0 = mkmap(
		mknod(10, black, nil, nil))

	var then = m0.treeString()
	//var dupM0 = m0.dup()

	var m1 = m0.Del(IntKey(10))

	if m1.NumEntries() != 0 {
		t.Fatal("m.NumEntries() != 0")
	}

	var now = m0.treeString()
	if then != now {
		log.Printf("origninal tree changeed:\nTHEN: %s\nNOW: %s",
			then, now)
		t.Fatal("The original tree changed.")
	}

	//if !m0.equiv(dupM0) {
	//	t.Fatal("The original tree changed.")
	//}
}

func TestBasicDelCase1Tree1(t *testing.T) {
	var m0 = mkmap(
		mknod(10, black,
			nil,
			mknod(20, red, nil, nil),
		))

	var then = m0.treeString()

	var m1 = m0.Del(IntKey(10))

	if m1.NumEntries() != 1 {
		t.Fatal("m.NumEntries() != 1")
	}

	var now = m0.treeString()
	if then != now {
		log.Printf("origninal tree changeed:\nTHEN: %s\nNOW: %s",
			then, now)
		t.Fatal("The original tree changed.")
	}
}

func TestBasicDelCase1Tree2(t *testing.T) {
	var m0 = mkmap(
		mknod(20, black,
			mknod(10, red, nil, nil),
			nil,
		))

	var then = m0.treeString()

	var m1 = m0.Del(IntKey(20))

	if m1.NumEntries() != 1 {
		t.Fatal("m.NumEntries() != 1")
	}

	var now = m0.treeString()
	if then != now {
		log.Printf("origninal tree changeed:\nTHEN: %s\nNOW: %s",
			then, now)
		t.Fatal("The original tree changed.")
	}
}

// DeleteCase1 is exhaustively tested.

func TestBasicDelCase2Tree0(t *testing.T) {
	var m = mkmap(
		mknod(20, black,
			nil,
			mknod(30, red, nil, nil),
		))

	//log.Printf("BEFORE REMOVE: Map m=\n%s", m.treeString())

	m = m.Del(IntKey(30))

	//log.Printf("AFTER REMOVE Map m=\n%s", m.treeString())

	if m.NumEntries() != 1 {
		t.Fatalf("m.NumEntries(),%d != 1", m.NumEntries())
	}

	if !m.root.isBlack() {
		t.Fatal("!m.root.isBlack()")
	}

	if m.root.ln != nil {
		t.Fatal("m.root.rn != nil")
	}

	if m.root.rn != nil {
		t.Fatal("m.root.ln != nil")
	}
}

func TestBasicDelCase2Tree1(t *testing.T) {
	var m = mkmap(
		mknod(20, black,
			mknod(10, red, nil, nil),
			nil,
		))

	//log.Printf("BEFORE REMOVE: Map m=\n%s", m.treeString())

	m = m.Del(IntKey(10))

	//log.Printf("AFTER REMOVE Map m=\n%s", m.treeString())

	if m.NumEntries() != 1 {
		t.Fatalf("m.NumEntries(),%d != 1", m.NumEntries())
	}

	if !m.root.isBlack() {
		t.Fatal("!m.root.isBlack()")
	}

	if m.root.ln != nil {
		t.Fatal("m.root.ln != nil")
	}

	if m.root.rn != nil {
		t.Fatal("m.root.rn != nil")
	}
}

func TestBasicDelCase3Tree0(t *testing.T) {
	var m = mkmap(
		mknod(20, black,
			mknod(10, black, nil, nil),
			mknod(30, black, nil, nil),
		))

	m = m.Del(IntKey(30))
	if m.NumEntries() != 2 {
		t.Fatalf("m.NumEntries(),%d != 2", m.NumEntries())
	}

	if !m.root.isBlack() {
		t.Fatal("!m.root.isBlack()")
	}

	if !m.root.ln.isRed() {
		t.Fatal("!m.root.ln.isRed()")
	}

	if m.root.rn != nil {
		t.Fatal("m.root.rn != nil")
	}
}

func TestBasicDelCase6Tree0(t *testing.T) {
	var m = mkmap(
		mknod(20, black,
			mknod(10, black, nil, nil),
			mknod(40, red,
				mknod(30, black, nil, nil),
				mknod(50, black,
					nil,
					mknod(60, red, nil, nil)))))

	var origM = m
	var dupOrigM = m.dup()
	var origMapStr0 = m.treeString()

	//log.Printf("origMapStr0 =\n%s", origMapStr0)

	m = m.Del(IntKey(30))

	if m.NumEntries() != 5 {
		t.Fatalf("m.NumEntries(),%d != 5", m.NumEntries())
	}

	var origMapStr1 = origM.treeString()
	if origMapStr0 != origMapStr1 {
		log.Printf("origMapStr0 != origMapStr1:\n"+
			"origMapStr0=\n%s\norigMapStr1=\n%s", origMapStr0, origMapStr1)
	}

	if !origM.equiv(dupOrigM) {
		t.Fatal("TestBasicPutCase4: orig Map and duplicate of orig Map are not identical.")
	}
}

func TestBasicDelTwoChildTree0(t *testing.T) {
	var m = mkmap(
		mknod(40, black,
			mknod(20, black,
				mknod(10, red, nil, nil),
				mknod(30, red, nil, nil)),
			mknod(70, red,
				mknod(50, black, nil, nil),
				mknod(80, black, nil, nil))))

	var shouldHaveKvs = []KeyVal{
		{IntKey(10), 10},
		{IntKey(30), 30},
		{IntKey(40), 40},
		{IntKey(50), 50},
		{IntKey(70), 70},
		{IntKey(80), 80},
	}

	var origM = m
	var dupOrigM = m.dup()
	var origMapStr0 = m.treeString()

	//log.Printf("BEFORE DEL m = \n%s", m.treeString())

	m = m.Del(IntKey(20))

	//log.Printf("AFTER DEL m = \n%s", m.treeString())

	if err := m.valid(); err != nil {
		t.Fatalf("INVALID TREE AFTER Del(IntKey(20)); err=%s\n", err)
	}

	if m.NumEntries() != 6 {
		t.Fatalf("m.NumEntries(),%d != 6", m.NumEntries())
	}

	for _, kv := range shouldHaveKvs {
		var val, found = m.Load(kv.Key)
		if !found {
			log.Printf("failed to find shouldHave key=%s in Tree=\n%s",
				kv.Key, m.treeString())
			t.Fatalf("Failed to find shouldHave kv.Key=%s", kv.Key)
		}
		if val != kv.Val {
			t.Fatalf("found val,%v != expected val,%v", val, kv.Val)
		}
	}

	var origMapStr1 = origM.treeString()
	if origMapStr0 != origMapStr1 {
		log.Printf("origMapStr0 != origMapStr1:\n"+
			"origMapStr0=\n%s\norigMapStr1=\n%s", origMapStr0, origMapStr1)
	}

	if !origM.equiv(dupOrigM) {
		t.Fatal("TestBasicDelTwoChildrenCase2: " +
			"orig Map and duplicate of orig Map are not identical.")
	}
}

//deleteCase4
func TestBasicDelTwoChildTree1(t *testing.T) {
	var m = mkmap(
		mknod(40, black,
			mknod(10, black,
				nil,
				mknod(30, red, nil, nil)),
			mknod(70, red,
				mknod(50, black, nil, nil),
				mknod(80, black, nil, nil))))

	//shouldHave after Del(70)
	var shouldHaveKvs = []KeyVal{
		{IntKey(10), 10},
		{IntKey(30), 30},
		{IntKey(40), 40},
		{IntKey(50), 50},
		{IntKey(80), 80},
	}

	var origM = m
	var dupOrigM = m.dup()
	var origMapStr0 = m.treeString()

	//log.Printf("BEFORE DEL m = \n%s", m.treeString())

	m = m.Del(IntKey(70))

	//log.Printf("AFTER DEL m = \n%s", m.treeString())

	if err := m.valid(); err != nil {
		t.Fatalf("INVALID TREE AFTER Del(IntKey(70)); err=%s\n", err)
	}

	if m.NumEntries() != 5 {
		t.Fatal("m.NumEntries(),%d != 5", m.NumEntries())
	}

	for _, kv := range shouldHaveKvs {
		var val, found = m.Load(kv.Key)
		if !found {
			log.Printf("failed to find shouldHave key=%s in Tree=\n%s",
				kv.Key, m.treeString())
			t.Fatalf("Failed to find shouldHave kv.Key=%s", kv.Key)
		}
		if val != kv.Val {
			t.Fatalf("found val,%v != expected val,%v", val, kv.Val)
		}
	}

	var origMapStr1 = origM.treeString()
	if origMapStr0 != origMapStr1 {
		log.Printf("origMapStr0 != origMapStr1:\n"+
			"origMapStr0=\n%s\norigMapStr1=\n%s", origMapStr0, origMapStr1)
	}

	if !origM.equiv(dupOrigM) {
		t.Fatal("TestBasicDelTwoChildrenCase2: " +
			"orig Map and duplicate of orig Map are not identical.")
	}
}

func TestBasicDelTwoChildTree2(t *testing.T) {
	var m = mkmap(
		mknod(40, black,
			mknod(10, black, nil, nil),
			mknod(70, red,
				mknod(50, black, nil, nil),
				mknod(80, black, nil, nil))))

	//shouldHave after Del(70)
	var shouldHaveKvs = []KeyVal{
		{IntKey(10), 10},
		{IntKey(50), 50},
		{IntKey(70), 70},
		{IntKey(80), 80},
	}

	var origM = m
	var dupOrigM = m.dup()
	var origMapStr0 = m.treeString()

	//log.Printf("BEFORE DEL m = \n%s", m.treeString())

	m = m.Del(IntKey(40))

	//log.Printf("AFTER DEL m = \n%s", m.treeString())

	if err := m.valid(); err != nil {
		t.Fatalf("INVALID TREE AFTER Del(IntKey(40)); err=%s\n", err)
	}

	if m.NumEntries() != 4 {
		t.Fatalf("m.NumEntries(),%d != 4", m.NumEntries())
	}

	for _, kv := range shouldHaveKvs {
		var val, found = m.Load(kv.Key)
		if !found {
			log.Printf("failed to find shouldHave key=%s in Tree=\n%s",
				kv.Key, m.treeString())
			t.Fatalf("Failed to find shouldHave kv.Key=%s", kv.Key)
		}
		if val != kv.Val {
			t.Fatalf("found val,%v != expected val,%v", val, kv.Val)
		}
	}

	var origMapStr1 = origM.treeString()
	if origMapStr0 != origMapStr1 {
		log.Printf("origMapStr0 != origMapStr1:\n"+
			"origMapStr0=\n%s\norigMapStr1=\n%s", origMapStr0, origMapStr1)
	}

	if !origM.equiv(dupOrigM) {
		t.Fatal("TestBasicDelTwoChildrenCase2: " +
			"orig Map and duplicate of orig Map are not identical.")
	}
}

func TestBasicDelTwoChildTree3(t *testing.T) {
	var m = mkmap(
		mknod(50, black,
			mknod(20, red,
				mknod(10, black, nil, nil),
				mknod(40, black,
					mknod(30, red, nil, nil),
					nil)),
			mknod(80, black, nil, nil)))

	//shouldHave after Del(20)
	var shouldHaveKvs = []KeyVal{
		{IntKey(10), 10},
		{IntKey(30), 30},
		{IntKey(40), 40},
		{IntKey(50), 50},
		{IntKey(80), 80},
	}

	var origM = m
	var dupOrigM = m.dup()
	var origMapStr0 = m.treeString()

	//log.Printf("BEFORE DEL m = \n%s", m.treeString())

	m = m.Del(IntKey(20))

	//log.Printf("AFTER DEL m = \n%s", m.treeString())

	if err := m.valid(); err != nil {
		t.Fatalf("INVALID TREE AFTER Del(IntKey(20)); err=%s\n", err)
	}

	if m.NumEntries() != 5 {
		t.Fatalf("m.NumEntries(),%d != 5", m.NumEntries())
	}

	for _, kv := range shouldHaveKvs {
		var val, found = m.Load(kv.Key)
		if !found {
			log.Printf("failed to find shouldHave key=%s in Tree=\n%s",
				kv.Key, m.treeString())
			t.Fatalf("Failed to find shouldHave kv.Key=%s", kv.Key)
		}
		if val != kv.Val {
			t.Fatalf("found val,%v != expected val,%v", val, kv.Val)
		}
	}

	var origMapStr1 = origM.treeString()
	if origMapStr0 != origMapStr1 {
		log.Printf("origMapStr0 != origMapStr1:\n"+
			"origMapStr0=\n%s\norigMapStr1=\n%s", origMapStr0, origMapStr1)
	}

	if !origM.equiv(dupOrigM) {
		t.Fatal("TestBasicDelTwoChildrenCase2: " +
			"orig Map and duplicate of orig Map are not identical.")
	}
}

func TestBasicDelTwoChildTree4(t *testing.T) {
	var m = mkmap(
		mknod(60, black,
			mknod(20, black,
				mknod(10, black, nil, nil),
				mknod(40, red,
					mknod(30, black, nil, nil),
					mknod(50, black, nil, nil))),
			mknod(80, black,
				mknod(70, black, nil, nil),
				mknod(90, black, nil, nil))))

	if err := m.valid(); err != nil {
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
	var dupOrigM = m.dup()
	var origMapStr0 = m.treeString()

	//log.Printf("BEFORE DEL m = \n%s", m.treeString())

	m = m.Del(IntKey(80))

	//log.Printf("AFTER DEL m = \n%s", m.treeString())

	if err := m.valid(); err != nil {
		t.Fatalf("INVALID TREE AFTER Del(IntKey(80)); err=%s\n", err)
	}

	if m.NumEntries() != 8 {
		t.Fatalf("m.NumEntries(),%d != 8", m.NumEntries())
	}

	for _, kv := range shouldHaveKvs {
		var val, found = m.Load(kv.Key)
		if !found {
			log.Printf("failed to find shouldHave key=%s in Tree=\n%s",
				kv.Key, m.treeString())
			t.Fatalf("Failed to find shouldHave kv.Key=%s", kv.Key)
		}
		if val != kv.Val {
			t.Fatalf("found val,%v != expected val,%v", val, kv.Val)
		}
	}

	var origMapStr1 = origM.treeString()
	if origMapStr0 != origMapStr1 {
		log.Printf("origMapStr0 != origMapStr1:\n"+
			"origMapStr0=\n%s\norigMapStr1=\n%s", origMapStr0, origMapStr1)
	}

	if !origM.equiv(dupOrigM) {
		t.Fatal("TestBasicDelTwoChildrenCase2: " +
			"orig Map and duplicate of orig Map are not identical.")
	}
}

func TestBasicRange(t *testing.T) {
	var m = mkmap(
		mknod(60, black,
			mknod(20, black,
				mknod(10, black, nil, nil),
				mknod(40, black,
					mknod(30, red, nil, nil),
					mknod(50, red, nil, nil))),
			mknod(100, black,
				mknod(80, black,
					mknod(70, red, nil, nil),
					mknod(90, red, nil, nil)),
				mknod(120, black,
					mknod(110, red, nil, nil),
					mknod(130, red, nil, nil)))))

	if err := m.valid(); err != nil {
		t.Fatalf("INVALID TREE; err=%s\n", err)
	}

	var shouldHaveKvs = []KeyVal{
		{IntKey(10), 10},
		{IntKey(20), 20},
		{IntKey(30), 30},
		{IntKey(40), 40},
		{IntKey(50), 50},
		{IntKey(60), 60},
		{IntKey(70), 70},
		{IntKey(80), 80},
		{IntKey(90), 90},
		{IntKey(100), 100},
		{IntKey(110), 110},
		{IntKey(120), 120},
		{IntKey(130), 130},
	}

	var i int
	var fn = func(k0 MapKey, v0 interface{}) bool {
		var k1 = shouldHaveKvs[i].Key
		var v1 = shouldHaveKvs[i].Val
		//log.Printf("k0=%s; v0=%v;", k0, v0)
		//log.Printf("k1=%s; v0=%v;", k1, v1)
		if k0.Less(k1) || k1.Less(k0) {
			t.Fatalf("i=%d; k0,%s != k1,%s", i, k0, k1)
		}
		if v0 != v1 {
			t.Fatalf("i=%d; v0,%d != v1,%d", i, v0, v1)
		}
		i++
		return true
	}
	m.Range(fn)
}

func TestBasicRangeForwBeg(t *testing.T) {
	var m = mkmap(
		mknod(60, black,
			mknod(20, black,
				mknod(10, black, nil, nil),
				mknod(40, black,
					mknod(30, red, nil, nil),
					mknod(50, red, nil, nil))),
			mknod(100, black,
				mknod(80, black,
					mknod(70, red, nil, nil),
					mknod(90, red, nil, nil)),
				mknod(120, black,
					mknod(110, red, nil, nil),
					mknod(130, red, nil, nil)))))

	if err := m.valid(); err != nil {
		t.Fatalf("INVALID TREE; err=%s\n", err)
	}

	var shouldHaveKvs = []KeyVal{
		{IntKey(10), 10},
		{IntKey(20), 20},
		{IntKey(30), 30},
		{IntKey(40), 40},
		{IntKey(50), 50},
		{IntKey(60), 60},
		{IntKey(70), 70},
		{IntKey(80), 80},
		{IntKey(90), 90},
		{IntKey(100), 100},
		{IntKey(110), 110},
		{IntKey(120), 120},
		{IntKey(130), 130},
	}

	//var numKeys = len(shouldHaveKvs)
	var eltOffset = 3
	var startKey = IntKey(eltOffset * 10)
	var endKey = InfKey(1) //positive infinity
	var keyRange = shouldHaveKvs[eltOffset-1:]
	var i = 0
	var fn = func(k0 MapKey, v0 interface{}) bool {
		if i >= len(keyRange) {
			t.Fatalf("i,%d >= len(keyRange),%d", i, len(keyRange))
		}
		var k1 = keyRange[i].Key
		var v1 = keyRange[i].Val
		//log.Printf("k0=%s; v0=%v;", k0, v0)
		//log.Printf("k1=%s; v0=%v;", k1, v1)
		if k0.Less(k1) || k1.Less(k0) {
			t.Fatalf("i=%d; k0,%s != k1,%s", i, k0, k1)
		}
		if v0 != v1 {
			t.Fatalf("i=%d; v0,%d != v1,%d", i, v0, v1)
		}
		i++
		return true
	}
	m.RangeLimit(startKey, endKey, fn)
	if i != len(keyRange) {
		t.Fatalf("after RangeLimit: i,%d != len(keyRange),%d", i, len(keyRange))
	}
}

func TestBasicRangeForwBegInexact(t *testing.T) {
	var m = mkmap(
		mknod(60, black,
			mknod(20, black,
				mknod(10, black, nil, nil),
				mknod(40, black,
					mknod(30, red, nil, nil),
					mknod(50, red, nil, nil))),
			mknod(100, black,
				mknod(80, black,
					mknod(70, red, nil, nil),
					mknod(90, red, nil, nil)),
				mknod(120, black,
					mknod(110, red, nil, nil),
					mknod(130, red, nil, nil)))))

	if err := m.valid(); err != nil {
		t.Fatalf("INVALID TREE; err=%s\n", err)
	}

	var shouldHaveKvs = []KeyVal{
		{IntKey(10), 10},
		{IntKey(20), 20},
		{IntKey(30), 30},
		{IntKey(40), 40},
		{IntKey(50), 50},
		{IntKey(60), 60},
		{IntKey(70), 70},
		{IntKey(80), 80},
		{IntKey(90), 90},
		{IntKey(100), 100},
		{IntKey(110), 110},
		{IntKey(120), 120},
		{IntKey(130), 130},
	}

	//var numKeys = len(shouldHaveKvs)
	var eltOffset = 2
	var startKey = IntKey((eltOffset * 10) - 5) //IntKey(15)
	var endKey = InfKey(1)                      //positive infinity
	var keyRange = shouldHaveKvs[eltOffset-1:]
	var i = 0
	var fn = func(k0 MapKey, v0 interface{}) bool {
		if i >= len(keyRange) {
			t.Fatalf("i,%d >= len(keyRange),%d", i, len(keyRange))
		}
		var k1 = keyRange[i].Key
		var v1 = keyRange[i].Val
		//log.Printf("k0=%s; v0=%v;", k0, v0)
		//log.Printf("k1=%s; v0=%v;", k1, v1)
		if k0.Less(k1) || k1.Less(k0) {
			t.Fatalf("i=%d; k0,%s != k1,%s", i, k0, k1)
		}
		if v0 != v1 {
			t.Fatalf("i=%d; v0,%d != v1,%d", i, v0, v1)
		}
		i++
		return true
	}
	m.RangeLimit(startKey, endKey, fn)
	if i != len(keyRange) {
		t.Fatalf("after RangeLimit: i,%d != len(keyRange),%d", i, len(keyRange))
	}
}

func TestBasicRangeForwEnd(t *testing.T) {
	var m = mkmap(
		mknod(60, black,
			mknod(20, black,
				mknod(10, black, nil, nil),
				mknod(40, black,
					mknod(30, red, nil, nil),
					mknod(50, red, nil, nil))),
			mknod(100, black,
				mknod(80, black,
					mknod(70, red, nil, nil),
					mknod(90, red, nil, nil)),
				mknod(120, black,
					mknod(110, red, nil, nil),
					mknod(130, red, nil, nil)))))

	if err := m.valid(); err != nil {
		t.Fatalf("INVALID TREE; err=%s\n", err)
	}

	var shouldHaveKvs = []KeyVal{
		{IntKey(10), 10},
		{IntKey(20), 20},
		{IntKey(30), 30},
		{IntKey(40), 40},
		{IntKey(50), 50},
		{IntKey(60), 60},
		{IntKey(70), 70},
		{IntKey(80), 80},
		{IntKey(90), 90},
		{IntKey(100), 100},
		{IntKey(110), 110},
		{IntKey(120), 120},
		{IntKey(130), 130},
	}

	var numKeys = len(shouldHaveKvs)
	var eltOffset = 3
	var startKey = InfKey(-1)                           //negative infinity
	var endKey = IntKey((numKeys - eltOffset) * 10)     //IntKey(100)
	var keyRange = shouldHaveKvs[:len(shouldHaveKvs)-3] //??
	var i = 0
	var fn = func(k0 MapKey, v0 interface{}) bool {
		if i >= len(keyRange) {
			t.Fatalf("i,%d >= len(keyRange),%d", i, len(keyRange))
		}
		var k1 = keyRange[i].Key
		var v1 = keyRange[i].Val
		//log.Printf("k0=%s; v0=%v;", k0, v0)
		//log.Printf("k1=%s; v0=%v;", k1, v1)
		if k0.Less(k1) || k1.Less(k0) {
			t.Fatalf("i=%d; k0,%s != k1,%s", i, k0, k1)
		}
		if v0 != v1 {
			t.Fatalf("i=%d; v0,%d != v1,%d", i, v0, v1)
		}
		i++
		return true
	}
	m.RangeLimit(startKey, endKey, fn)
	if i != len(keyRange) {
		t.Fatalf("after RangeLimit: i,%d != len(keyRange),%d", i, len(keyRange))
	}
}

func TestBasicRangeForwEndInexact(t *testing.T) {
	var m = mkmap(
		mknod(60, black,
			mknod(20, black,
				mknod(10, black, nil, nil),
				mknod(40, black,
					mknod(30, red, nil, nil),
					mknod(50, red, nil, nil))),
			mknod(100, black,
				mknod(80, black,
					mknod(70, red, nil, nil),
					mknod(90, red, nil, nil)),
				mknod(120, black,
					mknod(110, red, nil, nil),
					mknod(130, red, nil, nil)))))

	if err := m.valid(); err != nil {
		t.Fatalf("INVALID TREE; err=%s\n", err)
	}

	var shouldHaveKvs = []KeyVal{
		{IntKey(10), 10},
		{IntKey(20), 20},
		{IntKey(30), 30},
		{IntKey(40), 40},
		{IntKey(50), 50},
		{IntKey(60), 60},
		{IntKey(70), 70},
		{IntKey(80), 80},
		{IntKey(90), 90},
		{IntKey(100), 100},
		{IntKey(110), 110},
		{IntKey(120), 120},
		{IntKey(130), 130},
	}

	var numKeys = len(shouldHaveKvs)
	var eltOffset = 3
	var startKey = InfKey(-1)                           //negative infinity
	var endKey = IntKey((numKeys-eltOffset)*10 + 5)     //IntKey(105)
	var keyRange = shouldHaveKvs[:len(shouldHaveKvs)-3] //??
	var i = 0
	var fn = func(k0 MapKey, v0 interface{}) bool {
		if i >= len(keyRange) {
			t.Fatalf("i,%d >= len(keyRange),%d", i, len(keyRange))
		}
		var k1 = keyRange[i].Key
		var v1 = keyRange[i].Val
		//log.Printf("k0=%s; v0=%v;", k0, v0)
		//log.Printf("k1=%s; v0=%v;", k1, v1)
		if k0.Less(k1) || k1.Less(k0) {
			t.Fatalf("i=%d; k0,%s != k1,%s", i, k0, k1)
		}
		if v0 != v1 {
			t.Fatalf("i=%d; v0,%d != v1,%d", i, v0, v1)
		}
		i++
		return true
	}
	m.RangeLimit(startKey, endKey, fn)
	if i != len(keyRange) {
		t.Fatalf("after RangeLimit: i,%d != len(keyRange),%d", i, len(keyRange))
	}
}

func TestBasicRangeRevBeg(t *testing.T) {
	var m = mkmap(
		mknod(60, black,
			mknod(20, black,
				mknod(10, black, nil, nil),
				mknod(40, black,
					mknod(30, red, nil, nil),
					mknod(50, red, nil, nil))),
			mknod(100, black,
				mknod(80, black,
					mknod(70, red, nil, nil),
					mknod(90, red, nil, nil)),
				mknod(120, black,
					mknod(110, red, nil, nil),
					mknod(130, red, nil, nil)))))

	if err := m.valid(); err != nil {
		t.Fatalf("INVALID TREE; err=%s\n", err)
	}

	var shouldHaveKvs = []KeyVal{
		{IntKey(10), 10},
		{IntKey(20), 20},
		{IntKey(30), 30},
		{IntKey(40), 40},
		{IntKey(50), 50},
		{IntKey(60), 60},
		{IntKey(70), 70},
		{IntKey(80), 80},
		{IntKey(90), 90},
		{IntKey(100), 100},
		{IntKey(110), 110},
		{IntKey(120), 120},
		{IntKey(130), 130},
	}

	var numKeys = len(shouldHaveKvs)
	var eltOffset = 3
	var keyRange = shouldHaveKvs[:numKeys-eltOffset]
	var startKey = IntKey((numKeys - eltOffset) * 10) //IntKey(100)
	var endKey = InfKey(-1)                           //negative infinity
	var i = len(keyRange) - 1
	var fn = func(k0 MapKey, v0 interface{}) bool {
		if i < 0 {
			t.Fatalf("i,%d < 0", i)
		}
		var k1 = keyRange[i].Key
		var v1 = keyRange[i].Val
		//log.Printf("k0=%s; v0=%v;", k0, v0)
		//log.Printf("k1=%s; v0=%v;", k1, v1)
		if k0.Less(k1) || k1.Less(k0) {
			t.Fatalf("i=%d; k0,%s != k1,%s", i, k0, k1)
		}
		if v0 != v1 {
			t.Fatalf("i=%d; v0,%d != v1,%d", i, v0, v1)
		}
		i--
		return true
	}
	//log.Printf("m.RangeLimit(startKey,%s, endKey,%s, fn)", startKey, endKey)
	m.RangeLimit(startKey, endKey, fn)
	if i != -1 {
		t.Fatalf("after RangeLimit: i,%d != -1", i)
	}
}

func TestBasicRangeRevBegInexact(t *testing.T) {
	var m = mkmap(
		mknod(60, black,
			mknod(20, black,
				mknod(10, black, nil, nil),
				mknod(40, black,
					mknod(30, red, nil, nil),
					mknod(50, red, nil, nil))),
			mknod(100, black,
				mknod(80, black,
					mknod(70, red, nil, nil),
					mknod(90, red, nil, nil)),
				mknod(120, black,
					mknod(110, red, nil, nil),
					mknod(130, red, nil, nil)))))

	if err := m.valid(); err != nil {
		t.Fatalf("INVALID TREE; err=%s\n", err)
	}

	var shouldHaveKvs = []KeyVal{
		{IntKey(10), 10},
		{IntKey(20), 20},
		{IntKey(30), 30},
		{IntKey(40), 40},
		{IntKey(50), 50},
		{IntKey(60), 60},
		{IntKey(70), 70},
		{IntKey(80), 80},
		{IntKey(90), 90},
		{IntKey(100), 100},
		{IntKey(110), 110},
		{IntKey(120), 120},
		{IntKey(130), 130},
	}

	var numKeys = len(shouldHaveKvs)
	var eltOffset = 3
	var keyRange = shouldHaveKvs[:numKeys-eltOffset]
	var startKey = IntKey((numKeys-eltOffset)*10 + 5) //IntKey(105)
	var endKey = InfKey(-1)                           //negative infinity
	var i = len(keyRange) - 1
	var fn = func(k0 MapKey, v0 interface{}) bool {
		if i < 0 {
			t.Fatalf("i,%d < 0", i)
		}
		var k1 = keyRange[i].Key
		var v1 = keyRange[i].Val
		//log.Printf("k0=%s; v0=%v;", k0, v0)
		//log.Printf("k1=%s; v0=%v;", k1, v1)
		if k0.Less(k1) || k1.Less(k0) {
			t.Fatalf("i=%d; k0,%s != k1,%s", i, k0, k1)
		}
		if v0 != v1 {
			t.Fatalf("i=%d; v0,%d != v1,%d", i, v0, v1)
		}
		i--
		return true
	}
	//log.Printf("m.RangeLimit(startKey,%s, endKey,%s, fn)", startKey, endKey)
	m.RangeLimit(startKey, endKey, fn)
	if i != -1 {
		t.Fatalf("after RangeLimit: i,%d != -1", i)
	}
}

func TestBasicRangeRevEnd(t *testing.T) {
	var m = mkmap(
		mknod(60, black,
			mknod(20, black,
				mknod(10, black, nil, nil),
				mknod(40, black,
					mknod(30, red, nil, nil),
					mknod(50, red, nil, nil))),
			mknod(100, black,
				mknod(80, black,
					mknod(70, red, nil, nil),
					mknod(90, red, nil, nil)),
				mknod(120, black,
					mknod(110, red, nil, nil),
					mknod(130, red, nil, nil)))))

	if err := m.valid(); err != nil {
		t.Fatalf("INVALID TREE; err=%s\n", err)
	}

	var shouldHaveKvs = []KeyVal{
		{IntKey(10), 10},
		{IntKey(20), 20},
		{IntKey(30), 30},
		{IntKey(40), 40},
		{IntKey(50), 50},
		{IntKey(60), 60},
		{IntKey(70), 70},
		{IntKey(80), 80},
		{IntKey(90), 90},
		{IntKey(100), 100},
		{IntKey(110), 110},
		{IntKey(120), 120},
		{IntKey(130), 130},
	}

	//var numKeys = len(shouldHaveKvs)
	var eltOffset = 3
	var startKey = InfKey(1)                   //positive infinity
	var endKey = IntKey(eltOffset * 10)        //IntKey(30)
	var keyRange = shouldHaveKvs[eltOffset-1:] //??
	var i = len(keyRange) - 1
	var fn = func(k0 MapKey, v0 interface{}) bool {
		if i < 0 {
			t.Fatalf("i,%d < 0", i)
		}
		var k1 = keyRange[i].Key
		var v1 = keyRange[i].Val
		//log.Printf("k0=%s; v0=%v;", k0, v0)
		//log.Printf("k1=%s; v0=%v;", k1, v1)
		if k0.Less(k1) || k1.Less(k0) {
			t.Fatalf("i=%d; k0,%s != k1,%s", i, k0, k1)
		}
		if v0 != v1 {
			t.Fatalf("i=%d; v0,%d != v1,%d", i, v0, v1)
		}
		i--
		return true
	}
	m.RangeLimit(startKey, endKey, fn)
	if i != -1 {
		t.Fatalf("after RangeLimit: i,%d != -1", i)
	}
}

func TestBasicRangeRevEndInexact(t *testing.T) {
	var m = mkmap(
		mknod(60, black,
			mknod(20, black,
				mknod(10, black, nil, nil),
				mknod(40, black,
					mknod(30, red, nil, nil),
					mknod(50, red, nil, nil))),
			mknod(100, black,
				mknod(80, black,
					mknod(70, red, nil, nil),
					mknod(90, red, nil, nil)),
				mknod(120, black,
					mknod(110, red, nil, nil),
					mknod(130, red, nil, nil)))))

	if err := m.valid(); err != nil {
		t.Fatalf("INVALID TREE; err=%s\n", err)
	}

	var shouldHaveKvs = []KeyVal{
		{IntKey(10), 10},
		{IntKey(20), 20},
		{IntKey(30), 30},
		{IntKey(40), 40},
		{IntKey(50), 50},
		{IntKey(60), 60},
		{IntKey(70), 70},
		{IntKey(80), 80},
		{IntKey(90), 90},
		{IntKey(100), 100},
		{IntKey(110), 110},
		{IntKey(120), 120},
		{IntKey(130), 130},
	}

	//var numKeys = len(shouldHaveKvs)
	var eltOffset = 3
	var startKey = InfKey(1)                   //positive infinity
	var endKey = IntKey(eltOffset*10 - 5)      //IntKey(25)
	var keyRange = shouldHaveKvs[eltOffset-1:] //??
	var i = len(keyRange) - 1
	var fn = func(k0 MapKey, v0 interface{}) bool {
		if i < 0 {
			t.Fatalf("i,%d < 0", i)
		}
		var k1 = keyRange[i].Key
		var v1 = keyRange[i].Val
		//log.Printf("k0=%s; v0=%v;", k0, v0)
		//log.Printf("k1=%s; v0=%v;", k1, v1)
		if k0.Less(k1) || k1.Less(k0) {
			t.Fatalf("i=%d; k0,%s != k1,%s", i, k0, k1)
		}
		if v0 != v1 {
			t.Fatalf("i=%d; v0,%d != v1,%d", i, v0, v1)
		}
		i--
		return true
	}
	m.RangeLimit(startKey, endKey, fn)
	if i != -1 {
		t.Fatalf("after RangeLimit: i,%d != -1", i)
	}
}

// TestBasicRangeForwToSmall test a range which is between two valid keys
func TestBasicRangeForwToSmall(t *testing.T) {
	var m = mkmap(
		mknod(60, black,
			mknod(20, black,
				mknod(10, black, nil, nil),
				mknod(40, black,
					mknod(30, red, nil, nil),
					mknod(50, red, nil, nil))),
			mknod(100, black,
				mknod(80, black,
					mknod(70, red, nil, nil),
					mknod(90, red, nil, nil)),
				mknod(120, black,
					mknod(110, red, nil, nil),
					mknod(130, red, nil, nil)))))

	var startKey = IntKey(62)
	var endKey = IntKey(68)
	var fn = func(k MapKey, v interface{}) bool {
		t.Fatalf("node found where no node should be found.k=%s; v=%d", k, v)
		return true
	}
	m.RangeLimit(startKey, endKey, fn)
}

// TestBasicRangeRevToSmall test a range which is between two valid keys
func TestBasicRangeRevToSmall(t *testing.T) {
	var m = mkmap(
		mknod(60, black,
			mknod(20, black,
				mknod(10, black, nil, nil),
				mknod(40, black,
					mknod(30, red, nil, nil),
					mknod(50, red, nil, nil))),
			mknod(100, black,
				mknod(80, black,
					mknod(70, red, nil, nil),
					mknod(90, red, nil, nil)),
				mknod(120, black,
					mknod(110, red, nil, nil),
					mknod(130, red, nil, nil)))))

	var startKey = IntKey(58)
	var endKey = IntKey(52)
	var fn = func(k MapKey, v interface{}) bool {
		t.Fatalf("node found where no node should be found.k=%s; v=%d", k, v)
		return true
	}
	m.RangeLimit(startKey, endKey, fn)
}

// TestBasicRangeForwToFarAbove test a range which is above any valid keys
func TestBasicRangeForwToFarAbove(t *testing.T) {
	var m = mkmap(
		mknod(60, black,
			mknod(20, black,
				mknod(10, black, nil, nil),
				mknod(40, black,
					mknod(30, red, nil, nil),
					mknod(50, red, nil, nil))),
			mknod(100, black,
				mknod(80, black,
					mknod(70, red, nil, nil),
					mknod(90, red, nil, nil)),
				mknod(120, black,
					mknod(110, red, nil, nil),
					mknod(130, red, nil, nil)))))

	var startKey = IntKey(200)
	var endKey = IntKey(300)
	var fn = func(k MapKey, v interface{}) bool {
		t.Fatalf("node found where no node should be found.k=%s; v=%d", k, v)
		return true
	}
	m.RangeLimit(startKey, endKey, fn)
}

// TestBasicRangeRevToFarAbove test a range which is above any valid keys
func TestBasicRangeRevToFarAbove(t *testing.T) {
	var m = mkmap(
		mknod(60, black,
			mknod(20, black,
				mknod(10, black, nil, nil),
				mknod(40, black,
					mknod(30, red, nil, nil),
					mknod(50, red, nil, nil))),
			mknod(100, black,
				mknod(80, black,
					mknod(70, red, nil, nil),
					mknod(90, red, nil, nil)),
				mknod(120, black,
					mknod(110, red, nil, nil),
					mknod(130, red, nil, nil)))))

	var startKey = IntKey(300)
	var endKey = IntKey(200)
	var fn = func(k MapKey, v interface{}) bool {
		t.Fatalf("node found where no node should be found.k=%s; v=%d", k, v)
		return true
	}
	m.RangeLimit(startKey, endKey, fn)
}

// TestBasicRangeForwToFarBelow test a range which is below any valid keys
func TestBasicRangeForwToFarBelow(t *testing.T) {
	var m = mkmap(
		mknod(60, black,
			mknod(20, black,
				mknod(10, black, nil, nil),
				mknod(40, black,
					mknod(30, red, nil, nil),
					mknod(50, red, nil, nil))),
			mknod(100, black,
				mknod(80, black,
					mknod(70, red, nil, nil),
					mknod(90, red, nil, nil)),
				mknod(120, black,
					mknod(110, red, nil, nil),
					mknod(130, red, nil, nil)))))

	var startKey = IntKey(-100)
	var endKey = IntKey(0)
	var fn = func(k MapKey, v interface{}) bool {
		t.Fatalf("node found where no node should be found.k=%s; v=%d", k, v)
		return true
	}
	m.RangeLimit(startKey, endKey, fn)
}

// TestBasicRangeRevToFarBelow test a range which is below any valid keys
func TestBasicRangeRevToFarBelow(t *testing.T) {
	var m = mkmap(
		mknod(60, black,
			mknod(20, black,
				mknod(10, black, nil, nil),
				mknod(40, black,
					mknod(30, red, nil, nil),
					mknod(50, red, nil, nil))),
			mknod(100, black,
				mknod(80, black,
					mknod(70, red, nil, nil),
					mknod(90, red, nil, nil)),
				mknod(120, black,
					mknod(110, red, nil, nil),
					mknod(130, red, nil, nil)))))

	var startKey = IntKey(0)
	var endKey = IntKey(-100)
	var fn = func(k MapKey, v interface{}) bool {
		t.Fatalf("node found where no node should be found.k=%s; v=%d", k, v)
		return true
	}
	m.RangeLimit(startKey, endKey, fn)
}

func TestBasicRangeStop(t *testing.T) {
	var m = mkmap(
		mknod(60, black,
			mknod(20, black,
				mknod(10, black, nil, nil),
				mknod(40, black,
					mknod(30, red, nil, nil),
					mknod(50, red, nil, nil))),
			mknod(100, black,
				mknod(80, black,
					mknod(70, red, nil, nil),
					mknod(90, red, nil, nil)),
				mknod(120, black,
					mknod(110, red, nil, nil),
					mknod(130, red, nil, nil)))))

	var fn = func(k MapKey, v interface{}) bool {
		if cmp(k, IntKey(60)) == 0 {
			return false
		}
		if less(IntKey(60), k) {
			t.Fatal("encountered a key higher that the stop condition")
		}
		return true
	}
	m.Range(fn)
}

func TestBasicMapString(t *testing.T) {
	var m = mkmap(
		mknod(20, black,
			mknod(10, red, nil, nil),
			mknod(30, red, nil, nil)))

	var expectedStr = "{10: 10, 20: 20, 30: 30}"

	if m.String() != expectedStr {
		t.Fatalf("m.String() != %q", expectedStr)
	}
}

//func TestBasicKeys(t *testing.T) {
//	var m = mkmap(
//		mknod(60, black,
//			mknod(20, black,
//				mknod(10, black, nil, nil),
//				mknod(40, black,
//					mknod(30, red, nil, nil),
//					mknod(50, red, nil, nil))),
//			mknod(100, black,
//				mknod(80, black,
//					mknod(70, red, nil, nil),
//					mknod(90, red, nil, nil)),
//				mknod(120, black,
//					mknod(110, red, nil, nil),
//					mknod(130, red, nil, nil)))))
//
//	var shouldHaveKvs = []KeyVal{
//		{IntKey(10), 10},
//		{IntKey(20), 20},
//		{IntKey(30), 30},
//		{IntKey(40), 40},
//		{IntKey(50), 50},
//		{IntKey(60), 60},
//		{IntKey(70), 70},
//		{IntKey(80), 80},
//		{IntKey(90), 90},
//		{IntKey(100), 100},
//		{IntKey(110), 110},
//		{IntKey(120), 120},
//		{IntKey(130), 130},
//	}
//
//	var foundKeys = m.Keys()
//	for i, kv := range shouldHaveKvs {
//		if cmp(kv.Key, foundKeys[i]) != 0 {
//			t.Fatalf("kv.Key,%s != foundKeys[%d],%s", kv.Key, i, foundKeys[i])
//		}
//	}
//}
