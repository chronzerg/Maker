package tests

import (
	"bufio"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"testing"
)

type MakeExecution struct {
	mocks    []string
	dir      string
	cliExec  string
	makefile string
	argPort  int
}

func (m *MakeExecution) Run(t *testing.T) {
	cmd := exec.Command("make", "-f", m.makefile)
	cmd.Env = m.env()
	cmd.Dir = m.dir
	log.Println(cmd)

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		panic(errors.Wrap(err, "failed to open stdout pipe"))
	}

	stderrPipe, err := cmd.StderrPipe()
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

	err = cmd.Start()
	if err != nil {
		panic(errors.Wrap(err, "failed to run make"))
	}
	err = cmd.Wait()
	<-loggingDone
	if err != nil {
		t.Fatal(errors.Wrap(err, "make returned an error"))
	}
}

func (m *MakeExecution) env() []string {
	env := make([]string, len(m.mocks))
	for i, mock := range m.mocks {
		env[i] = fmt.Sprintf(`%s="%s %s %d"`,
			mock, m.cliExec, mock, m.argPort)
	}
	return env
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
