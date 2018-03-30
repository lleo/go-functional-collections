package set_test

import (
	//"testing"
	"log"
	"math/rand"
	"os"

	"github.com/lleo/go-functional-collections/key"
	"github.com/lleo/go-functional-collections/set"
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

//type StringKey = key.Str

func buildKeys(num int) []key.Hash {
	var keys = make([]key.Hash, num)

	var keyStr = "a"
	for i := 0; i < num; i++ {
		keys[i] = key.Str(keyStr)
		keyStr = Inc(keyStr)
	}

	return keys
}

func buildKeysByN(num int, n int) []key.Hash {
	var keys = make([]key.Hash, num)

	var keyStr = "a"
	for i := 0; i < num; i++ {
		keys[i] = key.Str(keyStr)
		for j := 0; j < n; j++ {
			keyStr = Inc(keyStr)
		}
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

func buildKeysFromStrings(strs []string) []key.Hash {
	var keys = make([]key.Hash, len(strs))

	for i := 0; i < len(strs); i++ {
		keys[i] = key.Str(strs[i])
	}

	return keys
}

func randomizeKeys(keys []key.Hash) []key.Hash {
	var randKeys = make([]key.Hash, len(keys))
	copy(randKeys, keys)

	//randomize keys
	//https://en.wikipedia.org/wiki/Fisher–Yates_shuffle#The_modern_algorithm
	for i := len(randKeys) - 1; i > 0; i-- {
		var j = rand.Intn(i + 1)
		randKeys[i], randKeys[j] = randKeys[j], randKeys[i]
	}

	return randKeys
}

func buildSet(keys []key.Hash) *set.Set {
	var m = set.New()
	for _, key := range keys {
		m = m.Set(key)
	}
	return m
}
