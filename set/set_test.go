package set_test

import (
	//"testing"
	"log"
	"math/rand"
	"os"

	"github.com/lleo/go-functional-collections/set"
	"github.com/lleo/go-functional-collections/set/hash"
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

func (sk StringKey) Equals(other set.SetKey) bool {
	var osk, ok = other.(StringKey)
	if !ok {
		return false
	}
	return sk == osk
}

func (sk StringKey) String() string {
	return string(sk)
}

func buildKeys(num int) []set.SetKey {
	var keys = make([]set.SetKey, num)

	var keyStr = "a"
	for i := 0; i < num; i++ {
		keys[i] = StringKey(keyStr)
		keyStr = Inc(keyStr)
	}

	return keys
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

func buildKeysFromStrings(strs []string) []set.SetKey {
	var keys = make([]set.SetKey, len(strs))

	for i := 0; i < len(strs); i++ {
		keys[i] = StringKey(strs[i])
	}

	return keys
}

func randomizeKeys(keys []set.SetKey) []set.SetKey {
	var randKeys = make([]set.SetKey, len(keys))
	copy(randKeys, keys)

	//randomize keys
	//https://en.wikipedia.org/wiki/Fisher–Yates_shuffle#The_modern_algorithm
	for i := len(randKeys) - 1; i > 0; i-- {
		var j = rand.Intn(i + 1)
		randKeys[i], randKeys[j] = randKeys[j], randKeys[i]
	}

	return randKeys
}

func buildSet(keys []set.SetKey) *set.Set {
	var m = set.New()
	for _, key := range keys {
		m = m.Set(key)
	}
	return m
}
