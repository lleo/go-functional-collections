package sortedSet

import (
	"log"
	"testing"

	"github.com/lleo/go-functional-collections/key"
)

func TestBasicIsSetPos(t *testing.T) {
	var s = mkset(
		mknod(20, black,
			mknod(10, red, nil, nil),
			mknod(30, red, nil, nil)))

	var isSet = s.IsSet(key.Int(20))

	if !isSet {
		t.Fatal("IsSet(key.Int(20)) returned false")
	}
}

func TestBasicIsSetNeg(t *testing.T) {
	var s = mkset(
		mknod(20, black,
			mknod(10, red, nil, nil),
			mknod(30, red, nil, nil)))

	var isSet = s.IsSet(key.Int(40))

	if isSet {
		t.Fatal("s.IsSet(key.Int(40)) returned true")
	}
}

func TestBasicAddReplace(t *testing.T) {
	var m0 = mkset(
		mknod(20, black,
			mknod(10, red, nil, nil),
			mknod(30, red, nil, nil)))

	var origM0 = m0
	var dupM0 = m0.DeepCopy()

	var m1, added = m0.Add(key.Int(30))

	if added {
		t.Fatal("Add added new entry when it should not")
	}

	if m1.NumEntries() != 3 {
		t.Fatalf("m1.NumEntries(),%d != 3", m1.NumEntries())
	}
	if !origM0.Equiv(dupM0) {
		t.Fatal("TestBasicStoreReplace: orig Set and duplicate of orig Set are not identical.")
	}
}

func TestBasicSetCase1(t *testing.T) {
	var s = mkset(nil)

	var origM = s
	var dupM = s.DeepCopy()

	//log.Printf("BEFORE Set s =\n%s", s.treeString())

	s = s.Set(key.Int(10))

	//log.Printf("AFTER Set s =\n%s", s.treeString())

	if s.NumEntries() != 1 {
		t.Fatal("s.NumEntries() != 1")
	}

	if err := s.valid(); err != nil {
		t.Fatalf("set not valid; err=%s", err)
	}

	if !origM.Equiv(dupM) {
		t.Fatal("TestBasicSetCase1: orig Set and duplicate of orig Set are not identical.")
	}
}

func TestBasicSetCase2(t *testing.T) {
	var s = mkset(
		mknod(20, black,
			mknod(10, red, nil, nil),
			nil))

	var origM = s
	var dupM = s.DeepCopy()

	//log.Printf("BEFORE Set s =\n%s", s.treeString())

	s = s.Set(key.Int(30))

	//log.Printf("AFTER Set s =\n%s", s.treeString())

	if s.NumEntries() != 3 {
		t.Fatal("s.NumEntries() != 1")
	}

	if err := s.valid(); err != nil {
		t.Fatalf("set not valid; err=%s", err)
	}

	if !origM.Equiv(dupM) {
		t.Fatal("TestBasicSetCase2: orig Set and duplicate of orig Set are not identical.")
	}
}

func TestBasicSetCase3(t *testing.T) {
	//insert order 10, 20, 50, 40, 30, 60
	var s = mkset(
		mknod(20, black,
			mknod(10, black, nil, nil),
			mknod(40, black,
				mknod(30, red, nil, nil),
				mknod(50, red, nil, nil),
			),
		))

	var origM = s
	var dupM = s.DeepCopy()

	//log.Printf("BEFORE Set s =\n%s", s.treeString())

	s = s.Set(key.Int(60))

	//log.Printf("AFTER Set s =\n%s", s.treeString())

	if s.NumEntries() != 6 {
		t.Fatal("s.NumEntries() != 6")
	}

	if err := s.valid(); err != nil {
		t.Fatalf("set not valid; err=%s", err)
	}

	if !origM.Equiv(dupM) {
		t.Fatal("TestBasicSetCase3: orig Set and duplicate of orig Set are not identical.")
	}
}

func TestBasicSetCase4(t *testing.T) {
	var s = mkset(
		mknod(50, black,
			mknod(20, black,
				nil,
				mknod(40, red, nil, nil)),
			mknod(60, black,
				nil,
				mknod(70, red, nil, nil)),
		))

	var origM = s           //copy the pointer
	var dupM = s.DeepCopy() //copy the value

	//log.Printf("BEFORE Set s =\n%s", s.treeString())

	s = s.Set(key.Int(30))

	//log.Printf("AFTER Set s =\n%s", s.treeString())

	if s.NumEntries() != 6 {
		t.Fatal("s.NumEntries() != 6")
	}

	if err := s.valid(); err != nil {
		t.Fatalf("set not valid; err=%s", err)
	}

	if !origM.Equiv(dupM) {
		t.Fatal("TestBasicSetCase4: orig Set and duplicate of orig Set are not identical.")
	}
}

func TestBasicRemoveNeg(t *testing.T) {
	var m0 = mkset(
		mknod(20, black,
			mknod(10, red, nil, nil),
			mknod(30, red, nil, nil)))

	var m1, found = m0.Remove(key.Int(40))

	if found {
		t.Fatal("found a key that does not exist")
	}

	if m1 != m0 {
		t.Fatal("returned set not the same as the original set")
	}
}

func TestBasicUnsetCase1Tree0(t *testing.T) {
	var m0 = mkset(
		mknod(10, black, nil, nil))

	var then = m0.treeString()
	//var dupM0 = m0.DeepCopy()

	var m1 = m0.Unset(key.Int(10))

	if m1.NumEntries() != 0 {
		t.Fatal("s.NumEntries() != 0")
	}

	var now = m0.treeString()
	if then != now {
		log.Printf("origninal tree changeed:\nTHEN: %s\nNOW: %s",
			then, now)
		t.Fatal("The original tree changed.")
	}

	//if !m0.Equiv(dupM0) {
	//	t.Fatal("The original tree changed.")
	//}
}

func TestBasicUnsetCase1Tree1(t *testing.T) {
	var m0 = mkset(
		mknod(10, black,
			nil,
			mknod(20, red, nil, nil),
		))

	var then = m0.treeString()

	var m1 = m0.Unset(key.Int(10))

	if m1.NumEntries() != 1 {
		t.Fatal("s.NumEntries() != 1")
	}

	var now = m0.treeString()
	if then != now {
		log.Printf("origninal tree changeed:\nTHEN: %s\nNOW: %s",
			then, now)
		t.Fatal("The original tree changed.")
	}
}

func TestBasicUnsetCase1Tree2(t *testing.T) {
	var m0 = mkset(
		mknod(20, black,
			mknod(10, red, nil, nil),
			nil,
		))

	var then = m0.treeString()

	var m1 = m0.Unset(key.Int(20))

	if m1.NumEntries() != 1 {
		t.Fatal("s.NumEntries() != 1")
	}

	var now = m0.treeString()
	if then != now {
		log.Printf("origninal tree changeed:\nTHEN: %s\nNOW: %s",
			then, now)
		t.Fatal("The original tree changed.")
	}
}

// DeleteCase1 is exhaustively tested.

func TestBasicUnsetCase2Tree0(t *testing.T) {
	var s = mkset(
		mknod(20, black,
			nil,
			mknod(30, red, nil, nil),
		))

	//log.Printf("BEFORE REMOVE: Set s=\n%s", s.treeString())

	s = s.Unset(key.Int(30))

	//log.Printf("AFTER REMOVE Set s=\n%s", s.treeString())

	if s.NumEntries() != 1 {
		t.Fatalf("s.NumEntries(),%d != 1", s.NumEntries())
	}

	if !s.root.isBlack() {
		t.Fatal("!s.root.isBlack()")
	}

	if s.root.ln != nil {
		t.Fatal("s.root.rn != nil")
	}

	if s.root.rn != nil {
		t.Fatal("s.root.ln != nil")
	}
}

func TestBasicUnsetCase2Tree1(t *testing.T) {
	var s = mkset(
		mknod(20, black,
			mknod(10, red, nil, nil),
			nil,
		))

	//log.Printf("BEFORE REMOVE: Set s=\n%s", s.treeString())

	s = s.Unset(key.Int(10))

	//log.Printf("AFTER REMOVE Set s=\n%s", s.treeString())

	if s.NumEntries() != 1 {
		t.Fatalf("s.NumEntries(),%d != 1", s.NumEntries())
	}

	if !s.root.isBlack() {
		t.Fatal("!s.root.isBlack()")
	}

	if s.root.ln != nil {
		t.Fatal("s.root.ln != nil")
	}

	if s.root.rn != nil {
		t.Fatal("s.root.rn != nil")
	}
}

func TestBasicUnsetCase3Tree0(t *testing.T) {
	var s = mkset(
		mknod(20, black,
			mknod(10, black, nil, nil),
			mknod(30, black, nil, nil),
		))

	s = s.Unset(key.Int(30))
	if s.NumEntries() != 2 {
		t.Fatalf("s.NumEntries(),%d != 2", s.NumEntries())
	}

	if !s.root.isBlack() {
		t.Fatal("!s.root.isBlack()")
	}

	if !s.root.ln.isRed() {
		t.Fatal("!s.root.ln.isRed()")
	}

	if s.root.rn != nil {
		t.Fatal("s.root.rn != nil")
	}
}

func TestBasicUnsetCase6Tree0(t *testing.T) {
	var s = mkset(
		mknod(20, black,
			mknod(10, black, nil, nil),
			mknod(40, red,
				mknod(30, black, nil, nil),
				mknod(50, black,
					nil,
					mknod(60, red, nil, nil)))))

	var origM = s
	var dupOrigM = s.DeepCopy()
	var origSetStr0 = s.treeString()

	//log.Printf("origSetStr0 =\n%s", origSetStr0)

	s = s.Unset(key.Int(30))

	if s.NumEntries() != 5 {
		t.Fatalf("s.NumEntries(),%d != 5", s.NumEntries())
	}

	var origSetStr1 = origM.treeString()
	if origSetStr0 != origSetStr1 {
		log.Printf("origSetStr0 != origSetStr1:\n"+
			"origSetStr0=\n%s\norigSetStr1=\n%s", origSetStr0, origSetStr1)
	}

	if !origM.Equiv(dupOrigM) {
		t.Fatal("TestBasicPutCase4: orig Set and duplicate of orig Set are not identical.")
	}
}

func TestBasicUnsetTwoChildTree0(t *testing.T) {
	var s = mkset(
		mknod(40, black,
			mknod(20, black,
				mknod(10, red, nil, nil),
				mknod(30, red, nil, nil)),
			mknod(70, red,
				mknod(50, black, nil, nil),
				mknod(80, black, nil, nil))))

	var shouldHaveKeys = []key.Sort{
		key.Int(10),
		key.Int(30),
		key.Int(40),
		key.Int(50),
		key.Int(70),
		key.Int(80),
	}

	var origM = s
	var dupOrigM = s.DeepCopy()
	var origSetStr0 = s.treeString()

	//log.Printf("BEFORE DEL s = \n%s", s.treeString())

	s = s.Unset(key.Int(20))

	//log.Printf("AFTER DEL s = \n%s", s.treeString())

	if err := s.valid(); err != nil {
		t.Fatalf("INVALID TREE AFTER Unset(key.Int(20)); err=%s\n", err)
	}

	if s.NumEntries() != 6 {
		t.Fatalf("s.NumEntries(),%d != 6", s.NumEntries())
	}

	for _, key := range shouldHaveKeys {
		var isSet = s.IsSet(key)
		if !isSet {
			log.Printf("failed to find shouldHave key=%s in Tree=\n%s",
				key, s.treeString())
			t.Fatalf("Failed to find shouldHave key=%s", key)
		}
	}

	var origSetStr1 = origM.treeString()
	if origSetStr0 != origSetStr1 {
		log.Printf("origSetStr0 != origSetStr1:\n"+
			"origSetStr0=\n%s\norigSetStr1=\n%s", origSetStr0, origSetStr1)
	}

	if !origM.Equiv(dupOrigM) {
		t.Fatal("TestBasicUnsetTwoChildrenCase2: " +
			"orig Set and duplicate of orig Set are not identical.")
	}
}

//deleteCase4
func TestBasicUnsetTwoChildTree1(t *testing.T) {
	var s = mkset(
		mknod(40, black,
			mknod(10, black,
				nil,
				mknod(30, red, nil, nil)),
			mknod(70, red,
				mknod(50, black, nil, nil),
				mknod(80, black, nil, nil))))

	//shouldHave after Unset(70)
	var shouldHaveKeys = []key.Sort{
		key.Int(10),
		key.Int(30),
		key.Int(40),
		key.Int(50),
		key.Int(80),
	}

	var origM = s
	var dupOrigM = s.DeepCopy()
	var origSetStr0 = s.treeString()

	//log.Printf("BEFORE DEL s = \n%s", s.treeString())

	s = s.Unset(key.Int(70))

	//log.Printf("AFTER DEL s = \n%s", s.treeString())

	if err := s.valid(); err != nil {
		t.Fatalf("INVALID TREE AFTER Del(key.Int(70)); err=%s\n", err)
	}

	if s.NumEntries() != 5 {
		t.Fatalf("s.NumEntries(),%d != 5", s.NumEntries())
	}

	for _, key := range shouldHaveKeys {
		var isSet = s.IsSet(key)
		if !isSet {
			log.Printf("failed to find shouldHave key=%s in Tree=\n%s",
				key, s.treeString())
			t.Fatalf("Failed to find shouldHave key=%s", key)
		}
	}

	var origSetStr1 = origM.treeString()
	if origSetStr0 != origSetStr1 {
		log.Printf("origSetStr0 != origSetStr1:\n"+
			"origSetStr0=\n%s\norigSetStr1=\n%s", origSetStr0, origSetStr1)
	}

	if !origM.Equiv(dupOrigM) {
		t.Fatal("TestBasicDelTwoChildrenCase2: " +
			"orig Set and duplicate of orig Set are not identical.")
	}
}

func TestBasicUnsetTwoChildTree2(t *testing.T) {
	var s = mkset(
		mknod(40, black,
			mknod(10, black, nil, nil),
			mknod(70, red,
				mknod(50, black, nil, nil),
				mknod(80, black, nil, nil))))

	//shouldHave after Unset(70)
	var shouldHaveKeys = []key.Sort{
		key.Int(10),
		key.Int(50),
		key.Int(70),
		key.Int(80),
	}

	var origM = s
	var dupOrigM = s.DeepCopy()
	var origSetStr0 = s.treeString()

	//log.Printf("BEFORE DEL s = \n%s", s.treeString())

	s = s.Unset(key.Int(40))

	//log.Printf("AFTER DEL s = \n%s", s.treeString())

	if err := s.valid(); err != nil {
		t.Fatalf("INVALID TREE AFTER Unset(key.Int(40)); err=%s\n", err)
	}

	if s.NumEntries() != 4 {
		t.Fatalf("s.NumEntries(),%d != 4", s.NumEntries())
	}

	for _, key := range shouldHaveKeys {
		var isSet = s.IsSet(key)
		if !isSet {
			log.Printf("failed to find shouldHave key=%s in Tree=\n%s",
				key, s.treeString())
			t.Fatalf("Failed to find shouldHave key=%s", key)
		}
	}

	var origSetStr1 = origM.treeString()
	if origSetStr0 != origSetStr1 {
		log.Printf("origSetStr0 != origSetStr1:\n"+
			"origSetStr0=\n%s\norigSetStr1=\n%s", origSetStr0, origSetStr1)
	}

	if !origM.Equiv(dupOrigM) {
		t.Fatal("TestBasicUnsetTwoChildrenCase2: " +
			"orig Set and duplicate of orig Set are not identical.")
	}
}

func TestBasicUnsetTwoChildTree3(t *testing.T) {
	var s = mkset(
		mknod(50, black,
			mknod(20, red,
				mknod(10, black, nil, nil),
				mknod(40, black,
					mknod(30, red, nil, nil),
					nil)),
			mknod(80, black, nil, nil)))

	//shouldHave after Unset(20)
	var shouldHaveKeys = []key.Sort{
		key.Int(10),
		key.Int(30),
		key.Int(40),
		key.Int(50),
		key.Int(80),
	}

	var origM = s
	var dupOrigM = s.DeepCopy()
	var origSetStr0 = s.treeString()

	//log.Printf("BEFORE DEL s = \n%s", s.treeString())

	s = s.Unset(key.Int(20))

	//log.Printf("AFTER DEL s = \n%s", s.treeString())

	if err := s.valid(); err != nil {
		t.Fatalf("INVALID TREE AFTER Unset(key.Int(20)); err=%s\n", err)
	}

	if s.NumEntries() != 5 {
		t.Fatalf("s.NumEntries(),%d != 5", s.NumEntries())
	}

	for _, key := range shouldHaveKeys {
		var isSet = s.IsSet(key)
		if !isSet {
			log.Printf("failed to find shouldHave key=%s in Tree=\n%s",
				key, s.treeString())
			t.Fatalf("Failed to find shouldHave key=%s", key)
		}
	}

	var origSetStr1 = origM.treeString()
	if origSetStr0 != origSetStr1 {
		log.Printf("origSetStr0 != origSetStr1:\n"+
			"origSetStr0=\n%s\norigSetStr1=\n%s", origSetStr0, origSetStr1)
	}

	if !origM.Equiv(dupOrigM) {
		t.Fatal("TestBasicUnsetTwoChildrenCase2: " +
			"orig Set and duplicate of orig Set are not identical.")
	}
}

func TestBasicUnsetTwoChildTree4(t *testing.T) {
	var s = mkset(
		mknod(60, black,
			mknod(20, black,
				mknod(10, black, nil, nil),
				mknod(40, red,
					mknod(30, black, nil, nil),
					mknod(50, black, nil, nil))),
			mknod(80, black,
				mknod(70, black, nil, nil),
				mknod(90, black, nil, nil))))

	if err := s.valid(); err != nil {
		t.Fatalf("INVALID TREE; err=%s\n", err)
	}

	//shouldHave after Unset(80)
	var shouldHaveKeys = []key.Sort{
		key.Int(10),
		key.Int(20),
		key.Int(30),
		key.Int(40),
		key.Int(50),
		key.Int(60),
		key.Int(70),
		//key.Int(80),
		key.Int(90),
	}

	var origM = s
	var dupOrigM = s.DeepCopy()
	var origSetStr0 = s.treeString()

	//log.Printf("BEFORE DEL s = \n%s", s.treeString())

	s = s.Unset(key.Int(80))

	//log.Printf("AFTER DEL s = \n%s", s.treeString())

	if err := s.valid(); err != nil {
		t.Fatalf("INVALID TREE AFTER Unset(key.Int(80)); err=%s\n", err)
	}

	if s.NumEntries() != 8 {
		t.Fatalf("s.NumEntries(),%d != 8", s.NumEntries())
	}

	for _, key := range shouldHaveKeys {
		var isSet = s.IsSet(key)
		if !isSet {
			log.Printf("failed to find shouldHave key=%s in Tree=\n%s",
				key, s.treeString())
			t.Fatalf("Failed to find shouldHave key=%s", key)
		}
	}

	var origSetStr1 = origM.treeString()
	if origSetStr0 != origSetStr1 {
		log.Printf("origSetStr0 != origSetStr1:\n"+
			"origSetStr0=\n%s\norigSetStr1=\n%s", origSetStr0, origSetStr1)
	}

	if !origM.Equiv(dupOrigM) {
		t.Fatal("TestBasicUnsetTwoChildrenCase2: " +
			"orig Set and duplicate of orig Set are not identical.")
	}
}

func TestBasicRange(t *testing.T) {
	var s = mkset(
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

	if err := s.valid(); err != nil {
		t.Fatalf("INVALID TREE; err=%s\n", err)
	}

	var shouldHaveKeys = []key.Sort{
		key.Int(10),
		key.Int(20),
		key.Int(30),
		key.Int(40),
		key.Int(50),
		key.Int(60),
		key.Int(70),
		key.Int(80),
		key.Int(90),
		key.Int(100),
		key.Int(110),
		key.Int(120),
		key.Int(130),
	}

	var i int
	var fn = func(k0 key.Sort) bool {
		var k1 = shouldHaveKeys[i]
		//log.Printf("k0=%s;", k0)
		//log.Printf("k1=%s;", k1)
		if key.Cmp(k0, k1) != 0 {
			t.Fatalf("i=%d; k0,%s != k1,%s", i, k0, k1)
		}
		i++
		return true
	}
	s.Range(fn)
}

func TestBasicRangeForwBeg(t *testing.T) {
	var s = mkset(
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

	if err := s.valid(); err != nil {
		t.Fatalf("INVALID TREE; err=%s\n", err)
	}

	var shouldHaveKeys = []key.Sort{
		key.Int(10),
		key.Int(20),
		key.Int(30),
		key.Int(40),
		key.Int(50),
		key.Int(60),
		key.Int(70),
		key.Int(80),
		key.Int(90),
		key.Int(100),
		key.Int(110),
		key.Int(120),
		key.Int(130),
	}

	//var numKeys = len(shouldHaveKeys)
	var eltOffset = 3
	var startKey = key.Int(eltOffset * 10)
	var endKey = key.Inf(1) //positive infinity
	var keyRange = shouldHaveKeys[eltOffset-1:]
	var i = 0
	var fn = func(k0 key.Sort) bool {
		if i >= len(keyRange) {
			t.Fatalf("i,%d >= len(keyRange),%d", i, len(keyRange))
		}
		var k1 = keyRange[i]
		//log.Printf("k0=%s;", k0)
		//log.Printf("k1=%s;", k1)
		if key.Cmp(k0, k1) != 0 {
			t.Fatalf("i=%d; k0,%s != k1,%s", i, k0, k1)
		}
		i++
		return true
	}
	s.RangeLimit(startKey, endKey, fn)
	if i != len(keyRange) {
		t.Fatalf("after RangeLimit: i,%d != len(keyRange),%d", i, len(keyRange))
	}
}

func TestBasicRangeForwBegInexact(t *testing.T) {
	var s = mkset(
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

	if err := s.valid(); err != nil {
		t.Fatalf("INVALID TREE; err=%s\n", err)
	}

	var shouldHaveKeys = []key.Sort{
		key.Int(10),
		key.Int(20),
		key.Int(30),
		key.Int(40),
		key.Int(50),
		key.Int(60),
		key.Int(70),
		key.Int(80),
		key.Int(90),
		key.Int(100),
		key.Int(110),
		key.Int(120),
		key.Int(130),
	}

	//var numKeys = len(shouldHaveKeys)
	var eltOffset = 2
	var startKey = key.Int((eltOffset * 10) - 5) //key.Int(15)
	var endKey = key.Inf(1)                      //positive infinity
	var keyRange = shouldHaveKeys[eltOffset-1:]
	var i = 0
	var fn = func(k0 key.Sort) bool {
		if i >= len(keyRange) {
			t.Fatalf("i,%d >= len(keyRange),%d", i, len(keyRange))
		}
		var k1 = keyRange[i]
		//log.Printf("k0=%s;", k0)
		//log.Printf("k1=%s;", k1)
		if key.Cmp(k0, k1) != 0 {
			t.Fatalf("i=%d; k0,%s != k1,%s", i, k0, k1)
		}
		i++
		return true
	}
	s.RangeLimit(startKey, endKey, fn)
	if i != len(keyRange) {
		t.Fatalf("after RangeLimit: i,%d != len(keyRange),%d", i, len(keyRange))
	}
}

func TestBasicRangeForwEnd(t *testing.T) {
	var s = mkset(
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

	if err := s.valid(); err != nil {
		t.Fatalf("INVALID TREE; err=%s\n", err)
	}

	var shouldHaveKeys = []key.Sort{
		key.Int(10),
		key.Int(20),
		key.Int(30),
		key.Int(40),
		key.Int(50),
		key.Int(60),
		key.Int(70),
		key.Int(80),
		key.Int(90),
		key.Int(100),
		key.Int(110),
		key.Int(120),
		key.Int(130),
	}

	var numKeys = len(shouldHaveKeys)
	var eltOffset = 3
	var startKey = key.Inf(-1)                            //negative infinity
	var endKey = key.Int((numKeys - eltOffset) * 10)      //key.Int(100)
	var keyRange = shouldHaveKeys[:len(shouldHaveKeys)-3] //??
	var i = 0
	var fn = func(k0 key.Sort) bool {
		if i >= len(keyRange) {
			t.Fatalf("i,%d >= len(keyRange),%d", i, len(keyRange))
		}
		var k1 = keyRange[i]
		//log.Printf("k0=%s;", k0)
		//log.Printf("k1=%s;", k1)
		if key.Cmp(k0, k1) != 0 {
			t.Fatalf("i=%d; k0,%s != k1,%s", i, k0, k1)
		}
		i++
		return true
	}
	s.RangeLimit(startKey, endKey, fn)
	if i != len(keyRange) {
		t.Fatalf("after RangeLimit: i,%d != len(keyRange),%d", i, len(keyRange))
	}
}

func TestBasicRangeForwEndInexact(t *testing.T) {
	var s = mkset(
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

	if err := s.valid(); err != nil {
		t.Fatalf("INVALID TREE; err=%s\n", err)
	}

	var shouldHaveKeys = []key.Sort{
		key.Int(10),
		key.Int(20),
		key.Int(30),
		key.Int(40),
		key.Int(50),
		key.Int(60),
		key.Int(70),
		key.Int(80),
		key.Int(90),
		key.Int(100),
		key.Int(110),
		key.Int(120),
		key.Int(130),
	}

	var numKeys = len(shouldHaveKeys)
	var eltOffset = 3
	var startKey = key.Inf(-1)                            //negative infinity
	var endKey = key.Int((numKeys-eltOffset)*10 + 5)      //key.Int(105)
	var keyRange = shouldHaveKeys[:len(shouldHaveKeys)-3] //??
	var i = 0
	var fn = func(k0 key.Sort) bool {
		if i >= len(keyRange) {
			t.Fatalf("i,%d >= len(keyRange),%d", i, len(keyRange))
		}
		var k1 = keyRange[i]
		//log.Printf("k0=%s;", k0)
		//log.Printf("k1=%s;", k1)
		if key.Cmp(k0, k1) != 0 {
			t.Fatalf("i=%d; k0,%s != k1,%s", i, k0, k1)
		}
		i++
		return true
	}
	s.RangeLimit(startKey, endKey, fn)
	if i != len(keyRange) {
		t.Fatalf("after RangeLimit: i,%d != len(keyRange),%d", i, len(keyRange))
	}
}

func TestBasicRangeRevBeg(t *testing.T) {
	var s = mkset(
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

	if err := s.valid(); err != nil {
		t.Fatalf("INVALID TREE; err=%s\n", err)
	}

	var shouldHaveKeys = []key.Sort{
		key.Int(10),
		key.Int(20),
		key.Int(30),
		key.Int(40),
		key.Int(50),
		key.Int(60),
		key.Int(70),
		key.Int(80),
		key.Int(90),
		key.Int(100),
		key.Int(110),
		key.Int(120),
		key.Int(130),
	}

	var numKeys = len(shouldHaveKeys)
	var eltOffset = 3
	var keyRange = shouldHaveKeys[:numKeys-eltOffset]
	var startKey = key.Int((numKeys - eltOffset) * 10) //key.Int(100)
	var endKey = key.Inf(-1)                           //negative infinity
	var i = len(keyRange) - 1
	var fn = func(k0 key.Sort) bool {
		if i < 0 {
			t.Fatalf("i,%d < 0", i)
		}
		var k1 = keyRange[i]
		//log.Printf("k0=%s;", k0)
		//log.Printf("k1=%s;", k1)
		if key.Cmp(k0, k1) != 0 {
			t.Fatalf("i=%d; k0,%s != k1,%s", i, k0, k1)
		}
		i--
		return true
	}
	//log.Printf("s.RangeLimit(startKey,%s, endKey,%s, fn)", startKey, endKey)
	s.RangeLimit(startKey, endKey, fn)
	if i != -1 {
		t.Fatalf("after RangeLimit: i,%d != -1", i)
	}
}

func TestBasicRangeRevBegInexact(t *testing.T) {
	var s = mkset(
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

	if err := s.valid(); err != nil {
		t.Fatalf("INVALID TREE; err=%s\n", err)
	}

	var shouldHaveKeys = []key.Sort{
		key.Int(10),
		key.Int(20),
		key.Int(30),
		key.Int(40),
		key.Int(50),
		key.Int(60),
		key.Int(70),
		key.Int(80),
		key.Int(90),
		key.Int(100),
		key.Int(110),
		key.Int(120),
		key.Int(130),
	}

	var numKeys = len(shouldHaveKeys)
	var eltOffset = 3
	var keyRange = shouldHaveKeys[:numKeys-eltOffset]
	var startKey = key.Int((numKeys-eltOffset)*10 + 5) //key.Int(105)
	var endKey = key.Inf(-1)                           //negative infinity
	var i = len(keyRange) - 1
	var fn = func(k0 key.Sort) bool {
		if i < 0 {
			t.Fatalf("i,%d < 0", i)
		}
		var k1 = keyRange[i]
		//log.Printf("k0=%s;", k0)
		//log.Printf("k1=%s;", k1)
		if key.Cmp(k0, k1) != 0 {
			t.Fatalf("i=%d; k0,%s != k1,%s", i, k0, k1)
		}
		i--
		return true
	}
	//log.Printf("s.RangeLimit(startKey,%s, endKey,%s, fn)", startKey, endKey)
	s.RangeLimit(startKey, endKey, fn)
	if i != -1 {
		t.Fatalf("after RangeLimit: i,%d != -1", i)
	}
}

func TestBasicRangeRevEnd(t *testing.T) {
	var s = mkset(
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

	if err := s.valid(); err != nil {
		t.Fatalf("INVALID TREE; err=%s\n", err)
	}

	var shouldHaveKeys = []key.Sort{
		key.Int(10),
		key.Int(20),
		key.Int(30),
		key.Int(40),
		key.Int(50),
		key.Int(60),
		key.Int(70),
		key.Int(80),
		key.Int(90),
		key.Int(100),
		key.Int(110),
		key.Int(120),
		key.Int(130),
	}

	//var numKeys = len(shouldHaveKeys)
	var eltOffset = 3
	var startKey = key.Inf(1)                   //positive infinity
	var endKey = key.Int(eltOffset * 10)        //key.Int(30)
	var keyRange = shouldHaveKeys[eltOffset-1:] //??
	var i = len(keyRange) - 1
	var fn = func(k0 key.Sort) bool {
		if i < 0 {
			t.Fatalf("i,%d < 0", i)
		}
		var k1 = keyRange[i]
		//log.Printf("k0=%s;", k0)
		//log.Printf("k1=%s;", k1)
		if key.Cmp(k0, k1) != 0 {
			t.Fatalf("i=%d; k0,%s != k1,%s", i, k0, k1)
		}
		i--
		return true
	}
	s.RangeLimit(startKey, endKey, fn)
	if i != -1 {
		t.Fatalf("after RangeLimit: i,%d != -1", i)
	}
}

func TestBasicRangeRevEndInexact(t *testing.T) {
	var s = mkset(
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

	if err := s.valid(); err != nil {
		t.Fatalf("INVALID TREE; err=%s\n", err)
	}

	var shouldHaveKeys = []key.Sort{
		key.Int(10),
		key.Int(20),
		key.Int(30),
		key.Int(40),
		key.Int(50),
		key.Int(60),
		key.Int(70),
		key.Int(80),
		key.Int(90),
		key.Int(100),
		key.Int(110),
		key.Int(120),
		key.Int(130),
	}

	//var numKeys = len(shouldHaveKeys)
	var eltOffset = 3
	var startKey = key.Inf(1)                   //positive infinity
	var endKey = key.Int(eltOffset*10 - 5)      //key.Int(25)
	var keyRange = shouldHaveKeys[eltOffset-1:] //??
	var i = len(keyRange) - 1
	var fn = func(k0 key.Sort) bool {
		if i < 0 {
			t.Fatalf("i,%d < 0", i)
		}
		var k1 = keyRange[i]
		//log.Printf("k0=%s;", k0)
		//log.Printf("k1=%s;", k1)
		if key.Cmp(k0, k1) != 0 {
			t.Fatalf("i=%d; k0,%s != k1,%s", i, k0, k1)
		}
		i--
		return true
	}
	s.RangeLimit(startKey, endKey, fn)
	if i != -1 {
		t.Fatalf("after RangeLimit: i,%d != -1", i)
	}
}

// TestBasicRangeForwToSmall test a range which is between two valid keys
func TestBasicRangeForwToSmall(t *testing.T) {
	var s = mkset(
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

	var startKey = key.Int(62)
	var endKey = key.Int(68)
	var fn = func(k key.Sort) bool {
		t.Fatalf("node found where no node should be found.k=%s;", k)
		return true
	}
	s.RangeLimit(startKey, endKey, fn)
}

// TestBasicRangeRevToSmall test a range which is between two valid keys
func TestBasicRangeRevToSmall(t *testing.T) {
	var s = mkset(
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

	var startKey = key.Int(58)
	var endKey = key.Int(52)
	var fn = func(k key.Sort) bool {
		t.Fatalf("node found where no node should be found.k=%s;", k)
		return true
	}
	s.RangeLimit(startKey, endKey, fn)
}

// TestBasicRangeForwToFarAbove test a range which is above any valid keys
func TestBasicRangeForwToFarAbove(t *testing.T) {
	var s = mkset(
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

	var startKey = key.Int(200)
	var endKey = key.Int(300)
	var fn = func(k key.Sort) bool {
		t.Fatalf("node found where no node should be found.k=%s;", k)
		return true
	}
	s.RangeLimit(startKey, endKey, fn)
}

// TestBasicRangeRevToFarAbove test a range which is above any valid keys
func TestBasicRangeRevToFarAbove(t *testing.T) {
	var s = mkset(
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

	var startKey = key.Int(300)
	var endKey = key.Int(200)
	var fn = func(k key.Sort) bool {
		t.Fatalf("node found where no node should be found.k=%s;", k)
		return true
	}
	s.RangeLimit(startKey, endKey, fn)
}

// TestBasicRangeForwToFarBelow test a range which is below any valid keys
func TestBasicRangeForwToFarBelow(t *testing.T) {
	var s = mkset(
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

	var startKey = key.Int(-100)
	var endKey = key.Int(0)
	var fn = func(k key.Sort) bool {
		t.Fatalf("node found where no node should be found.k=%s;", k)
		return true
	}
	s.RangeLimit(startKey, endKey, fn)
}

// TestBasicRangeRevToFarBelow test a range which is below any valid keys
func TestBasicRangeRevToFarBelow(t *testing.T) {
	var s = mkset(
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

	var startKey = key.Int(0)
	var endKey = key.Int(-100)
	var fn = func(k key.Sort) bool {
		t.Fatalf("node found where no node should be found.k=%s;", k)
		return true
	}
	s.RangeLimit(startKey, endKey, fn)
}

func TestBasicRangeStop(t *testing.T) {
	var s = mkset(
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

	var fn = func(k key.Sort) bool {
		if key.Cmp(k, key.Int(60)) == 0 {
			return false
		}
		if key.Less(key.Int(60), k) {
			t.Fatal("encountered a key higher that the stop condition")
		}
		return true
	}
	s.Range(fn)
}

func TestBasicString(t *testing.T) {
	var s = mkset(
		mknod(20, black,
			mknod(10, red, nil, nil),
			mknod(30, red, nil, nil)))

	var expectedStr = "{10, 20, 30}"

	if s.String() != expectedStr {
		t.Fatalf("s.String(),%s != %q", s.String(), expectedStr)
	}
}

func TestBasicKeys(t *testing.T) {
	var s = mkset(
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

	var shouldHaveKeys = []key.Sort{
		key.Int(10),
		key.Int(20),
		key.Int(30),
		key.Int(40),
		key.Int(50),
		key.Int(60),
		key.Int(70),
		key.Int(80),
		key.Int(90),
		key.Int(100),
		key.Int(110),
		key.Int(120),
		key.Int(130),
	}

	var foundKeys = s.Keys()
	for i, k := range shouldHaveKeys {
		if key.Cmp(k, foundKeys[i]) != 0 {
			t.Fatalf("key,%s != foundKeys[%d],%s", k, i, foundKeys[i])
		}
	}
}

func TestBasicCount(t *testing.T) {
	var size = 1000
	var keys = buildKeys(size)
	var s = buildSet(keys)

	var count = s.Count()
	if count != size {
		t.Fatalf("s.Count(),%d != size,%d", count, size)
	}

	var numEnts = s.NumEntries()
	if numEnts != size {
		t.Fatalf("s.NumEntries(),%d != size,%d", numEnts, size)
	}
}

func TestBasicNewFromList(t *testing.T) {
	var size = 10
	var keys = buildKeys(size)
	var s = NewFromList(keys)

	var count = s.Count()
	if count != size {
		t.Fatalf("s.Count(),%d != size,%d", count, size)
	}

	var numEnts = s.NumEntries()
	if numEnts != size {
		t.Fatalf("s.NumEntries(),%d != size,%d", numEnts, size)
	}
}

func TestBasicBulkInsert(t *testing.T) {}

func TestBasicMerge(t *testing.T) {}

func TestBasicBulkDelete(t *testing.T) {}

func TestBasicBulkDelete2(t *testing.T) {}

func TestBasicUnion(t *testing.T) {}

func TestBasicIntersect(t *testing.T) {}

func TestBasicDifference(t *testing.T) {}
