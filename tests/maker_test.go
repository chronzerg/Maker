package tests

// Tests panic with any inter-process problems. For instance,
// if RPC or forking fails. These actions facilitate the tests
// but are not the test invariants themselves. The testing.T
// methods are reserved for when test invariants are violated.

import (
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

const cliExec = "cli/cli.run"
const makefile = "../makefile"

func TestMain(m *testing.M) {
	log.SetFlags(0)
	err := exec.Command("go", "build", "-o", cliExec, "cli/cli.go").Run()
	if err != nil {
		panic(errors.Wrap(err, "failed to build CLI tool"))
	}
	os.Exit(m.Run())
}

func TestMaker(t *testing.T) {
	argListener := newArgListener()
	//defer argListener.close()

	e := MakeExecution{
		mocks:    []string{"cxx"},
		dir:      tempDir(),
		cliExec:  absPath(cliExec),
		makefile: absPath(makefile),
		argPort:  argListener.port,
	}

	e.Run(t)

	fmt.Println(argListener.args)
}

func absPath(path string) string {
	abs, err := filepath.Abs(path)
	if err != nil {
		panic(errors.Wrap(err, "failed to get absolute path"))
	}
	return abs
}

func tempDir() string {
	dir, err := ioutil.TempDir("", "")
	if err != nil {
		panic(errors.Wrap(err, "failed to create temp dir"))
	}
	return dir
}
