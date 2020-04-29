package game

import "github.com/DGHeroin/golua/lua"

func initModLogger(L *lua.State) {
    L.GetGlobal("BabyEngine")
    L.PushString("Logger")
    {
        // 创建子表
        L.CreateTable(0, 1)

        L.PushString("Debug")
        L.PushGoFunction(gLoggerDebug)
        L.SetTable(-3)

        L.PushString("Warn")
        L.PushGoFunction(gLoggerWarn)
        L.SetTable(-3)

        L.PushString("Error")
        L.PushGoFunction(gLoggerError)
        L.SetTable(-3)

    }
    L.SetTable(-3)
}