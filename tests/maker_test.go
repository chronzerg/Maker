package tests

// Tests panic with any inter-process problems. For instance,
// if RPC or forking fails. These actions facilitate the tests
// but are not the test invariants themselves. The testing.T
// methods are reserved for when test invariants are violated.

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"
)

func TestMain(m *testing.M) {
	log.SetFlags(0)
	buildCLI()
	os.Exit(m.Run())
}

func TestFramework(t *testing.T) {
	args := newArgListener()
	defer args.close()

	mock := newMock(args.port, nil, nil)

	const cfg = `
$(call exec, mod,
	# Dependencies
	,
	# Compile Flags
	,
	# Linking Flags
);`

	cfgPath := filepath.Join(mock.dir, "config.mkr")
	err := ioutil.WriteFile(cfgPath, []byte(cfg), 0644)
	if err != nil {
		panic(err)
	}

	modPath := filepath.Join(mock.dir, "mod")
	err = os.MkdirAll(modPath, os.ModePerm)
	if err != nil {
		panic(err)
	}

	srcPath := filepath.Join(modPath, "file.cpp")
	err = ioutil.WriteFile(srcPath, nil, 0644)
	if err != nil {
		panic(err)
	}

	mock.Run(t)

	log.Println(args.args)
}
