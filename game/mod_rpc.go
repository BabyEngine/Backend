package game

import (
    "errors"
    "github.com/BabyEngine/Backend/logger"
    "github.com/BabyEngine/Backend/rpc"
    "github.com/DGHeroin/golua/lua"
    "sync"
)

func initModRPC(L *lua.State) {
    L.GetGlobal("BabyEngine")
    L.PushString("JSONRPC")
    {
        // 创建子表
        L.CreateTable(0, 1)

        L.PushString("NewClient")
        L.PushGoFunction(gJSONRPCNewClient)
        L.SetTable(-3)

        L.PushString("Call")
        L.PushGoFunction(gJSONRPCClientCall)
        L.SetTable(-3)

        L.PushString("StopClient")
        L.PushGoFunction(gJSONRPCStopClient)
        L.SetTable(-3)

        L.PushString("NewServer")
        L.PushGoFunction(gJSONRPCNewServer)
        L.SetTable(-3)

        L.PushString("StopServer")
        L.PushGoFunction(gJSONRPCStopServer)
        L.SetTable(-3)

    }
    L.SetTable(-3)
}

type rpcServer struct {
    server *rpc.JSONRPCServer
}

func gJSONRPCNewServer(L *lua.State) int {
    address := L.ToString(1)
    ref := L.Ref(lua.LUA_REGISTRYINDEX)
    if L.Type(-1) == lua.LUA_TFUNCTION {
        logger.Debug("gJSONRPCNewServer args error")
        return 0
    }
    srv := rpcServer{}
    app := GetApplication(L)

    srv.server = rpc.NewJSONRPCServer(func(request rpc.Request, reply *rpc.Reply) error {
        var err error
        wg := sync.WaitGroup{}
        wg.Add(1)
        app.eventSys.OnMainThread(func() {
            L.RawGeti(lua.LUA_REGISTRYINDEX, ref)
            if L.Type(-1) == lua.LUA_TFUNCTION {
                L.PushString(request.Action)
                L.PushBytes(request.Data)
                L.PushGoFunction(func(L *lua.State) int {
                    code := L.ToInteger(1)
                    rep := L.ToBytes(2)
                    *reply = rpc.Reply{code, rep}
                    wg.Done()
                    return 0
                })

                if err := L.Call(3, 0); err != nil {
                    logger.Debug("rpc invoke error:", err)
                    wg.Done()
                }

            } else {
                wg.Done()
                err = errors.New("handle func error")
            }
        })
        wg.Wait()
        return err
    })
    go func() {
        if err := srv.server.ListenServe(address); err != nil {
            logger.Debug(err)
        }
    }()
    L.PushGoStruct(srv)
    return 1
}

func gJSONRPCStopServer(L *lua.State) int {
    ptr := L.ToGoStruct(1)
    if ptr, ok := ptr.(*rpcServer); ok {
        if ptr.server != nil {
            if err := ptr.server.Stop(); err != nil {
                L.PushString(err.Error())
                return 1
            }
        }
    }
    return 0
}

type rpcClient struct {
    client *rpc.JSONRPCClient
}

func gJSONRPCNewClient(L *lua.State) int {
    address := L.ToString(1)
    cli, err := rpc.NewJSONRPCClient(address)
    if err != nil {
        L.PushNil()
        L.PushString(err.Error())
        return 2
    }
    c := &rpcClient{client: cli}
    L.PushGoStruct(c)
    L.PushNil()
    return 2
}

func gJSONRPCStopClient(L *lua.State) int {
    ptr := L.ToGoStruct(1)
    if c, ok := ptr.(*rpcClient); ok {
        c.client.Stop()
    }
    return 0
}

func gJSONRPCClientCall(L *lua.State) int {
    ptr := L.ToGoStruct(1)
    action := L.ToString(2)
    data := L.ToBytes(3)
    cb := L.Ref(lua.LUA_REGISTRYINDEX)
    app := GetApplication(L)
    if c, ok := ptr.(*rpcClient); ok {
        go func() {
            r, err := c.client.Call(action, data)
            app.eventSys.OnMainThread(func() {
                L.RawGeti(lua.LUA_REGISTRYINDEX, cb)
                if L.Type(-1) == lua.LUA_TFUNCTION {
                    if err != nil {
                        L.PushInteger(-1)
                        L.PushNil()
                        L.PushString(err.Error())
                    } else {
                        L.PushInteger(int64(r.Code))
                        L.PushBytes(r.Data)
                        L.PushNil()
                    }
                    L.Call(3, 0)
                }
            })
        }()
    } else {
        logger.Debug("args err")
    }
    return 0
}
