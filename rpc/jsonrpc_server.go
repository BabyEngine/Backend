package rpc

import (
    "io"
    "net"
    "net/http"
    "github.com/gorilla/mux"
    "github.com/gorilla/rpc"
    "github.com/gorilla/rpc/json"
)

type JSONRPCServer struct {
    rpcServer  *rpc.Server
    httpServer *http.Server
    rpc        RPC
    closer     io.Closer
}

func (s *JSONRPCServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    io.WriteString(w, "rpc online!")
}

func NewJSONRPCServer(cb func(request Request, reply *Reply) error) *JSONRPCServer {
    s := &JSONRPCServer{}
    s.rpc.Handler = cb
    return s
}

func (s *JSONRPCServer) ListenServe(address string) error {
    return s.ListenServeTLS(address, "", "")
}

func (s *JSONRPCServer) ListenServeTLS(address string, crt string, key string) error {
    rpcServer := rpc.NewServer()
    rpcServer.RegisterCodec(json.NewCodec(), "application/json")
    rpcServer.RegisterCodec(json.NewCodec(), "application/json;charset=UTF-8")
    rpcServer.RegisterService(&s.rpc, "")
    r := mux.NewRouter()
    r.Handle("/jsonrpc", rpcServer)

    s.rpcServer = rpcServer
    ln, err := net.Listen("tcp", address)
    if err != nil {
        return err
    }
    s.closer = ln
    s.httpServer = &http.Server{
        Addr:    address,
        Handler: s.rpcServer,
    }

    if crt != "" && key != "" {
        return s.httpServer.ServeTLS(ln, crt, key)
    } else {
        return s.httpServer.Serve(ln)
    }
}

func (s *JSONRPCServer) Stop() error {
    return s.closer.Close()
}
