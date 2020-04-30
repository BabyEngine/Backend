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
type Client interface {
    Connect() error
    Call(string, []byte) (Reply, error)
    Disconnect() error
}
type Server interface {
    ListenServe(string) error
    Close() error
}

func (r *RPC) Call(req *http.Request, request *Request, reply *Reply) error {
    if r.Handler != nil {
        return r.Handler(*request, reply)
    }
    return ErrorRPCHandler
}

func NewClient(cType string, address string) Client {
    switch cType {
    case "jsonrpc":
        return NewJSONRPCClient(address)
    case "grpc":
        return NewGRPCClient(address)
    }
    return nil
}
func NewServer(cType string, cb func(request Request, reply *Reply) error) Server {
    switch cType {
    case "jsonrpc":
        return NewJSONRPCServer(cb)
    case "grpc":
        return NewGRPCServer(cb)
    }
    return nil
}
