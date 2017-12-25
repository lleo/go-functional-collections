package fmap_test

import (
	//"testing"
	"log"
	"math/rand"
	"os"

	"github.com/lleo/go-functional-collections/fmap"
	"github.com/lleo/go-functional-collections/fmap/hash"
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

type StringKey string

func (sk StringKey) Hash() hash.HashVal {
	return hash.CalcHash([]byte(sk))
}

func (sk StringKey) Equals(other fmap.MapKey) bool {
	var osk, ok = other.(StringKey)
	if !ok {
		return false
	}
	return sk == osk
}

func (sk StringKey) String() string {
	return string(sk)
}

type keyVal struct {
	Key fmap.MapKey
	Val interface{}
}

func buildKvs(num int) []keyVal {
	var kvs = make([]keyVal, num)

	var keyStr = "a"
	for i := 0; i < num; i++ {
		kvs[i].Key = StringKey(keyStr)
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

func buildKvsFromStrings(strs []string) []keyVal {
	var kvs = make([]keyVal, len(strs))

	for i := 0; i < len(strs); i++ {
		kvs[i].Key = StringKey(strs[i])
		kvs[i].Val = i
	}

	return kvs
}

func randomizeKvs(kvs []keyVal) []keyVal {
	var randKvs = make([]keyVal, len(kvs))
	copy(randKvs, kvs)

	//randomize kvs
	//https://en.wikipedia.org/wiki/Fisher–Yates_shuffle#The_modern_algorithm
	for i := len(randKvs) - 1; i > 0; i-- {
		var j = rand.Intn(i + 1)
		randKvs[i], randKvs[j] = randKvs[j], randKvs[i]
	}

	return randKvs
}

func buildMap(kvs []keyVal) *fmap.Map {
	var m = fmap.New()
	for _, kv := range kvs {
		m = m.Put(kv.Key, kv.Val)
	}
	return m
}
