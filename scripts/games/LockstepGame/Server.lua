local pb = require "pb"
local protoc = require "protobuf/protoc"

if not pb.type"GameMessage" then
    protoc:load(PBDefine)
end

NetService = {}
-- 游戏消息服务器
function NetService.startGameServer( addr )
    local server = net.NewKCPBinaryServer(addr, "帧同步服务器")
    NetService.server = server
    server.OnNew = function ( cli )
        NewPlayerActor(cli)
    end

    server.OnClose = function ( cli )
        local player = GetPlayer(cli)
        if player then
            player.Release()
            print('玩家释放', cli)
        end
    end

    server.OnData = function ( cli, data )
        print('收到数据', table.tostring(cli), tostring(data))
    end

    local handlers = {}
    -- 列出房间
    handlers['listRoom'] = function(req, cli, data, respFunc)
        local datas = ListRoom()
        local rooms = {}
        for k,v in pairs(datas) do
            table.insert(rooms, {id = v.id, players=v.GetPlayers()})
        end
        print('列出房间', table.tostring(rooms))
        local data = pb.encode('ListRoomResp', {rooms=rooms})
        respFunc(data)
    end
    -- 创建房间
    handlers['createRoom'] = function(req, cli, data, respFunc)
        local room = NewRoom()
        print('create room', room.id)
        local data = pb.encode('RespMessage', {code=0, msg='创建成功['..tostring(room.id)..']'})
        respFunc(data)
    end
    -- 加入房间
    handlers['joinRoom'] = function(req, cli, data, respFunc)
        local info = pb.decode('JoinRoom', req.data)
        local room = JoinRoom(info.roomId, info.playerId, cli)
        print('加入房间', table.tostring(info))
        local data
        if room then
            data = pb.encode('RespMessage', {code=0, msg='加入成功['..tostring(room.id)..']'})
        else
            data = pb.encode('RespMessage', {code=-1, msg='房间不存在'})
        end
        print('加入结果', table.tostring(room))
        respFunc(data)
    end
    -- 离开房间
    handlers['leaveRoom'] = function(req, cli, data, respFunc)
        local player = LeaveRoom(cli)
        local data
        if player then
            data = pb.encode('RespMessage', {code=0, msg='离开成功'})
        else
            data = pb.encode('RespMessage', {code=-1, msg='房间不存在'})
        end
        print('离开结果', table.tostring(room))
        respFunc(data)
    end
    -- 游戏消息
    handlers['playing'] = function(req, cli, data, respFunc)
        local msg = pb.decode('PlayingMessage', req.data)
        local player = GetPlayer(cli)
        if player then
            player.OnPlayingMessage(msg)
        end
        local data = pb.encode('RespMessage', {code=0, msg='ok'})
        respFunc(data)
    end

    server.OnRequest = function ( cli, data, respFunc )
        -- print('请求数据', table.tostring(cli), tostring(data))
        local req = pb.decode('GameMessage', data)
        local func = handlers[req.action]
        if func then
            func(req, cli, data, respFunc)
            return
        end
        print('不支持请求', req.action)
        local data = pb.encode('RespMessage', {code=-1, msg='不支持请求['..req.action..']'})
        respFunc(data)
        do return end
        if req.action == 'listRoom' then
            -- local datas = ListRoom()
            -- local rooms = {}
            -- for k,v in pairs(datas) do
            --     table.insert(rooms, {id = v.id, players=v.GetPlayers()})
            -- end
            -- print('列出房间', table.tostring(rooms))
            -- local data = pb.encode('ListRoomResp', {rooms=rooms})
            -- respFunc(data)
        elseif req.action == 'createRoom' then
            -- local room = NewRoom()
            -- print('create room', room.id)
            -- local data = pb.encode('RespMessage', {code=0, msg='创建成功['..tostring(room.id)..']'})
            -- respFunc(data)
        elseif req.action == 'joinRoom' then
            -- local info = pb.decode('JoinRoom', req.data)
            -- local room = JoinRoom(info.roomId, info.playerId, cli)

            -- local player = players[cli]
            -- if player then
            --     player.room = room
            -- end

            -- print('加入房间', table.tostring(info))
            -- local data
            -- if room then
            --     data = pb.encode('RespMessage', {code=0, msg='加入成功['..tostring(room.id)..']'})
            -- else
            --     data = pb.encode('RespMessage', {code=-1, msg='房间不存在'})
            -- end
            -- respFunc(data)
        elseif req.action == 'playing' then
            -- local msg = pb.decode('PlayingMessage', req.data)
            -- local player = players[cli]
            -- if player then
            --     player.OnPlayingMessage(msg)
            -- end
            -- local data = pb.encode('RespMessage', {code=0, msg='ok'})
            -- respFunc(data)
        else
            print('不支持请求', req.action)
            local data = pb.encode('RespMessage', {code=-1, msg='不支持请求['..req.action..']'})
            respFunc(data)
        end

    end

    server.Start()
end
