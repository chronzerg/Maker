package tests

import (
	"fmt"
	"github.com/pkg/errors"
	"os"
	"os/exec"
	"strconv"
	"testing"
)

const cliExec = "cli/cli.run"

func TestMain(m *testing.M) {
	err := exec.Command("go", "build", "-o", cliExec, "cli/cli.go").Run()
	if err != nil {
		panic(errors.Wrap(err, "failed to build CLI tool"))
	}
	os.Exit(m.Run())
}

func TestMaker(t *testing.T) {
	a, err := newArgListener()
	if err != nil {
		panic(errors.Wrap(err, "failed to construct ArgListener"))
	}
	defer func() {
		err := a.close()
		if err != nil {
			panic(errors.Wrap(err, "failed to close ArgListener"))
		}
	}()
	err = exec.Command(cliExec, strconv.Itoa(a.Port), "hello").Run()
	if err != nil {
		panic(errors.Wrap(err, "failed to run CLI tool"))
	}
	fmt.Println(a.get())
}
