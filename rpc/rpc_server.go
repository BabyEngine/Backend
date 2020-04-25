package rpc

import (
    "io"
    "net"
    "net/http"
    "net/rpc"
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
    if err := rpc.Register(&s.rpc); err != nil {
        return err
    }
    s.rpcServer = rpc.NewServer()
    s.rpcServer.Register(&s.rpc)
    s.rpcServer.HandleHTTP(rpc.DefaultRPCPath, rpc.DefaultDebugPath)
    ln, err := net.Listen("tcp", address)
    if err != nil {
        return err
    }
    s.closer = ln
    s.httpServer = &http.Server{Addr: address, Handler: s.rpcServer}
    s.httpServer.Serve(ln)

    if crt != "" && key != "" {
        return s.httpServer.ListenAndServeTLS(crt, key)
    } else {

        return s.httpServer.ListenAndServe()
    }
}

func (s *Server) Stop() error {
    return s.closer.Close()
}
