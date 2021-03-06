package sortedMap

import (
	"log"
	"math/rand"
	"os"

	"github.com/lleo/go-functional-collections/key"
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

func mkmap(r *node) *Map {
	var num = r.count()
	return &Map{num, r}
}

func mknod(i int, c colorType, ln, rn *node) *node {
	return &node{key.Int(i), i, c, ln, rn}
}

type KeyVal struct {
	Key key.Sort
	Val interface{}
}

func genIntKeyVals(n int) []KeyVal {
	var kvs = make([]KeyVal, n)

	for i := 0; i < n; i++ {
		var x = (i + 1) * 10
		var k = key.Int(x)
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
