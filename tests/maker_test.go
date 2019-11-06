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

	save = flag.Bool("save", false, "overwrites the test expectations")
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
$(call slib, d1,
	# Dependencies
	,
	# Compile Flags
	,
	# Linking Flags
	--kill
);
$(call exec, d2,
	# Dependencies
	d1,
	#Compile Flags
	-x {0}+,
	#Linking Flags
	--hope --love
);`)

	fs.dir("d1").
		file("f1.cpp", "#include\"f3.hpp\"").
		file("f2.cpp", "#include\"f3.hpp\"").
		file("f3.hpp", "").
		file("f4.hpp", "")

	fs.dir("d2").
		file("f1.cpp", "#include\"f2.hpp\"").
		file("f2.hpp", "")

	mock.Run(t)

	if *save {
		saveArgs("test", args.args)
	} else {
		checkArgs(t, "test", args.args)
	}
}
