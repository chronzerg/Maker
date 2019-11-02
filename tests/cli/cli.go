package main

import (
	"net/rpc"
	"os"
)

func main() {
	client, err := rpc.Dial("tcp", "127.0.0.1:"+os.Args[1])
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = client.Close()
	}()
	err = client.Call("Args.Put", os.Args[1:], nil)
	if err != nil {
		panic(err)
	}
}
