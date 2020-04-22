package networking

import (
    "encoding/base64"
    "fmt"
    "github.com/BabyEngine/Backend/debugging"
    "github.com/googollee/go-socket.io"
    "github.com/gorilla/websocket"
    "io"
    "net"
    "net/http"
    "sync"
    "sync/atomic"
)

type mSocketIOServer struct {
    opts     *Options
    wg       sync.WaitGroup
    totalQPS uint32
    tx       uint64 // transfer bytes
    rx       uint64 // receive bytes
    tp       uint64 // transfer packet
    rp       uint64 // receive packet
    server   *socketio.Server
    closer   io.Closer
}

func (s *mSocketIOServer) Init() {

}

func (s *mSocketIOServer) Serve(address string) error {
    server, err := socketio.NewServer(nil)
    if err != nil {
        debugging.Logf("%v", err)
        return err
    }
    ln, err := net.Listen("tcp", address)
    if err != nil {
        debugging.Logf("%v", err)
        return err
    }
    s.server = server
    s.closer = ln
    s.postServe(ln)
    return err
}

func (s *mSocketIOServer) handleKCPConn(conn *websocket.Conn) {
    client := &mWebsocketClient{
        conn: conn,
        opts: s.opts,
    }
    client.init()
    s.wg.Add(1)
    defer func() {
        s.wg.Done()
        s.opts.Handler.OnClose(client)
    }()

    client.Serve()
}

func (s *mSocketIOServer) checkClient() {
    var (
        checkList []*mWebsocketClient
        deathList []*mWebsocketClient
    )
    if s.opts == nil || s.opts.Handler == nil {
        return
    }
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
func (s *mSocketIOServer) onNetStat(aType int, n, l uint64) {
    if aType == 0 {
        atomic.AddUint64(&s.rp, n)
        atomic.AddUint64(&s.rx, l)
    } else if aType == 1 {
        atomic.AddUint64(&s.tp, n)
        atomic.AddUint64(&s.tx, l)
    }
}

func (s *mSocketIOServer) postServe(ln net.Listener) {
    server := s.server
    server.OnConnect("/", func(conn socketio.Conn) error {
        conn.SetContext("")
        client := &mSocketIOClient{
            conn:conn,
            server:s,
            opts:s.opts,
        }
        client.init()
        s.opts.Handler.Mapping(conn, client)
        go client.Serve()
        return nil
    })
    //server.OnEvent("/", "notice", func(conn socketio.Conn, msg string) {
    //    debugging.Println("notice:", msg)
    //    conn.Emit("reply", "have "+msg)
    //})
    server.OnEvent("/", "text", func(conn socketio.Conn, msg string) []byte {
        conn.SetContext(msg)
        data, err := base64.StdEncoding.DecodeString(msg)
        if err != nil {
            debugging.Log("socket.io data format error:", err)
            return nil
        }
        if cli := s.opts.Handler.GetClientByKey(conn); cli != nil {
            if cc, ok := cli.(*mSocketIOClient); ok {
                m := &mSocketIOClientMessage{
                    data:data,
                    reply:make(chan []byte),
                }
                cc.msgChan <- m
                reply := <-m.reply
                return reply
                //return string(reply)
            }
        }
        return nil
    })
    server.OnError("/", func(conn socketio.Conn, e error) {
        if cli := s.opts.Handler.GetClientByKey(conn); cli != nil {
            s.opts.Handler.OnError(cli, e)
        }
    })
    server.OnDisconnect("/", func(conn socketio.Conn, reason string) {
        if cli := s.opts.Handler.GetClientByKey(conn); cli != nil {
            s.opts.Handler.OnClose(cli)
            s.opts.Handler.Mapping(conn, nil)
        }
    })
    go server.Serve()
    defer server.Close()
    defer ln.Close()

    http.Handle("/socket.io/", server)
    http.Handle("/", http.FileServer(http.Dir("./asset")))

    if s.opts.TLSEnable {
        if err := http.ServeTLS(ln, server, s.opts.TLSCert, s.opts.TLSKey); err != nil {
            debugging.Log(err)
        }
    } else {
        if err := http.Serve(ln, server); err != nil {
            debugging.Log(err)
        }
    }
}
