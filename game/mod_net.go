package game

import "github.com/DGHeroin/golua/lua"

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

        L.PushString("Redirect")
        L.PushGoFunction(gRedirectNetClient)
        L.SetTable(-3)
    }
    L.SetTable(-3)
}