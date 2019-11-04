package tests

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

type Files string

func (b Files) dir(path string) Files {
	modPath := filepath.Join(string(b), path)
	err := os.MkdirAll(modPath, os.ModePerm)
	if err != nil {
		panic(err)
	}
	return Files(modPath)
}

func (b Files) file(name string, contents string) Files {
	srcPath := filepath.Join(string(b), name)
	err := ioutil.WriteFile(srcPath, []byte(contents), 0644)
	if err != nil {
		panic(err)
	}
	return b
}
