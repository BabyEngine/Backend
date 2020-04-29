package networking

import (
    "fmt"
    "github.com/BabyEngine/Backend/logger"
    "github.com/gorilla/websocket"
    "net/http"
    "sync"
    "sync/atomic"
)

type mWebsocketServer struct {
    opts     *Options
    wg       sync.WaitGroup
    totalQPS uint32
    tx       uint64 // transfer bytes
    rx       uint64 // receive bytes
    tp       uint64 // transfer packet
    rp       uint64 // receive packet
    ws       *websocket.Upgrader
}

func (s *mWebsocketServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    c, err := s.ws.Upgrade(w, r, nil)
    if err != nil {
        logger.Debugf("upgrade:%v", err)
        return
    }
    s.handleKCPConn(c)
}

func (s *mWebsocketServer) Init() {
    s.ws = &websocket.Upgrader{
        CheckOrigin: func(r *http.Request) bool {
            return true
        },
    }
}

func (s *mWebsocketServer) Serve(address string) error {
    if s.opts.TLSEnable {
        err := http.ListenAndServeTLS(address, s.opts.TLSCert, s.opts.TLSKey, s)
        if err != nil {
            logger.Debugf("%v", err)
        }
        return err
    } else {
        err := http.ListenAndServe(address, s)
        if err != nil {
            logger.Debugf("%v", err)
        }
        return err
    }
}

func (s *mWebsocketServer) handleKCPConn(conn *websocket.Conn) {
    client := &mWebsocketClient{
        conn:   conn,
        opts:   s.opts,
        server: s,
    }
    client.init()
    s.wg.Add(1)
    defer func() {
        s.wg.Done()
        s.opts.Handler.OnClose(client)
    }()

    client.Serve()
}

func (s *mWebsocketServer) checkClient() {
    var (
        checkList []*mWebsocketClient
        deathList []*mWebsocketClient
    )
    if s.opts == nil || s.opts.Handler == nil { return }
    clients := s.opts.Handler.GetAllClient()

    for _, cli := range clients {
        if cc, ok := cli.(*mWebsocketClient); ok {
            checkList = append(checkList, cc)
        }
    }

    for _, cli := range checkList {
        if !cli.IsAlive() {
            deathList = append(deathList, cli)
        }
    }

    if len(deathList) > 0 {
        for _, cli := range deathList {
            s.opts.Handler.OnClose(cli)
        }
    }
    totalCount := len(clients)
    deathCount := len(deathList)
    if deathCount > 0 {
        fmt.Printf("当前客户端 总数:%d 死亡:%d\n", totalCount, deathCount)
    }
}

// 网络统计
func (s *mWebsocketServer) onNetStat(aType int, n, l uint64) {
    if aType == 0 {
        atomic.AddUint64(&s.rp, n)
        atomic.AddUint64(&s.rx, l)
    } else if aType == 1 {
        atomic.AddUint64(&s.tp, n)
        atomic.AddUint64(&s.tx, l)
    }
}
