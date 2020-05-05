package game

import (
    "bufio"
    "bytes"
    "fmt"
    "github.com/BabyEngine/Backend/logger"
    "github.com/BabyEngine/Backend/events"
    "github.com/BabyEngine/Backend/networking"
    "github.com/DGHeroin/golua/lua"
    "os"
    "strings"
    "sync"
    "time"
)

func NewApp() *Application {
    ga := &Application{}
    return ga
}

type Application struct {
    L        *lua.State
    apiMap   map[string]func(state *lua.State) int
    servers  map[interface{}]networking.ClientHandler
    serverM  sync.RWMutex
    eventSys *events.EventSystem
}

func (app *Application) Init(L *lua.State) {
    app.L = L
    app.servers = make(map[interface{}]networking.ClientHandler)
    app.apiMap = make(map[string]func(state *lua.State) int)
    app.eventSys = events.NewEventSystem()
    L.PushGoStruct(app)
    L.SetGlobal("ApplicationContext")
    // 创建全局 BabyEngine 表
    L.CreateTable(0, 1)
    L.SetGlobal("BabyEngine")
    // init App
    initModApp(L)
    // Logger
    initModLogger(L)
    // KV 表
    initModKV(L)
    // Net 表
    initModNet(L)
    // Redis
    initModRedis(L)
    // Cipher
    initModCipher(L)
    // RPC
    initModRPC(L)
    // Cron
    initModCron(L)
    //
    injectArgs(L)

    // 导出接口
    for k, v := range app.apiMap {
        L.Register(k, v)
    }

    L.PushGoStruct(app)
    L.SetGlobal("AppContext")

    initLuaCode := `
package.path=package.path .. ';./framework/?.lua'
pcall = pcall or unsafe_pcall
pcall(require,'framework.init')
BabyEngine.App.AddUpdateFunc(function()
    if LooperManager then
        LooperManager.UpdateFunc()
        LooperManager.FixedUpdateFunc()
        LooperManager.LateUpdateFunc()
    end
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
    BabyEngine.Logger.Debug(msg)
end

`
    if err := L.DoString(initLuaCode); err != nil {
        fmt.Println(err)
    }
}

func injectArgs(L *lua.State) {
    L.PushString("Args")
    {
        // 创建子表
        L.CreateTable(0, 1)
        for k, v := range os.Args {
            L.PushInteger(int64(k+1))
            L.PushString(v)
            L.SetTable(-3)
        }
    }
    L.SetTable(-3)
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
            logger.Debug(err, "\n", luaCode)
        }
    }
    updateTimer()
    tickCounter := uint64(0)
    app.eventSys.OnUpdateFunc = func() {
        tickCounter++
        // set Time.time
        updateTimer()
        // invoke all callbacks
        refs := app.eventSys.AllRef()
        for _, ref := range refs {
            L.RawGeti(lua.LUA_REGISTRYINDEX, ref)
            if L.IsFunction(-1) {
                if err := L.Call(0, 0); err != nil {
                    logger.Debug(err)
                }
            }
        }
    }
    // runtime cmd
    go func() {
        reader := bufio.NewReader(os.Stdin)
        cmd := bytes.NewBufferString("")
        var wgCmd sync.WaitGroup
        for app.eventSys.IsRunning() {
            if len(cmd.String()) == 0 {
                fmt.Print("$: ")
            } else {
                fmt.Print(cmd.String())
            }

            text, _ := reader.ReadString('\n')
            text = strings.TrimSpace(text)
            if len(text) == 0 {
                continue
            }

            if strings.HasSuffix(text, "\\") {
                cmd.WriteString(strings.TrimSuffix(text, "\\") + "\n")
                continue
            }
            cmd.WriteString(text)
            wgCmd.Add(1)
            go func(cmdStr string) {
                app.eventSys.OnMainThread(func() {
                    if err := L.DoString(cmdStr); err !=nil {
                        logger.Debugf("%v", err)
                    }
                    wgCmd.Done()
                })
            }(cmd.String())
            wgCmd.Wait()
            cmd.Reset()
        }
    }()

    // Main Loop
    for app.eventSys.IsRunning() {
        app.eventSys.Update()
    }
}
func (app *Application) Stop() {
    app.eventSys.Stop()
    app.eventSys.Reset()
    for _, s := range app.servers {
        s.Stop()
    }
}

func (app *Application) SetNetServer(key interface{}, value networking.ClientHandler) {
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
    if e, ok := app.servers[key]; ok {
        return e
    }
    return nil
}
