package game

import (
    "fmt"
    "github.com/DGHeroin/golua/lua"
)

// 获取栈顶的table, 并转换成go map
func GetTableMap(L *lua.State) (result map[string]string) {
    result = make(map[string]string)
    if L.Type(-1) != lua.LUA_TTABLE {
        return
    }
    L.PushNil()
    for L.Next(-2) != 0 {
        key := L.ToString(-2)
        val := L.ToString(-1)
        L.Pop(1)
        result[key] = val
    }
    L.Pop(1)
    return
}

func CheckTypes(L *lua.State,types ...lua.LuaValType) error {
    for i, t := range types {
        idx := i + 1
        if t == lua.LUA_TNIL {
            continue
        } else {
            L.CheckType(idx, t)
            if L.Type(idx) != t {
                //return fmt.Errorf("args(%d) want:%v got: %v", i, L.CheckType)
                return fmt.Errorf("args error")
            }
        }
    }
    return nil
}