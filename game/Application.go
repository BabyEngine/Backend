package game

import (
    "fmt"
    "github.com/BabyEngine/Backend/Debug"
    "github.com/BabyEngine/Backend/events"
    "github.com/BabyEngine/Backend/networking"
    "github.com/DGHeroin/golua/lua"
    "sync"
    "time"
)

func NewApp() *Application {
    ga := &Application{}
    return ga
}

type Application struct {
    L       *lua.State
    apiMap  map[string]func(state *lua.State) int
    servers map[interface{}]networking.ClientHandler
    serverM sync.RWMutex
}

func (app *Application) Init(L *lua.State) {
    app.L = L
    app.servers = make(map[interface{}]networking.ClientHandler)
    app.apiMap = make(map[string]func(state *lua.State) int)
    app.apiMap["Invoke"] = gInvoke
    app.apiMap["Exit"] = gExit
    app.apiMap["StartNetServer"] = gStartNetServer
    app.apiMap["StopNetServer"] = gStopNetServer
    app.apiMap["BindNetServer"] = gBindNetServer
    app.apiMap["SendNetData"] = gSendNetData
    app.apiMap["CloseNetClient"] = gCloseNetClient
    app.apiMap["RedirectNetClient"] = gRedirectNetClient
    app.apiMap["AddUpdateFunc"] = gAddUpdateFunc
    app.apiMap["ApplicationSetFPS"] = gApplicationSetFPS
    //app.apiMap["print"] = gPrint
    app.apiMap["AppLog"] = gAppLog

    L.CreateTable(0, 1)
    L.SetGlobal("BabyEngine")
    L.GetGlobal("BabyEngine")
    // KV 表
    L.PushString("KV")
    {
        // 创建子表
        L.CreateTable(0, 1)
        // OpenFunc
        L.PushString("Open")
        L.PushGoFunction(gKVOpen)
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

    // 导出接口
    for k, v := range app.apiMap {
        L.Register(k, v)
    }

    L.PushGoStruct(app)
    L.SetGlobal("AppContext")

    initLuaCode := `
package.path=package.path .. ';./framework/?.lua'
require 'framework.init'
AddUpdateFunc(function()
    LooperManager.UpdateFunc()
    LooperManager.FixedUpdateFunc()
    LooperManager.LateUpdateFunc()
end)
__log_trace__ = false
local _print = print
print=function(...)
    local msg = ''

    for k,v in pairs({...}) do
        msg = msg .. tostring(v) .. ' ' 
    end
    if __log_trace__ then
        msg = msg .. '\n' .. debug.traceback()
    end
    AppLog(msg)
end

`
    if err := L.DoString(initLuaCode); err != nil {
        fmt.Println(err)
    }
}

func (app *Application) Start() {
    L := app.L
    // on update
    runStartTime := time.Now()  // 程序开始时间
    fpsCounter := 0             // FPS 计数器
    lastTickTime := time.Now()  // 上一帧时间
    tickStartTime := time.Now() // 帧时间
    lastFPS := 0                // 帧数
    updateTimer := func() {
        // FPS
        fpsCounter++
        // timer
        elapsedTime := float64(time.Now().Sub(runStartTime).Milliseconds()) / 1000.0
        deltaTime := float64(time.Now().Sub(lastTickTime).Milliseconds()) / 1000.0
        // make Time info
        luaCode := fmt.Sprintf("Time=Time or {}\nTime.time=%f\nTime.deltaTime=%f\nTime.FPS=%d\n", elapsedTime, deltaTime, lastFPS)
        // reset data
        lastTickTime = time.Now()
        if time.Now().Sub(tickStartTime) >= time.Second {
            lastFPS = fpsCounter
            fpsCounter = 0
            tickStartTime = time.Now()
        }
        L.SetTop(0)
        if err := L.DoString(luaCode); err != nil {
            Debug.Log(err, "\n", luaCode)
        }
    }
    updateTimer()
    tickCounter := uint64(0)
    events.DefaultEventSystem.OnUpdateFunc = func() {
        tickCounter++
        // set Time.time
        updateTimer()
        // invoke all callbacks
        refs := events.DefaultEventSystem.AllRef()
        for _, ref := range refs {
            L.RawGeti(lua.LUA_REGISTRYINDEX, ref)
            if L.IsFunction(-1) {
                L.Call(0, 0)
            }
        }
    }
    // Main Loop
    for events.DefaultEventSystem.IsRunning() {
        events.DefaultEventSystem.Update()
    }
}
func (app *Application) Stop() {
    events.DefaultEventSystem.Stop()
    events.DefaultEventSystem.Reset()
    for _, s := range app.servers {
        s.Stop()
    }
}

func (app *Application) SetNetServer(key interface{}, value networking.ClientHandler)  {
    app.serverM.Lock()
    if key != nil && value != nil { // add
        app.servers[key] = value
    } else { // del?
        delete(app.servers, key)
    }
    app.serverM.Unlock()
}

func (app *Application) GetNetServer(key interface{}) networking.ClientHandler {
    app.serverM.RLock()
    defer app.serverM.RUnlock()
    if e, ok := app.servers[key];ok {
        return e
    }
    return nil
}