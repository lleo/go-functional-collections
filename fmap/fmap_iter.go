package fmap

import (
	"fmt"

	"github.com/lleo/go-functional-collections/hash"
)

// Iter struct maintains the current state for walking the *Map data structure.
type Iter struct {
	kvIdx       int
	curLeaf     leafI
	tblNextNode tableIterFunc
	stack       *tableIterStack
}

func newIter(root tableI) *Iter {
	var it = new(Iter)
	//it.kvIdx = 0
	//it.curLeaf = nil
	it.tblNextNode = root.iter()
	it.stack = newTableIterStack()
	return it
}

// Next returns each sucessive key/value mapping in the *Map. When all enrties
// have been returned it will return a nil hash.Key.
func (it *Iter) Next() (hash.Key, interface{}) {
	//log.Printf("it.Next: called. it=%s", it)
	var key hash.Key
	var val interface{}

LOOP:
	for {
		switch x := it.curLeaf.(type) {
		case nil:
			key = nil //the end
			break LOOP
		case *flatLeaf:
			key = x.key
			val = x.val
			it.kvIdx = 0
			it.setNextNode() //ignore return false == the end
			//if !it.setNextNode() {
			//	log.Printf("it.Next: case *flatLeaf: it.setNextNode()==false")
			//}
			break LOOP
		case *collisionLeaf:
			if it.kvIdx >= len(x.kvs) {
				it.setNextNode() //ignore return false == the end
				//if !it.setNextNode() {
				//	log.Printf("it.Next: case *collisionLeaf: it.setNextNode()==false")
				//}
				continue LOOP
			}
			key = x.kvs[it.kvIdx].Key
			val = x.kvs[it.kvIdx].Val
			it.kvIdx++
			break LOOP
		default:
			panic("Set (*iter).Next(); it.curLeaf unknown type")
		}
	}
	//log.Printf("it.Next: key=%s; val=%v;", key, val)
	return key, val
}

//setNextNode() sets the iter struct pointing to the next node. If there is no
//next node it returns false, else it returns true.
func (it *Iter) setNextNode() bool {
	//log.Printf("it.setNextNode: called; it=%s", it)
LOOP:
	for {
		//log.Printf("it.setNextNode: it=%s", it)
		var cur = it.tblNextNode()
		//log.Printf("it.setNextNode: it.tblNextNode()=>cur=%s", cur)

		//if cur==nil pop stack and loop
		for cur == nil {
			it.tblNextNode = it.stack.pop()
			if it.tblNextNode == nil {
				it.curLeaf = nil
				return false
			}
			cur = it.tblNextNode()
		}
		//cur != nil
		switch x := cur.(type) {
		case nil:
			panic("WTF!!! cur == nil")
		case tableI:
			it.stack.push(it.tblNextNode)
			it.tblNextNode = x.iter()
			//break switch and LOOP
		case leafI:
			it.curLeaf = x
			break LOOP
		}
		//log.Println("it.setNextNode: looping for")
	}
	return true
}

func (it *Iter) String() string {
	return fmt.Sprintf("%#v", *it)
}
