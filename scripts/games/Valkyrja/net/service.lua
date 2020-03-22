NetService = {}
-- 区域
local zones = {
    lobby = {}
}

local UnAuthConn = {} -- 未认证玩家
local AuthConn   = {} -- 已经认证玩家
local lastCheckTime = -1 -- 时间标识

-- 游戏消息服务器
function NetService.startGameServer( addr )
    local server = net.NewKCPBinaryServer(addr, "ValkyrjaGameServer")
    NetService.server = server
    server.OnNew = function ( cli )
        -- 新连接进来的玩家放到未认证列表, 10秒内没有认证则踢出服务器
        local player = NewPlayerActor(cli)
        player.connectTime = Time.time
        UnAuthConn[cli] = player
    end

    server.OnClose = function ( cli )
        local player = nil
        player = UnAuthConn[cli]
        if player then
            UnAuthConn[cli] = nil
            player.Release()
        end
        player = AuthConn[cli]
        if player then
            AuthConn[cli] = nil
            player.Release()
        end
    end

    server.OnData = function ( cli, data )
        -- print('收到数据', table.tostring(cli), tostring(data))
        if UnAuthConn[cli] then
            UnAuthConn[cli].OnMessage(data)
            return
        end
        if AuthConn[cli] then
            AuthConn[cli].OnMessage(data)
            return
        end
    end

    server.OnRequest = function ( cli, data, respFunc )
        if UnAuthConn[cli] then
            UnAuthConn[cli].OnRequest(data, respFunc)
        end
        if AuthConn[cli] then
            AuthConn[cli].OnRequest(data, respFunc)
        end
    end

    server.Start()
end

function NetService.OnAuth ( cli, player )
    if UnAuthConn[cli] then
        UnAuthConn[cli] = nil
    end
    AuthConn[cli] = player
end

local function onUpdateCheck ()
    -- 每 3 秒做一次检查
    if Time.time - lastCheckTime < 3 then return end
    lastCheckTime = Time.time

    local keepList = {}
    for k,player in pairs(UnAuthConn) do
        if not player.isAuth and Time.time - player.connectTime > 10 then
            player.conn.Close()
            print('玩家过期了')
        else
            keepList[k] = player -- 还未过期, 留着以后验证
        end
    end
    UnAuthConn = keepList
end


function NetService.status(  )
    local c1 = 0
    local c2 = 0
    for k,v in pairs(UnAuthConn) do
        c1 = c1 + 1
    end
    for k,v in pairs(AuthConn) do
        c2 = c2 + 1
    end
    return string.format("auth: %d | unauth: %d", c2, c1, countP)
end

Looper.AddUpdate(onUpdateCheck)