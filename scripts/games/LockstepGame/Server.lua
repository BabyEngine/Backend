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
        -- 新连接进来的玩家放到未认证列表, 10秒内没有认证则踢出服务器
        local player = NewPlayerActor(cli)
    end

    server.OnClose = function ( cli )

    end

    server.OnData = function ( cli, data )
        print('收到数据', table.tostring(cli), tostring(data))
    end

    server.OnRequest = function ( cli, data, respFunc )
        -- print('请求数据', table.tostring(cli), tostring(data))
        local req = pb.decode('GameMessage', data)
        if req.action == 'listRoom' then
            local datas = ListRoom()
            local rooms = {}
            for k,v in pairs(datas) do
                table.insert(rooms, {id = v.id, v.GetPlayers()})
            end
            local data = pb.encode('ListRoomResp', {rooms=rooms})
            respFunc(data)
        elseif req.action == 'createRoom' then
            local room = NewRoom()
            local data = pb.encode('RespMessage', {code=0, msg='创建成功['..tostring(room.id)..']'})
            respFunc(data)
        elseif req.action == 'joinRoom' then
            local info = pb.decode('JoinRoom', req.data)
            local room = JoinRoom(info.roomId, info.playerId, cli)
            print('加入房间', table.tostring(info))
            local data
            if room then
                data = pb.encode('RespMessage', {code=0, msg='加入成功['..tostring(room.id)..']'})
            else
                data = pb.encode('RespMessage', {code=-1, msg='房间不存在'})
            end
            respFunc(data)
        else
            print('不支持请求', req.action)
            local data = pb.encode('RespMessage', {code=-1, msg='不支持请求['..req.action..']'})
            respFunc(data)
        end

    end

    server.Start()
end
