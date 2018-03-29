package key_test

import (
	"log"
	"testing"

	"github.com/lleo/go-functional-collections/key"
)

func TestIntCreate(t *testing.T) {
	var ik = key.Int(10)
	log.Printf("TestIntCreate: ik=%s\n", ik.String())
}

func TestIntSort(t *testing.T) {

}

func TestIntHash(t *testing.T) {

}
