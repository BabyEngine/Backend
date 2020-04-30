package rpc

import (
    "io"
    "net"
    "net/http"
    "github.com/gorilla/mux"
    "github.com/gorilla/rpc"
    "github.com/gorilla/rpc/json"
)

type Server struct {
    rpcServer  *rpc.Server
    httpServer *http.Server
    rpc        RPC
    closer     io.Closer
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    io.WriteString(w, "rpc online!")
}

func NewServer(cb func(request Request, reply *Reply) error) *Server {
    s := &Server{}
    s.rpc.Handler = cb
    return s
}

func (s *Server) ListenServe(address string) error {
    return s.ListenServeTLS(address, "", "")
}

func (s *Server) ListenServeTLS(address string, crt string, key string) error {
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

func (s *Server) Stop() error {
    return s.closer.Close()
}
