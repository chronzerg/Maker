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
		Module(tests.ModSpec{
			Type:      "exec",
			Deps:      "d2",
			CompFlags: "-X{0} -wasted",
			LinkFlags: "-redflag -nope",
		}).
		File("f1.cpp", "#include\"f3.hpp\"").
		File("f2.cpp", "#include\"f3.hpp\"").
		File("f3.hpp", "").
		File("f4.hpp", "")

	fs.Dir("d2").
		Module(tests.ModSpec{
			Type:      "slib",
			Deps:      "",
			CompFlags: "-Y{1} -sober",
			LinkFlags: "-greencard -yup",
		}).
		File("f1.cpp", "#include\"f2.hpp\"").
		File("f2.hpp", "")

	fs.Dir("d3").
		Module(tests.ModSpec{
			Type:      "exec",
			Deps:      "d2",
			CompFlags: "-C{9} -tired",
			LinkFlags: "-purpledrank -maybe",
		}).
		File("f1.cpp", "")

	makeCmd.Run()

	if save {
		tests.SaveArgs("test", args.Args)
	} else {
		os.Exit(tests.CheckArgs("test", args.Args))
	}
}
