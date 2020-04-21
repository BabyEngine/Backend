package networking

import (
    "fmt"
    "github.com/BabyEngine/Backend/debugging"
    "github.com/xtaci/kcp-go"
    "net"
    "sync"
    "sync/atomic"
    "time"
)

type mKCPServer struct {
    opts     *Options
    wg       sync.WaitGroup
    totalQPS uint32
    tx       uint64 // transfer bytes
    rx       uint64 // receive bytes
    tp       uint64 // transfer packet
    rp       uint64 // receive packet
}

func (s *mKCPServer) Init() {
}

func (s *mKCPServer) Serve(address string) error {
    ln, err := kcp.Listen(address)
    if err != nil {
        return err
    }
    ticker := time.NewTicker(time.Second)
    defer func() {
        if ln != nil {
            if err := ln.Close(); err != nil {
                debugging.Log(err)
            }
        }
        ticker.Stop()
    }()

    go func() {
        for {
            select {
            case <-s.opts.Ctx.Done():
                if ln != nil {
                    if err := ln.Close(); err != nil {
                        debugging.Log(err)
                    }
                }
                ln = nil
            case <-ticker.C:
                s.checkClient()
            }
        }
    }()

    for {
        conn, err := ln.Accept()
        if err != nil {
            if ne, ok := err.(net.Error); ok && ne.Temporary() {
                // ignore code here
                continue
            }
            return err
        }
        go s.handleKCPConn(conn)
    }
}

func (s *mKCPServer) handleKCPConn(conn net.Conn) {
    client := &mKCPClient{
        conn: conn,
        opts: s.opts,
        server:s,
    }
    client.init()
    s.wg.Add(1)
    defer func() {
        s.wg.Done()
        s.opts.Handler.OnClose(client)
    }()

    client.Serve()
}

func (s *mKCPServer) checkClient() {
    var (
        checkList []*mKCPClient
        deathList []*mKCPClient
    )
    if s.opts == nil || s.opts.Handler == nil { return }
    clients := s.opts.Handler.GetAllClient()
    for _, cli := range clients {
        if cc, ok := cli.(*mKCPClient); ok {
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
func (s *mKCPServer) onNetStat(aType int, n, l uint64)  {
    if aType == 0 {
        atomic.AddUint64(&s.rp, n)
        atomic.AddUint64(&s.rx, l)
    } else if aType == 1 {
        atomic.AddUint64(&s.tp, n)
        atomic.AddUint64(&s.tx, l)
    }
}
