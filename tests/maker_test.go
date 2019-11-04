package tests

// Tests panic with any inter-process problems. For instance,
// if RPC or forking fails. These actions facilitate the tests
// but are not the test invariants themselves. The testing.T
// methods are reserved for when test invariants are violated.

import (
	"flag"
	"log"
	"os"
	"testing"
)

var save *bool

func TestMain(m *testing.M) {
	log.SetFlags(0)
	flag.Parse()
	buildCLI()
	os.Exit(m.Run())
}

func TestFramework(t *testing.T) {
	args := newArgListener()
	defer args.close()

	mock := newMock(args.port, nil, nil)

	fs := Files(mock.dir)
	fs.file("config.mkr", `
$(call exec, mod,
	# Dependencies
	,
	# Compile Flags
	,
	# Linking Flags
);`)

	fs.dir("mod").file("file.cpp", "")

	mock.Run(t)

	log.Println(args.args)
}
