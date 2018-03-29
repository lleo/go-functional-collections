package fmap_test

import (
	"log"
	"math/rand"
	"os"

	"github.com/lleo/go-functional-collections/fmap"
	"github.com/lleo/go-functional-collections/key"
	"github.com/lleo/stringutil"
	"github.com/pkg/errors"
)

func init() {
	log.SetFlags(log.Lshortfile)

	var logFileName = "test.log"
	var logFile, err = os.Create(logFileName)
	if err != nil {
		log.Fatal(errors.Wrapf(err, "failed to os.Create(%q)", logFileName))
	}
	log.SetOutput(logFile)
}

var Inc = stringutil.Lower.Inc

type KeyVal = fmap.KeyVal

func buildKvs(num int) []KeyVal {
	var kvs = make([]KeyVal, num)

	var keyStr = "a"
	for i := 0; i < num; i++ {
		kvs[i].Key = key.Str(keyStr)
		kvs[i].Val = i
		keyStr = Inc(keyStr)
	}

	return kvs
}

func buildStrings(num int) []string {
	var strs = make([]string, num)

	var str = "a"
	for i := 0; i < num; i++ {
		strs[i] = str
		str = Inc(str)
	}

	return strs
}

func buildKeys(num int) []key.Hash {
	var keys = make([]key.Hash, num)
	for i, s := range buildStrings(num) {
		keys[i] = key.Str(s)
	}
	return keys
}

func buildKvsFromStrings(strs []string) []KeyVal {
	var kvs = make([]KeyVal, len(strs))

	for i := 0; i < len(strs); i++ {
		kvs[i].Key = key.Str(strs[i])
		kvs[i].Val = i
	}

	return kvs
}

func randomizeKvs(kvs []KeyVal) []KeyVal {
	var randKvs = make([]KeyVal, len(kvs))
	copy(randKvs, kvs)

	// randomize kvs
	// https://en.wikipedia.org/wiki/Fisher–Yates_shuffle#The_modern_algorithm
	for i := len(randKvs) - 1; i > 0; i-- {
		var j = rand.Intn(i + 1)
		randKvs[i], randKvs[j] = randKvs[j], randKvs[i]
	}

	return randKvs
}

func buildMap(kvs []KeyVal) *fmap.Map {
	var m = fmap.New()
	for _, kv := range kvs {
		m = m.Put(kv.Key, kv.Val)
	}
	return m
}
