package tests

import (
	"net"
	"net/rpc"
)

// Listens for args sent from the CLI client via RPC.
type ArgListener struct {
	listener net.Listener
	argsCh   chan []string
	server   *rpc.Server
	Port     int
}

func newArgListener() (*ArgListener, error) {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		return nil, err
	}
	a := &ArgListener{
		listener: listener,
		argsCh:   make(chan []string, 1),
		server:   rpc.NewServer(),
		Port:     listener.Addr().(*net.TCPAddr).Port,
	}
	if err := a.server.RegisterName("Args", a); err != nil {
		return nil, err
	}
	go a.server.Accept(listener)
	return a, nil
}

func (a *ArgListener) Put(args []string, _ *struct{}) error {
	a.argsCh <- args
	return nil
}

func (a *ArgListener) get() []string {
	return <-a.argsCh
}

func (a *ArgListener) close() error {
	return a.listener.Close()
}
