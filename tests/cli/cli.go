package main

import (
	"github.com/janderland/Maker/tests"
	"net/rpc"
	"os"
	"strings"
)

func main() {
	client, err := rpc.Dial("tcp", "127.0.0.1:"+os.Args[1])
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = client.Close()
	}()
	err = client.Call(tests.ListenerName+".Put", tests.Invocation{
		Name: os.Args[2],
		Args: strings.Join(os.Args[3:], " "),
	}, nil)
	if err != nil {
		panic(err)
	}
}
