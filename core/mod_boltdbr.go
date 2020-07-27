package core

import (
    "github.com/BabyEngine/Backend/logger"
    "github.com/BabyEngine/Backend/rpc"
    "github.com/DGHeroin/golua/lua"
)

func initModBoltDBR(L *lua.State) {
    L.GetGlobal("BabyEngine")
    L.PushString("BoltDBR")
    {
        // 创建子表
        L.CreateTable(0, 1)

        L.PushString("Create")
        L.PushGoFunction(gBoltDBRCreate)
        L.SetTable(-3)

        L.PushString("Get")
        L.PushGoFunction(gBoltDBRGet)
        L.SetTable(-3)

        L.PushString("Set")
        L.PushGoFunction(gBoltDBRSet)
        L.SetTable(-3)

        L.PushString("Delete")
        L.PushGoFunction(gBoltDBRDelete)
        L.SetTable(-3)
    }
    L.SetTable(-3)
}
func gBoltDBRCreate(L *lua.State) int {
    token := L.ToString(1)
    address := L.ToString(2)
    aFunc := L.Ref(lua.LUA_REGISTRYINDEX) // last one is function

    client := rpc.NewBoltDBR(token)
    app := GetApplication(L)
    client.Connect(address, func(err error) {
        app.eventSys.OnMainThread(func() {
            L.RawGeti(lua.LUA_REGISTRYINDEX, aFunc)
            if L.Type(-1) == lua.LUA_TFUNCTION {
                defer L.Unref(lua.LUA_REGISTRYINDEX, aFunc)
                if err != nil {
                    L.PushString(err.Error())
                    if err := L.Call(1, 0); err != nil {
                        logger.WarnIff(EnableDebug,"boltdbr create error:%s", err)
                        return
                    }
                } else {
                    if err := L.Call(0, 0); err != nil {
                        logger.WarnIff(EnableDebug,"boltdbr create error:%s", err)
                        return
                    }
                }
            }
        })
    })
    L.PushGoStruct(client)
    return 1
}
func gBoltDBRGet(L *lua.State) int {
    ptr := L.ToGoStruct(1)
    bucketName := L.ToString(2)
    key := L.ToString(3)
    aFunc := L.Ref(lua.LUA_REGISTRYINDEX) // last one is function
    app := GetApplication(L)

    client, ok := ptr.(*rpc.BoltDBRClient)
    if !ok {
        L.PushString("can not convert to go struct")
        return 1
    }
    client.Get(bucketName, key, func(s1 string, s2 string) {
        app.eventSys.OnMainThread(func() {
            L.RawGeti(lua.LUA_REGISTRYINDEX, aFunc)
            if L.Type(-1) == lua.LUA_TFUNCTION {
                L.PushString(s1)
                L.PushString(s2)
                if err := L.Call(2, 0); err != nil {
                    logger.WarnIff(EnableDebug,"boltdbr get error:%s", err)
                    return
                }
            }
        })
    })
    return 0
}

func gBoltDBRSet(L *lua.State) int {
    ptr := L.ToGoStruct(1)
    bucketName := L.ToString(2)
    key := L.ToString(3)
    value := L.ToString(4)
    aFunc := L.Ref(lua.LUA_REGISTRYINDEX) // last one is function
    app := GetApplication(L)

    client, ok := ptr.(*rpc.BoltDBRClient)
    if !ok {
        L.PushString("can not convert to go struct")
        return 1
    }
    client.Set(bucketName, key, value, func(s1 string) {
        app.eventSys.OnMainThread(func() {
            L.RawGeti(lua.LUA_REGISTRYINDEX, aFunc)
            if L.Type(-1) == lua.LUA_TFUNCTION {
                L.PushString(s1)
                if err := L.Call(1, 0); err != nil {
                    logger.WarnIff(EnableDebug,"boltdbr set error:%s", err)
                    return
                }
            }
        })
    })
    return 0
}

func gBoltDBRDelete(L *lua.State) int {
    ptr := L.ToGoStruct(1)
    bucketName := L.ToString(2)
    key := L.ToString(3)
    aFunc := L.Ref(lua.LUA_REGISTRYINDEX) // last one is function
    app := GetApplication(L)

    client, ok := ptr.(*rpc.BoltDBRClient)
    if !ok {
        L.PushString("can not convert to go struct")
        return 1
    }
    client.Delete(bucketName, key, func(s1 string) {
        app.eventSys.OnMainThread(func() {
            L.RawGeti(lua.LUA_REGISTRYINDEX, aFunc)
            if L.Type(-1) == lua.LUA_TFUNCTION {
                L.PushString(s1)
                if err := L.Call(1, 0); err != nil {
                    logger.WarnIff(EnableDebug,"boltdbr delete error:%s", err)
                    return
                }
            }
        })
    })
    return 0
}
