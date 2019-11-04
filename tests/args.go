package tests

import (
	"context"
	"github.com/pkg/errors"
	"net"
	"net/rpc"
)

const ListenerName = "ArgsListener"

type Invocation struct {
	Name string
	Args string
}

// Listens for args sent from the CLI client via RPC.
type ArgListener struct {
	cancel context.CancelFunc

	listener net.Listener
	port     int

	argsCh chan Invocation
	args   []Invocation
}

func newArgListener() *ArgListener {
	const errMsg = "failed to construct ArgListener"

	ctx, cancel := context.WithCancel(context.Background())

	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		panic(errors.Wrap(err, errMsg))
	}

	a := &ArgListener{
		cancel:   cancel,
		listener: listener,
		port:     listener.Addr().(*net.TCPAddr).Port,
		argsCh:   make(chan Invocation),
		args:     nil,
	}

	server := rpc.NewServer()
	if err := server.RegisterName(ListenerName, a); err != nil {
		panic(errors.Wrap(err, errMsg))
	}

	go func() {
		for {
			conn, err := listener.Accept()
			if ctx.Err() != nil {
				return
			}
			if err != nil {
				panic(errors.Wrap(err, "TCP listener failed"))
			}
			go server.ServeConn(conn)
		}
	}()

	go func() {
		for {
			select {
			case arg := <-a.argsCh:
				a.args = append(a.args, arg)
			case <-ctx.Done():
				return
			}
		}
	}()

	return a
}

func (a *ArgListener) Put(args Invocation, _ *struct{}) error {
	a.argsCh <- args
	return nil
}

func (a *ArgListener) close() {
	a.cancel()
	err := a.listener.Close()
	if err != nil {
		panic(errors.Wrap(err, "failed to close ArgListener"))
	}
}
