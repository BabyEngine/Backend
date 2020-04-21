package networking

import (
    "fmt"
    "io"
    "net"
    "net/http"
)

type mHTTPServer struct {
    opts         *Options
    //wg           sync.WaitGroup
    totalQPS     uint32
    tx           uint64 // transfer bytes
    rx           uint64 // receive bytes
    tp           uint64 // transfer packet
    rp           uint64 // receive packet
    //clients      map[int64]*mHTTPClient
    //clientM      sync.RWMutex
   // clientId     int64
    serverCloser io.Closer
}

func (s *mHTTPServer) Init() {
 //   s.clients = make(map[int64]*mHTTPClient)
}

func (s *mHTTPServer) Serve(addr string) error {
    closer, err := HTTPListenAndServeWithClose(addr, s)
    if err != nil {
        return err
    }
    s.serverCloser = closer
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
