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
	Port     int

	argsCh chan Invocation
	Args   []Invocation
}

func NewArgListener() *ArgListener {
	const errMsg = "failed to construct ArgListener"

	ctx, cancel := context.WithCancel(context.Background())

	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		panic(errors.Wrap(err, errMsg))
	}

	a := &ArgListener{
		cancel:   cancel,
		listener: listener,
		Port:     listener.Addr().(*net.TCPAddr).Port,
		argsCh:   make(chan Invocation),
		Args:     nil,
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
				a.Args = append(a.Args, arg)
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

func (a *ArgListener) Close() {
	a.cancel()
	if err := a.listener.Close(); err != nil {
		panic(err)
	}
}
