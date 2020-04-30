package rpc

import (
    "errors"
    "net/http"
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

func (r *RPC) Call(req *http.Request, request *Request, reply *Reply) error {
    if r.Handler != nil {
        return r.Handler(*request, reply)
    }
    return ErrorRPCHandler
}