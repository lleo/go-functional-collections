package sorted_map_test

import (
	"log"
	"os"
	"strconv"

	"github.com/lleo/go-functional-collections/sorted_map"
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

	log.Println("IT HAS STARTED...")
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
