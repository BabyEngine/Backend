package rpc

import (
    "errors"
)

var (
    ErrorRPCHandler          = errors.New("RPC handler not found")
    ErrorRPCClientInvalid    = errors.New("RPC client invalid")
    ErrorRPCClientConnecting = errors.New("RPC client connecting")
)

type Request struct {
    Action string
    Data   []byte
}

type Reply struct {
    Code int
    Data []byte
}

type RPC struct {
    Handler func(request Request, reply *Reply) error
}

func (r *RPC) Call(request Request, reply *Reply) error {
    if r.Handler != nil {
        return r.Handler(request, reply)
    }
    return ErrorRPCHandler
}
