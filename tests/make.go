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
	"testing"
)

const makefile = "../makefile"

var mocks = []string{"cxx", "ar"}

type MakeMock struct {
	command *exec.Cmd
	dir     string
}

func newMock(argPort int, targets []string, opts map[string]string) *MakeMock {
	makefile, err := filepath.Abs(makefile)
	if err != nil {
		panic(errors.Wrap(err, "failed to get makefile path"))
	}
	cmd := exec.Command("make", append([]string{"-f", makefile}, targets...)...)

	dir, err := ioutil.TempDir("", "")
	if err != nil {
		panic(errors.Wrap(err, "failed to create temp dir"))
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

	return &MakeMock{
		command: cmd,
		dir:     dir,
	}
}

func (m *MakeMock) Run(t *testing.T) {
	log.Println(m.command)

	stdoutPipe, err := m.command.StdoutPipe()
	if err != nil {
		panic(errors.Wrap(err, "failed to open stdout pipe"))
	}

	stderrPipe, err := m.command.StderrPipe()
	if err != nil {
		panic(errors.Wrap(err, "failed to open stderr pipe"))
	}

	stdoutCh := make(chan string)
	stderrCh := make(chan string)

	go doRead(stdoutPipe, stdoutCh)
	go doRead(stderrPipe, stderrCh)

	loggingDone := make(chan struct{})
	go func() {
		defer close(loggingDone)
		lgr := log.New(os.Stderr, "| ", 0)
		for stdoutCh != nil || stderrCh != nil {
			var line string
			var open bool

			select {
			case line, open = <-stdoutCh:
				if !open {
					stdoutCh = nil
					break
				}
			case line, open = <-stderrCh:
				if !open {
					stderrCh = nil
					break
				}
			}

			if len(line) > 0 {
				lgr.Println(strings.ReplaceAll(line, "\n", ""))
			}
		}
	}()

	err = m.command.Start()
	if err != nil {
		panic(errors.Wrap(err, "failed to run make"))
	}
	err = m.command.Wait()
	<-loggingDone
	if err != nil {
		t.Fatal(errors.Wrap(err, "make returned an error"))
	}
}

func doRead(inReader io.Reader, outCh chan<- string) {
	defer close(outCh)
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
