package set_test

import (
	//"testing"
	"log"
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
