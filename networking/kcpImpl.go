package networking

import (
    "context"
    "fmt"
    "github.com/BabyEngine/Backend/Debug"
    "github.com/BabyEngine/Backend/events"
    "github.com/DGHeroin/golua/lua"
    "sync"
)

type mKCPGameServerHandler struct {
    ctx        context.Context
    cancel     func()
    L          *lua.State
    refNew     int
    refClose   int
    refData    int
    refError   int
    refRequest int
    clients    map[int64]Client
}

var (
    enableDebug = true
)

func (m *mKCPGameServerHandler) Init() {
    m.ctx, m.cancel = context.WithCancel(context.Background())
    m.clients = make(map[int64]Client)
}

func (m *mKCPGameServerHandler) OnNew(client Client) {
    Debug.LogIff(enableDebug, "OnNew:%v", client)
    events.DefaultEventSystem.OnMainThread(func() {
        m.clients[client.Id()] = client
        L := m.L
        L.RawGeti(lua.LUA_REGISTRYINDEX, m.refNew)
        if L.Type(-1) == lua.LUA_TFUNCTION {
            L.PushInteger(client.Id())
            L.Call(1, 0)
        }
    })
}

func (m *mKCPGameServerHandler) OnData(client Client, data []byte) {
    Debug.LogIff(enableDebug, "OnData:%v %v", client, data)
    if data == nil || len(data) == 0 {
        return
    }
    events.DefaultEventSystem.OnMainThread(func() {
        L := m.L
        L.RawGeti(lua.LUA_REGISTRYINDEX, m.refData)
        if L.Type(-1) == lua.LUA_TFUNCTION {
            L.PushInteger(client.Id())
            L.PushBytes(data)

            L.Call(2, 0)
        }
    })
}

func (m *mKCPGameServerHandler) OnClose(client Client) {
    Debug.LogIff(enableDebug, "OnClose:%v", client)
    m.CloseClient(client.Id())
    events.DefaultEventSystem.OnMainThread(func() {
        L := m.L
        L.RawGeti(lua.LUA_REGISTRYINDEX, m.refClose)
        if L.Type(-1) == lua.LUA_TFUNCTION {
            L.PushInteger(client.Id())
            L.Call(1, 0)
        }
    })
}

func (m *mKCPGameServerHandler) OnError(client Client, err error) {
    Debug.LogIff(enableDebug, "OnError:%v %v", client, err)
    events.DefaultEventSystem.OnMainThread(func() {
        L := m.L
        L.RawGeti(lua.LUA_REGISTRYINDEX, m.refError)
        if L.Type(-1) == lua.LUA_TFUNCTION {
            L.PushInteger(client.Id())
            L.PushString(fmt.Sprint(err))
            L.Call(2, 0)
        }
    })

}
func (m *mKCPGameServerHandler) OnRequest(client Client, data []byte) []byte {
    Debug.LogIff(enableDebug, "OnRequest:%v %v", client, data)
    var (
        wg     sync.WaitGroup
        result []byte
    )
    if data == nil || len(data) == 0 {
        return nil
    }

    wg.Add(1)
    events.DefaultEventSystem.OnMainThread(func() {
        L := m.L
        L.RawGeti(lua.LUA_REGISTRYINDEX, m.refRequest)
        if L.Type(-1) == lua.LUA_TFUNCTION {
            L.PushInteger(client.Id())
            L.PushBytes(data)
            L.Call(2, 1)
            result = L.ToBytes(1)
        }
        wg.Done()
    })
    wg.Wait()
    //return []byte("收到你的请求啦")
    return result
}

func (h *mKCPGameServerHandler) Stop() {
    refs := []int{h.refNew, h.refClose, h.refError, h.refData, h.refRequest}
    for _, ref := range refs {
        if ref != 0 {
            h.L.Unref(lua.LUA_REGISTRYINDEX, ref)
        }
    }
    h.cancel()
}

func (h *mKCPGameServerHandler) BindFunc(name string, ref int) {
    switch name {
    case "new":
        h.refNew = ref
    case "close":
        h.refClose = ref
    case "data":
        h.refData = ref
    case "error":
        h.refError = ref
    case "request":
        h.refRequest = ref
    }
}

func (h *mKCPGameServerHandler) SendClientData(clientId int64, data []byte) {
    if cli, ok := h.clients[clientId]; ok {
        cli.SendData(data)
    }
}
func (h *mKCPGameServerHandler) SendClientRawData(clientId int64, op OpCode, data []byte) {
    Debug.Logf("重定向 %v %v %s", clientId, op, data)
    if cli, ok := h.clients[clientId]; ok {
        cli.SendRaw(op, data)
    }
}

func (h *mKCPGameServerHandler) CloseClient(clientId int64) {
    if cli, ok := h.clients[clientId]; ok {
        cli.Close()
    }
}

func newKCP(L *lua.State, address string, tag string) ClientHandler {
    // 主服务器
    h := &mKCPGameServerHandler{}
    h.L = L
    h.Init()

    go func() {
        if err := Listen(
            WithType("kcp"),
            WithTag(tag),
            WithAddress(address),
            WithContext(h.ctx),
            WithHandler(h)); err != nil {
        }
    }()

    return h
}

func Listen(options ...OptionFunc) error {
    opts := &Options{}
    for _, cb := range options {
        cb(opts)
    }
    if opts.Ctx == nil {
        opts.Ctx = context.TODO()
    }
    switch opts.Type {
    case "kcp":
        return mKCPListenAndServe(opts)
    }
    return ErrorOptionsInvalid
}

