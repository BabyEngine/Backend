package game

import (
    "github.com/BabyEngine/Backend/Debug"
    "github.com/BabyEngine/Backend/kv"
    "github.com/DGHeroin/golua/lua"
    "time"
)

func initModRedis(L *lua.State) {
    L.GetGlobal("BabyEngine")
    L.PushString("Redis")
    {
        // 创建子表
        L.CreateTable(0, 1)

        L.PushString("Open")
        L.PushGoFunction(gRedisOpen)
        L.SetTable(-3)

        L.PushString("Set")
        L.PushGoFunction(gRedisSet)
        L.SetTable(-3)

        L.PushString("Get")
        L.PushGoFunction(gRedisGet)
        L.SetTable(-3)
    }
    L.SetTable(-3)
}
// open redis
func gRedisOpen(L *lua.State) int {
    config := GetTableMap(L)
    if r, err := kv.OpenRedis(config["address"], config["password"], 0); err ==nil {
        L.PushGoStruct(r)
        return 1
    }
    return 0
}

func gRedisGet(L *lua.State) int {
    var r * kv.RedisDB
    ptr := L.ToGoStruct(1)
    if _r, ok := ptr.(*kv.RedisDB); !ok {
        return 0
    } else {
        r = _r
    }
    key := L.ToString(2)
    if val, err := r.Get(key); err != nil {
        Debug.Logf("%v", err)
        return 0
    } else {
        L.PushString(val)
        return 1
    }
}
func gRedisSet(L *lua.State) int {
    var r * kv.RedisDB
    ptr := L.ToGoStruct(1)
    if _r, ok := ptr.(*kv.RedisDB); !ok {
        return 0
    } else {
        r = _r
    }
    key := L.ToString(2)
    value := L.ToString(3)
    TTLms := float64(0)
    if L.Type(4) == lua.LUA_TNUMBER {
        TTLms = L.ToNumber(4)
    }

    duration := TTLms * float64(time.Millisecond)
    r.Set(key, value, time.Duration(duration))
    return 0
}
