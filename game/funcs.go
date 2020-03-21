package game

import (
    "bytes"
    "fmt"
    "github.com/BabyEngine/Backend/Debug"
    "github.com/BabyEngine/Backend/events"
    "github.com/BabyEngine/Backend/kv"
    "github.com/BabyEngine/Backend/networking"
    "github.com/DGHeroin/golua/lua"
)
// 往主线程塞一个回调
func gInvoke(L *lua.State) int {
    if L.GetTop() == 1 { // next tick
        if L.Type(-1) == lua.LUA_TFUNCTION {
            var ref int = 0
            ref = L.Ref(lua.LUA_REGISTRYINDEX)
            events.DefaultEventSystem.OnMainThread(func() {
                L.RawGeti(lua.LUA_REGISTRYINDEX, ref)
                if L.Type(-1) == lua.LUA_TFUNCTION {
                    L.Call(0, 0)
                }
                L.Unref(lua.LUA_REGISTRYINDEX, ref)
            })
        }
    } else if L.GetTop() == 2 { // delay
        delay := L.ToNumber(1)
        var ref int = 0
        ref = L.Ref(lua.LUA_REGISTRYINDEX)
        events.DefaultEventSystem.OnMainThreadDelay(delay, func() {
            L.RawGeti(lua.LUA_REGISTRYINDEX, ref)
            if L.Type(-1) == lua.LUA_TFUNCTION {
                L.Call(0, 0)
            }
            L.Unref(lua.LUA_REGISTRYINDEX, ref)
        })
    }
    return 0
}
// 启动服务器
func gStartNetServer(L *lua.State) int {
    netType := L.ToString(1)
    addr := L.ToString(2)
    serverTag := "Backend"
    if L.Type(3) == lua.LUA_TSTRING {
        serverTag = L.ToString(3)
    }
    L.GetGlobal("AppContext")
    _app := L.ToGoStruct(-1)
    app:= _app.(*Application)

    server := networking.StartNetServer(L, netType, addr, serverTag)
    app.SetNetServer(app, server)

    L.PushGoStruct(server)
    return 1
}
// 停止服务器
func gStopNetServer(L*lua.State) int {
    ptr := L.ToGoStruct(1)
    L.GetGlobal("AppContext")
    _app := L.ToGoStruct(-1)
    app:= _app.(*Application)
    app.SetNetServer(ptr, nil)
    return 0
}
// 绑定回调
func gBindNetServer(L*lua.State) int {
    ptr := L.ToGoStruct(1)
    name := L.ToString(2)
    aFunc := L.Ref(lua.LUA_REGISTRYINDEX) // last one is function

    L.GetGlobal("AppContext")
    _app := L.ToGoStruct(-1)
    app:= _app.(*Application)
    app.GetNetServer(ptr)

    networking.BindNetServerFunc(L, ptr, name, aFunc)
    return 0
}
// 给客户端发消息
func gSendNetData(L*lua.State) int {
    ptr := L.ToGoStruct(1)
    cliId := L.ToInteger(2)
    data := L.ToBytes(3)
    networking.SendNetData(L, ptr, int64(cliId), data)
    return 0
}
// 关闭客户端
func gCloseNetClient(L*lua.State) int {
    ptr := L.ToGoStruct(1)
    cliId := L.ToInteger(2)
    networking.CloseClient(L, ptr, int64(cliId))
    return 0
}
// 重定向客户端到其他服务器
func gRedirectNetClient(L*lua.State) int  {
    ptr := L.ToGoStruct(1)
    cliId := L.ToInteger(2)
    data := L.ToBytes(3)
    networking.SendNetRawData(L, ptr, int64(cliId), networking.OPCODE_TURN, data)
    return 0
}
// 退出app
func gExit(L *lua.State) int {
    events.DefaultEventSystem.Stop()
    return 0
}
// 往update加入一个回调
func gAddUpdateFunc(L *lua.State) int {
    ref := L.Ref(lua.LUA_REGISTRYINDEX)
    if L.IsFunction(-1) {
        events.DefaultEventSystem.AddRef(ref)
        L.PushInteger(int64(ref))
        return 1
    }
    return 0
}
// 设置update帧率
func gApplicationSetFPS(L *lua.State) int {
    fps := L.ToInteger(1)
    events.DefaultEventSystem.SetFPS(fps)
    return 0
}
// 覆盖lua print函数
func gPrint(L *lua.State) int {

    nargs := L.GetTop()
    buf := bytes.NewBufferString(L.StackTraceString())
    for i := 1; i <= nargs; i++ {
        msg := ""
        switch L.Type(i) {
        case lua.LUA_TNIL:
            msg = "nil"
        case lua.LUA_TNUMBER:
            msg = fmt.Sprintf("%v", L.ToNumber(i))
        case lua.LUA_TBOOLEAN:
            msg = fmt.Sprintf("%v", L.ToBoolean(i))
        case lua.LUA_TSTRING:
            msg = L.ToString(i)
        case lua.LUA_TTABLE:
            msg = fmt.Sprintf("table:0x%x", L.ToPointer(i))
        case lua.LUA_TFUNCTION:
            msg = fmt.Sprintf("function:0x%x", L.ToPointer(i))
        case lua.LUA_TUSERDATA:
            msg = fmt.Sprintf("udata:0x%x", L.ToPointer(i))
        case lua.LUA_TTHREAD:
            msg = fmt.Sprintf("thread:0x%x", L.ToPointer(i))
        case lua.LUA_TLIGHTUSERDATA:
            msg = fmt.Sprintf("ludatea:0x%x", L.ToUserdata(i))
        }
        if i < nargs - 1 {
            buf.WriteString(fmt.Sprintf("%v\t", msg))
        } else {
            buf.WriteString(fmt.Sprintf("%v", msg))
        }
    }
    Debug.Log(buf)
    return 0
}

// 打印日志
func gAppLog(L *lua.State) int  {
    msg := L.CheckString(-1)
    Debug.Log(msg)
    return 0
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
// 写db
func gKVPut(L *lua.State) int {
    ptr := L.ToGoStruct(1)
    bucketName := L.ToString(2)
    key := L.ToString(3)
    value := L.ToBytes(4)
    if db, ok := ptr.(*kv.DB); ok {
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
    ptr := L.ToGoStruct(1)
    bucketName := L.ToString(2)
    key := L.ToString(3)
    cbRef := L.Ref(lua.LUA_REGISTRYINDEX)
    if db, ok := ptr.(*kv.DB); ok {
        db.View(bucketName, key, func(i []byte, err error) {
            events.DefaultEventSystem.OnMainThread(func() {
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
                    L.Call(2, 0)
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
    if db, ok := ptr.(*kv.DB); ok {
        db.RemoveValue(bucketName, key)
    }
    return 0
}
func gKVRemoveBucket(L *lua.State) int {
    ptr := L.ToGoStruct(1)
    bucketName := L.ToString(2)
    if db, ok := ptr.(*kv.DB); ok {
        db.RemoveBucket(bucketName)
    }
    return 0
}