package tests

import (
	"github.com/pkg/errors"
	"os/exec"
)

const cliExec = "cli/cli.run"

func BuildCLI() {
	err := exec.Command("go", "build", "-o", cliExec, "cli/cli.go").Run()
	if err != nil {
		panic(errors.Wrap(err, "failed to build CLI tool"))
	}
}
