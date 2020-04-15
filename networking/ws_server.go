package networking

import (
    "fmt"
    "github.com/BabyEngine/Backend/Debug"
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
    clients  map[int64]*mWebsocketClient
    clientM  sync.RWMutex
    clientId int64
    ws       *websocket.Upgrader
}

func (s *mWebsocketServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    c, err := s.ws.Upgrade(w, r, nil)
    if err != nil {
        Debug.Logf("upgrade:%v", err)
        return
    }
    s.handleKCPConn(c)
}

func (s *mWebsocketServer) Init() {
    s.clients = make(map[int64]*mWebsocketClient)
    s.ws = &websocket.Upgrader{
        CheckOrigin: func(r *http.Request) bool {
            return true
        },
    }
}

func (s *mWebsocketServer) Serve(address string) error {
    err := http.ListenAndServe(address, s)
    if err != nil {
        Debug.Logf("%v", err)
    }
    return err
}

func (s *mWebsocketServer) handleKCPConn(conn *websocket.Conn) {
    client := &mWebsocketClient{
        conn:   conn,
        opts:   s.opts,
        server: s,
    }
    s.wg.Add(1)
    defer func() {
        s.wg.Done()
        s.RemoveClient(client)
        s.opts.Handler.OnClose(client)
    }()

    s.AddClient(client)
    client.Serve()
}

func (s *mWebsocketServer) AddClient(client *mWebsocketClient) {
    s.clientM.Lock()

    for {
        s.clientId++
        if _, exist := s.clients[client.id]; !exist {
            client.id = s.clientId
            s.clients[client.id] = client
            break
        }
    }
    s.clientM.Unlock()
}

func (s *mWebsocketServer) RemoveClient(client *mWebsocketClient) {
    s.clientM.Lock()
    if _, ok := s.clients[client.id]; ok {
        delete(s.clients, client.id)
    }
    s.clientM.Unlock()
}

func (s *mWebsocketServer) checkClient() {
    var (
        checkList []*mWebsocketClient
        deathList []*mWebsocketClient
    )
    s.clientM.RLock()
    for _, cli := range s.clients {
        checkList = append(checkList, cli)
    }
    s.clientM.RUnlock()

    for _, cli := range checkList {
        if !cli.IsAlive() {
            deathList = append(deathList, cli)
        }
    }

    if len(deathList) > 0 {
        for _, cli := range deathList {
            fmt.Println("移除客户端", cli)
            cli.Stop()
            s.RemoveClient(cli)
        }
    }
    totalCount := len(s.clients)
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
