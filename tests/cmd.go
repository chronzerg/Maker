package tests

import (
	"bufio"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const makefile = "../makefile"

var mocks = []string{"cxx", "ar"}

type MakeCmd struct{ *exec.Cmd }

func NewMakeCmd(argPort int, targets []string, opts map[string]string) MakeCmd {
	makefile, err := filepath.Abs(makefile)
	if err != nil {
		panic(errors.Wrap(err, "failed to get absolute path of makefile"))
	}
	cmd := exec.Command("make", append([]string{"-f", makefile}, targets...)...)

	dir, err := ioutil.TempDir("", "")
	if err != nil {
		panic(errors.Wrap(err, "failed to create working dir"))
	}
	cmd.Dir = dir

	cliExec, err := filepath.Abs(cliExec)
	if err != nil {
		panic(errors.Wrap(err, "failed to get CLI path"))
	}

	env := make([]string, len(mocks))
	for i, mock := range mocks {
		env[i] = fmt.Sprintf("%s=%s %d %s", mock, cliExec, argPort, mock)
	}
	for opt, val := range opts {
		env = append(env, fmt.Sprintf("%s=%s", opt, val))
	}
	cmd.Env = env

	return MakeCmd{cmd}
}

func (m *MakeCmd) Run() {
	log.Println(m)

	stdoutPipe, err := m.StdoutPipe()
	if err != nil {
		panic(errors.Wrap(err, "failed to open stdout pipe"))
	}

	stderrPipe, err := m.StderrPipe()
	if err != nil {
		panic(errors.Wrap(err, "failed to open stderr pipe"))
	}

	outCh := make(chan string)
	stdoutDone := make(chan struct{})
	stderrDone := make(chan struct{})
	go doRead(stdoutPipe, outCh, stdoutDone)
	go doRead(stderrPipe, outCh, stderrDone)
	go doLog(outCh)

	err = m.Start()
	if err != nil {
		panic(errors.Wrap(err, "failed to run make"))
	}

	err = m.Wait()
	<-stdoutDone
	<-stderrDone
	close(outCh)
	if err != nil {
		log.Fatal(errors.Wrap(err, "make returned an error"))
	}
}

func doRead(inReader io.Reader, outCh chan<- string, done chan<- struct{}) {
	defer close(done)
	reader := bufio.NewReader(inReader)
	for {
		line, err := reader.ReadString('\n')
		outCh <- line
		if err != nil {
			if err == io.EOF {
				return
			}
			panic(errors.Wrap(err, "failed while reading"))
		}
	}
}

func doLog(outputCh <-chan string) {
	lgr := log.New(os.Stderr, "| ", 0)
	for line := range outputCh {
		if len(line) > 0 {
			lgr.Println(strings.ReplaceAll(line, "\n", ""))
		}
	}
}
