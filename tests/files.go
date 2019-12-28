package tests

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Files string

func (f Files) Dir(path string) Files {
	modPath := filepath.Join(string(f), path)
	err := os.MkdirAll(modPath, os.ModePerm)
	if err != nil {
		panic(err)
	}
	return Files(modPath)
}

func (f Files) File(name string, contents string) Files {
	srcPath := filepath.Join(string(f), name)
	err := ioutil.WriteFile(srcPath, []byte(contents), 0644)
	if err != nil {
		panic(err)
	}
	return f
}

type ModSpec struct {
	Type      string
	Deps      string
	CompFlags string
	LinkFlags string
}

func (f Files) Module(spec ModSpec) Files {
	modVar := func(n string, v string) string {
		return fmt.Sprintf("%s = %s\n", n, v)
	}
	contents := modVar("moduleType", spec.Type)
	contents += modVar("moduleDeps", spec.Deps)
	contents += modVar("moduleCompFlags", spec.CompFlags)
	contents += modVar("moduleLinkFlags", spec.LinkFlags)
	return f.File("makefile", contents)
}
