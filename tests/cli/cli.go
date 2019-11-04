package main

import (
	"github.com/janderland/Maker/tests"
	"net/rpc"
	"os"
)

func main() {
	client, err := rpc.Dial("tcp", "127.0.0.1:"+os.Args[2])
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = client.Close()
	}()
	err = client.Call("Args.Put", tests.Args{
		Name: os.Args[1],
		Args: os.Args[3:],
	}, nil)
	if err != nil {
		panic(err)
	}
}
