package main

import (
    "fmt"
    "github.com/BabyEngine/Backend/rpc"
)

func main()  {
    cli := rpc.NewGRPCClient("127.0.0.1:9981")
    if err := cli.Connect(); err != nil {
        fmt.Println(err)
        return
    }
    if reply, err := cli.Call("api.echo", []byte("hello")); err != nil {
        fmt.Println(err)
        return
    } else {
        fmt.Println(reply.Code, string(reply.Data))
    }

}
