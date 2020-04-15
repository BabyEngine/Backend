package game

import (
    "github.com/DGHeroin/golua/lua"
    "os"
)

func initModApp(L *lua.State) {
    L.GetGlobal("BabyEngine")
    L.PushString("App")
    {
        // 创建子表
        L.CreateTable(0, 1)

        L.PushString("Invoke")
        L.PushGoFunction(gInvoke)
        L.SetTable(-3)

        L.PushString("Exit")
        L.PushGoFunction(gExit)
        L.SetTable(-3)

        L.PushString("AddUpdateFunc")
        L.PushGoFunction(gAddUpdateFunc)
        L.SetTable(-3)

        L.PushString("SetFPS")
        L.PushGoFunction(gSetFPS)
        L.SetTable(-3)

        L.PushString("Log")
        L.PushGoFunction(gAppLog)
        L.SetTable(-3)

        L.PushString("GetEnv")
        L.PushGoFunction(gGetEnv)
        L.SetTable(-3)

    }
    L.SetTable(-3)
}

func gGetEnv(L *lua.State) int {
    key := L.ToString(1)
    val := os.Getenv(key)
    L.PushString(val)
    return 1
}
