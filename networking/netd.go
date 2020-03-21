package networking

import (
    "github.com/BabyEngine/UnityConnector/common"
    "github.com/DGHeroin/golua/lua"

)

//var (
//    serverMap = make(map[interface{}]*mKCPGameServerHandler)
//    serverM   sync.RWMutex
//)

func StartNetServer(L *lua.State, netType string, address string, tag string) common.ClientHandler {
    switch netType {
    case "kcp":
        return newKCP(L, address, tag)
    }
    return nil
}

func BindNetServerFunc(L *lua.State, p interface{}, name string, ref int)  {
    s := p.(*mKCPGameServerHandler)
    L.GetGlobal("A")
    s.BindFunc(name, ref)
}

func SendNetData(L *lua.State, p interface{}, cliId int64, data []byte)  {
    s := p.(*mKCPGameServerHandler)
    if s == nil {
        return
    }
    s.SendClientData(cliId, data)
}

func CloseClient(L *lua.State, p interface{}, cliId int64)  {
    s := p.(*mKCPGameServerHandler)
    if s == nil {
        return
    }
    s.CloseClient(cliId)
}

func SendNetRawData(L *lua.State, p interface{}, cliId int64, op common.OpCode, data []byte)  {
    s := p.(*mKCPGameServerHandler)
    if s == nil {
        return
    }
    s.SendClientRawData(cliId, op, data)
}