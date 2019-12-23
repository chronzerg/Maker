package main

import (
	"flag"
	"github.com/janderland/Maker/tests"
	"log"
	"os"
)

func main() {
	log.SetFlags(0)

	save := flag.Bool("save", false, "overwrites the test expectations")
	flag.Parse()

	tests.BuildCLI()

	TestFramework(*save)
}

func TestFramework(save bool) {
	args := tests.NewArgListener()
	defer args.Close()

	makeCmd := tests.NewMakeCmd(args.Port, nil, nil)
	fs := tests.Files(makeCmd.Dir)

	fs.Dir("d1").
		File("f1.cpp", "#include\"f3.hpp\"").
		File("f2.cpp", "#include\"f3.hpp\"").
		File("f3.hpp", "").
		File("f4.hpp", "").
		File("makefile", `
moduleType = exec
moduleDeps = d2
moduleCompFlags = -X{0} -wasted
moduleLinkFlags = -redflag -nope
`)

	fs.Dir("d2").
		File("f1.cpp", "#include\"f2.hpp\"").
		File("f2.hpp", "").
		File("makefile", `
moduleType = slib
moduleDeps =
moduleCompFlags = -Y{1} -sober
moduleLinkFlags = -greencard -yup
`)

	fs.Dir("d3").
		File("f1.cpp", "").
		File("makefile", `
moduleType = exec
moduleDeps = d2
moduleCompFlags = -C{9} -tired
moduleLinkFlags = -purpledrank -maybe
`)

	makeCmd.Run()

	if save {
		tests.SaveArgs("test", args.Args)
	} else {
		os.Exit(tests.CheckArgs("test", args.Args))
	}
}
