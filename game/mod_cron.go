package game

import (
    "github.com/BabyEngine/Backend/jobs"
    "github.com/BabyEngine/Backend/logger"
    "github.com/DGHeroin/golua/lua"
)

func initModCron(L *lua.State) {
    L.GetGlobal("BabyEngine")
    L.PushString("Cron")
    {
        // 创建子表
        L.CreateTable(0, 1)

        L.PushString("New")
        L.PushGoFunction(gCronNew)
        L.SetTable(-3)

        L.PushString("Add")
        L.PushGoFunction(gCronAdd)
        L.SetTable(-3)

        L.PushString("Remove")
        L.PushGoFunction(gCronRemove)
        L.SetTable(-3)

        L.PushString("Stop")
        L.PushGoFunction(gCronStop)
        L.SetTable(-3)
    }
    L.SetTable(-3)
}

func gCronNew(L *lua.State) int {
    c := jobs.NewCron()
    L.PushGoStruct(c)
    return 1
}

func gCronAdd(L *lua.State) int {
    // check args
    if err := CheckTypes(L, lua.LUA_TUSERDATA, lua.LUA_TSTRING, lua.LUA_TFUNCTION); err != nil {
        L.PushNil()
        L.PushString(err.Error())
        return 2
    }

    c := L.ToGoStruct(1)
    spec := L.ToString(2)
    ref := L.Ref(lua.LUA_REGISTRYINDEX)
    app := GetApplication(L)
    if cc, ok := c.(*jobs.Context); ok {
        id, err := cc.Add(spec, func() {
            app.eventSys.OnMainThread(func() {
                L.RawGeti(lua.LUA_REGISTRYINDEX, ref)
                if err := L.Call(0, 0); err != nil {
                    logger.Warn(err)
                }
            })
        })
        if err != nil {
            L.PushNil()
            L.PushString(err.Error())
            return 2
        }
        L.PushInteger(int64(id))
        L.PushNil()
        return 2
    }
    L.PushNil()
    L.PushString("args error")
    return 2
}

func gCronRemove(L *lua.State) int {
    if err := CheckTypes(L, lua.LUA_TUSERDATA, lua.LUA_TNUMBER); err != nil {
        L.PushString(err.Error())
        return 1
    }
    c := L.ToGoStruct(1)
    id := L.ToInteger(2)
    if cc, ok := c.(*jobs.Context); ok {
        cc.Remove(int(id))
        L.PushNil()
        return 1
    }
    L.PushString("args error")
    return 1
}

func gCronStop(L *lua.State) int {
    if err := CheckTypes(L, lua.LUA_TUSERDATA ); err != nil {
        L.PushString(err.Error())
        return 1
    }
    c := L.ToGoStruct(1)
    if cc, ok := c.(*jobs.Context); ok {
        cc.Stop()
        return 0
    }
    L.PushString("args error")
    return 1
}
