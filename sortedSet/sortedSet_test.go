package sortedSet

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

func mkset(r *node) *Set {
	var num = r.count()
	return &Set{num, r}
}

func mknod(i int, c colorType, ln, rn *node) *node {
	return &node{key.Int(i), c, ln, rn}
}

func buildKeys(n int) []key.Sort {
	var keys = make([]key.Sort, n)

	for i := 0; i < n; i++ {
		var x = (i + 1) * 10
		keys[i] = key.Int(x)
	}

	return keys
}

func randomizeKeys(keys []key.Sort) []key.Sort {
	var randKeys = make([]key.Sort, len(keys))
	copy(randKeys, keys)
	//var randKeys = keys

	//From: https://en.wikipedia.org/wiki/Fisher%E2%80%93Yates_shuffle#The_modern_algorithm
	for i := len(randKeys) - 1; i > 0; i-- {
		var j = rand.Intn(i + 1)
		randKeys[i], randKeys[j] = randKeys[j], randKeys[i]
	}

	return randKeys
}

func buildSet(kvs []key.Sort) *Set {
	var s = New()
	for _, key := range kvs {
		s = s.Set(key)
	}
	return s
}
