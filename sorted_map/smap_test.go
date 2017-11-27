package sorted_map_test

import (
	"log"
	"math/rand"
	"os"
	"strconv"

	"github.com/lleo/go-functional-collections/sorted_map"
	"github.com/lleo/stringutil"
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

	log.Println("TESTING HAS STARTED...")
}

var Inc = stringutil.Lower.Inc

type StringKey string

func (sk StringKey) Less(o sorted_map.MapKey) bool {
	var osk, ok = o.(StringKey)
	if !ok {
		panic("o is not a StringKey")
	}
	return sk < osk
}

func (sk StringKey) String() string {
	return string(sk)
}

type IntKey int

func (ik IntKey) Less(o sorted_map.MapKey) bool {
	var oik, ok = o.(IntKey)
	if !ok {
		panic("o is not a IntKey")
	}
	return ik < oik
}

func (ik IntKey) String() string {
	return strconv.Itoa(int(ik))
}

type KeyVal struct {
	key sorted_map.MapKey
	val interface{}
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

func buildMap(kvs []KeyVal) *sorted_map.Map {
	var m = sorted_map.New()
	for _, kv := range kvs {
		m = m.Put(kv.key, kv.val)
	}
	return m
}
