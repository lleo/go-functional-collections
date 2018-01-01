package hash

import (
	"log"
	"os"

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

	log.Printf("NumIndexBits = %d\n", NumIndexBits)
	log.Printf("hashSize     = %d\n", hashSize)
	log.Printf("remainder    = %d\n", remainder)
	log.Printf("DepthLimit   = %d\n", DepthLimit)
	log.Printf("MaxDepth     = %d\n", MaxDepth)
	log.Printf("IndexLimit   = %d\n", IndexLimit)
	log.Printf("MaxIndex     = %d\n", MaxIndex)
}
