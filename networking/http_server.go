package networking

import (
    "fmt"
    "github.com/BabyEngine/Backend/logger"
    "io"
    "net"
    "net/http"
)

type mHTTPServer struct {
    opts         *Options
    totalQPS     uint32
    serverCloser io.Closer
}

func (s *mHTTPServer) Init() {
}

func (s *mHTTPServer) Serve(addr string) error {
    var (
        ln  net.Listener
        err error
    )
    srv := &http.Server{
        Handler: s,
    }
    if addr == "" {
        if s.opts.TLSEnable {
            addr = ":https"
        } else {
            addr = ":http"
        }
    }

    ln, err = net.Listen("tcp", addr)
    if err != nil {
        return err
    }

    go func() {
        if s.opts.TLSEnable {
            if err := srv.ServeTLS(ln, s.opts.TLSCert, s.opts.TLSKey); err != nil {
                logger.Debug(err)
            }
        } else {
            if err := srv.Serve(ln); err != nil {
                logger.Debug(err)
            }
        }
    }()

    s.serverCloser = ln
    return nil
}

func (s *mHTTPServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    client := &mHTTPClient{
        w:      w,
        r:      r,
        server: s,
        opts:   s.opts,
    }
    client.init()
    defer func() {
        s.opts.Handler.OnClose(client)
    }()
    client.Serve()
}
func HTTPListenAndServeWithClose(addr string, handler http.Handler) (io.Closer, error) {
    var (
        ln     net.Listener
        closer io.Closer
        err    error
    )
    s := &http.Server{
        Handler: handler,
    }
    if addr == "" {
        addr = ":http"
    }

    ln, err = net.Listen("tcp", addr)
    if err != nil {
        return nil, err
    }

    go func() {
        if err := s.Serve(ln); err != nil {
            fmt.Println(err)
        }
    }()
    closer = ln
    return closer, nil
}
