package main

import (
    "fmt"
    "github.com/BabyEngine/Backend/rpc"
)

func main()  {
    s := rpc.NewGRPCServer(func(request rpc.Request, reply *rpc.Reply) error {
        reply.Code = 11
        reply.Data = request.Data
        return nil
    })
    if err := s.ListenServe(":9981"); err != nil {
        fmt.Println(err)
    }
}
