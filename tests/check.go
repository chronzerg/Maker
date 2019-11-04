package tests

import (
	"encoding/gob"
	"github.com/pkg/errors"
	"io"
	"log"
	"os"
	"path/filepath"
	"testing"
)

const dataDir = "testdata"

func checkArgs(t *testing.T, name string, args []Invocation) {
	path := filepath.Join(dataDir, name)
	file, err := os.Open(path)
	if err != nil {
		panic(errors.Wrap(err, "failed to open arg data"))
	}
	dec := gob.NewDecoder(file)
	for i := 0; true; i++ {
		var inv Invocation
		if err := dec.Decode(&inv); err != nil {
			if err == io.EOF {
				return
			}
			panic(errors.Wrap(err, "failed to decode arg data"))
		}
		log.Printf("expected: %v", inv)
		if i >= len(args) {
			t.Fatal("actual: no more")
		}
		if inv != args[i] {
			t.Fatalf("actual: %v", args[i])
		}
	}
}

func saveArgs(name string, args []Invocation) {
	err := os.MkdirAll(dataDir, os.ModePerm)
	if err != nil {
		panic(errors.Wrap(err, "failed to create arg data dir"))
	}
	path := filepath.Join(dataDir, name)
	file, err := os.Create(path)
	if err != nil {
		panic(errors.Wrap(err, "failed to create arg data file"))
	}
	enc := gob.NewEncoder(file)
	for _, inv := range args {
		log.Printf("expected: %v", inv)
		if err := enc.Encode(&inv); err != nil {
			panic(errors.Wrap(err, "failed to encode arg data"))
		}
	}
}
