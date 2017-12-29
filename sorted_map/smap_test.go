package sorted_map

import (
	"log"
	"math/rand"
	"os"
	"strconv"

	"github.com/pkg/errors"
)

//Set up log file
func init() {
	log.SetFlags(log.Lshortfile)

	var logFileName = "test.log"
	var logFile, err = os.Create(logFileName)
	if err != nil {
		log.Fatal(errors.Wrapf(err, "failed to os.Create(%q)", logFileName))
	}
	log.SetOutput(logFile)

	//log.Println("TESTING HAS STARTED...")
}

type IntKey int

func (ik IntKey) Less(o MapKey) bool {
	var oik, ok = o.(IntKey)
	if !ok {
		panic("o is not a IntKey")
	}
	return ik < oik
}

func (ik IntKey) String() string {
	return strconv.Itoa(int(ik))
}

func mkmap(r *node) *Map {
	var num = uint(r.count())
	return &Map{num, r}
}

func mknod(i int, c colorType, ln, rn *node) *node {
	return &node{IntKey(i), i, c, ln, rn}
}

type KeyVal struct {
	Key MapKey
	Val interface{}
}

func genIntKeyVals(n int) []KeyVal {
	var kvs = make([]KeyVal, n)

	for i := 0; i < n; i++ {
		var x = (i + 1) * 10
		var k = IntKey(x)
		var v = x
		kvs[i] = KeyVal{k, v}
	}

	return kvs
}

func randomizeKeyVals(kvs []KeyVal) []KeyVal {
	var randKvs = make([]KeyVal, len(kvs))
	copy(randKvs, kvs)
	//var randKvs = kvs

	//From: https://en.wikipedia.org/wiki/Fisher%E2%80%93Yates_shuffle#The_modern_algorithm
	for i := len(randKvs) - 1; i > 0; i-- {
		var j = rand.Intn(i + 1)
		randKvs[i], randKvs[j] = randKvs[j], randKvs[i]
	}

	return randKvs
}

func buildMap(kvs []KeyVal) *Map {
	var m = New()
	for _, kv := range kvs {
		m = m.Put(kv.Key, kv.Val)
	}
	return m
}
