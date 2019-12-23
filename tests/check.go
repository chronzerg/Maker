package tests

import (
	"encoding/gob"
	"github.com/pkg/errors"
	"io"
	"log"
	"os"
	"path/filepath"
)

const dataDir = "testdata"

func CheckArgs(name string, args []Invocation) int {
	file, err := os.Open(filepath.Join(dataDir, name))
	if err != nil {
		panic(errors.Wrap(err, "failed to open arg data"))
	}

	dec := gob.NewDecoder(file)
	var invocations []Invocation
	for {
		var inv Invocation
		if err := dec.Decode(&inv); err != nil {
			if err == io.EOF {
				break
			}
			panic(errors.Wrap(err, "failed to decode arg data"))
		}
		invocations = append(invocations, inv)
	}

	exitCode := 0
	for i, inv := range invocations {
		if i >= len(args) {
			log.Printf("FAIL\n expect: %v\n actual:\n", inv)
			exitCode = 1
		} else if inv != args[i] {
			log.Printf("FAIL\n expect: %v\n actual: %v\n", inv, args[i])
			exitCode = 1
		} else {
			log.Printf("PASS %v", inv)
		}
	}
	for i := len(invocations); i < len(args); i++ {
		log.Printf("FAIL\n expect:\n actual: %v\n", args[i])
		exitCode = 1
	}

	return exitCode
}

func SaveArgs(name string, args []Invocation) {
	err := os.MkdirAll(dataDir, os.ModePerm)
	if err != nil {
		panic(errors.Wrap(err, "failed to create arg data dir"))
	}

	file, err := os.Create(filepath.Join(dataDir, name))
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
