package tests

import (
	"github.com/pkg/errors"
	"net"
	"net/rpc"
)

// Listens for args sent from the CLI client via RPC.
type ArgListener struct {
	listener net.Listener
	server   *rpc.Server
	port     int

	argsCh chan Args
	args   map[string][]string
}

type Args struct {
	Name string
	Args []string
}

func newArgListener() *ArgListener {
	const errMsg = "failed to construct ArgListener"
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		panic(errors.Wrap(err, errMsg))
	}
	a := &ArgListener{
		listener: listener,
		server:   rpc.NewServer(),
		port:     listener.Addr().(*net.TCPAddr).Port,
		argsCh:   make(chan Args),
		args:     make(map[string][]string),
	}
	if err := a.server.RegisterName("Args", a); err != nil {
		panic(errors.Wrap(err, errMsg))
	}
	go a.server.Accept(listener)
	go a.doPut()
	return a
}

func (a *ArgListener) doPut() {
	for {
		arg := <-a.argsCh
		a.args[arg.Name] = arg.Args
	}
}

func (a *ArgListener) Put(args Args, _ *struct{}) error {
	a.argsCh <- args
	return nil
}

func (a *ArgListener) close() {
	err := a.listener.Close()
	if err != nil {
		panic(errors.Wrap(err, "failed to close ArgListener"))
	}
}
