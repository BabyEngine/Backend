package game

import (
    "context"
    "fmt"
    "github.com/BabyEngine/Backend/logger"
    "github.com/BabyEngine/Backend/networking"
    "github.com/DGHeroin/golua/lua"
    "sync"
    "sync/atomic"
)

func initModNet(L *lua.State) {
    L.GetGlobal("BabyEngine")
    L.PushString("Net")
    {
        // 创建子表
        L.CreateTable(0, 1)

        L.PushString("Start")
        L.PushGoFunction(gStartNetServer)
        L.SetTable(-3)

        L.PushString("Stop")
        L.PushGoFunction(gStopNetServer)
        L.SetTable(-3)

        L.PushString("Bind")
        L.PushGoFunction(gBindNetServer)
        L.SetTable(-3)

        L.PushString("Send")
        L.PushGoFunction(gSendNetData)
        L.SetTable(-3)

        L.PushString("Close")
        L.PushGoFunction(gCloseNetClient)
        L.SetTable(-3)

        L.PushString("Redirect")
        L.PushGoFunction(gRedirectNetClient)
        L.SetTable(-3)

        L.PushString("SendRaw")
        L.PushGoFunction(gNetClientSendRaw)
        L.SetTable(-3)

        L.PushString("RunCmd")
        L.PushGoFunction(gRunCmd)
        L.SetTable(-3)
    }
    L.SetTable(-3)
}

func StartNetServer(L *lua.State, netType string, address string, flags map[string]string) networking.ClientHandler {
    app := GetApplication(L)
    tag := flags["tag"]
    isRawMode := flags["raw"] == "true"
    switch netType {
    case "kcp":
        h := &MessageServerHandler{}
        h.L = L
        h.Init(app)
        // 配置 Server
        go func() {
            if err := networking.ListenAndServe(
                networking.WithType(netType),
                networking.WithTag(tag),
                networking.WithAddress(address),
                networking.WithContext(h.ctx),
                networking.WithRawMode(isRawMode),
                networking.WithArgs(flags),
                networking.WithHandler(h)); err != nil {
            }
        }()

        return h
    case "ws":
        h := &MessageServerHandler{}
        h.L = L
        h.Init(app)
        // 配置 Server
        go func() {
            if err := networking.ListenAndServe(
                networking.WithType(netType),
                networking.WithTag(tag),
                networking.WithAddress(address),
                networking.WithContext(h.ctx),
                networking.WithRawMode(isRawMode),
                networking.WithTLS(flags["ssl_key"], flags["ssl_cert"]),
                networking.WithArgs(flags),
                networking.WithHandler(h)); err != nil {
            }
        }()

        return h
    case "http":
        h := &MessageServerHandler{}
        h.L = L
        h.Init(app)
        // 配置 Server
        go func() {
            if err := networking.ListenAndServe(
                networking.WithType(netType),
                networking.WithTag(tag),
                networking.WithAddress(address),
                networking.WithContext(h.ctx),
                networking.WithRawMode(isRawMode),
                networking.WithTLS(flags["ssl_key"], flags["ssl_cert"]),
                networking.WithArgs(flags),
                networking.WithHandler(h)); err != nil {
                logger.Debug(err)
            }
        }()
        return h
    case "socket.io":
        h := &MessageServerHandler{}
        h.L = L
        h.Init(app)
        // 配置 Server
        go func() {
            if err := networking.ListenAndServe(
                networking.WithType(netType),
                networking.WithTag(tag),
                networking.WithAddress(address),
                networking.WithContext(h.ctx),
                networking.WithRawMode(isRawMode),
                networking.WithTLS(flags["ssl_key"], flags["ssl_cert"]),
                networking.WithArgs(flags),
                networking.WithHandler(h)); err != nil {
                logger.Debug(err)
            }
        }()
        return h
    }
    return nil
}

func BindNetServerFunc(L *lua.State, p interface{}, name string, ref int) {
    s := p.(*MessageServerHandler)
    //L.GetGlobal("A")
    s.BindFunc(name, ref)
}

func SendNetData(L *lua.State, p interface{}, cliId int64, data []byte) {
    s := p.(*MessageServerHandler)
    if s == nil {
        return
    }
    s.SendClientData(cliId, data)
}

func CloseClient(L *lua.State, p interface{}, cliId int64) {
    s := p.(*MessageServerHandler)
    if s == nil {
        return
    }
    s.CloseClient(cliId)
}

//func SendNetRawData(L *lua.State, p interface{}, cliId int64, op networking.OpCode, data []byte) {
//    s := p.(*MessageServerHandler)
//    if s == nil {
//        return
//    }
//    s.SendClientRawData(cliId, op, data)
//}

func SendNetRawDataEvent(L *lua.State, p interface{}, cliId int64, event string, op networking.OpCode, data []byte) {
    s := p.(*MessageServerHandler)
    if s == nil {
        return
    }
    s.SendClientRawDataEvent(cliId, event, op, data)
}

func NetRunCmd(L *lua.State, p interface{}, cliId int64, cmd string, args []string) string {
    s := p.(*MessageServerHandler)
    if s == nil {
        return ""
    }
    if c := s.GetClient(cliId); c != nil {
        return c.RunCmd(cmd, args)
    }
    return ""
}

type MessageServerHandler struct {
    app        *Application
    ctx        context.Context
    cancel     func()
    L          *lua.State
    refNew     int
    refClose   int
    refData    int
    refError   int
    refRequest int
    clients    map[int64]networking.Client
    clientM    sync.RWMutex
    clientId   int64
    clientMap map[interface{}]networking.Client
}

func (h *MessageServerHandler) GetAllClient() []networking.Client {
    h.clientM.RLock()
    result := make([]networking.Client, len(h.clients))
    i := 0
    for _, v := range h.clients {
        result[i] = v
        i++
    }
    h.clientM.RUnlock()
    return result
}

func (h *MessageServerHandler) Mapping(key interface{}, client networking.Client) {
    h.clientM.Lock()
    if client == nil {
        delete(h.clientMap, key)
    } else {
        h.clientMap[key] = client
    }
    h.clientM.Unlock()
}
func (h *MessageServerHandler) GetClientByKey(key interface{}) networking.Client  {
    h.clientM.RLock()
    cli, _ := h.clientMap[key]
    h.clientM.RUnlock()
    return cli
}
var (
    EnableDebug = false
)

func (h *MessageServerHandler) Init(app *Application) {
    h.app = app
    h.ctx, h.cancel = context.WithCancel(context.Background())
    h.clients = make(map[int64]networking.Client)
    h.clientMap = make(map[interface{}]networking.Client)
}

func (m *MessageServerHandler) OnNew(client networking.Client) {
    id := atomic.AddInt64(&m.clientId, 1)
    client.SetId(id)
    logger.DebugIff(EnableDebug, "OnNew:%v", client)
    m.app.eventSys.
        OnMainThread(func() {
            m.clients[client.Id()] = client
            L := m.L
            L.RawGeti(lua.LUA_REGISTRYINDEX, m.refNew)
            if L.Type(-1) == lua.LUA_TFUNCTION {
                L.PushInteger(client.Id())
                if err := L.Call(1, 0); err != nil {
                    logger.Debug(err)
                }
            }
        })
}

func (h *MessageServerHandler) OnData(client networking.Client, data []byte) {
    logger.DebugIff(EnableDebug, "OnData:%v %v", client, data)
    if data == nil || len(data) == 0 {
        return
    }
    h.app.eventSys.
        OnMainThread(func() {
            L := h.L
            L.RawGeti(lua.LUA_REGISTRYINDEX, h.refData)
            if L.Type(-1) == lua.LUA_TFUNCTION {
                L.PushInteger(client.Id())
                L.PushBytes(data)

                if err := L.Call(2, 0); err != nil {
                    logger.Debug(err)
                }
            }
        })
}

func (h *MessageServerHandler) OnClose(client networking.Client) {
    logger.DebugIff(EnableDebug, "OnClose:%v", client)
    h.CloseClient(client.Id())
    h.app.eventSys.
        OnMainThread(func() {
            L := h.L
            L.RawGeti(lua.LUA_REGISTRYINDEX, h.refClose)
            if L.Type(-1) == lua.LUA_TFUNCTION {
                L.PushInteger(client.Id())
                if err := L.Call(1, 0); err != nil {
                    logger.Debug(err)
                }
            }
        })
}

func (h *MessageServerHandler) OnError(client networking.Client, err error) {
    logger.DebugIff(EnableDebug, "OnError:%v %v", client, err)
    h.app.eventSys.
        OnMainThread(func() {
            L := h.L
            L.RawGeti(lua.LUA_REGISTRYINDEX, h.refError)
            if L.Type(-1) == lua.LUA_TFUNCTION {
                L.PushInteger(client.Id())
                L.PushString(fmt.Sprint(err))
                if err := L.Call(2, 0); err != nil {
                    logger.Debug(err)
                }
            }
        })

}
func (h *MessageServerHandler) OnRequest(client networking.Client, data []byte) []byte {
    logger.DebugIff(EnableDebug, "OnRequest:%v %v", client, data)
    var (
        wg     sync.WaitGroup
        result []byte
    )
    if data == nil || len(data) == 0 {
        return nil
    }
    wg.Add(1)
    respFunc := func(L *lua.State) int {
        result = L.ToBytes(1)
        wg.Done()
        return 0
    }
    h.app.eventSys.
        OnMainThread(func() {
            L := h.L
            L.RawGeti(lua.LUA_REGISTRYINDEX, h.refRequest)
            if L.Type(-1) == lua.LUA_TFUNCTION {
                L.PushInteger(client.Id())
                L.PushBytes(data)
                L.PushGoFunction(respFunc)
                if err := L.Call(3, 0); err != nil {
                    logger.Debug(err)
                    wg.Done()
                }
            }
        })
    wg.Wait()
    return result
}

func (h *MessageServerHandler) Stop() {
    refs := []int{h.refNew, h.refClose, h.refError, h.refData, h.refRequest}
    for _, ref := range refs {
        if ref != 0 {
            h.L.Unref(lua.LUA_REGISTRYINDEX, ref)
        }
    }
    h.cancel()
}

func (h *MessageServerHandler) BindFunc(name string, ref int) {
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

func (h *MessageServerHandler) SendClientData(clientId int64, data []byte) {
    h.clientM.RLock()
    cli, ok := h.clients[clientId]
    h.clientM.RUnlock()
    if ok {
        if err := cli.SendData(data); err != nil {
            logger.Debug(err)
        }
    }
}
//func (h *MessageServerHandler) SendClientRawData(clientId int64, op networking.OpCode, data []byte) {
//    h.clientM.RLock()
//    cli, ok := h.clients[clientId]
//    h.clientM.RUnlock()
//    if ok {
//        if err := cli.SendRaw(op, data); err != nil {
//            logger.Debug(err)
//        }
//    }
//}

func (h *MessageServerHandler) SendClientRawDataEvent(clientId int64, e string, op networking.OpCode, data []byte) {
    h.clientM.RLock()
    cli, ok := h.clients[clientId]
    h.clientM.RUnlock()
    if ok {
        if err := cli.SendRawEvent(e, op, data); err != nil {
            logger.Debug(err)
        }
    }
}

func (h *MessageServerHandler) CloseClient(clientId int64) {
    h.clientM.RLock()
    cli, ok := h.clients[clientId]
    h.clientM.RUnlock()
    if ok {
        cli.Close()
        h.clientM.Lock()
        delete(h.clients, clientId)
        h.clientM.Unlock()
    }
}

func (h *MessageServerHandler) GetClient(id int64) networking.Client {
    h.clientM.RLock()
    defer h.clientM.RUnlock()
    if cli, ok := h.clients[id]; ok {
        return cli
    }
    return nil
}
