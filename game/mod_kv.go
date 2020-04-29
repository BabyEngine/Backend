package game

import (
    "fmt"
    "github.com/BabyEngine/Backend/logger"
    "github.com/BabyEngine/Backend/kv"
    "github.com/DGHeroin/golua/lua"
)

func initModKV(L *lua.State) {
    L.GetGlobal("BabyEngine")
    L.PushString("KV")
    {
        // 创建子表
        L.CreateTable(0, 1)

        L.PushString("Open")
        L.PushGoFunction(gKVOpen)
        L.SetTable(-3)

        L.PushString("Close")
        L.PushGoFunction(gKVClose)
        L.SetTable(-3)

        L.PushString("Get")
        L.PushGoFunction(gKVGet)
        L.SetTable(-3)

        L.PushString("Put")
        L.PushGoFunction(gKVPut)
        L.SetTable(-3)

        L.PushString("RemoveValue")
        L.PushGoFunction(gKVRemoveValue)
        L.SetTable(-3)

        L.PushString("RemoveBucket")
        L.PushGoFunction(gKVRemoveBucket)
        L.SetTable(-3)
    }
    L.SetTable(-3)
}


// 打开kvdb
func gKVOpen(L *lua.State) int {
    path := L.CheckString(-1)
    if db, err := kv.OpenKVDB(path); err != nil {
        return 0
    } else {
        L.PushGoStruct(db)
        return 1
    }
}
//关闭
func gKVClose(L *lua.State) int {
    ptr := L.ToGoStruct(1)
    if db, ok := ptr.(*kv.KVDB); ok {
        db.Close()
    }
    return 0
}

// 写db
func gKVPut(L *lua.State) int {
    ptr := L.ToGoStruct(1)
    bucketName := L.ToString(2)
    key := L.ToString(3)
    value := L.ToBytes(4)
    if db, ok := ptr.(*kv.KVDB); ok {
        if err := db.Update(bucketName, key, value); err == nil {
            L.PushNil()
            L.PushBoolean(true)
            return 2
        } else {
            L.PushString(fmt.Sprint(err))
            L.PushBoolean(false)
            return 2
        }
    }
    L.PushNil()
    L.PushBoolean(false)
    return 2
}
// 读db
func gKVGet(L *lua.State) int {
    app := GetApplication(L)
    ptr := L.ToGoStruct(1)
    bucketName := L.ToString(2)
    key := L.ToString(3)
    cbRef := L.Ref(lua.LUA_REGISTRYINDEX)
    if db, ok := ptr.(*kv.KVDB); ok {
        db.View(bucketName, key, func(i []byte, err error) {
            app.eventSys.OnMainThread(func() {
                defer L.Unref(lua.LUA_REGISTRYINDEX, cbRef)
                L.RawGeti(lua.LUA_REGISTRYINDEX, cbRef)
                if L.Type(-1) == lua.LUA_TFUNCTION {
                    if err != nil {
                        L.PushNil()
                        L.PushString(fmt.Sprint(err))
                    } else {
                        if i == nil {
                            L.PushNil()
                        } else {
                            L.PushBytes(i)
                        }
                        L.PushNil()
                    }
                    if err := L.Call(2, 0); err != nil {
                        logger.Debug(err)
                    }
                }
            })
        })
    }
    return 0
}

func gKVRemoveValue(L *lua.State) int {
    ptr := L.ToGoStruct(1)
    bucketName := L.ToString(2)
    key := L.ToString(3)
    if db, ok := ptr.(*kv.KVDB); ok {
        db.RemoveValue(bucketName, key)
    }
    return 0
}
func gKVRemoveBucket(L *lua.State) int {
    ptr := L.ToGoStruct(1)
    bucketName := L.ToString(2)
    if db, ok := ptr.(*kv.KVDB); ok {
        db.RemoveBucket(bucketName)
    }
    return 0
}